package config

import (
	"fmt"
	"github.com/Dreamacro/clash/component/geodata"
	"io"
	"net/http"
	"os"

	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
)

func downloadGeoIP(path string) (err error) {
	resp, err := http.Get(C.GeoIpUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)

	return err
}

func initGeoIP() error {
	if _, err := os.Stat(C.Path.GeoIP()); os.IsNotExist(err) {
		log.Infoln("Can't find GeoIP.dat, start download")
		if err := downloadGeoIP(C.Path.GeoIP()); err != nil {
			return fmt.Errorf("can't download GeoIP.dat: %s", err.Error())
		}
		log.Infoln("Download GeoIP.dat finish")
	}

	if err := geodata.Verify(C.GeoipName); err != nil {
		log.Warnln("GeoIP.dat invalid, remove and download: %s", err)
		if err := os.Remove(C.Path.GeoIP()); err != nil {
			return fmt.Errorf("can't remove invalid GeoIP.dat: %s", err.Error())
		}
		if err := downloadGeoIP(C.Path.GeoIP()); err != nil {
			return fmt.Errorf("can't download GeoIP.dat: %s", err.Error())
		}
	}
	return nil
}

// Init prepare necessary files
func Init(dir string) error {
	// initial homedir
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0o777); err != nil {
			return fmt.Errorf("can't create config directory %s: %s", dir, err.Error())
		}
	}

	// initial config.yaml
	if _, err := os.Stat(C.Path.Config()); os.IsNotExist(err) {
		log.Infoln("Can't find config, create a initial config file")
		f, err := os.OpenFile(C.Path.Config(), os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return fmt.Errorf("can't create file %s: %s", C.Path.Config(), err.Error())
		}
		f.Write([]byte(`mixed-port: 7890`))
		f.Close()
	}
	buf, _ := os.ReadFile(C.Path.Config())
	rawCfg, err := UnmarshalRawConfig(buf)
	if err != nil {
		log.Errorln(err.Error())
		fmt.Printf("configuration file %s test failed\n", C.Path.Config())
		os.Exit(1)
	}
	C.GeoIpUrl = rawCfg.GeoXUrl.GeoIp
	C.GeoSiteUrl = rawCfg.GeoXUrl.GeoSite
	// initial GeoIP
	if err := initGeoIP(); err != nil {
		return fmt.Errorf("can't initial GeoIP: %w", err)
	}

	return nil
}
