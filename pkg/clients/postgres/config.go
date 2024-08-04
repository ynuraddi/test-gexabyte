package postgres

import "time"

const (
	pingTimeout = 5 * time.Second
)

// TODO: max idle/open/lifetime connection...
type Config struct {
	DSN string
}
