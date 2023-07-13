package hotel

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/YunchengHua/hotels-data-merge/internal/log"
	"github.com/YunchengHua/hotels-data-merge/internal/util"
)

const (
	Url1 = "http://www.mocky.io/v2/5ebbea002e000054009f3ffc"
	Url2 = "http://www.mocky.io/v2/5ebbea102e000029009f3fff"
	Url3 = "http://www.mocky.io/v2/5ebbea1f2e00002b009f4000"
)

//go:generate mockgen -destination ./repo.mock.gen.go -source repo.go -package hotel
type Repo interface {
	GetFromSource1(ctx context.Context) (map[string]*Hotel, error)
	GetFromSource2(ctx context.Context) (map[string]*Hotel, error)
	GetFromSource3(ctx context.Context) (map[string]*Hotel, error)
}

func NewRepo(httpService util.HttpService) *repo {
	return &repo{
		httpService: httpService,
	}
}

type repo struct {
	httpService util.HttpService
}

func (r *repo) GetFromSource1(ctx context.Context) (map[string]*Hotel, error) {
	raw, err := r.httpService.Get(ctx, Url1)
	if err != nil {
		return nil, err
	}
	fmt.Println(hex.EncodeToString(raw))

	var responseHotels []hotelFromSource1
	err = json.Unmarshal(raw, &responseHotels)
	if err != nil {
		log.Error(ctx, "fail to unmarshl response: %v", err.Error())
		return nil, err
	}

	res := make(map[string]*Hotel, len(responseHotels))
	for _, responseHotel := range responseHotels {
		// Convert the response hotel data to the target hotel struct
		hotel := &Hotel{
			ID:            responseHotel.Id,
			Name:          responseHotel.Name,
			DestinationID: responseHotel.DestinationId,
			Description:   responseHotel.Description,
			Amenities:     responseHotel.Facilities,
			Images: Images{
				Amenities: []Image{},
				Rooms:     []Image{},
				Site:      []Image{},
			},
			Location: Location{
				Address:    responseHotel.Address,
				City:       responseHotel.City,
				Country:    responseHotel.Country,
				Latitude:   float64(responseHotel.Latitude),
				Longitude:  float64(responseHotel.Longitude),
				PostalCode: responseHotel.PostalCode,
			},
		}
		res[hotel.ID] = hotel
	}

	return res, nil
}

func (r *repo) GetFromSource2(ctx context.Context) (map[string]*Hotel, error) {
	raw, err := r.httpService.Get(ctx, Url2)
	if err != nil {
		return nil, err
	}

	var responseHotels []hotelFromSource2
	err = json.Unmarshal(raw, &responseHotels)
	if err != nil {
		log.Error(ctx, "fail to unmarshl response: %v", err.Error())
		return nil, err
	}

	res := make(map[string]*Hotel, len(responseHotels))
	for _, hotel := range responseHotels {
		var amenities []string
		amenities = append(amenities, append(hotel.Amenities.General, hotel.Amenities.Room...)...)

		var roomImages, siteImages []Image
		for _, img := range hotel.Images.Rooms {
			roomImages = append(roomImages, Image{URL: img.Link, Caption: img.Caption})
		}
		for _, img := range hotel.Images.Site {
			siteImages = append(siteImages, Image{URL: img.Link, Caption: img.Caption})
		}

		res[hotel.HotelID] = &Hotel{
			ID:                hotel.HotelID,
			Name:              hotel.HotelName,
			DestinationID:     hotel.DestinationID,
			Description:       hotel.Details,
			BookingConditions: hotel.BookingConditions,
			Amenities:         amenities,
			Images: Images{
				Rooms: roomImages,
				Site:  siteImages,
			},
			Location: Location{
				Address: hotel.Location.Address,
				Country: hotel.Location.Country,
			},
		}
	}
	return res, nil
}

func (r *repo) GetFromSource3(ctx context.Context) (map[string]*Hotel, error) {
	raw, err := r.httpService.Get(ctx, Url3)
	if err != nil {
		return nil, err
	}

	var responseHotels []hotelFromSource3
	err = json.Unmarshal(raw, &responseHotels)
	if err != nil {
		log.Error(ctx, "fail to unmarshl response: %v", err.Error())
		return nil, err
	}

	res := make(map[string]*Hotel, len(responseHotels))
	for _, hotel := range responseHotels {
		var rooms []Image
		for _, roomImg := range hotel.Images.Rooms {
			rooms = append(rooms, Image{URL: roomImg.URL, Caption: roomImg.Description})
		}

		var amenityImgs []Image
		for _, amenityImg := range hotel.Images.Amenities {
			amenityImgs = append(amenityImgs, Image{URL: amenityImg.URL, Caption: amenityImg.Description})
		}

		res[hotel.ID] = &Hotel{
			ID:            hotel.ID,
			Name:          hotel.Name,
			DestinationID: hotel.Destination,
			Description:   getFromStringPointer(hotel.Info),
			Amenities:     hotel.Amenities,
			Images: Images{
				Amenities: amenityImgs,
				Rooms:     rooms,
			},
			Location: Location{
				Address:   getFromStringPointer(hotel.Address),
				Latitude:  hotel.Lat,
				Longitude: hotel.Lng,
			},
		}
	}

	return res, nil
}
