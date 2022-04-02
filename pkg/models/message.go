package models

import "google.golang.org/api/gmail/v1"

type Message struct {
	ID           string             `json:"id"`
	HistoryID    uint64             `json:"historyId"`
	InternalDate int64              `json:"internalDate"`
	LabelIDS     []string           `json:"labelIds"`
	Payload      *gmail.MessagePart `json:"payload"`
	SizeEstimate int64              `json:"sizeEstimate"`
	Snippet      string             `json:"snippet"`
	ThreadID     string             `json:"threadId"`
	HTML         string             `json:"html"`
}
