package subscriber

import (
	"context"
	"sync"
	"time"

	"github.com/lib/pq"

	"github.com/charliemcelfresh/event_worker/internal/config"

	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/rabbitpubsub"
)

func Subscribe() {
	ctx := context.Background()
	ch := make(chan *pubsub.Message)
	var messages []*pubsub.Message
	ticker := time.NewTicker(500 * time.Millisecond)
	subscription, err := pubsub.OpenSubscription(ctx, "rabbit://events_development")
	defer subscription.Shutdown(ctx)
	if err != nil {
		config.Logger.Panicf("could not open topic subscription: %v", err)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			msg, err := subscription.Receive(ctx)
			if err != nil {
				msg.Nack()
			}
			ch <- msg
		}

	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case t := <-ticker.C:
				config.Logger.Debug(t)
				if len(messages) >= 0 {
					config.Logger.Debug(("-------- inside ticker -----------"))
					print(messages)
					err := write(ctx, messages)
					if err != nil {
						config.Logger.Error(err)
					}
					messages = []*pubsub.Message{}
				}
				messages = []*pubsub.Message{}
			case m := <-ch:
				messages = append(messages, m)
				if len(messages) == 10000 {
					config.Logger.Debug(("-------- inside length -----------"))
					print(messages)
					err := write(ctx, messages)
					if err != nil {
						config.Logger.Error(err)
					}
					messages = []*pubsub.Message{}
				}
			}
		}
	}()
	wg.Wait()
}

func print(messages []*pubsub.Message) {
	config.Logger.Debugf("length is %d\n", len(messages))
	for _, m := range messages {
		config.Logger.Debugf("length is %s\n", string(m.Body))
	}
}

func write(ctx context.Context, messages []*pubsub.Message) error {
	db := config.DBPool
	txn, err := db.Begin()
	if err != nil {
		return err
	}
	statement, err := txn.PrepareContext(ctx, pq.CopyIn("event", "event"))
	if err != nil {
		return err
	}
	for _, m := range messages {
		_, err := statement.Exec(string(m.Body))
		if err != nil {
			return err
		}
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}
	err = statement.Close()
	if err != nil {
		return err
	}
	err = txn.Commit()
	if err != nil {
		return err
	}
	for _, m := range messages {
		m.Ack()
	}
	return nil
}
