package structs

type GeoData struct {
	Standard struct {
		City        string `json:"city"`
		Countryname string `json:"countryname"`
	} `json:"standard"`
	Longt string `json:"longt"`
	Latt  string `json:"latt"`
}
