package session

type sessionKey int

const (
	loggingKey sessionKey = iota
	requestIdKey
)
