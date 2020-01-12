package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Login struct {
	Username    string `json:"usernam"`
	Password string `json:"password"`
}

type ResponseBool struct {
	Error string `json:"error"`
	ErrorCode int `json:"code"`
	ErrorText string `json:"text"`
	Result bool `json:"result"`
}

type ResponseString struct {
	Error string `json:"error"`
	ErrorCode int `json:"code"`
	ErrorText string `json:"text"`
	Result string `json:"result"`
}


const url = "http://localhost:8000/api"

func apiClear() (success bool) {
	resp, e := http.Get( url + "/clear")
	if e != nil {
		return false
	}
	defer resp.Body.Close()
	bodytext, e := ioutil.ReadAll(resp.Body)
	log.Print(string(bodytext))
	return true
}


func apiLogin(username string, password string) (token string) {
	buf := new(bytes.Buffer)
	body := &Login{username, password}
	json.NewEncoder(buf).Encode(body)
	req, e := http.NewRequest("POST", url + "/users/login",  buf)
	if e != nil {
		return ""
	}
	client := &http.Client{}
	resp, e := client.Do(req)
	if e != nil {
		return ""
	}
	defer resp.Body.Close()

	bodytext, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return ""
	}

	log.Print(string(bodytext))

	var response ResponseString
	json.NewDecoder(strings.NewReader(string(bodytext))).Decode(&response)

	if len(response.Error) > 0 {
		return ""
	}
	return response.Result
}

func apiRegister(username string, password string) (result bool) {
	buf := new(bytes.Buffer)
	body := &Login{username, password}
	json.NewEncoder(buf).Encode(body)
	req, e := http.NewRequest("POST", url + "/users/register",  buf)
	if e != nil {
		return false
	}
	client := &http.Client{}
	resp, e := client.Do(req)
	if e != nil {
		return false
	}
	defer resp.Body.Close()

	bodytext, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return false
	}

	log.Print(string(bodytext))

	var response ResponseBool
	json.NewDecoder(strings.NewReader(string(bodytext))).Decode(&response)

	if len(response.Error) > 0 {
		return false
	}

	return response.Result
}

func apiCreateAccount(username string, password string) (token string) {
	buf := new(bytes.Buffer)
	body := &Login{username, password}
	json.NewEncoder(buf).Encode(body)
	req, e := http.NewRequest("POST", url + "/users/login",  buf)
	if e != nil {
		return ""
	}
	client := &http.Client{}
	resp, e := client.Do(req)
	if e != nil {
		return ""
	}
	defer resp.Body.Close()

	bodytext, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return ""
	}

	log.Print(string(bodytext))

	var response ResponseString
	json.NewDecoder(strings.NewReader(string(bodytext))).Decode(&response)

	if len(response.Error) > 0 {
		return ""
	}
	return response.Result
}

func main() {
	result := apiClear()
	log.Println("Clear :", result)

	result = apiRegister("user1", "pwd1")
	log.Println("Register : ", result)

	token := apiLogin("user1", "pwd1")
	log.Println("Login token : ",  token)
}
