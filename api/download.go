package api

import (
	"geoip2-http/geoip"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func DownloadHander(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if geoip.DB == nil {
		http.Error(w, "DB Not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename=GeoLite2-City.mmdb")
	w.Header().Set("Content-type", "text/mmdb")
	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	_, err := w.Write([]byte(geoip.DBRaw))
	if err != nil {
		http.Error(w, "Could not write data", http.StatusInternalServerError)
		return
	}
}
