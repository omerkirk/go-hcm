package hcm

import (
	"errors"
	"strings"
)

var (
	// ErrInvalidMessage occurs if push notitication message is nil.
	ErrInvalidMessage = errors.New("message is invalid")

	// ErrInvalidTarget occurs if message topic is empty.
	ErrInvalidTarget = errors.New("topic is invalid or registration ids are not set")

	// ErrToManyRegIDs occurs when registration ids more then 1000.
	ErrToManyRegIDs = errors.New("too many registrations ids")

	// ErrInvalidTimeToLive occurs if TimeToLive more then 2419200.
	ErrInvalidTimeToLive = errors.New("messages time-to-live is invalid")
)

type notification struct {
	Title        string `json:"title,omitempty"`
	Body         string `json:"body,omitempty"`
	ChannelID    string `json:"android_channel_id,omitempty"`
	Icon         string `json:"icon,omitempty"`
	Sound        string `json:"sound,omitempty"`
	Badge        string `json:"badge,omitempty"`
	Tag          string `json:"tag,omitempty"`
	Color        string `json:"color,omitempty"`
	ClickAction  string `json:"click_action,omitempty"`
	BodyLocKey   string `json:"body_loc_key,omitempty"`
	BodyLocArgs  string `json:"body_loc_args,omitempty"`
	TitleLocKey  string `json:"title_loc_key,omitempty"`
	TitleLocArgs string `json:"title_loc_args,omitempty"`
}

type message struct {
	Token        []string       `json:"token,omitempty"`
	Data         string         `json:"data,omitempty"`
	Notification *notification  `json:"notification,omitempty"`
	Android      *androidConfig `json:"android,omitempty"`
	Topic        string         `json:"topic,omitempty"`
	Condition    string         `json:"condition,omitempty"`
}

type androidConfig struct {
	CollapseKey  string        `json:"collapse_key,omitempty"`
	TimeToLive   string        `json:"ttl,omitempty"`
	Notification *notification `json:"notification,omitempty"`
}

type Message struct {
	ValidateOnly bool     `json:"validate_only"`
	Message      *message `json:"message"`

	extra map[string]interface{}
}

func NewMessage(token []string, data string, ttl string, isProd bool, extra map[string]interface{}) *Message {
	msg := message{
		Token:   token,
		Data:    data,
		Android: &androidConfig{TimeToLive: ttl}}

	return &Message{
		ValidateOnly: !isProd,
		Message:      &msg,
		extra:        extra,
	}

}

func (msg *Message) SetExtra(extra map[string]interface{}) {
	msg.extra = extra
}

func (msg *Message) Extra() map[string]interface{} {
	return msg.extra
}

// Validate returns an error if the message is not well-formed.
func (msg *Message) Validate() error {
	if msg == nil {
		return ErrInvalidMessage
	}

	// validate target identifier: `to` or `condition`, or `registration_ids`
	opCnt := strings.Count(msg.Message.Condition, "&&") + strings.Count(msg.Message.Condition, "||")
	if (msg.Message.Condition == "" || opCnt > 2) && len(msg.Message.Token) == 0 {
		return ErrInvalidTarget
	}

	if len(msg.Message.Token) > 1000 {
		return ErrToManyRegIDs
	}
	return nil
}
