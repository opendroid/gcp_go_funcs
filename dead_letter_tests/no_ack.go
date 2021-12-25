package dead_letter_tests

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Example: In pub/sub
// https://cloud.google.com/functions/docs/calling/pubsub
// LogEntryPubSub defines structure of a pub/sub message

// AckPubMessage test a message
func AckPubMessage(_ context.Context, m *pubsub.Message) error {
	if string(m.Data) != "" {
		m := fmt.Sprintf(`{"severity": "INFO", "method": "AckPubMessage", "data": %s}`, m.Data)
		fmt.Println(m)
	} else {
		fmt.Println(`{"severity"": "INFO", "method": "AckPubMessage", "data": "no-data"}`)
	}
	return nil
}

// pushMessageRequest sample message
type pushMessageRequest struct {
	Message      pubsub.Message
	Subscription string
}

// BadAckFunc a cloud func that returns 'internal server' error
func BadAckFunc(w http.ResponseWriter, r *http.Request) {
	var msg pushMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		m := fmt.Sprintf(`{"severity": "ERROR", "message": "Request not formatted", "method": "BadAckFunc", "error": "%s"}`, err.Error())
		fmt.Println(m)
		return
	}
	_ = r.Body.Close()
	if data := string(msg.Message.Data); data != "" {
		if msg.Message.DeliveryAttempt != nil { // Note that this is not set yet.
			data += fmt.Sprintf(`, "attempt": %d`, *msg.Message.DeliveryAttempt)
		}
		m := fmt.Sprintf(`{"severity": "INFO", "method": "BadAckFunc", "subscription": "%s", "data": %s}`, msg.Subscription, data)
		fmt.Println(m)
	} else {
		m := fmt.Sprintf(`{"severity": "INFO", "method": "BadAckFunc", "subscription": "%s"}`, msg.Subscription)
		fmt.Println(m)
	}
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}
