package handler

import (
	"context"
	"testing"

	"github.com/YunchengHua/hotels-data-merge/internal/domain/hotel"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHotelsHandler_Execute_HotelID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHotelService := hotel.NewMockService(ctrl)
	handler := &HotelsHandler{
		hotelService: mockHotelService,
	}

	testHotel := &hotel.Hotel{ID: "testHotel"}
	mockHotelService.EXPECT().GetByHotelIDs(gomock.Any(), []string{"testHotel"}).Return([]*hotel.Hotel{testHotel}).Times(1)

	params := map[string]string{ParamNameHotelIDs: `"testHotel"`}

	resp, err := handler.Execute(context.Background(), params)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	hotelsResponse, ok := resp.(*HotelsResponse)
	assert.True(t, ok)

	assert.Equal(t, 1, len(hotelsResponse.Hotels))
	assert.Equal(t, "testHotel", hotelsResponse.Hotels[0].ID)
}

func TestHotelsHandler_Execute_DestinationIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHotelService := hotel.NewMockService(ctrl)
	handler := &HotelsHandler{
		hotelService: mockHotelService,
	}

	testHotel := &hotel.Hotel{DestinationID: 123}
	mockHotelService.EXPECT().GetByDestinationIDs(gomock.Any(), []int64{123}).Return([]*hotel.Hotel{testHotel}).Times(1)

	params := map[string]string{ParamNameDestinationIDs: `123`}

	resp, err := handler.Execute(context.Background(), params)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	hotelsResponse, ok := resp.(*HotelsResponse)
	assert.True(t, ok)

	assert.Equal(t, 1, len(hotelsResponse.Hotels))
	assert.Equal(t, int64(123), hotelsResponse.Hotels[0].DestinationID)
}
