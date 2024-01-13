package main

import "time"

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	IdToken          string `json:"id_token"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	Scope            string `json:"scope"`
	SessionState     string `json:"session_state"`
	TokenType        string `json:"token_type"`
}

type Profile struct {
	Registration DefaultGeschaeftspartnerRegistration `json:"defaultGeschaeftspartnerRegistration"`
}

type DefaultGeschaeftspartnerRegistration struct {
	Geschaeftspartner string `json:"geschaeftspartner"`
	Zaehlpunkt        string `json:"zaehlpunkt"`
}

type BewegungsDaten struct {
	Descriptor BewegungsDatenDescriptor `json:"descriptor"`
	Values     []BewegungsdatenValue
}

type BewegungsDatenDescriptor struct {
	Geschaeftspartnernummer string `json:"geschaeftspartnernummer"`
	Zaehlpunktnummer        string `json:"zaehlpunktnummer"`
	Rolle                   string `json:"rolle"`
	Aggregat                string `json:"aggregat"`
	Granularitaet           string `json:"granularitaet"`
	Einheit                 string `json:"einheit"`
}

type BewegungsdatenValue struct {
	Wert         float64                `json:"wert"`
	ZeitpunktVon BewegungsdatenDateTime `json:"zeitpunktVon"`
	ZeitpunktBis BewegungsdatenDateTime `json:"zeitpunktBis"`
	Geschaetzt   bool                   `json:"geschaetzt"`
}

type BewegungsdatenDateTime struct {
	time.Time
}

func (t *BewegungsdatenDateTime) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"2006-01-02T15:04:05Z"`, string(b))
	if err != nil {
		return err
	}
	t.Time = date
	return
}
