package main

import (
	"log"

	"github.com/NeoJRotary/GCB-bridge/app"
	"github.com/NeoJRotary/GCB-bridge/gcloud"
	"github.com/NeoJRotary/GCB-bridge/webhook"
	D "github.com/NeoJRotary/describe-go"
	"github.com/NeoJRotary/describe-go/dhttp"
)

func init() {
	app.Init()
	gcloud.Init()
}

func main() {
	addr := D.GetENV("LISTEN_ON", "0.0.0.0:8080")
	server := dhttp.Server().ListenOn(addr)
	webhook.InitWebhookRoute(server)

	log.Println("Server listen on 8080")
	err := server.Start()
	log.Fatal(err)
}
