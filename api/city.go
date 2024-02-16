package api

import (
	"encoding/json"
	"geoip2-http/geoip"
	"net"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func CityHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	println("RESPONSE", p.ByName("ip"), "FROM", r.RemoteAddr)
	w.Header().Set("Content-Type", "application/json")
	ip := p.ByName("ip")
	if strings.ToLower(ip) == "me" {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}
	city, err := geoip.DB.City(net.ParseIP(ip))
	if err != nil || city.Country.GeoNameID == 0 {
		w.Write([]byte("{\"Code\": \"404\"}"))
		return
	}
	response := geoip.CityResponse{Code: http.StatusOK, Response: city}
	json.NewEncoder(w).Encode(response)
}
