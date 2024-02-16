package geoip

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

func findDB(stream io.Reader) ([]byte, error) {
	uncompressedStream, err := gzip.NewReader(stream)
	if err != nil {
		return nil, err
	}
	tarReader := tar.NewReader(uncompressedStream)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}
		if header.Typeflag == tar.TypeReg {
			if strings.HasSuffix(header.Name, ".mmdb") {
				buf := new(bytes.Buffer)
				buf.ReadFrom(tarReader)
				return buf.Bytes(), nil
			}
		}
	}
	return nil, fmt.Errorf("database not found")
}

func NewDB() (*geoip2.Reader, error) {
	accountID, ok := os.LookupEnv("GEOIP2_ACCOUNT_ID")
	if !ok {
		return nil, fmt.Errorf("could not find environment variable GEOIP2_ACCOUNT_ID")
	}
	licenseKey, ok := os.LookupEnv("GEOIP2_LICENSE_KEY")
	if !ok {
		return nil, fmt.Errorf("could not find environment variable GEOIP2_LICENSE_KEY")
	}
	authorizationBase64 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", accountID, licenseKey)))
	req, _ := http.NewRequest("GET", "https://download.maxmind.com/geoip/databases/GeoLite2-City/download?suffix=tar.gz", nil)
	req.Header.Add("Authorization", "Basic "+authorizationBase64)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:122.0) Gecko/20100101 Firefox/122.0")
	cli := &http.Client{}
	resp, err := cli.Do(req)

	if err != nil {
		return nil, err
	}
	dbData, err := findDB(resp.Body)
	if err != nil {
		return nil, err
	}
	db, err := geoip2.FromBytes(dbData)
	if err != nil {
		return nil, err
	}
	return db, nil
}
