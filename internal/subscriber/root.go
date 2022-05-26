package subscriber

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/charliemcelfresh/event_worker/internal/config"

	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/rabbitpubsub"
)

//Store interface to write into database
type Store interface {
	InsertEvents(ctx context.Context, events []json.RawMessage) error
}

//subscriber --
type subscriber struct {
	Store     Store
	WaitTime  time.Duration
	BatchSize int
}

func New(s Store, wt time.Duration, bs int) subscriber {
	return subscriber{
		Store:     s,
		WaitTime:  wt,
		BatchSize: bs,
	}
}

func (s subscriber) Subscribe() {
	ctx := context.Background()
	subscription, err := pubsub.OpenSubscription(ctx, "rabbit://events_development")
	defer subscription.Shutdown(ctx)
	if err != nil {
		config.Logger.Panicf("could not open topic subscription: %v", err)
	}

	mch := make(chan json.RawMessage)
	go func() {
		for {
			msg, err := subscription.Receive(ctx)
			if err != nil {
				config.Logger.Panicf("could not receive message: %v", err)
			}
			fmt.Printf("Got message: %q\n", msg.Body)
			mch <- msg.Body
			msg.Ack()
		}
	}()

	toWrite := make([]json.RawMessage, 0, s.BatchSize)

	for {
		select {
		case <-time.After(s.WaitTime):
			//write db
			if len(toWrite) == 0 {
				config.Logger.Debugf("nothing to write")
				continue
			}
			config.Logger.Debugf("Write process")
			_ = s.Store.InsertEvents(ctx, toWrite)
			toWrite = make([]json.RawMessage, 0, s.BatchSize)
		case m := <-mch:
			toWrite = append(toWrite, m)
			if len(toWrite) == s.BatchSize {
				//write db
				config.Logger.Debugf("write process buffer is full")
				_ = s.Store.InsertEvents(ctx, toWrite)
				toWrite = make([]json.RawMessage, 0, s.BatchSize)
			}
		}
	}
}
