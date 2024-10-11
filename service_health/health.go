package servicehealth

import (
	"encoding/json"
	"monitor/config"
	tg "monitor/tg_alert"
	"net/http"
)

type ServiceHeltRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Apidomain  []interface{} `json:"apidomain"`
		Rootdomain []interface{} `json:"rootdomain"`
	} `json:"data"`
}

func GetServiceHealth() error {
	var Config config.Conf
	Config.GetConfig()
	urls := Config.ServiceApi.Urls
	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == 404 || resp.StatusCode == 500 {
			// Send alert (implementation depends on your alerting mechanism)
			tg.SendAlert("Service is down: " + url)
		}
		var serviceHeltRes ServiceHeltRes
		err = json.NewDecoder(resp.Body).Decode(&serviceHeltRes)
		if err != nil {
			return err
		}
		if serviceHeltRes.Status != "1" {
			// Send alert (implementation depends on your alerting mechanism)
			tg.SendAlert("Service is down: " + url + "\n" + serviceHeltRes.Message)
		}
	}
	return nil
}
