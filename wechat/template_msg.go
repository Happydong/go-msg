package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-msg/redis"
	"io/ioutil"
	"net/http"
)

// 发送模板消息
var (
	send_template_url        = "https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=%s"
	send_wxopen_template_url = "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s"
	get_access_token_url     = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
)

// 发送微信小程序消息
type TemplateMiniMsg struct {
	Touser           string        `json:"touser"`      // 接收者的OpenID
	TemplateID       string        `json:"template_id"` // 模板消息ID
	Page             string        `json:"page"`        // 点击后跳转链接
	Data             interface{} `json:"data"`
	MiniprogramState string        `json:"miniprogram_state"` // developer为开发版；trial为体验版；formal为正式版；默认为正式版
}

// 发送微信公众号消息
type TemplateWxOpenMsg struct {
	Touser      string       `json:"touser"`      // 接收者的OpenID
	TemplateID  string       `json:"template_id"` // 模板消息ID
	Url         string       `json:"url"`         // 模板跳转链接
	Miniprogram *Miniprogram `json:"miniprogram"` // 跳小程序所需数据，不需跳小程序可不用传该数据
	Data        interface{}  `json:"data"`        // 模拟数据
	Color       string       `json:"color"`       //  模板内容字体颜色，不填默认为黑色
}

type Miniprogram struct {
	AppID    string `json:"appid"`    // 所需跳转到的小程序appid（该小程序appid必须与发模板消息的公众号是绑定关联关系，暂不支持小游戏）
	Pagepath string `json:"pagepath"` // 所需跳转到小程序的具体页面路径，支持带参数,（示例index?foo=bar），要求该小程序已发布，暂不支持小游戏
}


type KeyWordData struct {
	Value string `json:"value"`
}

type KeyWordOpenData struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
}

type SendTemplateResponse struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	MsgID   string `json:"msgid"`
}

type GetAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// SendTemplate 发送模板消息
func SendMiniTemplate(msg *TemplateMiniMsg, miniAppID string, miniSecret string) (*SendTemplateResponse, error) {
	// 1. 获取token
	accessToken, err := getAccessToken(miniAppID, miniSecret, true)
	if err != nil {
		return nil, err
	}
	// 2. 参数格式化
	url := fmt.Sprintf(send_template_url, accessToken.AccessToken)
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	// 3. 发送请求
	return  sendUrl(url, data)
}

// SendTemplate 发送模板消息
func SendWxOpenTemplate(msg *TemplateWxOpenMsg, openAppID string, openSecret string) (*SendTemplateResponse, error) {
	// 1. 获取token
	accessToken, err := getAccessToken(openAppID, openSecret, false)
	if err != nil {
		return nil, err
	}
	// 2. 格式化请求url 和 请求参数
	url := fmt.Sprintf(send_wxopen_template_url, accessToken.AccessToken)
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	// 3. 发送请求
	return sendUrl(url, data)
}

// 发送请求
func sendUrl(url string, data []byte) (*SendTemplateResponse, error) {
	client := http.Client{}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("网络错误，发送模板消息失败", "err", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var templateResponse SendTemplateResponse
	err = json.Unmarshal(body, &templateResponse)
	if err != nil {
		fmt.Println("解析responseBody错误", "err", err.Error())
		return nil, err
	}
	return &templateResponse, nil
}

// 获取token
func getAccessToken(appID string, secret string, isMini bool) (*GetAccessTokenResponse, error) {
	var accessTokenResponse GetAccessTokenResponse
	// 先从redis中拿
	accessToken, err := getAccessTokenFromRedis(appID, isMini)
	if accessToken != "" && err == nil {
		accessTokenResponse = GetAccessTokenResponse{AccessToken: accessToken}
		fmt.Println("从redis中获取到access_token", "access_token", accessToken)
		return &accessTokenResponse, nil
	}

	url := fmt.Sprintf(get_access_token_url, appID, secret)
	client := http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("获取access_toke网络异常", "err", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &accessTokenResponse)
	if err != nil {
		fmt.Println("解析AccessToken失败", "err", err.Error())
		return nil, err
	}
	// 存到redis中
	if err := setAccessTokenToRedis(appID, accessTokenResponse.AccessToken, isMini); err != nil {
		fmt.Println("将access_token存储到redis中失败", "err", err.Error())
	}
	return &accessTokenResponse, nil
}

// 从redis中取access_token
func getAccessTokenFromRedis(AppId string, isMini bool) (string, error) {
	key := fmt.Sprintf("wechatMinToken:%s", AppId)
	if !isMini {
		key = fmt.Sprintf("wechatOpenToken:%s", AppId)
	}

	accessToken, err := redis.Get(key)
	actStr := ""
	if accessToken != nil {
		actStr = string(accessToken.([]uint8))
	}
	return actStr, err
}

// 将access_token存储到redis中
func setAccessTokenToRedis(AppId, token string, isMini bool) error {
	key := fmt.Sprintf("wechatMinToken:%s", AppId)
	if !isMini {
		key = fmt.Sprintf("wechatOpenToken:%s", AppId)
	}
	err := redis.Set(key, token, 7000)
	return err
}
