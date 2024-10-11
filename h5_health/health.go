package servicehealth

import (
	"monitor/config"
	tg "monitor/tg_alert"
	"os/exec"
)

func GetH5Health() error {
	var Config config.Conf
	Config.GetConfig()
	urls := Config.H5.Urls
	for _, url := range urls {
		_, err := exec.Command("curl", "--location", url).Output()
		if err != nil {
			tg.SendAlert("H5 is down: " + url)
			continue
		}

	}
	return nil
}
