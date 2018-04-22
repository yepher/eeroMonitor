package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"
)

// LoginRequest is use to login to Eero. Set the account email address
type LoginRequest struct {
	Login string `json:"login"`
}

// LoginResponse response to the login request
// user_token will be used as `Cookie: s={{user_token}}`
// Code looks similar to the HTTP response code
type LoginResponse struct {
	Meta struct {
		Code       int       `json:"code"`
		ServerTime time.Time `json:"server_time"`
	} `json:"meta"`
	Data struct {
		UserToken string `json:"user_token"`
	} `json:"data"`
}

// LoginVerifyRequest sends the code from email
type LoginVerifyRequest struct {
	Code string `json:"code"`
}

// LoginVerifyResponse Returns details about your network
type LoginVerifyResponse struct {
	Meta struct {
		Code       int       `json:"code"`
		ServerTime time.Time `json:"server_time"`
	} `json:"meta"`
	Data struct {
		Name  string `json:"name"`
		Phone struct {
			Value    string `json:"value"`
			Verified bool   `json:"verified"`
		} `json:"phone"`
		Email struct {
			Value    string `json:"value"`
			Verified bool   `json:"verified"`
		} `json:"email"`
		LogID    string `json:"log_id"`
		Networks struct {
			Count int `json:"count"`
			Data  []struct {
				URL     string    `json:"url"`
				Name    string    `json:"name"`
				Created time.Time `json:"created"`
			} `json:"data"`
		} `json:"networks"`
		Role          string `json:"role"`
		CanTransfer   bool   `json:"can_transfer"`
		IsProOwner    bool   `json:"is_pro_owner"`
		PremiumStatus string `json:"premium_status"`
		PushSettings  struct {
			NetworkOffline bool `json:"networkOffline"`
			NodeOffline    bool `json:"nodeOffline"`
		} `json:"push_settings"`
		TrustCertificatesEtag string `json:"trust_certificates_etag"`
	} `json:"data"`
}

// LogoutResponse session logout
type LogoutResponse struct {
	Meta struct {
		Code       int       `json:"code"`
		ServerTime time.Time `json:"server_time"`
	} `json:"meta"`
}

func main() {

	sessionKey := flag.String("sessionKey", "", "Eero session key")
	loginID := flag.String("loginID", "", "Eero loginId")
	verificationKey := flag.String("verificationKey", "", "Eero verification key")
	flag.Parse()

	if *verificationKey != "" && *sessionKey != "" {
		verifyKey(verificationKey, sessionKey)
	} else if *loginID != "" {
		sessionKey := login(loginID)
		fmt.Printf("sessionKey=%s\n", sessionKey)
	} else if *sessionKey != "" {
		monitor(sessionKey)
	} else {
		fmt.Printf("Unknow set of arguments...")
		return
	}
}

func login(loginID *string) string {
	fmt.Printf("Login: %s\n", *loginID)

	loginRequest := LoginRequest{Login: *loginID}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(loginRequest)

	r, err := http.Post("https://api-user.e2ro.com/2.2/login?", "application/json; charset=utf-8", b)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	var login LoginResponse
	if r.Body == nil {
		fmt.Printf("Body was empty.\n")
		return ""
	}
	err = json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		panic(err)
	}

	return login.Data.UserToken
}

func verifyKey(verificationKey *string, sessionKey *string) string {
	fmt.Printf("Verify: %s, %s\n", *verificationKey, *sessionKey)

	verifyRequest := LoginVerifyRequest{Code: *verificationKey}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(verifyRequest)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api-user.e2ro.com/2.2/login/verify?", b)
	if err != nil {
		panic(err)
	}

	sessionString := fmt.Sprintf("s=%s", *sessionKey)
	fmt.Printf("Session Key: %s\n", sessionString)
	req.Header.Add("Cookie", sessionString)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	r, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()

	var verifyResponse LoginVerifyResponse
	if r.Body == nil {
		fmt.Printf("Body was empty.\n")
		return ""
	}
	err = json.NewDecoder(r.Body).Decode(&verifyResponse)
	if err != nil {
		panic(err)
	}

	fmt.Println(verifyResponse)
	return ""
}

func monitor(sessionKey *string) {

}
