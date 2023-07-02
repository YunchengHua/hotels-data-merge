package hotel

import (
	"strconv"

	"github.com/YunchengHua/hotels-data-merge/internal/config"
)

type Hotel struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	DestinationID     int64    `json:"destination_id"`
	Description       string   `json:"description"`
	BookingConditions []string `json:"booking_conditions"`
	Amenities         []string `json:"amenities"`
	Images            Images   `json:"images"`
	Location          Location `json:"location"`
	Missing           bool     `json:"missing"`
}

type Images struct {
	Amenities []Image `json:"amenities"`
	Rooms     []Image `json:"rooms"`
	Site      []Image `json:"site"`
}

type Image struct {
	URL     string `json:"url"`
	Caption string `json:"caption"`
}

type Location struct {
	Address    string  `json:"address"`
	City       string  `json:"city"`
	Country    string  `json:"country"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	PostalCode string  `json:"postal_code"`
}

func MergeHotelsMap(hotels1 map[string]*Hotel, hotels2 map[string]*Hotel) map[string]*Hotel {
	mergedHotels := make(map[string]*Hotel)

	// iterate over first map
	for id, hotel := range hotels1 {
		if _, exists := hotels2[id]; exists {
			// if hotel exists in second map, merge
			mergedHotels[id] = MergeHotels(hotel, hotels2[id])
		} else {
			// if hotel doesn't exist in second map, add to merged map
			mergedHotels[id] = hotel
		}
	}

	// iterate over second map to find hotels not in first map
	for id, hotel := range hotels2 {
		if _, exists := mergedHotels[id]; !exists {
			// if hotel doesn't exist in merged map, add it
			mergedHotels[id] = hotel
		}
	}

	return mergedHotels
}

// function to merge two hotels
func MergeHotels(hotel1 *Hotel, hotel2 *Hotel) *Hotel {
	// create a new hotel struct
	mergedHotel := &Hotel{}

	// id: no merging necessary
	mergedHotel.ID = hotel1.ID
	// destinationId: no merging necessary
	mergedHotel.DestinationID = hotel1.DestinationID

	// name: the longer name wins
	mergedHotel.Name = chooseString(hotel1.Name, hotel2.Name)

	// description: based on config
	if config.GetConfig().MergeDescriptions {
		mergedHotel.Description = hotel1.Description + "\n" + hotel2.Description
	} else {
		mergedHotel.Description = chooseString(hotel1.Description, hotel2.Description)
	}

	// bookingConditions: remove duplicates
	mergedHotel.BookingConditions = removeDuplicates(append(hotel1.BookingConditions, hotel2.BookingConditions...))

	// amenities: remove duplicates after trimming
	mergedHotel.Amenities = removeDuplicates(append(hotel1.Amenities, hotel2.Amenities...))

	// images: merge and remove duplicates
	mergedHotel.Images = mergeImages(hotel1.Images, hotel2.Images)

	// location: non-null values win, prefer 2 character country codes
	mergedHotel.Location = mergeLocations(hotel1.Location, hotel2.Location)

	return mergedHotel
}

// A function that removes duplicate strings from a slice.
func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// A function that merges two sets of images and removes duplicates.
func mergeImages(images1, images2 Images) Images {
	// Using map to remove duplicates.
	amenitiesMap := make(map[string]Image)
	for _, item := range append(images1.Amenities, images2.Amenities...) {
		amenitiesMap[item.URL] = item
	}
	roomsMap := make(map[string]Image)
	for _, item := range append(images1.Rooms, images2.Rooms...) {
		roomsMap[item.URL] = item
	}
	siteMap := make(map[string]Image)
	for _, item := range append(images1.Site, images2.Site...) {
		siteMap[item.URL] = item
	}

	// Converting map back to slice.
	amenities := make([]Image, 0, len(amenitiesMap))
	for _, item := range amenitiesMap {
		amenities = append(amenities, item)
	}
	rooms := make([]Image, 0, len(roomsMap))
	for _, item := range roomsMap {
		rooms = append(rooms, item)
	}
	site := make([]Image, 0, len(siteMap))
	for _, item := range siteMap {
		site = append(site, item)
	}

	return Images{Amenities: amenities, Rooms: rooms, Site: site}
}

func mergeLocations(loc1, loc2 Location) Location {
	return Location{
		Address:    chooseString(loc1.Address, loc2.Address),
		City:       chooseString(loc1.City, loc2.City),
		Country:    chooseString(loc1.Country, loc2.Country),
		Latitude:   chooseFloat(loc1.Latitude, loc2.Latitude),
		Longitude:  chooseFloat(loc1.Longitude, loc2.Longitude),
		PostalCode: chooseString(loc1.PostalCode, loc2.PostalCode),
	}
}

// A function that chooses the non-null or longer string.
func chooseString(str1, str2 string) string {
	if len(str1) > len(str2) {
		return str1
	} else if str2 != "" {
		return str2
	} else {
		return str1
	}
}

// A function that chooses the non-zero float.
func chooseFloat(val1, val2 float64) float64 {
	if val1 != 0.0 {
		return val1
	} else {
		return val2
	}
}

type hotelFromSource1 struct {
	Id            string         `json:"Id"`
	Name          string         `json:"Name"`
	DestinationId int64          `json:"DestinationId"`
	Description   string         `json:"Description"`
	Facilities    []string       `json:"Facilities"`
	Address       string         `json:"Address"`
	City          string         `json:"City"`
	Country       string         `json:"Country"`
	Latitude      StringOrNumber `json:"Latitude"`
	Longitude     StringOrNumber `json:"Longitude"`
	PostalCode    string         `json:"PostalCode"`
}

type StringOrNumber float64

func (son *StringOrNumber) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s[0] == '"' {
		s = s[1 : len(s)-1]
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	*son = StringOrNumber(v)
	return nil
}

type hotelFromSource2 struct {
	HotelID           string            `json:"hotel_id"`
	DestinationID     int64             `json:"destination_id"`
	HotelName         string            `json:"hotel_name"`
	Location          LocationResponse2 `json:"location"`
	Details           string            `json:"details"`
	Amenities         AmenityResponse2  `json:"amenities"`
	Images            ImagesResponse2   `json:"images"`
	BookingConditions []string          `json:"booking_conditions"`
}

type AmenityResponse2 struct {
	General []string `json:"general"`
	Room    []string `json:"room"`
}

type ImagesResponse2 struct {
	Rooms []ImageResponse2 `json:"rooms"`
	Site  []ImageResponse2 `json:"site"`
}

type ImageResponse2 struct {
	Link    string `json:"link"`
	Caption string `json:"caption"`
}

type LocationResponse2 struct {
	Address string `json:"address"`
	Country string `json:"country"`
}

type hotelFromSource3 struct {
	ID          string          `json:"id"`
	Destination int64           `json:"destination"`
	Name        string          `json:"name"`
	Lat         float64         `json:"lat"`
	Lng         float64         `json:"lng"`
	Address     *string         `json:"address"`
	Info        *string         `json:"info"`
	Amenities   []string        `json:"amenities"`
	Images      ImagesResponse3 `json:"images"`
}

type ImagesResponse3 struct {
	Rooms     []ImageResponse3 `json:"rooms"`
	Amenities []ImageResponse3 `json:"amenities"`
}

type ImageResponse3 struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

func getFromStringPointer(t *string) string {
	if t == nil {
		return ""
	}
	return *t
}
