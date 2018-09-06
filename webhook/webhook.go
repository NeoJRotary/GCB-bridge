package webhook

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"log"
	"strings"

	D "github.com/NeoJRotary/describe-go"
	"github.com/NeoJRotary/describe-go/dhttp"
)

// InitWebhookRoute init http server route of webhook
func InitWebhookRoute(server *dhttp.TypeServer) {
	server.Route("/webhook").POST(webhookHandler)
}

func webhookHandler(w *dhttp.ResponseWriter, r *dhttp.Request) {
	log.Println("----- Webhook Receive -----\nX-GitHub-Delivery " + r.Header.Get("X-GitHub-Delivery"))

	body := r.ReadAllBody()
	signature := r.Header.Get("X-Hub-Signature")

	if !verifyWebhook(body, signature) {
		if w != nil {
			w.WriteHeader(403)
		}
		return
	}

	if w != nil {
		go eventHandler(r.Header.Get("X-GitHub-Event"), body)
		w.WriteHeader(200)
	}
}

var webhookSecret = []byte(D.GetENV("GITHUB_APP_WEBHOOK_SECRET", ""))

func verifyWebhook(body []byte, signature string) bool {
	// 40 with "sha1=" prefix
	if len(signature) != 45 {
		return false
	}

	// check prefix
	if !strings.HasPrefix(signature, "sha1=") {
		return false
	}

	hash := hmac.New(sha1.New, webhookSecret)
	hash.Write(body)
	mac1 := hash.Sum(nil)
	mac2 := make([]byte, 20)
	hex.Decode(mac2, []byte(signature[5:]))

	return hmac.Equal(mac1, mac2)
}
