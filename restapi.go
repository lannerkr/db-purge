package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func updateUser(d []bson.D) []DisabledUsers {

	var disabledUsers []DisabledUsers
	for _, data := range d {
		un := data.Map()
		// for key, val := range un {
		// 	fmt.Println(key, val)
		// }

		realm := fmt.Sprint(un["realm"])
		var auth string
		switch realm {
		case "store":
			auth = "02.store.local"
		case "partner":
			auth = "03.partner.local"
		case "EMP-GOTP":
			auth = "emp.local"
		case "emp":
			continue
		default:
			auth = "02.store.local"
		}

		val := un["user_name"]
		user := fmt.Sprint(val)

		// if user doesn't exist in pulse user db
		// make user's enabled status to "Deleted"
		// and continue for next user
		if !checkUser(user, auth) {
			users := DisabledUsers{user, "Deleted"}
			disabledUsers = append(disabledUsers, users)
			continue
		}

		pulse := configuration.PulseUri
		pUri := pulse + "/api/v1/configuration/authentication/auth-servers/auth-server/" + auth + "/local/users/user/" + user + "/enabled"

		status := map[string]string{
			"enabled": "false",
		}
		pbytes, _ := json.Marshal(status)
		buff := bytes.NewBuffer(pbytes)

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}

		req, err := http.NewRequest(http.MethodPut, pUri, buff)
		if err != nil {
			panic(err)
		}

		apikey := configuration.PulseApiKey
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(apikey, "")

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		log.Println("[SSLVPN] username: ", user, ", ", status, ", result-code: ", resp.StatusCode)
		if resp.StatusCode == 200 {
			users := DisabledUsers{user, "False"}
			disabledUsers = append(disabledUsers, users)
		}
	}
	//fmt.Println(disabledUsers)
	return disabledUsers
}

func checkUser(user, auth string) bool {
	pulse := configuration.PulseUri
	pUri := pulse + "/api/v1/configuration/authentication/auth-servers/auth-server/" + auth + "/local/users/user/" + user + "/enabled"

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest(http.MethodGet, pUri, nil)
	if err != nil {
		panic(err)
	}

	apikey := configuration.PulseApiKey
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(apikey, "")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		log.Println("[SSLVPN] user: " + user + " is checked OK.")
		return true
	} else {
		log.Println("[SSLVPN] user: " + user + " is not existed in realm " + auth)
		return false
	}
}
