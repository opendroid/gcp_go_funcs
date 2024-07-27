package df_v2

import (
	"context"
	"fmt"

	"cloud.google.com/go/logging"
)

// Note that to pipe zap to GCP logger, see [gcloudzap]
// [gcloudzap]: https://pkg.go.dev/github.com/jonstaryuk/gcloudzap

var (
	logger *logging.Logger // GCP Logger
)

func init() {
	// Create a client to gcp project "the-gpl" that logs to "df-v2" logName
	client, err := logging.NewClient(context.Background(), "the-gpl")
	if err != nil {
		panic(fmt.Errorf("logging.NewClient: %v", err))
	}
	log_name := "df-v2"
	logger = client.Logger(log_name)
}
