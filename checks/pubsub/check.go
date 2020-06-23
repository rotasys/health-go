package pubsub

import (
	"context"

	wErrors "github.com/pkg/errors"
	"gocloud.dev/pubsub"

	// Required to sucessfully healthcheck this in nats or in mem
	_ "gocloud.dev/pubsub/mempubsub"
	_ "gocloud.dev/pubsub/natspubsub"
)

// Config is the pubsub checker configuration settings container.
type Config struct {
	// DSN is the pubsub instance connection DSN. Required.
	DSN string
}

// New creates new pubsub health check that verifies the following:
// - connection establishing
func New(config Config) func() error {
	return func() (checkErr error) {
		topic, err := pubsub.OpenTopic(context.Background(), config.DSN)
		if err != nil {
			checkErr = wErrors.Wrap(err, "pubsub health check failed on connect")
		}

		defer func() {
			// override checkErr only if there were no other errors
			if err = topic.Shutdown(context.Background()); err != nil && checkErr == nil {
				checkErr = wErrors.Wrap(err, "pubsub health check failed on connection closing")
			}
		}()

		return
	}
}
