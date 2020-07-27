package utils

// RichError is a common interface for error types that can provide more detail
type RichError interface {
	error
	Code() string
	Extra() map[string]string
}
