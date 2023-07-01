package util

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
)

var (
	successHandlerFunc = func(context.Context, map[string][]string) (interface{}, error) {
		return nil, nil
	}
)

func TestGetHandlerWrapper_HappyCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fn := GetHandlerWrapper(successHandlerFunc)
	httpWriter := NewMockResponseWriter(ctrl)
	request := &http.Request{
		Method: "GET",
		URL:    &url.URL{},
	}
	httpWriter.EXPECT().Header().Return(http.Header{})
	httpWriter.EXPECT().Write(gomock.Any())
	fn(httpWriter, request)
}
