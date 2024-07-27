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
func ackPubMessage(_ context.Context, e event.Event) error {
	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %w", err)
	}
	data := string(msg.Message.Data) // Automatically decoded from base64.
	fmt.Println()                    // Fix error in printing
	// WARNING: failed to extract Pub/Sub topic name from the URL request path: "/",
	// configure your subscription's push endpoint to use the following path pattern: 'projects/PROJECT_NAME/topics/TOPIC_NAME'
	if string(data) != "" {
		m := fmt.Sprintf(`{"severity": "INFO", "method": "AckPubMessage",  "subscription": "%s", "data": %s}`, msg.Subscription, data)
		fmt.Println(m)
	} else {
		m := fmt.Sprintf(`{"severity"": "INFO", "method": "AckPubMessage", "subscription": "%s", "data": "no-data"}`, msg.Subscription)
		fmt.Println(m)
	}
	return nil
}

// badAckFunc a cloud func that returns 'internal server' error
func badAckFunc(w http.ResponseWriter, r *http.Request) {
	var msg MessagePublishedData
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		m := fmt.Sprintf(`{"severity": "ERROR", "message": "Request not formatted", "method": "BadAckFunc", "error": "%s"}`, err.Error())
		fmt.Println(m)
		return
	}
	_ = r.Body.Close()
	if data := string(msg.Message.Data); data != "" {
		deliveryAttempt := msg.Message.Attributes["googclient_deliveryattempt"]
		if deliveryAttempt != "" { // Note that this is not set yet.
			data += fmt.Sprintf(`, "attempt": %s`, deliveryAttempt)
		}
		m := fmt.Sprintf(`{"severity": "INFO", "method": "BadAckFunc", "subscription": "%s", "data": %s}`, msg.Subscription, data)
		fmt.Println(m)
	} else {
		m := fmt.Sprintf(`{"severity": "INFO", "method": "BadAckFunc", "subscription": "%s"}`, msg.Subscription)
		fmt.Println(m)
	}
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}
