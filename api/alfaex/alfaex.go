package alfaex

import (
	"bytes"
	"crypto/hmac"
	"encoding/hex"
	"encoding/json"
	"golang.org/x/crypto/sha3"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AlfaEXAPI struct {
	public_key []byte
	secret_key []byte
	url        string
	client     *http.Client
	nonce      int
}
type WithdrawOperationReplyArguments struct {
	ExternalTransaction string `json:"externaltransaction"`
	OperationType       string `json:"operation_type"`
	Currency            string `json:"currency"`
	Time                int    `json:"time"`
	Amount              string `json:"amount"`
	Commission          string `json:"commission"`
	Account             string `json:"account"`
	AccountType         string `json:"account_type"`
	AccountWallet       string `json:"account_wallet"`
	DestinationWallet   string `json:"destination_wallet"`
	DestinationTag      string `json:"destination_tag"`
}

type WithdrawOperationReplyPagination struct {
	Limit    int `json:limit`
	Offset   int `json:"offset"`
	SumCount int `json:"sum_count"`
}

type WithdrawOperationReply struct {
	Id         string                           `json:"id"`
	Arguments  WithdrawOperationReplyArguments  `json:"arguments"`
	Signature  string                           `json:"signature"`
	Pagination WithdrawOperationReplyPagination `json:"page"`
}

type WithdrawOperationsReply struct {
	Operation []WithdrawOperationReply `json:"operation"`
}

type WithdrawOperation struct {
	OperationID           string `json:"operation_id,omitempty"`
	ExternalTransactionID string `json:"external_transaction_id,omitempty"`
	Cause                 string `json:"cause,omitempty"`
	State                 int    `json:"state,omitempty"`
}

type WithdrawOperations []WithdrawOperation

func (a *AlfaEXAPI) Init(public_key string, secret_key string, url string) {
	a.url = url
	a.public_key, _ = hex.DecodeString(public_key)
	a.secret_key, _ = hex.DecodeString(secret_key)
	a.client = &http.Client{}

}

func (a *AlfaEXAPI) Request(method string, endpoint string, data interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if data != nil {
		json.NewEncoder(buf).Encode(data)
		buf.Truncate(buf.Len() - 1)
	}
	req, e := http.NewRequest(method, a.url+endpoint, buf)

	if e != nil {
		return []byte{}, e
	}

	now := (int)(time.Now().UTC().Unix())
	if now <= a.nonce {
		a.nonce += 1
	} else {
		a.nonce = now
	}

	signature := hmac.New(sha3.New512, ([]byte)(hex.EncodeToString(a.secret_key)))
	var sign []byte

	signature.Write(([]byte)(endpoint))

	if data != nil {
		signature.Write(([]byte)(":"))
		signature.Write(([]byte)(buf.Bytes()))
	}

	signature.Write(([]byte)(":"))
	signature.Write(([]byte)(strconv.Itoa(a.nonce)))
	sign = signature.Sum(nil)

	req.Header.Add("Alfex-Auth-Signature", hex.EncodeToString(sign))
	req.Header.Add("Alfex-Nonce", strconv.Itoa(a.nonce))
	req.Header.Add("Alfex-API-Key", hex.EncodeToString(a.public_key))

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

func (a *AlfaEXAPI) GetOperations() WithdrawOperationsReply {
	endpoint := "/bookkeeping-api/v1/service_level/transaction/withdraw/wait_for_process/wallet"

	resp, err := a.Request("GET", endpoint, nil)
	log.Println(err, (string)(resp))
	var ret WithdrawOperationsReply
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&ret)

	return ret
}

func (a *AlfaEXAPI) PutOperation(operation WithdrawOperation) {
	endpoint := "/bookkeeping-api/v1/service_level/transaction/withdraw/wait_for_process/wallet"
	resp, err := a.Request("PUT", endpoint, operation)
	log.Println(err, (string)(resp))
}
