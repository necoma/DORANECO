package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bytes"
)

// MewPostData $B$O(B mew $B$9$k;~$K;H$&%G!<%?$G$9!#(B
type MewPostData struct {
	UserName  string `json:"user_name"`
	APIKey    string `json:"api_key"`
	Text      string `json:"text"`
}

// MewUserData $B$O(B mew $B$9$k;~$K;H$o$l$k%f!<%6(BID$B$H(BAPI key$B$rJ];}$7$^$9(B
type MewUserData struct {
	UserName string `json:"userName"`
	APIKey   string `json:"apiKey"`
}

// MewResultData $B$O(B mew $B$N7k2L$H$7$FF@$i$l$k9=B$BN$G$9(B
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

// Mew $B$O(B NECOMAtter $B$K=q$-9~$_$^$9!#(B
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


