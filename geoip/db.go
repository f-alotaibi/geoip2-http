package geoip

import "github.com/oschwald/geoip2-golang"

var DB *geoip2.Reader

func ConnectDB() error {
	db, err := NewDB()
	if err != nil {
		return err
	}
	DB = db
	return nil
}
