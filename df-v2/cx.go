package df_v2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

type fulfillmentInfo struct {
	Tag string `json:"tag"`
}

type sessionInfo struct {
	Session    string         `json:"session,omitempty"`
	Parameters map[string]any `json:"parameters,omitempty"`
}

type text struct {
	Text []string `json:"text"`
}

type responseMessage struct {
	Text text `json:"text"`
}

type fulfillmentResponse struct {
	Messages []responseMessage `json:"messages,omitempty"`
}

type webhookRequest struct {
	FulfillmentInfo fulfillmentInfo `json:"fulfillmentInfo"`
	SessionInfo     sessionInfo     `json:"sessionInfo"`
}

type webhookResponse struct {
	FulfillmentResponse *fulfillmentResponse `json:"fulfillmentResponse,omitempty"`
	SessionInfo         *sessionInfo         `json:"sessionInfo,omitempty"`
}

type tagHandler func(webhookRequest) (webhookResponse, error)

var tagHandlers = map[string]tagHandler{
	"welcome": welcome,
	"confirm": confirm,
}

func init() {
	functions.HTTP("HandleWebhookRequest", handleWebhookRequest)
}

func welcome(_ webhookRequest) (webhookResponse, error) {
	return textResponse("Welcome to the Dialogflow webhook!"), nil
}

func confirm(request webhookRequest) (webhookResponse, error) {
	size, err := stringParameter(request.SessionInfo.Parameters, "size")
	if err != nil {
		return webhookResponse{}, err
	}

	color, err := stringParameter(request.SessionInfo.Parameters, "color")
	if err != nil {
		return webhookResponse{}, err
	}

	response := textResponse(fmt.Sprintf(
		"You can pick up your order for a %s %s shirt in 5 days.",
		size,
		color,
	))

	response.SessionInfo = &sessionInfo{
		Parameters: map[string]any{
			"cancel-period": 2,
		},
	}

	return response, nil
}

func textResponse(message string) webhookResponse {
	return webhookResponse{
		FulfillmentResponse: &fulfillmentResponse{
			Messages: []responseMessage{
				{
					Text: text{
						Text: []string{message},
					},
				},
			},
		},
	}
}

func stringParameter(parameters map[string]any, name string) (string, error) {
	value, ok := parameters[name]
	if !ok {
		return "", fmt.Errorf("missing session parameter %q", name)
	}

	s, ok := value.(string)
	if !ok || strings.TrimSpace(s) == "" {
		return "", fmt.Errorf(
			"session parameter %q must be a non-empty string",
			name,
		)
	}

	return s, nil
}

func handleWebhookRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request webhookRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.WarnContext(
			r.Context(),
			"invalid Dialogflow webhook request",
			"error", err,
		)

		http.Error(w, "invalid JSON request", http.StatusBadRequest)
		return
	}

	tag := strings.TrimSpace(request.FulfillmentInfo.Tag)
	if tag == "" {
		http.Error(w, "missing fulfillment tag", http.StatusBadRequest)
		return
	}

	logger.InfoContext(
		r.Context(),
		"Dialogflow webhook request",
		"tag", tag,
		"session", request.SessionInfo.Session,
	)

	handler, ok := tagHandlers[tag]
	if !ok {
		logger.WarnContext(
			r.Context(),
			"unknown Dialogflow fulfillment tag",
			"tag", tag,
		)

		http.Error(w, "unknown fulfillment tag", http.StatusBadRequest)
		return
	}

	response, err := handler(request)
	if err != nil {
		logger.WarnContext(
			r.Context(),
			"Dialogflow webhook request rejected",
			"tag", tag,
			"error", err,
		)

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.ErrorContext(
			r.Context(),
			"failed to encode Dialogflow webhook response",
			"tag", tag,
			"error", err,
		)
	}
}
