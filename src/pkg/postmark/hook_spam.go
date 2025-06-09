package postmark

type SpamComplaintData struct {
	RecordType    string            `json:"RecordType"`
	MessageStream string            `json:"MessageStream"`
	ID            int               `json:"ID"`
	Type          string            `json:"Type"`
	TypeCode      int               `json:"TypeCode"`
	Name          string            `json:"Name"`
	Tag           string            `json:"Tag"`
	MessageID     string            `json:"MessageID"`
	Metadata      map[string]string `json:"Metadata"`
	ServerID      int               `json:"ServerID"`
	Description   string            `json:"Description"`
	Details       string            `json:"Details"`
	Email         string            `json:"Email"`
	From          string            `json:"From"`
	BouncedAt     string            `json:"BouncedAt"`
	DumpAvailable bool              `json:"DumpAvailable"`
	Inactive      bool              `json:"Inactive"`
	CanActivate   bool              `json:"CanActivate"`
	Subject       string            `json:"Subject"`
	Content       string            `json:"Content"`
}
