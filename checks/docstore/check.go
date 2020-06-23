package docstore

import (
	"context"

	wErrors "github.com/pkg/errors"
	"gocloud.dev/docstore"

	//Required to test the specified docstore
	_ "gocloud.dev/docstore/memdocstore"
	_ "gocloud.dev/docstore/mongodocstore"
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
		coll, err := docstore.OpenCollection(context.Background(), config.DSN)
		if err != nil {
			checkErr = wErrors.Wrap(err, "docstore health check failed on connect")
		}

		defer func() {
			// override checkErr only if there were no other errors
			if err = coll.Close(); err != nil && checkErr == nil {
				checkErr = wErrors.Wrap(err, "pubsub health check failed on connection closing")
			}
		}()

		return
	}
}
