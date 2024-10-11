package navigation

import (
	"bytes"
	"encoding/json"
	"monitor/config"
	tg "monitor/tg_alert"
	"net/http"
)

const (
	loginUrl      string = "/user/login"
	navigationUrl string = "/system/config_navigation?token="
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

type NavigationRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Home     []interface{} `json:"home"`
		User     []interface{} `json:"user"`
		UserList []interface{} `json:"user_list"`
	} `json:"data"`
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
		resp, err := http.Get(navigation.Url + navigationUrl + token)
		if err != nil {
			tg.SendAlert("Failed to get navigation: " + navigation.Url)
		}
		defer resp.Body.Close()
		var navigationRes NavigationRes
		err = json.NewDecoder(resp.Body).Decode(&navigationRes)
		if err != nil {
			tg.SendAlert("Error response from navigation: " + navigation.Url + "\n" + err.Error())
		}
		if navigationRes.Status != "1" {
			tg.SendAlert("Navigation is down: " + navigation.Url + "\n" + "Message: " + navigationRes.Message)
		}

		tg.SendAlert("Navigation is up: " + navigation.Url + "\n" + "Message: " + navigationRes.Message)
	}
	return nil
}
