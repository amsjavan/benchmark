package standard

import (
	"fmt"
	"net/http"
	"time"
)

func init() {

}

func RunStandardServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(400 * time.Millisecond)
		fmt.Fprintf(w, "Welcome to my website!")
	})

	http.ListenAndServe(":9090", nil)
}
