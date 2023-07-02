package handler

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/YunchengHua/hotels-data-merge/internal/domain/hotel"
	"github.com/YunchengHua/hotels-data-merge/internal/util"
)

type HotelsResponse struct {
	Hotels []*hotel.Hotel `json:"hotels"`
}

const (
	paramNameHotelIDs       = "hotel_ids"
	paramNameDestinationIDs = "destination_ids"
)

func NewDefaultHotelsHandler() *HotelsHandler {
	return &HotelsHandler{
		hotelService: hotel.NewService(hotel.NewRepo(util.NewHttpService())),
	}
}

type HotelsHandler struct {
	hotelService hotel.Service
}

func (h *HotelsHandler) Execute(
	ctx context.Context, params map[string]string,
) (interface{}, error) {
	hotelIDs, destinationIDs, err := h.getFromRequest(ctx, params)
	if err != nil {
		return nil, err
	}
	if len(hotelIDs) != 0 {
		return &HotelsResponse{
			Hotels: h.hotelService.GetByHotelIDs(ctx, hotelIDs),
		}, nil
	}
	if len(destinationIDs) != 0 {
		return &HotelsResponse{
			Hotels: h.hotelService.GetByDestinationIDs(ctx, destinationIDs),
		}, nil
	}
	return nil, nil
}

func (h *HotelsHandler) getFromRequest(
	ctx context.Context, params map[string]string) ([]string, []int64, error) {
	if raw, exist := params[paramNameHotelIDs]; exist {
		raw = strings.ReplaceAll(raw, `"`, "") // Remove quotes
		hotelIds := strings.Split(raw, `,`)    // Split into array
		return hotelIds, nil, nil
	}
	if raw, exist := params[paramNameDestinationIDs]; exist {
		raws := strings.Split(raw, `,`) // Split into array
		res := make([]int64, 0, len(raws))
		for _, raw := range raws {
			n, err := strconv.ParseInt(raw, 10, 64)
			if err != nil {
				return nil, nil, err
			}
			res = append(res, n)
		}
		return nil, res, nil
	}
	return nil, nil, errors.New("both hotel_ids and destination_ids not provided")
}
