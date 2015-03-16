package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bytes"
)

// MewPostData は mew する時に使うデータです。
type MewPostData struct {
	UserName  string `json:"user_name"`
	APIKey    string `json:"api_key"`
	Text      string `json:"text"`
}

// MewUserData は mew する時に使われるユーザIDとAPI keyを保持します
type MewUserData struct {
	UserName string `json:"userName"`
	APIKey   string `json:"apiKey"`
}

// MewResultData は mew の結果として得られる構造体です
type MewResultData struct {
	UserName         string     `json:"user_name"`
	ID               string     `json:"id"`
	Text             string     `json:"text"`
	Time             string     `json:"time"`
	UnixTime         int        `json:"unix_time"`
	IconURL          string     `json:"icon_url"`
	OwnStard         bool       `json:"own_stard"`
	OwnRetweeted     bool       `json:"own_retweeted"`
	IsRetweet        bool       `json:"is_retweet"`
	RetweetUserName  string     `json:"retweet_user_name"`
	RetweetTime      string     `json:"retweet_time"`
	RetweetUnixTime  string     `json:"retweet_unix_time"`
	ListName         string     `json:"list_name"`
	ListOwnerName    string     `json:"list_owner_name"`
	Result           string     `json:"result"`
}

// Mew は NECOMAtter に書き込みます。
func Mew(text string, userData *MewUserData) (*MewResultData, error) {
	fmt.Println("mew: ", text)

	url := "http://necomatter.necoma-project.jp/post.json"
	mewData := MewPostData{
		UserName: userData.UserName,
		APIKey: userData.APIKey,
		Text: text}
	postData, err := json.Marshal(mewData)
	if err != nil {
		fmt.Println("mew error: ", err)
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
	if err != nil {
		fmt.Println("mew error: ", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("mew error: ", err)
		return nil, err
	}
	defer resp.Body.Close()
	mewResultJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("mew error: ", err)
		return nil, err
	}
	mewResult := new(MewResultData)
	err = json.Unmarshal(mewResultJSON, mewResult)
	if err != nil {
		fmt.Println("mew result unmarshal error: ", err)
		return nil, err
	}

	return mewResult, nil
}


