package main

import (
	"bytes"
	"encoding/json"
	"github.com/antchfx/htmlquery"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func updateEnergyConsumption() {
	log.Println("Start update consumption")
	var err error

	retryCount := 5
	for i := 0; i < retryCount; i++ {
		responseData, err := fetchEnergyConsumption()
		if err == nil {
			log.Println(responseData)
			return
		}
		time.Sleep(time.Hour)
	}

	if err != nil {
		log.Println("Fetching energy consumption")
	}

}

func fetchEnergyConsumption() (*BewegungsDaten, error) {

	redirectUri := "https://smartmeter-web.wienernetze.at/"

	authParams := url.Values{}
	authParams.Add("client_id", "wn-smartmeter")
	authParams.Add("redirect_uri", redirectUri)
	authParams.Add("response_mode", "fragment")
	authParams.Add("response_type", "code")
	authParams.Add("scope", "openid")
	authParams.Add("nonce", "")

	authUrl := "https://log.wien/auth/realms/logwien/protocol/openid-connect/"
	loginPageUrl := authUrl + "auth?" + authParams.Encode()

	loginPage, err := http.Get(loginPageUrl)

	document, err := htmlquery.Parse(loginPage.Body)
	if err != nil {
		return nil, err
	}

	form, _ := htmlquery.Query(document, "//*[@id=\"kc-login-form\"]")

	var loginEndpoint string

	for _, attr := range form.Attr {
		if attr.Key == "action" {
			loginEndpoint = attr.Val
			break
		}
	}

	loginValues := url.Values{}
	loginValues.Add("username", os.Getenv("SMARTMETER_USER"))
	loginValues.Add("password", os.Getenv("SMARTMETER_PASSWORD"))

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	cookies := loginPage.Cookies()

	loginRequest, err := http.NewRequest("POST", loginEndpoint, bytes.NewBuffer([]byte(loginValues.Encode())))
	if err != nil {
		return nil, err
	}

	loginRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	for _, cookie := range cookies {
		loginRequest.AddCookie(cookie)
	}

	res, err := client.Do(loginRequest)
	if err != nil {
		return nil, err
	}

	location, err := res.Location()
	if err != nil {
		return nil, err
	}

	locationQuery, err := url.ParseQuery(location.Fragment)
	if err != nil {
		return nil, err
	}

	code := locationQuery.Get("code")

	tokenValues := url.Values{}
	tokenValues.Add("grant_type", "authorization_code")
	tokenValues.Add("client_id", "wn-smartmeter")
	tokenValues.Add("redirect_uri", redirectUri)
	tokenValues.Add("code", code)

	tokenRequest, err := http.NewRequest("POST", authUrl+"token", bytes.NewBuffer([]byte(tokenValues.Encode())))

	tokenRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	for _, cookie := range cookies {
		tokenRequest.AddCookie(cookie)
	}

	tokenResponse, err := client.Do(tokenRequest)
	if err != nil {
		return nil, err
	}

	tokenBody, err := io.ReadAll(tokenResponse.Body)
	if err != nil {
		return nil, err
	}

	var tokens TokenResponse
	err = json.Unmarshal(tokenBody, &tokens)
	if err != nil {
		return nil, err
	}

	//TODO check if accesstoken is present

	profileRequest, err := http.NewRequest("GET", "https://service.wienernetze.at/sm/api/user/profile", nil)
	if err != nil {
		return nil, err
	}

	for _, cookie := range cookies {
		profileRequest.AddCookie(cookie)
	}

	profileRequest.Header.Add("Authorization", "Bearer "+tokens.AccessToken)

	profileResponse, err := client.Do(profileRequest)
	if err != nil {
		return nil, err
	}

	profileBody, err := io.ReadAll(profileResponse.Body)
	if err != nil {
		return nil, err
	}

	var profile Profile
	err = json.Unmarshal(profileBody, &profile)
	if err != nil {
		return nil, err
	}

	today := time.Now().Truncate(time.Hour * 24).Format("2006-01-02T15:04:05Z")
	oneMonthAgo := time.Now().Truncate(time.Hour*24).AddDate(0, -1, 0).Format("2006-01-02T15:04:05Z")

	bewegungsdatenParams := url.Values{}
	bewegungsdatenParams.Add("geschaeftspartner", profile.Registration.Geschaeftspartner)
	bewegungsdatenParams.Add("zaehlpunktnummer", profile.Registration.Zaehlpunkt)
	bewegungsdatenParams.Add("rolle", "V002")
	bewegungsdatenParams.Add("zeitpunktVon", oneMonthAgo)
	bewegungsdatenParams.Add("zeitpunktBis", today)
	bewegungsdatenParams.Add("aggregat", "NONE")

	bewegungsdatenRequest, err := http.NewRequest("GET", "https://service.wienernetze.at/sm/api/user/messwerte/bewegungsdaten?"+bewegungsdatenParams.Encode(), nil)
	if err != nil {
		return nil, err
	}

	for _, cookie := range cookies {
		bewegungsdatenRequest.AddCookie(cookie)
	}

	bewegungsdatenRequest.Header.Add("Authorization", "Bearer "+tokens.AccessToken)

	bewegungsdatenResponse, err := client.Do(bewegungsdatenRequest)
	if err != nil {
		return nil, err
	}

	bewegungsdatenBody, err := io.ReadAll(bewegungsdatenResponse.Body)
	if err != nil {
		return nil, err
	}

	var bewegungsdaten BewegungsDaten

	err = json.Unmarshal(bewegungsdatenBody, &bewegungsdaten)
	if err != nil {
		return nil, err
	}

	return &bewegungsdaten, nil
}
