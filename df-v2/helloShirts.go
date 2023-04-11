package df_v2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"cloud.google.com/go/logging"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("HelloShirts", helloShirts)
	functions.HTTP("Echo", echo)
	functions.HTTP("ManojS", manojS)
}

// helloShirts is an HTTP Cloud Function.
func helloShirts(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		_, _ = fmt.Fprintf(w, "Error reading request body: %v", err)
		return
	}
	fmt.Printf(`{"method": "helloShirts", "message": "Hello, %s!"}`, body)
	_, _ = fmt.Fprintf(w, "echo: %s", body)
}

// echo is an HTTP Cloud Function.
func echo(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()
	defer logger.Flush() // Ensure the entry is written.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		_, _ = fmt.Fprintf(w, "Error reading request body: %v", err)
		return
	}
	logger.Log(logging.Entry{Payload: json.RawMessage(string(body)),
		Labels: map[string]string{"method": "echo", "bot": "echo"},
	})
	_, _ = fmt.Fprintf(w, "%s", body)
}

// manojS is an HTTP Cloud Function.
func manojS(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		_, _ = fmt.Fprintf(w, "Error reading request body: %v", err)
		return
	}
	fmt.Printf(`{message: "manojS", "body": %s}`, body)
	_, _ = fmt.Fprintf(w, "%s", body)
}
