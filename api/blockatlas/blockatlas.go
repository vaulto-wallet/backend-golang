package blockatlas

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type BlockAtlasAPI struct {
	url    string
	client *http.Client
}

func (a *BlockAtlasAPI) Init(url string) {
	a.url = url
	a.client = &http.Client{}
}

func (a *BlockAtlasAPI) Request(method string, endpoint string, data interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(data)
	req, e := http.NewRequest(method, a.url+endpoint, buf)
	if e != nil {
		return []byte{}, e
	}
	resp, e := a.client.Do(req)
	if e != nil {
		return []byte{}, e
	}
	defer resp.Body.Close()
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return []byte{}, e
	}
	return body, nil
}

func (a *BlockAtlasAPI) GetTXs(asset string, address string) ([]Tx, error) {
	resp, err := a.Request("GET", "/v1/"+asset+"/"+address, nil)
	var response TxResponse
	if err != nil {
		return response.Docs, err
	}
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)
	return response.Docs, nil
}

func (a *BlockAtlasAPI) EstimateGas(asset string, data interface{}) (int64, error) {
	resp, err := a.Request("POST", "/v1/"+asset+"/"+"transaction/estimate", data)
	if err != nil {
		return 0, err
	}
	var response StringResponse

	err = json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)
	if err != nil {
		return 0, err
	}
	ret, err := strconv.ParseInt(response.Result, 10, 64)

	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (a *BlockAtlasAPI) GasPrice(asset string) (int64, error) {
	resp, err := a.Request("GET", "/v1/"+asset+"/"+"gas/price", nil)
	if err != nil {
		return 0, err
	}
	var response StringResponse
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)
	if err != nil {
		return 0, err
	}
	ret, err := strconv.ParseInt(response.Result, 10, 64)

	if err != nil {
		return 0, err
	}
	return ret, nil
}
