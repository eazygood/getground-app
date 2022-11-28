package domain

type Table struct {
	ID      int64  `json:"id" db:"id"`
	Seats   uint16 `json:"seats" db:"seats"`
	GuestID *int64 `json:"guest_id" db:"guest_id"`
	Guest   Guest  `json:"-" gorm:"foreignKey:ID;references:GuestID"`
}
