package domain

type GuestList struct {
	Table Table `json:"table"`
	Guest Guest `json:"guest"`
}
