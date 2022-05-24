package publisher

import (
	"context"
	"encoding/json"

	"github.com/bxcodec/faker/v3"
	"github.com/charliemcelfresh/event_worker/internal/config"
	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/rabbitpubsub"
)

type Event struct {
	Email string `faker:"email"`
}

func Publish(eventCount int) error {
	ctx := context.Background()
	topic, err := pubsub.OpenTopic(ctx, "rabbit://events_development_fanout")
	if err != nil {
		return err
	}
	defer topic.Shutdown(ctx)
	for i := 0; i < eventCount; i++ {
		e := Event{}
		err := faker.FakeData(&e)
		if err != nil {
			config.Logger.Error(err)
		}
		msgAsJson, err := json.Marshal(e)
		if err != nil {
			config.Logger.Error(err)
		}
		err = topic.Send(ctx, &pubsub.Message{Body: msgAsJson})
		if err != nil {
			config.Logger.Error(err)
		}
	}
	return nil
}
