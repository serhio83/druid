package structs

//easyjson:json
type Message struct {
	Receiver string  `json:"receiver"`
	Alerts   []Alert `json:"alerts"`
}

//easyjson:json
type Alert struct {
	Status      string     `json:"status"`
	Labels      Label      `json:"labels"`
	Annotations Annotation `json:"annotations"`
}

//easyjson:json
type Label struct {
	Alertname string `json:"alertname"`
	Instance  string `json:"instance"`
	Job       string `json:"job"`
	Monitor   string `json:"monitor"`
	Severity  string `json:"severity"`
}

//easyjson:json
type Annotation struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

//easyjson:json
type Responce struct {
	Status string `json:"status"`
}
