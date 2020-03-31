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

// Notification specifies the predefined, user-visible key-value pairs of the
// notification payload.
type Notification struct {
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

// Message represents list of targets, options, and payload for HTTP JSON
// messages.
type Message struct {
	Token        []string               `json:"token,omitempty"`
	Data         map[string]interface{} `json:"data,omitempty"`
	Notification *Notification          `json:"notification,omitempty"`
	Android      *AndroidConfig         `json:"android,omitempty"`
	Topic        string                 `json:"condition,omitempty"`
	Condition    string                 `json:"condition,omitempty"`

	extra map[string]interface{}
}

// Message represents list of targets, options, and payload for HTTP JSON
// messages.
type AndroidConfig struct {
	CollapseKey  string        `json:"collapse_key,omitempty"`
	TimeToLive   *uint         `json:"ttl,omitempty"`
	Notification *Notification `json:"notification,omitempty"`
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
	opCnt := strings.Count(msg.Condition, "&&") + strings.Count(msg.Condition, "||")
	if (msg.Condition == "" || opCnt > 2) && len(msg.Token) == 0 {
		return ErrInvalidTarget
	}

	if len(msg.Token) > 1000 {
		return ErrToManyRegIDs
	}

	if msg.Android != nil && *msg.Android.TimeToLive > uint(2419200) {
		return ErrInvalidTimeToLive
	}
	return nil
}
