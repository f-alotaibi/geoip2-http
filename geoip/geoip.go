package geoip

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

var dbSha256 string

func NewDB() (*geoip2.Reader, error) {
	accountID, ok := os.LookupEnv("GEOIP2_ACCOUNT_ID")
	if !ok {
		return nil, fmt.Errorf("could not find environment variable GEOIP2_ACCOUNT_ID")
	}
	licenseKey, ok := os.LookupEnv("GEOIP2_LICENSE_KEY")
	if !ok {
		return nil, fmt.Errorf("could not find environment variable GEOIP2_LICENSE_KEY")
	}
	authKey = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", accountID, licenseKey)))

	resp, err := getFile("https://download.maxmind.com/geoip/databases/GeoLite2-City/download?suffix=tar.gz")
	if err != nil {
		return nil, err
	}
	sha256resp, err := getFile("https://download.maxmind.com/geoip/databases/GeoLite2-City/download?suffix=tar.gz.sha256")
	if err != nil {
		return nil, err
	}
	v, err := io.ReadAll(sha256resp.Body)
	if err != nil {
		return nil, err
	}
	dbSha256 = strings.Split(string(v), " ")[0]
	dbData, err := findDB(resp.Body)
	if err != nil {
		return nil, err
	}
	db, err := geoip2.FromBytes(dbData)
	if err != nil {
		return nil, err
	}
	DBRaw = dbData
	runTicker()
	return db, nil
}
