package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Denuwan-Kalubowila/mentor-finder/model"
	"github.com/Denuwan-Kalubowila/mentor-finder/utils"
	"github.com/labstack/echo/v4"
)

type RequestHandeler struct {
	DB *sql.DB
	Es *utils.EmailService
}

// request a mentor
func (db *RequestHandeler) RequestToMentor(c echo.Context) error {

	var mentorreq model.MentorRequest

	err := c.Bind(&mentorreq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
		})
	}

	if err := utils.MentorRequestValidation(&mentorreq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	mentorreq.Status = "Pending"
	mentorreq.CreateAt = time.Now()
	mentorreq.UpdatedAt = time.Now()

	query := `INSERT INTO MentorRequest(StudentID, MentorID, Project, Description,Status, CreatedAt, UpdatedAt) values (?,?,?,?,?,?,?)`

	result, err := db.DB.Exec(query, mentorreq.StudentID, mentorreq.MentorID, mentorreq.Project, mentorreq.Description, mentorreq.Status, mentorreq.CreateAt, mentorreq.UpdatedAt)

	if err != nil {
		log.Printf("Error requesting mentor: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error while requesting mentor",
		})
	}

	stdName, _, err := utils.GetStudentData(db.DB, mentorreq.StudentID)
	if err != nil {
		log.Printf("Error getting student data: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error while getting student data",
		})
	}
	mentorName, mentorMail, err := utils.GetMentorData(db.DB, mentorreq.MentorID)
	if err != nil {
		log.Printf("Error getting mentor data: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error while getting mentor data",
		})
	}

	id, _ := result.LastInsertId()
	link := "http://localhost:5000/mentorequest/" + strconv.FormatInt(id, 10)
	// send email to mentor

	go func() {
		es := utils.NewEmailService()
		err = es.SendEmail(mentorMail, "Mentor Request", utils.EmailToMentor(mentorName, stdName, mentorreq.Project, link))
		if err != nil {
			log.Printf("Error sending email: %v", err)
		}
	}()

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Request sent successfully",
		"id":      id,
	})
}

// accept mentor request
func (db *RequestHandeler) AcceptStudentRequest(c echo.Context) error {
	var AcceptedStatus struct {
		Accepeted bool `json:"accepted"`
	}

	var status string

	err := c.Bind(&AcceptedStatus)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	userParams := c.Param("id")

	reqID, err := strconv.ParseInt(userParams, 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid mentor ID"})
	}

	if !AcceptedStatus.Accepeted {
		status = "Rejected"
	} else {
		status = "Accepted"
	}
	query := `UPDATE mentorrequest SET Status=? WHERE  ID=?`
	_, err = db.DB.Exec(query, status, reqID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to accept mentor request"})
	}

	var (
		studentEmail string
		studentName  string
		projectTitle string
	)
	err = db.DB.QueryRow("SELECT name,email,Project FROM mentorrequest as mr INNER JOIN students as s on mr.StudentID = s.id WHERE mr.ID = ?", reqID).Scan(&studentName, &studentEmail, &projectTitle)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to fetch student email",
		})
	}

	go func() {
		es := utils.NewEmailService()
		//link := "http://localhost:5000/mentorequest/" + strconv.FormatInt(reqID, 10)
		err = es.SendEmail(studentEmail, "Request Confirmation", utils.EmailtoStudent(studentName, studentName, projectTitle, status))
		if err != nil {
			log.Printf("Error sending email: %v", err)
		}
	}()

	return c.JSON(http.StatusOK, map[string]string{"message": "Mentor request" + status + " successfully"})
}

func (db *RequestHandeler) GetRequestDataById(c echo.Context) error {
	var mentorreq *model.MentorRequest = &model.MentorRequest{}

	reqID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid mentor request ID"})
	}

	query := "SELECT ID, StudentID, MentorId, Status FROM mentorrequest WHERE ID = ?"
	err = db.DB.QueryRow(query, reqID).Scan(&mentorreq.ID, &mentorreq.StudentID, &mentorreq.MentorID, &mentorreq.Status)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "Mentor request not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get mentor request data"})
	}

	return c.JSON(http.StatusOK, mentorreq)
}
