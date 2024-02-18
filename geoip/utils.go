package geoip

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/oschwald/geoip2-golang"
)

var authKey string

func getFile(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Basic "+authKey)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:122.0) Gecko/20100101 Firefox/122.0")
	cli := &http.Client{}
	resp, err := cli.Do(req)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

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

func runTicker() {
	ticker := time.NewTicker((7 * 24) * time.Hour) // A week
	go func() {
		for {
			t, ok := <-ticker.C
			if !ok {
				fmt.Println("Error, ticker stopped")
				break
			}
			fmt.Println("Checking for updates at", t)
			// check for updates using the db's sha256
			sha256resp, err := getFile("https://download.maxmind.com/geoip/databases/GeoLite2-City/download?suffix=tar.gz.sha256")
			if err != nil {
				fmt.Println("An error occured while checking for updates (could not get sha256)")
				break
			}
			v, err := io.ReadAll(sha256resp.Body)
			if err != nil {
				fmt.Println("An error occured while checking for updates (could not read sha256 body)")
				break
			}
			sha256file := strings.Split(string(v), " ")[0]
			if dbSha256 == sha256file {
				continue
			}
			// update
			fmt.Println("An update found at", t, "updating...")
			resp, err := getFile("https://download.maxmind.com/geoip/databases/GeoLite2-City/download?suffix=tar.gz")
			if err != nil {
				fmt.Println("An error occured while updating (could not read body)")
				break
			}
			dbData, err := findDB(resp.Body)
			if err != nil {
				fmt.Println("An error occured while updating (db not found)")
				break
			}
			db, err := geoip2.FromBytes(dbData)
			if err != nil {
				fmt.Println("An error occured while updating (invalid db)")
				break
			}
			fmt.Println("Updated")
			DBRaw = dbData
			DB = db
		}
	}()
}
