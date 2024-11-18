package model

type Student struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Password   string `json:"password"`
	University string `json:"university"`
	Degree     string `json:"degree"`
}
