package places

import (
	"fmt"
	"testing"
)

func TestPlacesSearchA(t *testing.T) {

	const (
		lat = "55.5711163"
		lng = "12.8757088"
	)

	query := SearchQuery{
		LocationBias: fmt.Sprintf("point:%s,%s", lat, lng),
		Input:        "TYGELSJÖVÄGEN 129",
		Language:     "sv",
	}

	result := query.Do()
	if result.Status != "OK" {
		t.Error(result.Status)
	}

	for _, r := range result.Candidates {
		fmt.Println(r.Geometry.Location.AsLatLng().String())
		fmt.Println(r.GetMunicipally())
	}

}

func TestPlacesSearchB(t *testing.T) {

	const (
		lat = "55.5711163"
		lng = "12.8757088"
	)

	query := SearchQuery{
		LocationBias: fmt.Sprintf("point:%s,%s", lat, lng),
		Input:        "Håkanstorpsvägen 31",
		Language:     "sv",
	}

	result := query.Do()
	if result.Status != "OK" {
		t.Error(result.Status)
	}

	for _, r := range result.Candidates {
		fmt.Println(r.Geometry.Location.AsLatLng().String())
		fmt.Println(r.GetMunicipally())
	}

}

func TestPlacesDetails(t *testing.T) {
	query := DetailsQuery{
		PlaceID: "ChIJWUpMSZugU0YReEaxgilNfdM",
	}

	result := query.Do()
	fmt.Println(result)
}
