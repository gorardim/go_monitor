package navigation

import (
	"bytes"
	"encoding/json"
	"monitor/config"
	tg "monitor/tg_alert"
	"net/http"
)

const (
	loginUrl string = "/user/login"
)

type AuthReq struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	FromTo   string `json:"from_to"`
	FormTo   string `json:"form_to"`
	Token    string `json:"token"`
}

type AuthRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	MsgCode string `json:"msgCode"`
	Data    struct {
		Token  string `json:"token"`
		Mobile string `json:"mobile"`
		Uid    int    `json:"uid"`
		Email  string `json:"email"`
	} `json:"data"`
}

type PostReq struct {
	Token string `json:"token"`
}

type NavigationRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func GetAuthToken(url string, authReq AuthReq) (string, error) {
	body, err := json.Marshal(authReq)
	if err != nil {
		return "", err
	}
	res, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		tg.SendAlert("Failed to login: " + url)
	}
	defer res.Body.Close()
	var authRes AuthRes
	err = json.NewDecoder(res.Body).Decode(&authRes)
	if err != nil {
		tg.SendAlert("Error response from login: " + url + "\n" + err.Error())
	}
	return authRes.Data.Token, nil
}

func GetNavigationHealth() error {
	var Config config.Conf
	Config.GetConfig()
	navigations := Config.Navigation.List
	for _, navigation := range navigations {
		authReq := AuthReq{
			Mobile:   navigation.Mobile,
			Password: navigation.Password,
			FromTo:   navigation.FromTo,
			FormTo:   navigation.FormTo,
		}
		token, err := GetAuthToken(navigation.Url+loginUrl, authReq)
		if err != nil {
			continue
		}
		// Call All GetUrls
		for _, getUrl := range Config.Navigation.GetUrls {
			url := navigation.Url + getUrl + token
			resp, err := http.Get(url)
			if err != nil {
				tg.SendAlert("Failed to get navigation: " + url)
			}
			defer resp.Body.Close()
			var navigationRes NavigationRes
			err = json.NewDecoder(resp.Body).Decode(&navigationRes)
			if err != nil {
				tg.SendAlert("Error response from navigation: " + url + "\n" + err.Error())
			}
			if navigationRes.Status != "1" {
				tg.SendAlert("Navigation is down: " + url + "\n" + "Message: " + navigationRes.Message)
			}
			tg.SendAlert("Navigation is up: " + url + "\n" + "Message: " + navigationRes.Message)
		}

		// Call All PostUrls
		for _, postUrl := range Config.Navigation.PostUrls {
			url := navigation.Url + postUrl
			v := PostReq{
				Token: token,
			}
			body, err := json.Marshal(v)
			if err != nil {
				tg.SendAlert("Failed to post navigation: " + url)
			}
			resp, err := http.Post(url, "application/json", bytes.NewReader(body))
			if err != nil {
				tg.SendAlert("Failed to post navigation: " + url)
			}
			defer resp.Body.Close()
			var navigationRes NavigationRes
			err = json.NewDecoder(resp.Body).Decode(&navigationRes)
			if err != nil {
				tg.SendAlert("Error response from navigation: " + url + "\n" + err.Error())
			}
			if navigationRes.Status != "1" {
				tg.SendAlert("Navigation is down: " + url + "\n" + "Message: " + navigationRes.Message)
			}
			tg.SendAlert("Navigation is up: " + url + "\n" + "Message: " + navigationRes.Message)
		}

		// Call All PutUrls
		for _, putUrl := range Config.Navigation.PutUrls {
			url := navigation.Url + putUrl
			v := PostReq{
				Token: token,
			}
			body, err := json.Marshal(v)
			if err != nil {
				tg.SendAlert("Failed to put navigation: " + url)
			}
			req, err := http.NewRequest("PUT", url, bytes.NewReader(body))
			if err != nil {
				tg.SendAlert("Failed to put navigation: " + url)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				tg.SendAlert("Failed to put navigation: " + url)
			}
			defer resp.Body.Close()
			var navigationRes NavigationRes
			err = json.NewDecoder(resp.Body).Decode(&navigationRes)
			if err != nil {
				tg.SendAlert("Error response from navigation: " + url + "\n" + err.Error())
			}
			if navigationRes.Status != "1" {
				tg.SendAlert("Navigation is down: " + url + "\n" + "Message: " + navigationRes.Message)
			}
			tg.SendAlert("Navigation is up: " + url + "\n" + "Message: " + navigationRes.Message)
		}
	}
	return nil
}
