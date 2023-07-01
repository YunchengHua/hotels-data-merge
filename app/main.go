package main

import (
	"net/http"

	"github.com/YunchengHua/hotels-data-merge/internal/context"
	"github.com/YunchengHua/hotels-data-merge/internal/handler"
	"github.com/YunchengHua/hotels-data-merge/internal/log"
	"github.com/YunchengHua/hotels-data-merge/internal/util"
)

func main() {
	ctx := context.NewContextWithTraceID()
	log.MustInit()
	log.Info(ctx, "System starting...")

	hotelsHandler := handler.NewDefaultHotelsHandler()
	http.HandleFunc("/hotels", util.GetHandlerWrapper(hotelsHandler.Execute))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Error(ctx, "http serve error: %v", err.Error())
	}

	log.Info(ctx, "System shutdown")
}
