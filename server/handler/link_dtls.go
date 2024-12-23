package handler

import (
	"net"
	"time"

	"github.com/bjdgyc/anylink/base"
	"github.com/bjdgyc/anylink/dbdata"
	"github.com/bjdgyc/anylink/pkg/utils"
	"github.com/bjdgyc/anylink/sessdata"
)

func LinkDtls(conn net.Conn, cSess *sessdata.ConnSession) {
	base.Debug("LinkDtls connect ip:", cSess.IpAddr, "user:", cSess.Username, "udp-rip:", conn.RemoteAddr())
	dSess := cSess.NewDtlsConn()
	if dSess == nil {
		// Creation failed, close the link directly
		_ = conn.Close()
		return
	}

	defer func() {
		base.Debug("LinkDtls return", cSess.Username, cSess.IpAddr)
		_ = conn.Close()
		dSess.Close()
	}()

	var (
		err  error
		n    int
		dead = time.Duration(cSess.CstpDpd+5) * time.Second
	)

	go dtlsWrite(conn, dSess, cSess)

	for {
		err = conn.SetReadDeadline(utils.NowSec().Add(dead))
		if err != nil {
			base.Error("SetDeadline: ", cSess.Username, cSess.IpAddr, err)
			return
		}

		pl := getPayload()
		n, err = conn.Read(pl.Data)
		if err != nil {
			base.Warn("read hdata: ", cSess.Username, cSess.IpAddr, err)
			return
		}

		// Current limiting settings
		err = cSess.RateLimit(n, true)
		if err != nil {
			base.Error(err)
		}

		switch pl.Data[0] {
		case 0x07: // KEEPALIVE
			// do nothing
			base.Trace("recv LinkDtls Keepalive", cSess.Username, cSess.IpAddr, conn.RemoteAddr())
		case 0x05: // Disconnect
			cSess.UserLogoutCode = dbdata.UserLogoutClient
			base.Info("DISCONNECT DTLS", cSess.Username, cSess.IpAddr, conn.RemoteAddr())
			return
		case 0x03: // Dpd req
			base.Trace("recv LinkDtls DPD-REQ", cSess.Username, cSess.IpAddr, conn.RemoteAddr(), n)
			pl.PType = 0x04
			// Starting from zero, you can assign values ​​directly
			pl.Data = pl.Data[:n]
			if payloadOutDtls(cSess, dSess, pl) {
				return
			}
		case 0x04:
			base.Trace("recv LinkDtls DPD-RESP", cSess.Username, cSess.IpAddr, conn.RemoteAddr())
		case 0x08: // Decompress
			if cSess.DtlsPickCmp == nil {
				continue
			}
			dst := getByteFull()
			nn, err := cSess.DtlsPickCmp.Uncompress(pl.Data[1:], *dst)
			if err != nil {
				putByte(dst)
				base.Error("dtls decompress error", err, n)
				continue
			}
			pl.Data = append(pl.Data[:1], (*dst)[:nn]...)
			putByte(dst)
			n = nn + 1
			fallthrough
		case 0x00: // DATA
			// Remove header
			// copy(pl.Data, pl.Data[1:n])
			// Update slice length
			// pl.Data = pl.Data[:n-1]
			pl.Data = append(pl.Data[:0], pl.Data[1:n]...)
			if payloadIn(cSess, pl) {
				return
			}
			// Only record the time when correct data is returned
			cSess.LastDataTime.Store(utils.NowSec().Unix())
		}

	}
}

func dtlsWrite(conn net.Conn, dSess *sessdata.DtlsSession, cSess *sessdata.ConnSession) {
	defer func() {
		base.Debug("dtlsWrite return", cSess.Username, cSess.IpAddr)
		_ = conn.Close()
		dSess.Close()
	}()

	var (
		pl *sessdata.Payload
	)

	for {
		// Dtls pushes data first
		select {
		case pl = <-cSess.PayloadOutDtls:
		case <-dSess.CloseChan:
			return
		}

		if pl.LType != sessdata.LTypeIPData {
			continue
		}

		// header = []byte{payload.PType}
		if pl.PType == 0x00 { // Data
			isCompress := false
			if cSess.DtlsPickCmp != nil && len(pl.Data) > base.Cfg.NoCompressLimit {
				dst := getByteFull()
				size, err := cSess.DtlsPickCmp.Compress(pl.Data, (*dst)[1:])
				if err == nil && size < len(pl.Data) {
					(*dst)[0] = 0x08
					pl.Data = append(pl.Data[:0], (*dst)[:size+1]...)
					isCompress = true
				}
				putByte(dst)
			}
			// uncompressed
			if !isCompress {
				// Get data length
				l := len(pl.Data)
				// Expand capacity first +1
				pl.Data = pl.Data[:l+1]
				// Data move back
				copy(pl.Data[1:], pl.Data)
				// Add header information
				pl.Data[0] = pl.PType
			}
		} else {
			// Set header type
			if pl.PType == 0x04 {
				pl.Data[0] = pl.PType
			} else {
				pl.Data = append(pl.Data[:0], pl.PType)
			}
		}
		n, err := conn.Write(pl.Data)
		if err != nil {
			base.Warn("write err", cSess.Username, cSess.IpAddr, err)
			return
		}

		putPayload(pl)

		// Current limiting settings
		err = cSess.RateLimit(n, false)
		if err != nil {
			base.Error(err)
		}
	}
}
