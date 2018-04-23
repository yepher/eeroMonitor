package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
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

// NetworkDeviceResponse Details about devices on network
type NetworkDeviceResponse struct {
	Meta struct {
		Code       int       `json:"code"`
		ServerTime time.Time `json:"server_time"`
	} `json:"meta"`
	Data []struct {
		URL            string      `json:"url"`
		Mac            string      `json:"mac"`
		Eui64          string      `json:"eui64"`
		Manufacturer   string      `json:"manufacturer"`
		IP             string      `json:"ip"`
		Ips            []string    `json:"ips"`
		Nickname       interface{} `json:"nickname"`
		Hostname       string      `json:"hostname"`
		Connected      bool        `json:"connected"`
		Wireless       bool        `json:"wireless"`
		ConnectionType string      `json:"connection_type"`
		Source         struct {
			Location string `json:"location"`
		} `json:"source"`
		LastActive   time.Time `json:"last_active"`
		FirstActive  time.Time `json:"first_active"`
		Connectivity struct {
			RxBitrate string  `json:"rx_bitrate"`
			Signal    string  `json:"signal"`
			SignalAvg string  `json:"signal_avg"`
			Score     float64 `json:"score"`
			ScoreBars int     `json:"score_bars"`
		} `json:"connectivity"`
		Interface struct {
			Frequency     string `json:"frequency"`
			FrequencyUnit string `json:"frequency_unit"`
		} `json:"interface"`
		Usage struct {
			DownMbps float64 `json:"down_mbps"`
			UpMbps   float64 `json:"up_mbps"`
		} `json:"usage"`
		Profile struct {
			URL    string `json:"url"`
			Name   string `json:"name"`
			Paused bool   `json:"paused"`
		} `json:"profile"`
		DeviceType string `json:"device_type"`
	} `json:"data"`
}

func main() {

	sessionKey := flag.String("sessionKey", "", "Eero session key")
	loginID := flag.String("loginID", "", "Eero loginId")
	verificationKey := flag.String("verificationKey", "", "Eero verification key")
	networkID := flag.String("networkID", "", "Network ID to monitor")
	command := flag.String("command", "monitor", "Actions to perform")
	flag.Parse()

	if *verificationKey != "" && *sessionKey != "" {
		verifyKey(verificationKey, sessionKey)
		fmt.Printf("Next monitor netowork with networkID (/2.2/networks/[ID]):\n")
		fmt.Printf("\t./eeroMonitor -sessionKey=\"%s\" --networkID=\n", *sessionKey)

	} else if *loginID != "" {
		sessionKey := login(loginID)
		//fmt.Printf("sessionKey=%s\n", sessionKey)
		fmt.Printf("Next verify session with verification code: \n")
		fmt.Printf("\t./eeroMonitor -sessionKey=\"%s\" -verificationKey=\n", sessionKey)

	} else if *sessionKey != "" && *networkID != "" {
		if *command == "" || *command == "monitor" {
			monitor(sessionKey, networkID)
		} else {
			fmt.Printf("Unknown command: %s\n", *command)

		}

	} else {
		fmt.Printf("Unknow set of arguments...\n")
		fmt.Printf("\t./eeroMonitor -loginID=[YOUR_LOGIN_ID]\n")
		return
	}
}

func login(loginID *string) string {
	fmt.Printf("Login: %s\n", *loginID)
	url := "https://api-user.e2ro.com/2.2/login?"

	loginRequest := LoginRequest{Login: *loginID}
	var loginResponse LoginResponse

	err := doRequest(url, nil, &loginRequest, &loginResponse)
	//r, err := http.Post("https://api-user.e2ro.com/2.2/login?", "application/json; charset=utf-8", b)
	if err != nil {
		panic(err)
	}

	return loginResponse.Data.UserToken
}

func verifyKey(verificationKey *string, sessionKey *string) string {
	fmt.Printf("Verify: %s, %s\n", *verificationKey, *sessionKey)

	verifyRequest := LoginVerifyRequest{Code: *verificationKey}
	var verifyResponse LoginVerifyResponse
	url := "https://api-user.e2ro.com/2.2/login/verify?"
	err := doRequest(url, sessionKey, &verifyRequest, &verifyResponse)

	if err != nil {
		panic(err)
	}

	//fmt.Println(verifyResponse)

	networks := verifyResponse.Data.Networks.Data
	for _, network := range networks {
		fmt.Printf("%s - %s\n", network.Name, network.URL)
	}
	return ""
}

func monitor(sessionKey *string, networkID *string) {
	fmt.Printf("Monitoring Network: %s\n", *networkID)
	url := fmt.Sprintf("https://api-user.e2ro.com%s/devices?thread=true", *networkID)
	fmt.Printf("URL: %s\n", url)

	// TODO: Query /2.2/networks/[NETWORK_ID]/burst_reporters?
	// to schedule next time to query usage

	for {
		var networkDeviceResponse NetworkDeviceResponse
		err := doRequest(url, sessionKey, nil, &networkDeviceResponse)

		if err != nil {
			panic(err)
		}

		networks := networkDeviceResponse.Data
		foundResult := false
		for _, device := range networks {
			up := device.Usage.UpMbps
			down := device.Usage.DownMbps
			if up > 0 || down > 0 {
				foundResult = true
				fmt.Printf("%s - %s (%f Mbps, %f Mbps)\n", device.Hostname, device.DeviceType, device.Usage.DownMbps, device.Usage.UpMbps)
			}
		}

		if foundResult {
			fmt.Printf("\n\n\n\n")
		}
		time.Sleep(60 * time.Second)
	}
}

func doRequest(url string, token *string, reqObj interface{}, respObj interface{}) error {
	method := "GET"

	b := new(bytes.Buffer)
	if reqObj != nil {
		method = "POST"
		json.NewEncoder(b).Encode(reqObj)
	} else {
		b = nil
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")

	if token != nil {
		sessionString := fmt.Sprintf("s=%s", *token)
		//fmt.Printf("Session Key: %s\n", sessionString)
		req.Header.Add("Cookie", sessionString)
	}

	if req != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	r, err := client.Do(req)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	if r.StatusCode != 200 {
		return errors.Errorf("Request failed: (%d) - %s\nURL: %s %s\nRequest: %s", r.StatusCode, r.Status, method, url, reqObj)
	}

	if r.Body == nil && respObj == nil {
		return nil
	}

	err = json.NewDecoder(r.Body).Decode(respObj)
	if err != nil {
		return err
	}

	return nil
}
