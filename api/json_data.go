package api

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AssetRequest struct {
	Name      string `json:"name,omitempty"`
	CoinIndex int    `json:"coinindex,omitempty"`
	Symbol    string `json:"symbol,omitempty"`
	Type      int    `json:"type,omitempty"`
	Decimals  int    `json:"decimals,omitempty"`
	Rounding  int    `json:"rounding,omitempty"`
}

type SeedRequest struct {
	Name     string `json:"name,omitempty"`
	Mnemonic string `json:"mnemonic,omitempty"`
}

type WalletRequest struct {
	Name        string
	NetworkType string
	SeedID      int
	AssetID     int
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
