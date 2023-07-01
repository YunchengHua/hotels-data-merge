package util

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	icontext "github.com/YunchengHua/hotels-data-merge/internal/context"
	"github.com/YunchengHua/hotels-data-merge/internal/log"
)

//go:generate mockgen -destination ./util.mock.gen.go -source util.go -package util
type ResponseWriter interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

func GetHandlerWrapper(
	exec func(context.Context, map[string][]string) (interface{}, error),
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var res interface{}
		var reqParam = r.URL.Query()
		ctx := icontext.NewContextWithTraceID()
		defer func() {
			log.Info(ctx, "Incoming_Get_Request: req: %v, resp: %v, err: %v", reqParam, res, err)
		}()

		if r.Method != http.MethodGet {
			err = errors.New("invalid request")
			log.Error(ctx, "Invalid request method: %v", r.Method)
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		res, err = exec(ctx, reqParam)
		if err != nil {
			log.Error(ctx, "Execute has error: %v", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the response header content type to JSON
		w.Header().Set("Content-Type", "application/json")
		var jsonResponse []byte
		// Convert the request parameters to JSON
		jsonResponse, err = json.Marshal(res)
		if err != nil {
			log.Error(ctx, "Unable to marshal the response: %v", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the JSON response to the response writer
		_, err = w.Write(jsonResponse)
		if err != nil {
			log.Error(ctx, "Unable to send the response: %v", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
