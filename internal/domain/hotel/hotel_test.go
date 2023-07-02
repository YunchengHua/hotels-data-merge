package hotel

import "testing"

func TestRemoveDuplicates(t *testing.T) {
	input := []string{"a", "b", "a", "c", "b", "c", "d"}
	expected := []string{"a", "b", "c", "d"}
	result := removeDuplicates(input)
	if !slicesEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestChooseString(t *testing.T) {
	str1 := "short"
	str2 := "longer string"
	if chooseString(str1, str2) != str2 {
		t.Errorf("Expected %v, but got %v", str2, chooseString(str1, str2))
	}
}

func TestChooseFloat(t *testing.T) {
	val1 := 0.0
	val2 := 1.1
	if chooseFloat(val1, val2) != val2 {
		t.Errorf("Expected %v, but got %v", val2, chooseFloat(val1, val2))
	}
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestMergeHotels(t *testing.T) {
	hotel1 := &Hotel{
		ID:                "1",
		Name:              "Hotel1",
		DestinationID:     1,
		Description:       "Hotel1 Description",
		BookingConditions: []string{"Condition1", "Condition2"},
		Amenities:         []string{"Amenity1", "Amenity2"},
		// Add more fields as needed
	}
	hotel2 := &Hotel{
		ID:                "1",
		Name:              "Hotel One", // Longer name
		DestinationID:     1,
		Description:       "Hotel One Description",              // Different description
		BookingConditions: []string{"Condition1", "Condition3"}, // Different booking conditions
		Amenities:         []string{"Amenity2", "Amenity3"},     // Different amenities
		// Add more fields as needed
	}
	merged := MergeHotels(hotel1, hotel2)
	// Assert that merged fields are as expected.
	if merged.Name != "Hotel One" {
		t.Errorf("Expected Name %v, but got %v", "Hotel One", merged.Name)
	}
	if merged.Description != "Hotel One Description" {
		t.Errorf("Expected Description %v, but got %v", "Hotel One Description", merged.Description)
	}
	expectedConditions := []string{"Condition1", "Condition2", "Condition3"}
	if !slicesEqual(merged.BookingConditions, expectedConditions) {
		t.Errorf("Expected BookingConditions %v, but got %v", expectedConditions, merged.BookingConditions)
	}
	expectedAmenities := []string{"Amenity1", "Amenity2", "Amenity3"}
	if !slicesEqual(merged.Amenities, expectedAmenities) {
		t.Errorf("Expected Amenities %v, but got %v", expectedAmenities, merged.Amenities)
	}
	// Add more assertions as needed
}

func TestMergeHotelsMap(t *testing.T) {
	hotel1 := &Hotel{
		ID:                "1",
		Name:              "Hotel1",
		DestinationID:     1,
		Description:       "Hotel1 Description",
		BookingConditions: []string{"Condition1", "Condition2"},
		Amenities:         []string{"Amenity1", "Amenity2"},
		// Add more fields as needed
	}
	hotel2 := &Hotel{
		ID:                "2",
		Name:              "Hotel Two",
		DestinationID:     2,
		Description:       "Hotel Two Description",
		BookingConditions: []string{"Condition3", "Condition4"},
		Amenities:         []string{"Amenity3", "Amenity4"},
		// Add more fields as needed
	}

	hotels1 := map[string]*Hotel{
		"1": hotel1,
	}
	hotels2 := map[string]*Hotel{
		"1": hotel1,
		"2": hotel2,
	}

	merged := MergeHotelsMap(hotels1, hotels2)

	if len(merged) != 2 {
		t.Errorf("Expected length %v, but got %v", 2, len(merged))
	}

	if _, exists := merged["2"]; !exists {
		t.Errorf("Expected hotel with ID 2 to exist in merged map, but it did not")
	}

}
