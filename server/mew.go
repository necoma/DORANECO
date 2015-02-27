package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bytes"
)

type MewPostData struct {
	UserName  string `json:"user_name"`
	ApiKey    string `json:"api_key"`
	Text      string `json:"text"`
}

type MewUserData struct {
	UserName string `json:"userName"`
	ApiKey   string `json:"apiKey"`
}

func Mew(text string, userData *MewUserData) error {
	fmt.Println("mew: ", text)

	url := "http://necomatter.necoma-project.jp/post.json"
	mewData := MewPostData{
		UserName: userData.UserName,
		ApiKey: userData.ApiKey,
		Text: text}
	postData, err := json.Marshal(mewData)
	if err != nil {
		fmt.Println("mew error: ", err)
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
	if err != nil {
		fmt.Println("mew error: ", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("mew error: ", err)
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("mew error: ", err)
		return err
	}

	return nil
}


