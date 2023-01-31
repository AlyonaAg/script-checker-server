package producer

// Message ...
type Message struct {
	ID     int64  `json:"id"`
	Script string `json:"script"`
}
