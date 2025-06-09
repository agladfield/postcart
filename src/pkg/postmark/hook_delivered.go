package postmark

type DeliveredData struct {
	RecordType    string            `json:"RecordType"`
	ServerID      int               `json:"ServerID"`
	MessageStream string            `json:"MessageStream"`
	MessageID     string            `json:"MessageID"`
	Recipient     string            `json:"Recipient"`
	Tag           string            `json:"Tag"`
	DeliveredAt   string            `json:"DeliveredAt"`
	Details       string            `json:"Details"`
	Metadata      map[string]string `json:"Metadata"`
}
