package dbdata

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/go-ldap/ldap"
)

type AuthLdap struct {
	Addr        string `json:"addr"`
	Tls         bool   `json:"tls"`
	BindName    string `json:"bind_name"`
	BindPwd     string `json:"bind_pwd"`
	BaseDn      string `json:"base_dn"`
	ObjectClass string `json:"object_class"`
	SearchAttr  string `json:"search_attr"`
	MemberOf    string `json:"member_of"`
}

func init() {
	authRegistry["ldap"] = reflect.TypeOf(AuthLdap{})
}

func (auth AuthLdap) checkData(authData map[string]interface{}) error {
	authType := authData["type"].(string)
	bodyBytes, err := json.Marshal(authData[authType])
	if err != nil {
		return errors.New("The LDAP configuration is filled in incorrectly")
	}
	json.Unmarshal(bodyBytes, &auth)
	// Support domain name and IP, port must be filled in
	if !ValidateIpPort(auth.Addr) && !ValidateDomainPort(auth.Addr) {
		return errors.New("LDAP server address(Including port) filled in incorrectly")
	}
	if auth.BindName == "" {
		return errors.New("LDAP administrator DN cannot be empty")
	}
	if auth.BindPwd == "" {
		return errors.New("The LDAP administrator password cannot be empty")
	}
	if auth.BaseDn == "" || !ValidateDN(auth.BaseDn) {
		return errors.New("LThe Base DN of DAP is filled in incorrectly")
	}
	if auth.ObjectClass == "" {
		return errors.New("The user object class of LDAP is filled in incorrectly")
	}
	if auth.SearchAttr == "" {
		return errors.New("LDAP user unique ID cannot be empty")
	}
	if auth.MemberOf != "" && !ValidateDN(auth.MemberOf) {
		return errors.New("The restricted user group of LDAP is filled in incorrectly")
	}
	return nil
}

func (auth AuthLdap) checkUser(name, pwd string, g *Group) error {
	pl := len(pwd)
	if name == "" || pl < 1 {
		return fmt.Errorf("%s %s", name, "Wrong password")
	}
	authType := g.Auth["type"].(string)
	if _, ok := g.Auth[authType]; !ok {
		return fmt.Errorf("%s %s", name, "The ldap value for LDAP does not exist")
	}
	bodyBytes, err := json.Marshal(g.Auth[authType])
	if err != nil {
		return fmt.Errorf("%s %s", name, "LDAP Marshal error occurred")
	}
	err = json.Unmarshal(bodyBytes, &auth)
	if err != nil {
		return fmt.Errorf("%s %s", name, "LDAP Unmarshal error occurred")
	}
	// Check server and port availability
	con, err := net.DialTimeout("tcp", auth.Addr, 3*time.Second)
	if err != nil {
		return fmt.Errorf("%s %s", name, "LDAP server connection exception, Please check the server and port")
	}
	defer con.Close()
	// connectldap
	l, err := ldap.Dial("tcp", auth.Addr)
	if err != nil {
		return fmt.Errorf("LDAP connection failed %s %s", auth.Addr, err.Error())
	}
	defer l.Close()
	if auth.Tls {
		err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return fmt.Errorf("%s LDAP TLS connection failed %s", name, err.Error())
		}
	}
	err = l.Bind(auth.BindName, auth.BindPwd)
	if err != nil {
		return fmt.Errorf("%s LDAP the administrator DN or password is incorrectly filled in %s", name, err.Error())
	}
	if auth.ObjectClass == "" {
		auth.ObjectClass = "person"
	}
	filterAttr := "(objectClass=" + auth.ObjectClass + ")"
	filterAttr += "(" + auth.SearchAttr + "=" + name + ")"
	if auth.MemberOf != "" {
		filterAttr += "(memberOf:=" + auth.MemberOf + ")"
	}
	searchRequest := ldap.NewSearchRequest(
		auth.BaseDn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 3, false,
		fmt.Sprintf("(&%s)", filterAttr),
		[]string{},
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		return fmt.Errorf("%s LDAP Query failed %s %s %s", name, auth.BaseDn, filterAttr, err.Error())
	}
	if len(sr.Entries) != 1 {
		if len(sr.Entries) == 0 {
			return fmt.Errorf("LDAP User %s not found, please check user or LDAP configuration parameters", name)
		}
		return fmt.Errorf("LDAP found that user %s has multiple accounts.", name)
	}
	err = parseEntries(sr)
	if err != nil {
		return fmt.Errorf("LDAP %s user %s", name, err.Error())
	}
	userDN := sr.Entries[0].DN
	err = l.Bind(userDN, pwd)
	if err != nil {
		return fmt.Errorf("%s LDAP Login failed, please check the login account or password %s", name, err.Error())
	}
	return nil
}

func parseEntries(sr *ldap.SearchResult) error {
	for _, attr := range sr.Entries[0].Attributes {
		switch attr.Name {
		case "shadowExpire":
			// -1 enabled, 1 disabled, >1 Number of days from 1970-01-01 to expiration date
			val, _ := strconv.ParseInt(attr.Values[0], 10, 64)
			if val == -1 {
				return nil
			}
			if val == 1 {
				return fmt.Errorf("Account has been deactivated")
			}
			if val > 1 {
				expireTime := time.Unix(val*86400, 0)
				t := time.Date(expireTime.Year(), expireTime.Month(), expireTime.Day(), 23, 59, 59, 0, time.Local)
				if t.Before(time.Now()) {
					return fmt.Errorf("Account has expired(Expiration date: %s)", t.Format("2006-01-02"))
				}
				return nil
			}
			return fmt.Errorf("Account shadowExpire value is abnormal: %d", val)
		}
	}
	return nil
}

func ValidateDomainPort(addr string) bool {
	re := regexp.MustCompile(`^([a-zA-Z0-9][-a-zA-Z0-9]{0,62}\.)+[A-Za-z]{2,18}\:([0-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-5]{2}[0-3][0-5])$`)
	return re.MatchString(addr)
}

func ValidateDN(dn string) bool {
	re := regexp.MustCompile(`^(?:(?:CN|cn|OU|ou|DC|dc)\=[^,'"]+,)*(?:CN|cn|OU|ou|DC|dc)\=[^,'"]+$`)
	return re.MatchString(dn)
}
