package domain

import "time"

type Guest struct {
	ID                 int64      `json:"id" db:"id"`
	Name               string     `json:"name" db:"name"`
	AccompanyingGuests uint16     `json:"accompanying_guests" db:"accompanying_guests"`
	TimeArrived        *time.Time `json:"time_arrived" db:"time_arrived"`
}
