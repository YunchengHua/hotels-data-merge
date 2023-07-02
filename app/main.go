package main

import (
	"fmt"
	"net/http"

	"github.com/YunchengHua/hotels-data-merge/internal/config"
	"github.com/YunchengHua/hotels-data-merge/internal/context"
	"github.com/YunchengHua/hotels-data-merge/internal/domain/hotel"
	"github.com/YunchengHua/hotels-data-merge/internal/handler"
	"github.com/YunchengHua/hotels-data-merge/internal/log"
	"github.com/YunchengHua/hotels-data-merge/internal/util"
)

func main() {
	ctx := context.NewContextWithTraceID()
	log.MustInit()
	fmt.Println("System starting...")
	config.MustInit()
	hotelService := hotel.NewService(hotel.NewRepo(util.NewHttpService()))
	hotel.MustInit(ctx, hotelService)

	hotelsHandler := handler.NewDefaultHotelsHandler()
	http.HandleFunc("/hotels", util.GetHandlerWrapper(hotelsHandler.Execute))
	fmt.Println("Preparation Done")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("http serve error: %v\n", err.Error())
	}

	fmt.Println("System shutdown")
}
