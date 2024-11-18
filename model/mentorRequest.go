package model

import "time"

type MentorRequest struct {
	ID          int64     `json:"id"`
	StudentID   int64     `json:"student_id"`
	MentorID    int64     `json:"mentor_id"`
	Project     string    `json:"project_title"`
	Description string    `json:"project_description"`
	Status      string    `json:"status"`
	CreateAt    time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
