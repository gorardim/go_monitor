package serverhealth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"monitor/config"
	tg "monitor/tg_alert"
	"net/http"
	"os/exec"
	"time"
)

type ServerListData struct {
	Platform string `json:"platform"`
	Note     string `json:"note"`
	ServerIP string `json:"server_ip"`
	Expire   string `json:"expire"`
	ServerID string `json:"server_id"`
}

type ServerListRes struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data []ServerListData `json:"data"`
}

// {
// 	"code": 1,
// 	"msg": "操作成功",
// 	"data": [
// 		{
// 			"platform": "diyvm",
// 			"note": "华耀后端 4 CPUs, 8192M Ram, 70G Disk, 156.254.127.18, CentOS-7.8.2003-x64",
// 			"server_ip": "156.254.127.18",
// 			"expire": "2024-10-28 22:47:43",
// 			"server_id": "122660"
// 		}
// 	]
// }

func GetServerList() (*ServerListRes, error) {
	var Config config.Conf
	Config.GetConfig()
	apiUrl := Config.ServerListApi.Url
	resp, err := http.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var serverListRes ServerListRes
	err = json.Unmarshal(body, &serverListRes)
	if err != nil {
		return nil, err
	}
	return &serverListRes, nil
}

func CheckServerHealth() {
	serverListRes, err := GetServerList()
	if err != nil {
		tg.SendAlert("GetServerList error: " + err.Error())
		return
	}
	for _, server := range serverListRes.Data {
		expireTime, err := time.Parse("2006-01-02 15:04:05", server.Expire)
		if err != nil {
			tg.SendAlert("Parse expire time error: " + err.Error())
			continue
		}
		if time.Now().After(expireTime) {
			tg.SendAlert(
				"Server is expired: " + server.ServerIP + "\n" +
					"Platform: " + server.Platform + "\n" +
					"Note: " + server.Note + "\n" +
					"Expire: " + server.Expire + "\n" +
					"ServerID: " + server.ServerID,
			)
		}

		if time.Now().AddDate(0, 0, 7).After(expireTime) {
			tg.SendAlert(
				fmt.Sprintf("Server is going to be expired in %d days: %s\n", int(time.Until(expireTime).Hours()/24), server.ServerIP) +
					"Platform: " + server.Platform + "\n" +
					"Note: " + server.Note + "\n" +
					"Expire: " + server.Expire + "\n" +
					"ServerID: " + server.ServerID,
			)
		}

		out, err := exec.Command("ping", "-c", "1", server.ServerIP).Output()
		if err != nil {
			// tg.SendAlert("Ping command error: " + err.Error())
			log.Printf("Ping command error: %v", err)
			continue
		}
		// log.Printf("Ping command output for server %s: %s", server.ServerIP, string(out))
		if string(out) == "" {
			tg.SendAlert(
				"Server is down: " + server.ServerIP + "\n" +
					"Platform: " + server.Platform + "\n" +
					"Note: " + server.Note + "\n" +
					"Expire: " + server.Expire + "\n" +
					"ServerID: " + server.ServerID,
			)
		}
	}
}
