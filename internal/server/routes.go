package server

import (
	"geoip2-http/api"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := httprouter.New()
	r.GET("/city/:ip", api.CityHandler)
	r.GET("/country/:ip", api.CountryHandler)
	r.GET("/download", api.DownloadHander)

	return r
}
