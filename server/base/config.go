package base

const (
	cfgStr = iota
	cfgInt
	cfgBool

	defaultJwt = "abcdef.0123456789.abcdef"
	defaultPwd = "$2a$10$UQ7C.EoPifDeJh6d8.31TeSPQU7hM/NOM2nixmBucJpAuXDQNqNke"
)

type config struct {
	Typ     int
	Name    string
	Short   string
	Usage   string
	ValStr  string
	ValInt  int
	ValBool bool
}

var configs = []config{
	{Typ: cfgStr, Name: "conf", Usage: "config file", ValStr: "./conf/server.toml", Short: "c"},
	{Typ: cfgStr, Name: "profile", Usage: "profile.xml file", ValStr: "./conf/profile.xml"},
	{Typ: cfgStr, Name: "profile_name", Usage: "profile name(Configuration used to distinguish different servers)", ValStr: "anylink"},
	{Typ: cfgStr, Name: "server_addr", Usage: "TCP service listening address(any port)", ValStr: ":443"},
	{Typ: cfgBool, Name: "server_dtls", Usage: "开启DTLS", ValBool: false},
	{Typ: cfgStr, Name: "server_dtls_addr", Usage: "DTLS listening address(any port)", ValStr: ":443"},
	{Typ: cfgStr, Name: "admin_addr", Usage: "Background service listening address", ValStr: ":8800"},
	{Typ: cfgBool, Name: "proxy_protocol", Usage: "TCP proxy protocol", ValBool: false},
	{Typ: cfgStr, Name: "db_type", Usage: "Database type [sqlite3 mysql postgres]", ValStr: "sqlite3"},
	{Typ: cfgStr, Name: "db_source", Usage: "Database source", ValStr: "./conf/anylink.db"},
	{Typ: cfgStr, Name: "cert_file", Usage: "certificate file", ValStr: "./conf/vpn_cert.pem"},
	{Typ: cfgStr, Name: "cert_key", Usage: "Certificate key", ValStr: "./conf/vpn_cert.key"},
	{Typ: cfgStr, Name: "files_path", Usage: "External download file path", ValStr: "./conf/files"},
	{Typ: cfgStr, Name: "log_path", Usage: "Log file path,Default standard output", ValStr: ""},
	{Typ: cfgStr, Name: "log_level", Usage: "Log level [debug info warn error]", ValStr: "debug"},
	{Typ: cfgBool, Name: "http_server_log", Usage: "Enable go standard library http.Server log", ValBool: false},
	{Typ: cfgBool, Name: "pprof", Usage: "Turn on pprof", ValBool: true},
	{Typ: cfgStr, Name: "issuer", Usage: "System name", ValStr: "XX company VPN"},
	{Typ: cfgStr, Name: "admin_user", Usage: "Manage username", ValStr: "admin"},
	{Typ: cfgStr, Name: "admin_pass", Usage: "Manage user passwords", ValStr: defaultPwd},
	{Typ: cfgStr, Name: "admin_otp", Usage: "Manage user otp, Build command ./anylink tool -o", ValStr: ""},
	{Typ: cfgStr, Name: "jwt_secret", Usage: "JWT key", ValStr: defaultJwt},
	{Typ: cfgStr, Name: "link_mode", Usage: "Virtual network type[tun tap macvtap ipvtap]", ValStr: "tun"},
	{Typ: cfgStr, Name: "ipv4_master", Usage: "ipv4 main network card name", ValStr: "eth0"},
	{Typ: cfgStr, Name: "ipv4_cidr", Usage: "IP address network segment", ValStr: "192.168.90.0/24"},
	{Typ: cfgStr, Name: "ipv4_gateway", Usage: "ipvch_guess", ValStr: "192.168.90.1"},
	{Typ: cfgStr, Name: "ipv4_start", Usage: "IPv4 start address", ValStr: "192.168.90.100"},
	{Typ: cfgStr, Name: "ipv4_end", Usage: "IPv4 ends", ValStr: "192.168.90.200"},
	{Typ: cfgStr, Name: "default_group", Usage: "Default user group", ValStr: "one"},
	{Typ: cfgStr, Name: "default_domain", Usage: "Default search domain for client dns", ValStr: ""},

	{Typ: cfgInt, Name: "ip_lease", Usage: "IP lease period(Second)", ValInt: 86400},
	{Typ: cfgInt, Name: "max_client", Usage: "Max user connections", ValInt: 200},
	{Typ: cfgInt, Name: "max_user_client", Usage: "Maximum single user connections", ValInt: 3},
	{Typ: cfgInt, Name: "cstp_keepalive", Usage: "keepalive time(Second)", ValInt: 3},
	{Typ: cfgInt, Name: "cstp_dpd", Usage: "Dead link detection time(Second)", ValInt: 20},
	{Typ: cfgInt, Name: "mobile_keepalive", Usage: "Mobile terminal keepalive connection detection time(Second)", ValInt: 4},
	{Typ: cfgInt, Name: "mobile_dpd", Usage: "Mobile terminal dead link detection time(Second)", ValInt: 60},
	{Typ: cfgInt, Name: "mtu", Usage: "Maximum transmission unit MTU", ValInt: 1460},
	{Typ: cfgInt, Name: "idle_timeout", Usage: "Idle link timeout(Second)-Disconnect the link after timeout, 0 turns off this function", ValInt: 0},
	{Typ: cfgInt, Name: "session_timeout", Usage: "session expiration time(Second)-Used for disconnection and reconnection, 0 will never expire", ValInt: 3600},
	// {Typ: cfgInt, Name: "auth_timeout", Usage: "auth_timeout", ValInt: 0},
	{Typ: cfgInt, Name: "audit_interval", Usage: "Audit deduplication interval(Second),-1 Close", ValInt: 600},

	{Typ: cfgBool, Name: "show_sql", Usage: "Display sql statements for debugging", ValBool: false},
	{Typ: cfgBool, Name: "iptables_nat", Usage: "Whether to automatically add NAT", ValBool: true},
	{Typ: cfgBool, Name: "compression", Usage: "Enable compression", ValBool: false},
	{Typ: cfgInt, Name: "no_compress_limit", Usage: "Below and equal to how many bytes are not compressed", ValInt: 256},

	{Typ: cfgBool, Name: "display_error", Usage: "Client displays detailed error message(Open the online environment with caution)", ValBool: false},
	{Typ: cfgBool, Name: "exclude_export_ip", Usage: "Exclude export ip routing(Export IP is not encrypted for transmission)", ValBool: true},
}

var envs = map[string]string{}
