package dead_letter_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

// Example: In pub/sub
// https://cloud.google.com/functions/docs/calling/pubsub
// LogEntryPubSub defines structure of a pub/sub message

func init() {
	functions.CloudEvent("AckPubMessage", ackPubMessage)
	functions.HTTP("BadAckFunc", badAckFunc)
}

// MessagePublishedData contains the full Pub/Sub message
// See the [pubsub-documentation] for more details:
//
// [pubsub-documentation]: https://cloud.google.com/eventarc/docs/cloudevents#pubsub
type MessagePublishedData struct {
	Message      PubSubMessage `json:"message,omitempty"`
	Subscription string        `json:"subscription,omitempty"`
}

// PubSubMessage is the payload of a Pub/Sub event.
// See the [event-documentation] for more details:
//
// [event-documentation]: https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data       []byte            `json:"data,omitempty"`
	ID         string            `json:"id"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// ackPubMessage test a message
func ackPubMessage(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		logger.ErrorContext(
			ctx,
			"failed to decode event data",
			"method", "AckPubMessage",
			"error", err,
		)
		return fmt.Errorf("event.DataAs: %w", err)
	}

	data := string(msg.Message.Data) // Automatically decoded from base64.

	if data != "" {
		logger.InfoContext(
			ctx,
			"AckPubMessage received message",
			"method", "AckPubMessage",
			"subscription", msg.Subscription,
			"data", data,
			"message_id", msg.Message.ID,
		)
	} else {
		logger.InfoContext(
			ctx,
			"AckPubMessage received message with no data",
			"method", "AckPubMessage",
			"subscription", msg.Subscription,
			"message_id", msg.Message.ID,
		)
	}
	return nil
}

// badAckFunc a cloud func that returns 'internal server' error
func badAckFunc(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()

	var msg MessagePublishedData
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		logger.ErrorContext(
			r.Context(),
			"request not formatted",
			"method", "BadAckFunc",
			"error", err.Error(),
		)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	data := string(msg.Message.Data)
	deliveryAttempt := msg.Message.Attributes["googclient_deliveryattempt"]

	if data != "" {
		if deliveryAttempt != "" {
			logger.InfoContext(
				r.Context(),
				"BadAckFunc received message",
				"method", "BadAckFunc",
				"subscription", msg.Subscription,
				"data", data,
				"attempt", deliveryAttempt,
				"message_id", msg.Message.ID,
			)
		} else {
			logger.InfoContext(
				r.Context(),
				"BadAckFunc received message",
				"method", "BadAckFunc",
				"subscription", msg.Subscription,
				"data", data,
				"message_id", msg.Message.ID,
			)
		}
	} else {
		logger.InfoContext(
			r.Context(),
			"BadAckFunc received message with no data",
			"method", "BadAckFunc",
			"subscription", msg.Subscription,
			"message_id", msg.Message.ID,
		)
	}

	http.Error(w, "Internal server error", http.StatusInternalServerError)
}
