package main

import (
	"flag"
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

	if *verificationKey != "" {
		verifyKey(verificationKey)
	} else if *loginID != "" {
		login(loginID)
	} else if *sessionKey != "" {
		monitor(sessionKey)
	}

}

func login(loginID *string) {
	resp, err := http.Get("https://www.merriam-webster.com/word-of-the-day")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

}

func verifyKey(key *string) {

}

func monitor(sessionKey *string) {

}
