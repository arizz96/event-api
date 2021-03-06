package v1

import (
	"encoding/json"
	"strings"

	"github.com/arizz96/event-api/config"
)

// NewMessage returns a new message of the given type
func NewMessage(kind string, raw []byte) (Message, error) {
	// Create the concrete message
	var msg Message
	switch kind {
	case "alias":
		msg = new(Alias)
	case "group":
		msg = new(Group)
	case "identify":
		msg = new(Identify)
	case "page":
		msg = new(Page)
	case "screen":
		msg = new(Screen)
	case "track":
		msg = new(Track)
	}

	// Unmarshal to type
	err := json.Unmarshal(raw, &msg)
	if err != nil {
		return nil, err
	}

	// Validate message
	err = msg.Validate()
	if err != nil {
		return nil, err
	}

	msg.WithType(kind)

	return msg, nil
}

func Authorize(writeKey string) bool {
	authorizedWriteKeys := strings.Split(config.AppConfig.AuthorizedWriteKeys, ",")

	return stringInSlice(writeKey, authorizedWriteKeys)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
