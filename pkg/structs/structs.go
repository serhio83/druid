package structs

type Message struct {
	Receiver string  `json:"receiver"`
	Alerts   []Alert `json:"alerts"`
}

type Alert struct {
	Status      string     `json:"status"`
	Labels      Label      `json:"labels"`
	Annotations Annotation `json:"annotations"`
}

type Label struct {
	Alertname string `json:"alertname"`
	Instance  string `json:"instance"`
	Job       string `json:"job"`
	Monitor   string `json:"monitor"`
	Severity  string `json:"severity"`
}

type Annotation struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

type Responce struct {
	Status string `json:"status"`
}
