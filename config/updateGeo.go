package config

import (
	"fmt"
	"github.com/Dreamacro/clash/component/geodata"
	_ "github.com/Dreamacro/clash/component/geodata/standard"
	C "github.com/Dreamacro/clash/constant"
	"io/ioutil"
	"net/http"
	"runtime"
)

func UpdateGeoDatabases() error {
	defer runtime.GC()
	geoLoader, err := geodata.GetGeoDataLoader("standard")
	if err != nil {
		return err
	}

	geoip, err := downloadForBytes(C.GeoIpUrl)
	if err != nil {
		return fmt.Errorf("can't download GeoIP database file: %w", err)
	}

	if _, err = geoLoader.LoadIPByBytes(geoip, "cn"); err != nil {
		return fmt.Errorf("invalid GeoIP database file: %s", err)
	}

	if saveFile(geoip, C.Path.GeoIP()) != nil {
		return fmt.Errorf("can't save GeoIP database file: %w", err)
	}

	geosite, err := downloadForBytes(C.GeoSiteUrl)
	if err != nil {
		return fmt.Errorf("can't download GeoSite database file: %w", err)
	}

	if _, err = geoLoader.LoadSiteByBytes(geosite, "cn"); err != nil {
		return fmt.Errorf("invalid GeoSite database file: %s", err)
	}

	if saveFile(geosite, C.Path.GeoSite()) != nil {
		return fmt.Errorf("can't save GeoSite database file: %w", err)
	}

	return nil
}

func downloadForBytes(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func saveFile(bytes []byte, path string) error {
	return ioutil.WriteFile(path, bytes, 0o644)
}
