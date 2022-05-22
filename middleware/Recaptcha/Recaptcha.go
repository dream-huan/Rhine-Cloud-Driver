package Recaptcha

import (
	"encoding/json"
	"golandproject/Class"
	"io/ioutil"
	"net/http"
	"net/url"
)

const privatekey = "6LdQ2vsfAAAAAN1e4mUhc9j4-vZd0k0iUHaNIgKR"
const recaptchaServerName = "https://recaptcha.net/recaptcha/api/siteverify"

func VerifyToken(token string) bool {
	resp, _ := http.PostForm(recaptchaServerName,
		url.Values{"secret": {privatekey}, "response": {token}})
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result Class.RecaptchaToken
	_ = json.Unmarshal(body, &result)
	return result.Success
}
