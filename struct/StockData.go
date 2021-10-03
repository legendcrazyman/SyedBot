package structs

type StockData struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Symbol             string  `json:"symbol"`
				RegularMarketPrice float64 `json:"regularMarketPrice"`
				PreviousClose      float64 `json:"previousClose"`
			} `json:"meta"`
		} `json:"result"`
	} `json:"chart"`
}