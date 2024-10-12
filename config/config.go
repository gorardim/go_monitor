package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	TgBot struct {
		ChatID string `yaml:"chatID"`
		Token  string `yaml:"token"`
	} `yaml:"tgBot"`

	ServerListApi struct {
		Url string `yaml:"url"`
	} `yaml:"serverListApi"`

	ServiceApi struct {
		Urls []string `yaml:"urls"`
	} `yaml:"serviceApi"`

	H5 struct {
		Urls []string `yaml:"urls"`
	} `yaml:"h5"`

	// navigation:
	// 	list: [
	// 		{
	// 		url: "https://api.js85881.com/apicom/",
	// 		{
	// 			"mobile": "13207863154",
	// 			"password": "123456",
	// 			"from_to": "https://hduegs-4hdshja.com",
	// 			"form_to": "https://api.js85881.com/",
	// 			"token": ""
	// 		}
	// 		}
	// 	]
	Navigation struct {
		List []struct {
			Url      string `yaml:"url"`
			Mobile   string `yaml:"mobile"`
			Password string `yaml:"password"`
			FromTo   string `yaml:"from_to"`
			FormTo   string `yaml:"form_to"`
			Token    string `yaml:"token"`
		} `yaml:"list"`
		GetUrls  []string `yaml:"get_urls"`
		PostUrls []string `yaml:"post_urls"`
		PutUrls  []string `yaml:"put_urls"`
	} `yaml:"navigation"`
}

func (c *Conf) GetConfig() *Conf {
	yamlFile, err := os.ReadFile("config.yml")
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}
