package consumer

import jsoniter "github.com/json-iterator/go"

type message struct {
	ID     int64  `json:"id"`
	Script string `json:"script"`
}

func (m *message) Unmarshall(data []byte) error {
	if err := jsoniter.Unmarshal(data, &m); err != nil {
		return err
	}

	return nil
}
