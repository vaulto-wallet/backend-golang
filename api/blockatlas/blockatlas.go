package blockatlas

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
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
