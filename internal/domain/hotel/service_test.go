package hotel

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestService_FetchAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepo(ctrl)

	mockRepo.EXPECT().GetFromSource1(gomock.Any()).Return(getMockHotels1(), nil)
	mockRepo.EXPECT().GetFromSource2(gomock.Any()).Return(getMockHotels2(), nil)
	mockRepo.EXPECT().GetFromSource3(gomock.Any()).Return(getMockHotels3(), nil)

	service := NewService(mockRepo)
	hotels, err := service.FetchAll(context.Background())

	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if len(hotels) != 4 {
		t.Errorf("expected length of hotels to be 4, but got %d", len(hotels))
	}

	if hotels["1"].Name != "Hotel One" {
		t.Errorf("expected name of hotel 1 to be 'Hotel One', but got '%s'", hotels["1"].Name)
	}

	if hotels["2"].Description != "Hotel Two Updated Description" {
		t.Errorf("expected description of hotel 2 to be 'Hotel Two Updated Description', but got '%s'", hotels["2"].Description)
	}

	if hotels["3"].Name != "Hotel Three" {
		t.Errorf("expected name of hotel 3 to be 'Hotel Three', but got '%s'", hotels["3"].Name)
	}

	if hotels["4"].Description != "Hotel Four Description" {
		t.Errorf("expected description of hotel 4 to be 'Hotel Four Description', but got '%s'", hotels["4"].Description)
	}

}

func TestService_GetByHotelIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepo(ctrl)
	ctx := context.Background()
	s := NewService(mockRepo)

	hotelsPool = getMockHotels1()
	expectedHotels := getMockHotels1()
	hotelIDs := []string{"1", "2"}
	hotels := s.GetByHotelIDs(ctx, hotelIDs)
	// Assertions
	if len(hotels) != len(hotelIDs) {
		t.Errorf("Expected length of hotels to be %d, but got %d", len(hotelIDs), len(hotels))
	}

	for i, hotel := range hotels {
		if hotel.ID != hotelIDs[i] {
			t.Errorf("Expected hotel ID to be %s, but got %s", hotelIDs[i], hotel.ID)
		}

		if hotel.Name != expectedHotels[hotel.ID].Name {
			t.Errorf("Expected hotel name to be %s, but got %s", expectedHotels[hotel.ID].Name, hotel.Name)
		}

		if hotel.Description != expectedHotels[hotel.ID].Description {
			t.Errorf("Expected hotel description to be %s, but got %s", expectedHotels[hotel.ID].Description, hotel.Description)
		}
	}
}

func TestService_GetByDestinationIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepo(ctrl)
	ctx := context.Background()
	s := NewService(mockRepo)

	// Assuming hotelDestinations is global and already populated
	hotelDestinations = map[int64][]string{
		1: {"1", "2"},
		2: {"3", "4"},
	}

	expectedHotels := MergeHotelsMap(getMockHotels1(), getMockHotels2())
	expectedHotels = MergeHotelsMap(expectedHotels, getMockHotels3())
	hotelsPool = expectedHotels

	destinationIDs := []int64{1, 2}
	hotels := s.GetByDestinationIDs(ctx, destinationIDs)

	// Assertions
	if len(hotels) != len(destinationIDs)*2 {
		t.Errorf("Expected length of hotels to be %d, but got %d", len(destinationIDs)*2, len(hotels))
	}

	for _, hotel := range hotels {
		if expectedHotels[hotel.ID].Name != hotel.Name {
			t.Errorf("Expected hotel name to be %s, but got %s", expectedHotels[hotel.ID].Name, hotel.Name)
		}

		if expectedHotels[hotel.ID].Description != hotel.Description {
			t.Errorf("Expected hotel description to be %s, but got %s", expectedHotels[hotel.ID].Description, hotel.Description)
		}
	}
}

func TestFetchFn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepo(ctrl)
	mockRepo.EXPECT().GetFromSource1(gomock.Any()).Return(getMockHotels1(), nil).Times(1)
	mockRepo.EXPECT().GetFromSource2(gomock.Any()).Return(getMockHotels2(), nil).Times(1)
	mockRepo.EXPECT().GetFromSource3(gomock.Any()).Return(getMockHotels3(), nil).Times(1)

	service := NewService(mockRepo)

	err := fetchFn(context.Background(), service)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	mu.RLock()
	defer mu.RUnlock()
	if len(hotelsPool) == 0 || len(hotelDestinations) == 0 {
		t.Fatalf("Global state was not updated correctly")
	}
}

func getMockHotels1() map[string]*Hotel {
	return map[string]*Hotel{
		"1": {ID: "1", Name: "Hotel One", Description: "Hotel One Description"},
		"2": {ID: "2", Name: "Hotel Two", Description: "Hotel Two Description"},
	}
}

func getMockHotels2() map[string]*Hotel {
	return map[string]*Hotel{
		"2": {ID: "2", Name: "Hotel Two", Description: "Hotel Two Updated Description"},
		"3": {ID: "3", Name: "Hotel Three", Description: "Hotel Three Description"},
	}
}

func getMockHotels3() map[string]*Hotel {
	return map[string]*Hotel{
		"3": {ID: "3", Name: "Hotel Three", Description: "Hotel Three Updated Description"},
		"4": {ID: "4", Name: "Hotel Four", Description: "Hotel Four Description"},
	}
}
