package model

type Mentor struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	Experience int64  `json:"experience"`
}
