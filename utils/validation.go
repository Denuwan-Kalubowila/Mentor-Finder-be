package utils

import (
	"errors"

	"github.com/Denuwan-Kalubowila/mentor-finder/model"
)

// StudentValidation validates the student data
func MentorRequestValidation(request *model.MentorRequest) error {
	if request.MentorID <= 0 {
		return errors.New("mentor id is required")
	}
	if request.StudentID <= 0 {
		return errors.New("student id is required")
	}
	if request.Project == "" {
		return errors.New("project title is required")
	}
	if request.Description == "" {
		return errors.New("description is required")
	}
	return nil
}
