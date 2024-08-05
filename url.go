package gostrava

type URL struct {
	URL       string `json:"url"`
	RetinaURL string `json:"retirna_url"`
	DarkURL   string `json:"dark_url"`
	LightURL  string `json:"light_url"`
}
