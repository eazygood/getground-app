package domain

type Table struct {
	ID     int64  `json:"id" db:"id"`
	Seats  uint16 `json:"seats" db:"seats"`
	IsBusy bool   `json:"is_busy" db:"is_busy"`
}
