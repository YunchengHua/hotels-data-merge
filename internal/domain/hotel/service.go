package hotel

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/YunchengHua/hotels-data-merge/internal/config"
	"github.com/YunchengHua/hotels-data-merge/internal/log"
)

var (
	hotelsPool        map[string]*Hotel
	hotelDestinations map[int64][]string
	mu                = &sync.RWMutex{}
)

func MustInit(ctx context.Context, hotelService Service) {
	err := FetchFn(ctx, hotelService)
	if err != nil {
		fmt.Println("Fail to fetch hotels", err.Error())
		log.Fatal(ctx, "Fail to fetch hotels: %v", err.Error())
	}

	ticker := time.NewTicker(time.Duration(config.GetConfig().DataSourcesRefreshIntervalSec) * time.Second)
	go func() {
		for {
			t := <-ticker.C
			err := FetchFn(ctx, hotelService)
			if err != nil {
				log.Error(ctx, "fetch error, skip this time: %v: %v", t, err.Error())
			}
		}
	}()
}

func FetchFn(ctx context.Context, hotelService Service) error {
	var err error
	var tempHotels map[string]*Hotel
	tempHotelDestinations := make(map[int64][]string)
	tempHotels, err = hotelService.FetchAll(ctx)
	if err != nil {
		return err
	}

	for _, hotel := range tempHotels {
		if _, exist := tempHotelDestinations[hotel.DestinationID]; !exist {
			tempHotelDestinations[hotel.DestinationID] = make([]string, 0)
		}
		tempHotelDestinations[hotel.DestinationID] = append(tempHotelDestinations[hotel.DestinationID], hotel.ID)
		log.Info(ctx, "Fetch the hotels info: hotel: %v", hotel)
	}
	mu.Lock()
	defer mu.Unlock()
	hotelsPool = tempHotels
	hotelDestinations = tempHotelDestinations

	log.Info(ctx, "Fetch the hotels info: hotel_destinations: %v", tempHotelDestinations)
	return nil
}

func NewService(repo Repo) *service {
	return &service{
		repo: repo,
	}
}

//go:generate mockgen -destination ./service.mock.gen.go -source service.go -package hotel
type Service interface {
	FetchAll(ctx context.Context) (map[string]*Hotel, error)
	GetByHotelIDs(ctx context.Context, hotelIDs []string) []*Hotel
	GetByDestinationIDs(ctx context.Context, destinationIDs []int64) []*Hotel
}

type service struct {
	repo Repo
}

func (s *service) FetchAll(ctx context.Context) (map[string]*Hotel, error) {
	wg := sync.WaitGroup{}
	var hotels1 map[string]*Hotel
	var hotels2 map[string]*Hotel
	var hotels3 map[string]*Hotel

	go func() {
		var err error
		wg.Add(1)
		defer func() {
			wg.Done()
		}()
		hotels1, err = s.repo.GetFromSource1(ctx)
		if err != nil {
			log.Error(ctx, "fetch err: %v", err)
		}
	}()
	go func() {
		var err error
		wg.Add(1)
		defer func() {
			wg.Done()
		}()
		hotels2, err = s.repo.GetFromSource2(ctx)
		if err != nil {
			log.Error(ctx, "fetch err: %v", err)
		}
	}()
	go func() {
		var err error
		wg.Add(1)
		defer func() {
			wg.Done()
		}()
		hotels3, err = s.repo.GetFromSource3(ctx)
		if err != nil {
			log.Error(ctx, "fetch err: %v", err)
		}
	}()
	wg.Wait()

	return MergeHotelsMap(MergeHotelsMap(hotels1, hotels2), hotels3), nil
}

func (s *service) GetByHotelIDs(ctx context.Context, hotelIDs []string) []*Hotel {
	mu.RLock()
	defer mu.RUnlock()
	res := make([]*Hotel, 0, len(hotelIDs))
	for _, id := range hotelIDs {
		hotel, exist := hotelsPool[id]
		if !exist {
			hotel = &Hotel{ID: id, Missing: true}
		}
		res = append(res, hotel)
	}
	return res
}

func (s *service) GetByDestinationIDs(ctx context.Context, destinationIDs []int64) []*Hotel {
	mu.RLock()
	defer mu.RUnlock()
	res := make([]*Hotel, 0, len(destinationIDs))
	for _, id := range destinationIDs {
		hotelIDs, exist := hotelDestinations[id]
		if !exist {
			continue
		}
		res = append(res, s.GetByHotelIDs(ctx, hotelIDs)...)
	}
	return res
}
