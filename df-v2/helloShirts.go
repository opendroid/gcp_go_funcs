package df_v2

import (
	"io"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("HelloShirts", helloShirts)
	functions.HTTP("Echo", echo)
	functions.HTTP("LogTestV2", logTestV2)
}

func helloShirts(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	logger.Info(
		"helloShirts request",
		"body", string(body),
	)

	_, _ = w.Write([]byte("echo: " + string(body)))
}

func echo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	logger.Info(
		"echo request",
		"method", "echo",
		"bot", "echo",
		"body", string(body),
	)

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(body)
}

func logTestV2(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	logger.Info(
		"logTestV2 request",
		"body", string(body),
	)

	_, _ = w.Write(body)
}
