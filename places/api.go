//Package places is a simple wrapper for Google Places APIs
package places

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"google.golang.org/genproto/googleapis/type/latlng"
)

const (
	//GooglePlacesURL is the Places API URL
	GooglePlacesURL = "https://maps.googleapis.com/maps/api/place/findplacefromtext/json"
	//GooglePlacesDetailsURL is Details API URL
	GooglePlacesDetailsURL = "https://maps.googleapis.com/maps/api/place/details/json"
)

//GooglePlacesAPIKey is the API KEY
var GooglePlacesAPIKey = os.Getenv("GOOGLE_APIKEY")

//Location see https://developers.google.com/places/web-service/search#FindPlaceRequests
type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

//AsLatLng returns Location as google.golang.org/genproto/googleapis/type/latlng
func (loc Location) AsLatLng() *latlng.LatLng {
	return &latlng.LatLng{Latitude: loc.Latitude, Longitude: loc.Longitude}
}

//Geometry see https://developers.google.com/places/web-service/search#FindPlaceRequests
type Geometry struct {
	Location Location `json:"location"`
}

// SearchResult see https://developers.google.com/places/web-service/search#FindPlaceRequests
type SearchResult struct {
	FormattedAddress string   `json:"formatted_address"`
	Name             string   `json:"name"`
	PlaceID          string   `json:"place_id"`
	Geometry         Geometry `json:"geometry"`
}

//GetMunicipally returns the city
func (r SearchResult) GetMunicipally() string {

	parts := strings.Split(r.FormattedAddress, ",")
	if len(parts) != 3 {
		return ""
	}
	city := regexp.MustCompile(`[A-ZÅÄÖa-zåäö]+`)
	return city.FindString(parts[1])
}

// SearchResults is https://developers.google.com/places/web-service/search#FindPlaceRequests
type SearchResults struct {
	Candidates []SearchResult `json:"candidates"`
	Status     string         `json:"status"`
}

// SearchQuery is a search query
type SearchQuery struct {
	Language, LocationBias, Input string
}

//DetailsQuery from Google places
type DetailsQuery struct {
	PlaceID string
}

//AddressComponent from Google places
type AddressComponent struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

//DetailsResult from Google places
type DetailsResult struct {
	FormattedAddress  string             `json:"formatted_address"`
	Name              string             `json:"name"`
	PlaceID           string             `json:"place_id"`
	Geometry          Geometry           `json:"geometry"`
	AddressComponents []AddressComponent `json:"address_components"`
}

//DetailsResults from Google places
type DetailsResults struct {
	Result DetailsResult
}

//Do executes a DetailsQuery query
func (q DetailsQuery) Do() DetailsResults {
	params := url.Values{}
	params.Set("key", GooglePlacesAPIKey)
	params.Set("place_id", q.PlaceID)
	params.Set("region", "se")
	params.Set("fields", "address_component,vicinity")

	u := fmt.Sprintf("%s?%s", GooglePlacesDetailsURL, params.Encode())
	resp, err := http.Get(u)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)

	var result DetailsResults
	json.Unmarshal(data, &result)

	return result
}

//Do executes a SearchQuery
func (q SearchQuery) Do() SearchResults {

	params := url.Values{}
	params.Set("key", GooglePlacesAPIKey)
	params.Set("inputtype", "textquery")
	params.Set("language", q.Language)
	params.Set("locationbias", q.LocationBias)
	params.Set("input", q.Input)
	params.Set("fields", "formatted_address,name,geometry,place_id")

	u := fmt.Sprintf("%s?%s", GooglePlacesURL, params.Encode())
	resp, err := http.Get(u)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)

	var result SearchResults
	json.Unmarshal(data, &result)

	return result

}
