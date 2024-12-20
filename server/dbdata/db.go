package dbdata

import (
	"net/http"
	"time"

	"github.com/bjdgyc/anylink/base"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

var (
	xdb *xorm.Engine
)

func GetXdb() *xorm.Engine {
	return xdb
}

func initDb() {
	var err error
	xdb, err = xorm.NewEngine(base.Cfg.DbType, base.Cfg.DbSource)
	// Initialize xorm time zone
	xdb.DatabaseTZ = time.Local
	xdb.TZLocation = time.Local
	if err != nil {
		base.Fatal(err)
	}

	if base.Cfg.ShowSQL {
		xdb.ShowSQL(true)
	}

	// Initialize database
	err = xdb.Sync2(&User{}, &Setting{}, &Group{}, &IpMap{}, &AccessAudit{}, &Policy{}, &StatsNetwork{}, &StatsCpu{}, &StatsMem{}, &StatsOnline{}, &UserActLog{})
	if err != nil {
		base.Fatal(err)
	}
	// fmt.Println("s1=============", err)
}

func initData() {
	var (
		err error
	)

	// Determine whether to use it for the first time
	install := &SettingInstall{}
	err = SettingGet(install)

	if err == nil && install.Installed {
		// Already installed
		return
	}

	// An error occurred
	if err != ErrNotFound {
		base.Fatal(err)
	}

	err = addInitData()
	if err != nil {
		base.Fatal(err)
	}

}

func addInitData() error {
	var (
		err error
	)

	sess := xdb.NewSession()
	defer sess.Close()

	err = sess.Begin()
	if err != nil {
		return err
	}

	// Setting smtp
	smtp := &SettingSmtp{
		Host:       "127.0.0.1",
		Port:       25,
		From:       "vpn@xx.com",
		Encryption: "None",
	}
	err = SettingSessAdd(sess, smtp)
	if err != nil {
		return err
	}

	// Setting audit log
	auditLog := SettingGetAuditLogDefault()
	err = SettingSessAdd(sess, auditLog)
	if err != nil {
		return err
	}

	// Setting dns provider
	provider := &SettingLetsEncrypt{
		Domain:      "vpn.xxx.com",
		Legomail:    "legomail",
		Name:        "aliyun",
		Renew:       false,
		DNSProvider: DNSProvider{},
	}
	err = SettingSessAdd(sess, provider)
	if err != nil {
		return err
	}
	// Lego user
	legouser := &LegoUserData{}
	err = SettingSessAdd(sess, legouser)
	if err != nil {
		return err
	}
	// Setting other
	other := &SettingOther{
		LinkAddr:    "vpn.xx.com",
		Banner:      "You have connected to the company network, please use it in accordance with company regulations. \nPlease do not perform non-work downloading and video activities! ",
		Homecode:    http.StatusOK,
		Homeindex:   "AnyLink is an enterprise-level remote office sslvpn software that can support multiple people using it online at the same time.",
		AccountMail: accountMail,
	}
	err = SettingSessAdd(sess, other)
	if err != nil {
		return err
	}

	// Install
	install := &SettingInstall{Installed: true}
	err = SettingSessAdd(sess, install)
	if err != nil {
		return err
	}

	err = sess.Commit()
	if err != nil {
		return err
	}

	g1 := Group{
		Name:         "all",
		AllowLan:     true,
		ClientDns:    []ValData{{Val: "114.114.114.114"}},
		RouteInclude: []ValData{{Val: ALL}},
		Status:       1,
	}
	err = SetGroup(&g1)
	if err != nil {
		return err
	}

	g2 := Group{
		Name:         "ops",
		AllowLan:     true,
		ClientDns:    []ValData{{Val: "114.114.114.114"}},
		RouteInclude: []ValData{{Val: "10.0.0.0/8"}},
		Status:       1,
	}
	err = SetGroup(&g2)
	if err != nil {
		return err
	}

	return nil
}

func CheckErrNotFound(err error) bool {
	return err == ErrNotFound
}

// base64 images
// User dynamic code (please keep it properly):<br/>
// <img src="{{.OtpImgBase64}}"/><br/>
const accountMail = `<p>Hello:</p>
<p>&nbsp;&nbsp;your{{.Issuer}} account has been reviewed and opened.</p>
<p>
    Login address: <b>{{.LinkAddr}}</b> <br/>
    User group: <b>{{.Group}}</b> <br/>
    username: <b>{{.Username}}</b> <br/>
    User PIN code: <b>{{.PinCode}}</b> <br/>
    User expiration time: <b>{{.LimitTime}}</b> <br/>
    {{if .DisableOtp}}
    <!-- nothing -->
    {{else}}
	
    <!-- 
    User dynamic code (expires after 3 days):<br/>
    <img src="{{.OtpImg}}"/><br/>
    -->
    User dynamic code (please keep it properly):<br/>
    <img src="cid:userOtpQr.png" alt="userOtpQr" /><br/>

    {{end}}
</p>
<div>
    Instructions for use:
    <ul>
        <li>Please use OTP software to scan the dynamic code QR code</li>
        <li>Then use anyconnect client to log in</li>
        <li>The login password is [PIN code + dynamic code] (no + sign in the middle)</li>
    </ul>
</div>
<p>
    软件下载地址: https://{{.LinkAddr}}/files/info.txt
</p>`
