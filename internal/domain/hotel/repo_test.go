package hotel

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/YunchengHua/hotels-data-merge/internal/util"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRepo_GetFromSource1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHttpService := util.NewMockHttpService(ctrl)

	testHotel := []hotelFromSource1{
		{
			Id:            "test_id",
			Name:          "test_name",
			DestinationId: 123,
			// initialize the rest of the fields with test data
		},
	}

	data, _ := json.Marshal(testHotel) // handle error separately

	mockHttpService.EXPECT().Get(gomock.Any(), url1).Return(data, nil)

	repo := NewRepo(mockHttpService)

	resp, err := repo.GetFromSource1(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	// perform more specific assertions here
}
func TestRepo_GetFromSource2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHttpService := util.NewMockHttpService(ctrl)

	testHotel := []hotelFromSource2{
		{
			HotelID:       "test_id",
			HotelName:     "test_name",
			DestinationID: 123,
			// initialize the rest of the fields with test data
		},
	}

	data, _ := json.Marshal(testHotel) // handle error separately

	mockHttpService.EXPECT().Get(gomock.Any(), url2).Return(data, nil)

	repo := NewRepo(mockHttpService)

	resp, err := repo.GetFromSource2(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	// perform more specific assertions here
}

func TestRepo_GetFromSource3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHttpService := util.NewMockHttpService(ctrl)

	testHotel := []hotelFromSource3{
		{
			ID:          "test_id",
			Name:        "test_name",
			Destination: 123,
			// initialize the rest of the fields with test data
		},
	}

	data, _ := json.Marshal(testHotel) // handle error separately

	mockHttpService.EXPECT().Get(gomock.Any(), url3).Return(data, nil)

	repo := NewRepo(mockHttpService)

	resp, err := repo.GetFromSource3(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	// perform more specific assertions here
}
