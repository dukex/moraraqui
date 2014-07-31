package models

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// type Contract int

const (
	Undefined int = iota
	House
	Apartament
)

// const (
//   Sale = Contract iota
//   Rent= Contract iota
// )

type Property struct {
	Id           int64  `json:"id"`
	Title        string `sql:"not null" json:"title"`
	Address      string `sql:"not null" json:"address"`
	City         string
	State        string
	Lat          float64 `json:"lat"`
	Lng          float64 `json:"lng"`
	Url          string  `sql:"not null;unique" json:"url"`
	RealStateId  int64   `json:"-"`
	Type         int     `sql:"type(integer);" json:"type"`
	Size         int     `json:"size"`
	Bedroom      int     `json:"bedroom"`
	Value        float64 `json:"value"`
	Neighborhood string  `json:"neighborhood"`
}

func (p Property) TableName() string {
	return "properties"
}

func (p *Property) FullAddress() string {
	return p.Address + "," + p.City + "," + p.State
}

func (p *Property) BeforeCreate() error {
	if p.Address != "" {
		var addressLocation []float64

		l := strings.ToLower(p.Address)

		if !(strings.Contains(l, "rua") ||
			strings.Contains(l, "r ") ||
			strings.Contains(l, "avenida") ||
			strings.Contains(l, "av") ||
			strings.Contains(l, "jardim") ||
			strings.Contains(l, "praca") ||
			strings.Contains(l, "pç") ||
			strings.Contains(l, "vila")) {
			p.Address = "Rua " + p.Address
		}

		if strings.Contains(l, "sob consulta") ||
			strings.Contains(l, "nao informado") ||
			strings.Contains(l, "não informado") ||
			p.Address == "Rua "+p.Neighborhood {
			p.Address = p.Neighborhood
		}

		addressLocationI, err := CacheGet(p.FullAddress(), addressLocation)
		addressLocation, ok := addressLocationI.([]float64)

		if err != nil || !ok {
			geoUrlRoot := "https://maps.googleapis.com/maps/api/geocode/json"
			geoUrl, _ := url.Parse(geoUrlRoot)
			geoParams := url.Values{}
			geoParams.Add("address", p.FullAddress())
			geoUrl.RawQuery = geoParams.Encode()
			resGeoFetch, _ := http.Get(geoUrl.String())
			defer resGeoFetch.Body.Close()
			contents, _ := ioutil.ReadAll(resGeoFetch.Body)
			var geo Geolocation
			json.Unmarshal(contents, &geo)

			if len(geo.Results) > 0 {
				location := geo.Results[0].Geometry.Location
				addressLocation = []float64{location.Lat, location.Lng}
				CacheSet(p.FullAddress(), addressLocation)
			}
		}

		if len(addressLocation) == 2 {
			p.Lat = addressLocation[0]
			p.Lng = addressLocation[1]
		}
	} else {
		p.Address = p.Neighborhood
	}

	return nil
}

type Geolocation struct {
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Bounds struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"bounds"`
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport     struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		Types []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}
