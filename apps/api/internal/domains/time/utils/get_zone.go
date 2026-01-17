package utils

import "net/http"

func GetGeoLocationZone(r *http.Request) *string {
	zone := r.URL.Query().Get("zone")
	if zone == "" {
		zone = "Europe/London"
	}
	return &zone
}