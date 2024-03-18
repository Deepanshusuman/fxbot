package capital

type Transactions struct {
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	DateUtc         string `json:"dateUtc"`
	InstrumentName  string `json:"instrumentName"`
	TransactionType string `json:"transactionType"`
	Reference       string `json:"reference"`
	Size            string `json:"size"`
}

type Positions struct {
	Positions []Position `json:"positions"`
}
type Position struct {
	Position PositionData `json:"position"`
	Market   MarketData   `json:"market"`
}

type Activities struct {
	Activities []Activity `json:"activities"`
}
type Activity struct {
	Date    string `json:"date"`
	DateUTC string `json:"dateUTC"`
	Epic    string `json:"epic"`
	DealId  string `json:"dealId"`
	Source  string `json:"source"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	Details Detail `json:"details"`
}
type Detail struct {
	DealReference  string  `json:"dealReference"`
	MarketName     string  `json:"marketName"`
	Currency       string  `json:"currency"`
	Size           float64 `json:"size"`
	Direction      string  `json:"direction"`
	Level          float64 `json:"level"`
	GuaranteedStop bool    `json:"guaranteedStop"`
	StopLevel      float64 `json:"stopLevel"`
}

type PositionData struct {
	ContractSize   int     `json:"contractSize"`
	CreatedDate    string  `json:"createdDate"`
	CreatedDateUTC string  `json:"createdDateUTC"`
	DealId         string  `json:"dealId"`
	DealReference  string  `json:"dealReference"`
	WorkingOrderId string  `json:"workingOrderId"`
	Size           float64 `json:"size"`
	Direction      string  `json:"direction"`
	Level          float64 `json:"level"`
	StopLevel      float64 `json:"stopLevel"`
	ProfitLevel    float64 `json:"profitLevel"`
	Currency       string  `json:"currency"`
}
type MarketData struct {
	InstrumentName   string  `json:"instrumentName"`
	Expiry           string  `json:"expiry"`
	MarketStatus     string  `json:"marketStatus"`
	Epic             string  `json:"epic"`
	InstrumentType   string  `json:"instrumentType"`
	LotSize          int     `json:"lotSize"`
	High             float64 `json:"high"`
	Low              float64 `json:"low"`
	PercentageChange float64 `json:"percentageChange"`
	NetChange        float64 `json:"netChange"`
	Bid              float64 `json:"bid"`
	Offer            float64 `json:"offer"`
	UpdateTime       string  `json:"updateTime"`
	UpdateTimeUTC    string  `json:"updateTimeUTC"`
}
type Candle struct {
	SnapshotTime    string `json:"snapshotTime"`
	SnapshotTimeUTC string `json:"snapshotTimeUTC"`
	OpenPrice       struct {
		Bid float64 `json:"bid"`
		Ask float64 `json:"ask"`
	} `json:"openPrice"`
	ClosePrice struct {
		Bid float64 `json:"bid"`
		Ask float64 `json:"ask"`
	} `json:"closePrice"`
	HighPrice struct {
		Bid float64 `json:"bid"`
		Ask float64 `json:"ask"`
	} `json:"highPrice"`
	LowPrice struct {
		Bid float64 `json:"bid"`
		Ask float64 `json:"ask"`
	} `json:"lowPrice"`
	LastTradedVolume int64 `json:"lastTradedVolume"`
}
type Candles struct {
	Candles        []Candle `json:"prices"`
	InstrumentType string   `json:"instrumentType"`
	ErrorCode      string   `json:"errorCode"`
}

// // result or error //{
// "errorCode": "validation.stoploss.maxvalue: 24.302"
// }
type PositionResponse struct {
	DealReference string `json:"dealReference"`
	ErrorCode     string `json:"errorCode"`
}

type Session struct {
	AccountType string `json:"accountType"`
	AccountInfo struct {
		Balance    float64 `json:"balance"`
		Deposit    float64 `json:"deposit"`
		ProfitLoss float64 `json:"profitLoss"`
		Available  float64 `json:"available"`
	} `json:"accountInfo"`
	CurrencyIsoCode string `json:"currencyIsoCode"`
	CurrencySymbol  string `json:"currencySymbol"`
	CurrentAccount  string `json:"currentAccountId"`
	StreamingHost   string `json:"streamingHost"`
	Accounts        []struct {
		AccountID   string `json:"accountId"`
		AccountName string `json:"accountName"`
		Preferred   bool   `json:"preferred"`
		AccountType string `json:"accountType"`
	} `json:"accounts"`
	ClientID              string `json:"clientId"`
	TimezoneOffset        int    `json:"timezoneOffset"`
	HasActiveDemoAccounts bool   `json:"hasActiveDemoAccounts"`
	HasActiveLiveAccounts bool   `json:"hasActiveLiveAccounts"`
	TrailingStopsEnabled  bool   `json:"trailingStopsEnabled"`
	Cst                   string `json:"cst"`
	SecurityToken         string `json:"securityToken"`
}

// {"status":"OK","destination":"quote","payload":{"epic":"EURUSD","product":"CFD","bid":1.06333,"bidQty":2.589885E7,"ofr":1.06339,"ofrQty":2.589885E7,"timestamp":1671005949822}}
type Quote struct {
	Status      string `json:"status"`
	Destination string `json:"destination"`
	Payload     struct {
		Epic    string  `json:"epic"`
		Product string  `json:"product"`
		Bid     float64 `json:"bid"`

		BidQty    float64 `json:"bidQty"`
		Ofr       float64 `json:"ofr"`
		OfrQty    float64 `json:"ofrQty"`
		Timestamp int64   `json:"timestamp"`
	} `json:"payload"`
}
