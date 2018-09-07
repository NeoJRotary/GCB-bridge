package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"
	"time"

	D "github.com/NeoJRotary/describe-go"
)

type eventPayload struct {
	Installation struct {
		ID int `json:"id"`
	} `json:"installation"`
}

func TestToken(t *testing.T) {
	InitToken()

	b, err := ioutil.ReadFile("../github-event-payload.json")
	if D.IsErr(err) {
		t.Fatal(err)
	}

	var payload eventPayload
	err = json.Unmarshal(b, &payload)
	if D.IsErr(err) {
		t.Fatal(err)
	}
	installationID := strconv.Itoa(payload.Installation.ID)

	tkn1 := getToken(installationID, time.Now())
	if tkn1 == "" {
		t.Fatal("should not be empty")
	}

	fmt.Println("tkn1", tkn1)

	tkn2 := getToken(installationID, time.Now().Add(time.Minute*30))
	if tkn2 == "" {
		t.Fatal("should not be empty")
	}

	fmt.Println("tkn2", tkn2)

	if tkn1 != tkn2 {
		t.Fatal("should be same")
	}

	tkn2 = getToken(installationID, time.Now().Add(time.Minute*62))
	if tkn2 == "" {
		t.Fatal("should not be empty")
	}

	fmt.Println("tkn2", tkn2)

	if tkn1 == tkn2 {
		t.Fatal("should not be same")
	}
}
