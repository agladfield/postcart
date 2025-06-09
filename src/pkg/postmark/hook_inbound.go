package postmark

type InboundData struct {
	FromName          string             `json:"FromName"`
	From              string             `json:"From"`
	FromFull          EmailAddressFull   `json:"FromFull"`
	To                string             `json:"To"`
	ToFull            []EmailAddressFull `json:"ToFull"`
	Cc                string             `json:"Cc,omitempty"`
	CcFull            []EmailAddressFull `json:"CcFull,omitempty"`
	Bcc               string             `json:"Bcc,omitempty"`
	BccFull           []EmailAddressFull `json:"BccFull,omitempty"`
	OriginalRecipient string             `json:"OriginalRecipient"`
	Subject           string             `json:"Subject"`
	MessageID         string             `json:"MessageID"`
	ReplyTo           string             `json:"ReplyTo"`
	MailboxHash       string             `json:"MailboxHash"`
	Date              string             `json:"Date"`
	TextBody          string             `json:"TextBody"`
	HTMLBody          string             `json:"HtmlBody"`
	StrippedTextReply string             `json:"StrippedTextReply"`
	Tag               string             `json:"Tag"`
	MessageStream     string             `json:"MessageStream"`
	Headers           []emailHeader      `json:"Headers"`
	Attachments       []EmailAttachment  `json:"Attachments,omitempty"`
}

type EmailAttachment struct {
	Name          string `json:"Name"`
	Content       string `json:"Content"`
	ContentID     string `json:"ContentID"`
	ContentType   string `json:"ContentType"`
	ContentLength int    `json:"ContentLength"`
}

type EmailAddressFull struct {
	Email       string `json:"Email"`
	Name        string `json:"Name"`
	MailboxHash string `json:"MailboxHash"`
}

type emailHeader struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}
