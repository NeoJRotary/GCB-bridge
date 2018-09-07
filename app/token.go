package app

import (
	"crypto/rsa"
	"encoding/json"
	"log"
	"time"

	D "github.com/NeoJRotary/describe-go"
	"github.com/NeoJRotary/describe-go/dhttp"
	jwt "github.com/dgrijalva/jwt-go"
)

type accessToken struct {
	Token       string `json:"token"`
	ExpiresAt   string `json:"expires_at"`
	ExpiredTime time.Time
}

var appID string
var rsaKey *rsa.PrivateKey

var tokenCache = map[string]*accessToken{}

// InitToken read RSAKey and App ID
func InitToken() {
	keyENV := D.GetENV("GITHUB_APP_PRIVATE_KEY", "")
	if keyENV == "" {
		log.Fatal("GITHUB_APP_PRIVATE_KEY ENV should not be empty")
	}
	var err error
	rsaKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(keyENV))
	if D.IsErr(err) {
		log.Fatal(err)
	}

	appID = D.GetENV("GITHUB_APP_ID", "")
	if appID == "" {
		log.Fatal("GITHUB_APP_ID should not be empty")
	}
}

// GetAccessToken get access token by installation ID
func GetAccessToken(installationID string) string {
	// return empty if empty
	if installationID == "" {
		return ""
	}
	return getToken(installationID, time.Now().Add(time.Minute*-1))
}

func getToken(installationID string, current time.Time) string {
	var tkn *accessToken

	// reuse token if there is and before expired time
	tkn, ok := tokenCache[installationID]
	if ok {
		if tkn.ExpiredTime.After(current) {
			return tkn.Token
		}
	}
	tkn = &accessToken{}

	// start to request new access token

	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		Issuer:    appID,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(time.Minute).Unix(),
	})

	tokenString, err := token.SignedString(rsaKey)
	if D.IsErr(err) {
		log.Println("Github GetAccessToken SignedString failed", err)
		return ""
	}

	res, err := dhttp.Client(dhttp.TypeClient{
		Method: "POST",
		URL:    "https://api.github.com/installations/" + installationID + "/access_tokens",
		Header: map[string]string{
			"Accept":        "application/vnd.github.machine-man-preview+json",
			"Authorization": "Bearer " + tokenString,
		},
	}).Do()

	if D.IsErr(err) {
		log.Println("Github GetAccessToken request error", err)
		return ""
	}

	if res.StatusCode != 201 {
		log.Println("Github GetAccessToken request failed with", res.StatusCode)
		return ""
	}

	err = json.Unmarshal(res.ReadAllBody(), tkn)
	if D.IsErr(err) {
		log.Println("Github GetAccessToken Unmarshal error", err)
		return ""
	}

	tkn.ExpiredTime, err = time.Parse(time.RFC3339, tkn.ExpiresAt)
	if D.IsErr(err) {
		log.Println("Github GetAccessToken time.Parse error", err)
		return ""
	}

	tokenCache[installationID] = tkn

	return tkn.Token
}
