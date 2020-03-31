package hcm

import (
	"errors"
)

const (
	RespCodeSuccess        = "80000000"
	RespCodePartialSuccess = "80100000"
)

var (
	ErrOAuth         = errors.New("oauth authentication error")
	ErrTokenExpired  = errors.New("token expired")
	ErrInvalidMsg    = errors.New("incorrect message structure.")
	ErrMessageTooBig = errors.New("message body size exceeded the default value (4 KB)")
	ErrCannotSend    = errors.New("messages cannot be sent to the app")
	ErrInvalidToken  = errors.New("invalid token")
	ErrUnknown       = errors.New("unknown error")
)

var (
	errMap = map[string]error{
		"80200001": ErrOAuth,
		"80200003": ErrTokenExpired,
		"80100003": ErrInvalidMsg,
		"80300008": ErrMessageTooBig,
		"80300002": ErrCannotSend,
		"80300007": ErrInvalidToken,
	}
)

// connectionError represents connection errors such as timeout error, etc.
// Implements `net.Error` interface.
type connectionError string

func (err connectionError) Error() string {
	return string(err)
}

func (err connectionError) Temporary() bool {
	return true
}

func (err connectionError) Timeout() bool {
	return true
}

// serverError represents internal server errors.
// Implements `net.Error` interface.
type serverError string

func (err serverError) Error() string {
	return string(err)
}

func (serverError) Temporary() bool {
	return true
}

func (serverError) Timeout() bool {
	return false
}

// Response represents the FCM server's response to the application
// server's sent message.
type Response struct {
	Code      string `json:"code"`
	Message   string `json:"msg"`
	RequestID string `json:"request_id"`
}
