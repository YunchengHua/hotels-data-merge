package util

import (
	"context"
	"encoding/json"
	"errors"
	"io"
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

type HttpService interface {
	Get(ctx context.Context, url string) ([]byte, error)
}

func GetHandlerWrapper(
	exec func(context.Context, map[string]string) (interface{}, error),
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var res interface{}
		var jsonResponse []byte
		var rawReqParam = r.URL.Query()
		ctx := icontext.NewContextWithTraceID()
		defer func() {
			log.Info(ctx, "Incoming_Get_Request: req: %v, resp: %v, err: %v", rawReqParam, string(jsonResponse), err)
		}()

		if r.Method != http.MethodGet {
			err = errors.New("invalid request")
			log.Error(ctx, "Invalid request method: %v", r.Method)
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		reqParam := make(map[string]string, len(rawReqParam))
		for key := range rawReqParam {
			reqParam[key] = rawReqParam.Get(key)
		}
		res, err = exec(ctx, reqParam)
		if err != nil {
			log.Error(ctx, "Execute has error: %v", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the response header content type to JSON
		w.Header().Set("Content-Type", "application/json")
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

type httpService struct{}

func NewHttpService() *httpService {
	return &httpService{}
}

func (s httpService) Get(ctx context.Context, url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Error(ctx, "http get error:%v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(ctx, "io readAll error:%v", err)
		return nil, err
	}
	return body, nil
}
