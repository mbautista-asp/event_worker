package subscriber

import (
	"context"
	"fmt"

	"github.com/charliemcelfresh/event_worker/internal/config"

	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/rabbitpubsub"
)

func Subscribe() {
	ctx := context.Background()
	subscription, err := pubsub.OpenSubscription(ctx, "rabbit://events_development")
	defer subscription.Shutdown(ctx)
	if err != nil {
		config.Logger.Panicf("could not open topic subscription: %v", err)
	}
	for {
		msg, err := subscription.Receive(ctx)
		if err != nil {
			config.Logger.Panicf("could not receive message: %v", err)
		}
		// put code here that 1) writes the first 10 messages to the db, or 2) if 1 second passes, and there
		// are not ten messages, write those to the db
		fmt.Printf("Got message: %q\n", msg.Body)

		msg.Ack()
	}
}
