package Recaptcha

import (
	"encoding/json"
	"github.com/dream-huan/Rhine-Cloud-Driver/Class"
	"github.com/dream-huan/Rhine-Cloud-Driver/config"
	logger "github.com/dream-huan/Rhine-Cloud-Driver/middleware/Log"
	"io/ioutil"
	"net/http"
	"net/url"
)

//const privatekey = "6LdQ2vsfAAAAAN1e4mUhc9j4-vZd0k0iUHaNIgKR"

// const privatekey = "6LdBFXIgAAAAAMam2T8Gih9gCOl0GhhBthRuSH3R"
const recaptchaServerName = "https://recaptcha.net/recaptcha/api/siteverify"

func VerifyToken(token string) bool {
	resp, err := http.PostForm(recaptchaServerName,
		url.Values{"secret": {config.GetPrivateKey()}, "response": {token}})
	if err != nil {
		logger.Errorf("httppost错误:%#v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("ioutil的read方法错误:%#v", err)
	}
	var result Class.RecaptchaToken
	err = json.Unmarshal(body, &result)
	if err != nil {
		logger.Errorf("对recaptcha结果处理错误:%#v", err)
	}
	return result.Success
}
