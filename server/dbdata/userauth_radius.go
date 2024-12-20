package dbdata

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"reflect"
	"time"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

type AuthRadius struct {
	Addr   string `json:"addr"`
	Secret string `json:"secret"`
	Nasip  string `json:"nasip"`
}

func init() {
	authRegistry["radius"] = reflect.TypeOf(AuthRadius{})
}

func (auth AuthRadius) checkData(authData map[string]interface{}) error {
	authType := authData["type"].(string)
	bodyBytes, err := json.Marshal(authData[authType])
	if err != nil {
		return errors.New("Radius key/The server address is entered incorrectly")
	}
	json.Unmarshal(bodyBytes, &auth)
	if !ValidateIpPort(auth.Addr) {
		return errors.New("The Radius server address is incorrectly entered.")
	}
	// freeradius official website has a maximum of 8000 characters, here the limit is 200
	if len(auth.Secret) < 8 || len(auth.Secret) > 200 {
		return errors.New("Radius key length needs to be between 8 and 200 characters")
	}
	return nil
}

func (auth AuthRadius) checkUser(name, pwd string, g *Group) error {
	pl := len(pwd)
	if name == "" || pl < 1 {
		return fmt.Errorf("%s %s", name, "Wrong password")
	}
	authType := g.Auth["type"].(string)
	if _, ok := g.Auth[authType]; !ok {
		return fmt.Errorf("%s %s", name, "The radius value of Radius does not exist")
	}
	bodyBytes, err := json.Marshal(g.Auth[authType])
	if err != nil {
		return fmt.Errorf("%s %s", name, "Radius Marshal error")
	}
	err = json.Unmarshal(bodyBytes, &auth)
	if err != nil {
		return fmt.Errorf("%s %s", name, "Radius An error occurred in Unmarshal")
	}
	// During Radius authentication, set the timeout to 3 seconds.
	packet := radius.New(radius.CodeAccessRequest, []byte(auth.Secret))
	err = rfc2865.UserName_SetString(packet, name)
	if err != nil {
		return fmt.Errorf("%s %s", name, "Radius set name An error occurred")
	}
	err = rfc2865.UserPassword_SetString(packet, pwd)
	if err != nil {
		return fmt.Errorf("%s %s", name, "Radius set pwd An error occurred")
	}
	if auth.Nasip != "" {
		nasip := net.ParseIP(auth.Nasip)
		err = rfc2865.NASIPAddress_Set(packet, nasip)
		if err != nil {
			return fmt.Errorf("%s %s", name, "Radius set nasip An error occurred")
		}
	}

	ctx, done := context.WithTimeout(context.Background(), 3*time.Second)
	defer done()
	response, err := radius.Exchange(ctx, packet, auth.Addr)
	if err != nil {
		return fmt.Errorf("%s %s", name, "Radius server connection exception, Please check the server and port")
	}
	if response.Code != radius.CodeAccessAccept {
		return fmt.Errorf("%s %s", name, "Radius: Wrong username or password")
	}
	return nil
}
