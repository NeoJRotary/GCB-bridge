package gcloud

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
)

// initPubSub init gcp pub/sub
func initPubSub() {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, ProjectID)
	if err != nil {
		log.Fatal(err)
	}

	sub := client.Subscription("cloud-builds")
	go subReceive(ctx, sub)

	log.Println("Cloud PubSub Listening")
}

func subReceive(ctx context.Context, sub *pubsub.Subscription) {
	err := sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		messageHandler(m)
		m.Ack()
	})

	panic(err)
}
