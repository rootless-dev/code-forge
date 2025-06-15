package entity

import "time"

type StateData struct {
	ID        string
	Timestamp time.Time
	IP        string
	UserAgent string
}
