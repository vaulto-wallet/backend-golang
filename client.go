package main

import (
	m "./models"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseBool struct {
	Error     string `json:"error"`
	ErrorCode int    `json:"code"`
	ErrorText string `json:"text"`
	Result    bool   `json:"result"`
}

type ResponseString struct {
	Error     string `json:"error"`
	ErrorCode int    `json:"code"`
	ErrorText string `json:"text"`
	Result    string `json:"result"`
}

type ResponseInterface struct {
	Error     string      `json:"error"`
	ErrorCode int         `json:"code"`
	ErrorText string      `json:"text"`
	Result    interface{} `json:"result"`
}

const url = "http://localhost:8000/api"

func apiClear() (success bool) {
	resp, e := http.Get(url + "/clear")
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
	req, e := http.NewRequest("POST", url+"/users/login", buf)
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
	req, e := http.NewRequest("POST", url+"/users/register", buf)
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
	req, e := http.NewRequest("POST", url+"/users/login", buf)
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

func apiCreateAsset(token string, name string, symbol string, index int, decimals int, rounding int) (result bool) {
	buf := new(bytes.Buffer)
	asset := m.Asset{
		Name:      name,
		CoinIndex: index,
		Symbol:    symbol,
		Type:      1,
		Decimals:  decimals,
		Rounding:  rounding,
	}
	json.NewEncoder(buf).Encode(asset)

	req, e := http.NewRequest("POST", url+"/assets", buf)
	if e != nil {
		return false
	}
	req.Header.Add("Authorization", "Bearer "+token)
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
	return true
}

func apiGetAssets(token string) (result interface{}) {
	buf := new(bytes.Buffer)
	req, e := http.NewRequest("GET", url+"/assets", buf)
	if e != nil {
		return false
	}
	req.Header.Add("Authorization", "Bearer "+token)
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

	var response ResponseInterface
	json.NewDecoder(strings.NewReader(string(bodytext))).Decode(&response)

	if len(response.Error) > 0 {
		return false
	}
	return response.Result
}

func apiCreateSeed(token string, name string) (result bool) {
	buf := new(bytes.Buffer)
	buf.Write(([]byte)("{\"name\":\"" + name + "\"}"))
	req, e := http.NewRequest("POST", url+"/seed", buf)
	if e != nil {
		return false
	}
	req.Header.Add("Authorization", "Bearer "+token)
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
	return true
}

func apiGetSeeds(token string) (result interface{}) {
	buf := new(bytes.Buffer)
	req, e := http.NewRequest("GET", url+"/seed", buf)
	if e != nil {
		return false
	}
	req.Header.Add("Authorization", "Bearer "+token)
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

func apiCreateWallet(token string, name string, network string, seed_id int, asset_id int) (result bool) {
	buf := new(bytes.Buffer)
	wallet := m.Wallet{
		Name:        name,
		NetworkType: network,
		AssetID:     asset_id,
		SeedID:      seed_id,
	}
	json.NewEncoder(buf).Encode(wallet)
	req, e := http.NewRequest("POST", url+"/wallets", buf)
	if e != nil {
		return false
	}
	req.Header.Add("Authorization", "Bearer "+token)
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
	return true
}

func apiGetWallets(token string) (result interface{}) {
	buf := new(bytes.Buffer)
	req, e := http.NewRequest("GET", url+"/wallets", buf)
	if e != nil {
		return false
	}
	req.Header.Add("Authorization", "Bearer "+token)
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

	var response ResponseInterface
	json.NewDecoder(strings.NewReader(string(bodytext))).Decode(&response)

	if len(response.Error) > 0 {
		return false
	}
	return response.Result
}

func apiCreateAddress(token string, wallet_id int, comment string) (result bool) {
	buf := new(bytes.Buffer)
	address := m.Address{
		Comment:  comment,
		WalletID: wallet_id,
	}
	json.NewEncoder(buf).Encode(address)
	req, e := http.NewRequest("POST", url+"/address", buf)
	if e != nil {
		return false
	}
	req.Header.Add("Authorization", "Bearer "+token)
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
	return true
}

func main() {
	result := apiClear()
	log.Println("Clear :", result)

	result = apiRegister("user1", "pwd1")
	log.Println("Register : ", result)

	token := apiLogin("user1", "pwd1")
	log.Println("Login token : ", token)

	asset_result := apiCreateAsset(token, "Ethereum", "ETH", 60, 18, 6)
	log.Println("Create asset : ", asset_result)

	assets := apiGetAssets(token)
	log.Println("Assets : ", assets)

	seed_result := apiCreateSeed(token, "Seed1")
	log.Println("Create seed : ", seed_result)

	seeds := apiGetSeeds(token)
	log.Println("Load seed : ", seeds)

	wallet_result := apiCreateWallet(token, "ETH wallet", "main", 1, 1)
	log.Println("Create wallet : ", wallet_result)

	wallets := apiGetWallets(token)
	log.Println("Load wallets : ", wallets)

	address_result := apiCreateAddress(token, 1, "New address")
	log.Println("Create address : ", address_result)

}
