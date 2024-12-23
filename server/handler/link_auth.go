package handler

import (
	"bytes"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"text/template"

	"github.com/bjdgyc/anylink/base"
	"github.com/bjdgyc/anylink/dbdata"
	"github.com/bjdgyc/anylink/sessdata"
)

var (
	profileHash = ""
	certHash    = ""
)

func LinkAuth(w http.ResponseWriter, r *http.Request) {
	// TODO Debug information output
	if base.GetLogLevel() == base.LogLevelTrace {
		hd, _ := httputil.DumpRequest(r, true)
		base.Trace("LinkAuth: ", string(hd))
	}
	// Determine anyconnect client
	userAgent := strings.ToLower(r.UserAgent())
	xAggregateAuth := r.Header.Get("X-Aggregate-Auth")
	xTranscendVersion := r.Header.Get("X-Transcend-Version")
	if !((strings.Contains(userAgent, "anyconnect") || strings.Contains(userAgent, "openconnect") || strings.Contains(userAgent, "anylink")) &&
		xAggregateAuth == "1" && xTranscendVersion == "1") {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "error request")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	cr := ClientRequest{}
	err = xml.Unmarshal(body, &cr)
	if err != nil {
		base.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	base.Trace(fmt.Sprintf("%+v \n", cr))
	// setCommonHeader(w)
	if cr.Type == "logout" {
		// Exit and delete session information
		if cr.SessionToken != "" {
			sessdata.DelSessByStoken(cr.SessionToken)
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	if cr.Type == "init" {
		w.WriteHeader(http.StatusOK)
		data := RequestData{Group: cr.GroupSelect, Groups: dbdata.GetGroupNamesNormal()}
		tplRequest(tpl_request, w, data)
		return
	}

	// Login parameter judgment
	if cr.Type != "auth-reply" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// User activity log
	ua := dbdata.UserActLog{
		Username:        cr.Auth.Username,
		GroupName:       cr.GroupSelect,
		RemoteAddr:      r.RemoteAddr,
		Status:          dbdata.UserAuthSuccess,
		DeviceType:      cr.DeviceId.DeviceType,
		PlatformVersion: cr.DeviceId.PlatformVersion,
	}
	// TODO User password verification
	err = dbdata.CheckUser(cr.Auth.Username, cr.Auth.Password, cr.GroupSelect)
	if err != nil {
		base.Warn(err, r.RemoteAddr)
		ua.Info = err.Error()
		ua.Status = dbdata.UserAuthFail
		dbdata.UserActLogIns.Add(ua, userAgent)

		w.WriteHeader(http.StatusOK)
		data := RequestData{Group: cr.GroupSelect, Groups: dbdata.GetGroupNamesNormal(), Error: "Wrong username or password"}
		if base.Cfg.DisplayError {
			data.Error = err.Error()
		}
		tplRequest(tpl_request, w, data)
		return
	}
	dbdata.UserActLogIns.Add(ua, userAgent)
	// if !ok {
	//	w.WriteHeader(http.StatusOK)
	//	data := RequestData{Group: cr.GroupSelect, Groups: base.Cfg.UserGroups, Error: "Please activate the user first"}
	//	tplRequest(tpl_request, w, data)
	//	return
	// }

	// Create new session information
	sess := sessdata.NewSession("")
	sess.Username = cr.Auth.Username
	sess.Group = cr.GroupSelect
	oriMac := cr.MacAddressList.MacAddress
	sess.UniqueIdGlobal = cr.DeviceId.UniqueIdGlobal
	sess.UserAgent = userAgent
	sess.DeviceType = ua.DeviceType
	sess.PlatformVersion = ua.PlatformVersion
	sess.RemoteAddr = r.RemoteAddr
	// Get the client mac address
	sess.UniqueMac = true
	macHw, err := net.ParseMAC(oriMac)
	if err != nil {
		var sum [16]byte
		if sess.UniqueIdGlobal != "" {
			sum = md5.Sum([]byte(sess.UniqueIdGlobal))
		} else {
			sum = md5.Sum([]byte(sess.Token))
			sess.UniqueMac = false
		}
		macHw = sum[0:5] // 5 bytes
		macHw = append([]byte{0x02}, macHw...)
		sess.MacAddr = macHw.String()
	}
	sess.MacHw = macHw
	// Unify the format of macAddr
	sess.MacAddr = macHw.String()

	other := &dbdata.SettingOther{}
	_ = dbdata.SettingGet(other)
	rd := RequestData{SessionId: sess.Sid, SessionToken: sess.Sid + "@" + sess.Token,
		Banner: other.Banner, ProfileName: base.Cfg.ProfileName, ProfileHash: profileHash, CertHash: certHash}
	w.WriteHeader(http.StatusOK)
	tplRequest(tpl_complete, w, rd)
	base.Info("login", cr.Auth.Username, userAgent)
}

const (
	tpl_request = iota
	tpl_complete
)

func tplRequest(typ int, w io.Writer, data RequestData) {
	if typ == tpl_request {
		t, _ := template.New("auth_request").Parse(auth_request)
		_ = t.Execute(w, data)
		return
	}

	if data.Banner != "" {
		buf := new(bytes.Buffer)
		_ = xml.EscapeText(buf, []byte(data.Banner))
		data.Banner = buf.String()
	}

	t, _ := template.New("auth_complete").Parse(auth_complete)
	_ = t.Execute(w, data)
}

// Set output information
type RequestData struct {
	Groups []string
	Group  string
	Error  string

	// complete
	SessionId    string
	SessionToken string
	Banner       string
	ProfileName  string
	ProfileHash  string
	CertHash     string
}

var auth_request = `<?xml version="1.0" encoding="UTF-8"?>
<config-auth client="vpn" type="auth-request" aggregate-auth-version="2">
    <opaque is-for="sg">
        <tunnel-group>{{.Group}}</tunnel-group>
        <group-alias>{{.Group}}</group-alias>
        <aggauth-handle>168179266</aggauth-handle>
        <config-hash>1595829378234</config-hash>
        <auth-method>multiple-cert</auth-method>
        <auth-method>single-sign-on-v2</auth-method>
    </opaque>
    <auth id="main">
        <title>Login</title>
        <message>Please enter your username and password</message>
        <banner></banner>
        {{if .Error}}
        <error id="88" param1="{{.Error}}" param2="">Login failed:  %s</error>
        {{end}}
        <form>
            <input type="text" name="username" label="Username:"></input>
            <input type="password" name="password" label="Password:"></input>
            <select name="group_list" label="GROUP:">
                {{range $v := .Groups}}
                <option {{if eq $v $.Group}} selected="true"{{end}}>{{$v}}</option>
                {{end}}
            </select>
        </form>
    </auth>
</config-auth>
`

var auth_complete = `<?xml version="1.0" encoding="UTF-8"?>
<config-auth client="vpn" type="complete" aggregate-auth-version="2">
    <session-id>{{.SessionId}}</session-id>
    <session-token>{{.SessionToken}}</session-token>
    <auth id="success">
        <banner>{{.Banner}}</banner>
        <message id="0" param1="" param2=""></message>
    </auth>
    <capabilities>
        <crypto-supported>ssl-dhe</crypto-supported>
    </capabilities>
    <config client="vpn" type="private">
        <vpn-base-config>
            <server-cert-hash>{{.CertHash}}</server-cert-hash>
        </vpn-base-config>
        <opaque is-for="vpn-client"></opaque>
        <vpn-profile-manifest>
            <vpn rev="1.0">
                <file type="profile" service-type="user">
                    <uri>/profile_{{.ProfileName}}.xml</uri>
                    <hash type="sha1">{{.ProfileHash}}</hash>
                </file>
            </vpn>
        </vpn-profile-manifest>
    </config>
</config-auth>
`

// var auth_profile = `<?xml version="1.0" encoding="UTF-8"?>
// <AnyConnectProfile xmlns="http://schemas.xmlsoap.org/encoding/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://schemas.xmlsoap.org/encoding/ AnyConnectProfile.xsd">

// 	<ClientInitialization>
// 		<UseStartBeforeLogon UserControllable="false">false</UseStartBeforeLogon>
// 		<StrictCertificateTrust>false</StrictCertificateTrust>
// 		<RestrictPreferenceCaching>false</RestrictPreferenceCaching>
// 		<RestrictTunnelProtocols>IPSec</RestrictTunnelProtocols>
// 		<BypassDownloader>true</BypassDownloader>
// 		<WindowsVPNEstablishment>AllowRemoteUsers</WindowsVPNEstablishment>
// 		<CertEnrollmentPin>pinAllowed</CertEnrollmentPin>
// 		<CertificateMatch>
// 			<KeyUsage>
// 				<MatchKey>Digital_Signature</MatchKey>
// 			</KeyUsage>
// 			<ExtendedKeyUsage>
// 				<ExtendedMatchKey>ClientAuth</ExtendedMatchKey>
// 			</ExtendedKeyUsage>
// 		</CertificateMatch>

// 		<BackupServerList>
// 	            <HostAddress>localhost</HostAddress>
// 		</BackupServerList>
// 	</ClientInitialization>

//	<ServerList>
//		<HostEntry>
//	            <HostName>VPN Server</HostName>
//	            <HostAddress>localhost</HostAddress>
//		</HostEntry>
//	</ServerList>
//
// </AnyConnectProfile>
// `
var ds_domains_xml = `
<?xml version="1.0" encoding="UTF-8"?>
<config-auth client="vpn" type="complete" aggregate-auth-version="2">
    <config client="vpn" type="private">
        <opaque is-for="vpn-client">
            <custom-attr>
            {{if .DsExcludeDomains}}
               <dynamic-split-exclude-domains><![CDATA[{{.DsExcludeDomains}},]]></dynamic-split-exclude-domains>
            {{else if .DsIncludeDomains}}
               <dynamic-split-include-domains><![CDATA[{{.DsIncludeDomains}}]]></dynamic-split-include-domains>
            {{end}}
            </custom-attr>
        </opaque>
    </config>
</config-auth>
`
