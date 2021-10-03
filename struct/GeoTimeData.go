package structs

type GeoTimeData struct {
	ZoneName     string `json:"zoneName"`
	Abbreviation string `json:"abbreviation"`
	GmtOffset    int    `json:"gmtOffset"`
	Dst          string `json:"dst"`
	Timestamp    int    `json:"timestamp"`
}