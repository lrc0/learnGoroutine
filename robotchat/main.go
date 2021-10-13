package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/eatMoreApple/openwechat"
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/logger.v1"
)

const (
	// 图灵机器人
	//apiKey  = "a866329e0339488a93d4080414e4f521"
	//baseURL = "http://www.tuling123.com/openapi/api"

	// 思知机器人
	baseURL = "https://api.ownthink.com/bot"
	appId   = "0b2b9c8b8842ee5b7b4fa77170ea9f0a"
	userId = "40E6DJsJ"
)

func main() {
	Login()
}

// 图灵 Answer .
//type Answer struct {
//	Code    int    `json:"code"`
//	Message string `json:"text"`
//}

// 思知 Answer
type Answer struct {
	Message string `json:"message"`
	Data    struct {
		Type int `json:"type"`
		Info struct {
			Text string `json:"text"`
		} `json:"info"`
	} `json:"data"`
}

func keepTalk(in string) {
	client := new(http.Client)

	answer, err := chat(in, client)
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Printf("answer: %s", string(answer))

	var an Answer
	err = json.Unmarshal(answer, &an)
	if err != nil {
		log.Error(err)
		return
	}

	l := len(an.Data.Info.Text)

	time.Sleep(time.Millisecond * time.Duration(l))

	sendMsg(an.Data.Info.Text)
}

func chat(speak string, client *http.Client) ([]byte, error) {
	// 图灵
	//robotURL := fmt.Sprintf("%s?key=%s&info=%s", baseURL, apiKey, speak)
	//req, err := http.NewRequest("GET", robotURL, nil)

	// 思知
	//var body Robot
	body := make(map[string]string)
	body["userid"] = userId
	body["appid"] = appId
	body["spoken"] = speak

	r, _ := json.Marshal(body)
	b := bytes.NewBuffer(r)

	req, err := http.NewRequest("POST", baseURL, b)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return respBody, nil
}

var bot = new(openwechat.Bot)

func Login() {
	//bot := openwechat.DefaultBot()
	bot = openwechat.DefaultBot(openwechat.Desktop) // 桌面模式，上面登录不上的可以尝试切换这种模式

	// 注册消息处理函数
	bot.MessageHandler = func(msg *openwechat.Message) {
		if msg.IsText() && msg.Content == "ping" {
			msg.ReplyText("pong")
		}
	}
	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 登陆
	if err := bot.Login(); err != nil {
		fmt.Println(err)
		return
	}

	receiveMsg()

	// 阻塞主goroutine, 知道发生异常或者用户主动退出
	bot.Block()
}

func sendMsg(in string) {
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}
	friends, err := self.Friends()
	if err != nil {
		fmt.Println(err)
		return
	}

	f := friends.Search(1, func(friend *openwechat.Friend) bool { return friend.User.RemarkName == "长虹同事陈强" })
	f.SendText(in)
}

func receiveMsg() {
	bot.MessageHandler = func(msg *openwechat.Message) {
		sender, _ := msg.Sender()
		fmt.Printf("%+v", sender)
		if sender.RemarkName == "长虹同事陈强" {
			keepTalk(msg.Content)
		}
	}
}
