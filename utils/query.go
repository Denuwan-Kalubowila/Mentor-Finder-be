package utils

import "database/sql"

// get student data
func GetStudentData(db *sql.DB, stdID int64) (stdName string, stdEmail string, error error) {
	row := db.QueryRow("SELECT name,email FROM students WHERE id=?", stdID)
	err := row.Scan(&stdName, &stdEmail)
	if err != nil {
		return "", "", err
	}
	return stdName, stdEmail, nil
}

// get mentor data
func GetMentorData(db *sql.DB, mentorID int64) (mentoName string, mentoEmail string, error error) {
	row := db.QueryRow("SELECT name,email FROM mentors WHERE id=?", mentorID)
	err := row.Scan(&mentoName, &mentoEmail)
	if err != nil {
		return "", "", err
	}
	return mentoName, mentoEmail, nil
}
