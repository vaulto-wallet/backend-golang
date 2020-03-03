package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type VaultoError struct {
	error string
}

func (e *VaultoError) Error() string {
	return e.error
}

type VaultoAPI struct {
	token  string
	url    string
	client *http.Client
	Error  error
}

func (a *VaultoAPI) Init(url string) {
	a.url = url
	a.client = &http.Client{}
}

func (a *VaultoAPI) Request(method string, endpoint string, data interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(data)
	req, e := http.NewRequest(method, a.url+endpoint, buf)
	if e != nil {
		return []byte{}, e
	}
	if len(a.token) > 0 {
		req.Header.Add("Authorization", "Bearer "+a.token)
	}

	resp, e := a.client.Do(req)
	if e != nil {
		return []byte{}, e
	}
	defer resp.Body.Close()
	bodytext, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return []byte{}, e
	}
	return bodytext, nil
}

func (a *VaultoAPI) ApiClear() (success bool) {
	resp, e := http.Get(a.url + "/clear")
	if e != nil {
		return false
	}
	defer resp.Body.Close()
	bodytext, e := ioutil.ReadAll(resp.Body)
	log.Print(string(bodytext))
	return true
}

func (a *VaultoAPI) ApiLogin(username string, password string) (bool, error) {

	resp, err := a.Request("POST", "/users/login", LoginRequest{username, password})
	if err != nil {
		return false, err
	}

	log.Print(string(resp))

	var response ResponseString
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return false, errors.New(response.ErrorText)
	}
	a.token = response.Result
	return true, nil
}

func (a *VaultoAPI) ApiRegister(username string, password string) (bool, error) {
	resp, err := a.Request("POST", "/user/register", LoginRequest{username, password})
	if err != nil {
		return false, err
	}
	var response ResponseBool
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return false, errors.New(response.ErrorText)
	}

	return response.Result, nil
}

func (a *VaultoAPI) ApiCreateAsset(name string, symbol string, index int, decimals int, rounding int) (bool, error) {
	resp, err := a.Request("POST", "/assets", AssetRequest{
		Name:      name,
		CoinIndex: index,
		Symbol:    symbol,
		Type:      1,
		Decimals:  decimals,
		Rounding:  rounding,
	})
	if err != nil {
		return false, err
	}

	log.Print(string(resp))

	var response ResponseBool
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return false, errors.New(response.ErrorText)
	}
	return true, nil
}

func (a *VaultoAPI) ApiGetAssets() (interface{}, error) {
	resp, err := a.Request("GET", "/assets", nil)
	if err != nil {
		return nil, err
	}

	log.Print(string(resp))

	var response ResponseInterface
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return nil, errors.New(response.ErrorText)
	}
	return response.Result, nil
}

func (a *VaultoAPI) ApiCreateSeed(name string, seed string) (bool, error) {
	resp, err := a.Request("POST", "/seed", SeedRequest{name, ""})
	if err != nil {
		return false, err
	}
	log.Print(string(resp))

	var response ResponseBool
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return false, errors.New(response.ErrorText)
	}
	return true, nil
}

func (a *VaultoAPI) ApiGetSeeds(token string) (interface{}, error) {
	resp, err := a.Request("GET", "/seeds", nil)
	if err != nil {
		return false, err
	}

	var response ResponseInterface
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return false, errors.New(response.ErrorText)
	}
	return response.Result, nil
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
