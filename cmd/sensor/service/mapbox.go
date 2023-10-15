package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

const baseURL = "https://api.mapbox.com/geocoding/v5/mapbox.places/"

type MapBox interface {
	FindLatLon(location string) (*Location, error)
}

type mapBox struct {
	apiKey string
}

func NewMapBox(apiKey string) *mapBox {
	return &mapBox{
		apiKey: apiKey,
	}
}

func (m mapBox) FindLatLon(location string) (*Location, error) {
	path := fmt.Sprintf("%s.json?access_token=%s", location, m.apiKey)
	resp, getErr := http.Get(baseURL + path)
	if getErr != nil {
		log.Fatal(getErr)
	}
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	return parseMapboxGeocode(string(body))
}

func parseMapboxGeocode(body string) (*Location, error) {
	var result map[string]any
	err := json.Unmarshal([]byte(body), &result)
	if err != nil {
		return nil, err
	}
	// The object stored in the "birds" key is also stored as
	// a map[string]any type, and its type is asserted from
	// the `any` type
	features := result["features"].([]any)
	if len(features) > 0 {
		feature := features[0].(map[string]any)
		center := feature["center"].([]any)
		if len(center) != 2 {
			return nil, errors.New("Not found")
		}
		return &Location{
			Lat: fmt.Sprint(center[1].(float64)),
			Lon: fmt.Sprint(center[0].(float64)),
		}, nil
	} else {
		return nil, errors.New("Not found")
	}
}
