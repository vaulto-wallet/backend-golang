package vaulto

type LoginRequest struct {
	Username string `json:"username, omitempty"`
	Password string `json:"password, omitempty"`
}

/*
type AssetRequest struct {
	Name      string `json:"name,omitempty"`
	CoinIndex int    `json:"coinindex,omitempty"`
	Symbol    string `json:"symbol,omitempty"`
	Type      int    `json:"type,omitempty"`
	Decimals  int    `json:"decimals,omitempty"`
	Rounding  int    `json:"rounding,omitempty"`
}

type AssetsResponse []AssetRequest
*/
type SeedRequest struct {
	Name     string `json:"name,omitempty"`
	Mnemonic string `json:"mnemonic,omitempty"`
}

/*
type WalletRequest struct {
	Id      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	SeedId  uint   `json:"seed_id,omitempty"`
	AssetId uint   `json:"asset_id,omitempty"`
	Symbol  string `json:"asset_symbol,omitempty"`
}


type WalletsResponse []WalletRequest
*/
/*type OrderRequest struct {
	Id        int     `json:"id,omitempty"`
	AssetId   uint    `json:"asset_id,omitempty"`
	Symbol    string  `json:"symbol"`
	WalletId  uint    `json:"wallet_id,omitempty"`
	AddressTo string  `json:"address_to,omitempty"`
	Amount    float64 `json:"amount,omitempty"`
	Comment   string  `json:"comment,omitempty"`
	Status    int     `json:"status,omitempty"`
	Tx        string  `json:"tx,omitempty"`
	TxHash    string  `json:"tx_hash,omitempty"`
}

type OrdersResponse []OrderRequest
*/
/*
type AddressRequest struct {
	Name     string `json:"name,omitempty"`
	WalletID int    `json:"wallet_id,omitempty"`
}
*/
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

type ResponseError struct {
	Error     string `json:"error"`
	ErrorCode int    `json:"code"`
	ErrorText string `json:"text"`
}

/*
func (ar *AssetsResponse) Find(symbol string) (request *AssetRequest) {
	for i, s := range ([]AssetRequest)(*ar) {
		if s.Symbol == symbol {
			return &([]AssetRequest)(*ar)[i]
		}
	}
	return nil
}

func (a *AssetRequest) ToBigInt(value float64) (ret *big.Int) {
	e := new(big.Int)
	r := new(big.Float)
	e.Exp(big.NewInt(10), big.NewInt( int64(a.Decimals) ), nil)
	r.SetFloat64(value).Mul( r, new(big.Float).SetInt(e) )
	ret, _ = r.Int(nil)
	return
}

func (a *AssetRequest) ToFloat(value *big.Int) (ret float64) {
	e := new(big.Int)
	r := new(big.Float)
	e.Exp(big.NewInt(10), big.NewInt( int64(a.Decimals) ), nil)
	r.SetInt(value)
	ret, _ = new(big.Float).Quo(r, new(big.Float).SetInt(e) ).Float64()
	return
}

*/
