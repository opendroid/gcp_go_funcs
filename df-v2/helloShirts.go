package df_v2

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("HelloShirts", helloShirts)
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
