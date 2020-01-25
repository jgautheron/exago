package model

import "time"

type Metadata struct {
	Image       string    `json:"image"`
	Description string    `json:"description"`
	Stars       int       `json:"stars"`
	LastPush    time.Time `json:"last_push"`
}
