package geoip

import "github.com/oschwald/geoip2-golang"

type CityResponse struct {
	Code     int          `json:"Code"`
	Response *geoip2.City `json:"Response"`
}

type CountryResponse struct {
	Code     int             `json:"Code"`
	Response *geoip2.Country `json:"Response"`
}
