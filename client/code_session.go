package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hust-tianbo/go_lib/log"
)

type CodeSession struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

const (
	AppID     = ""
	AppSecret = ""
)

func CodeToSession(code string) (*CodeSession, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s"+
			"&grant_type=authorization_code", AppID, AppSecret, code)
	resp, err := client.Get(url)
	if err != nil {
		log.Errorf("[IsCodeValid]get failed:%+v,%+v", url, err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[IsCodeValid]get failed:%+v,%+v", url, err)
		return nil, err
	}

	codeSession := &CodeSession{}
	json.Unmarshal(body, codeSession)
	if codeSession.ErrCode != 0 {
		log.Errorf("[IsCodeValid]codeSession code  err:%+v,%+v", url, codeSession)
		return nil, fmt.Errorf("code to session err:%+v", codeSession.ErrCode)
	}

	return codeSession, nil
}
