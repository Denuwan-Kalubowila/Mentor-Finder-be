package utils

import (
	"fmt"
	"net/smtp"
)

type EmailService struct {
	From     string
	Password string
	Host     string
	Port     string
}

func NewEmailService() *EmailService {
	return &EmailService{
		From:     "vishwa.caprisious@gmail.com",
		Password: "ufgmhadikmzbmbtu",
		Host:     "smtp.gmail.com",
		Port:     "587",
	}
}

func (emailservice *EmailService) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", emailservice.From, emailservice.Password, emailservice.Host)
	err := smtp.SendMail(emailservice.Host+":"+emailservice.Port, auth, emailservice.From, []string{to}, []byte("Subject: "+subject+"\r\n"+body))

	if err != nil {
		return err
	}
	fmt.Println("Email sent successfully")
	return nil
}

func EmailToMentor(studentName, mentorName, projectTitle, link string) string {
	return fmt.Sprintf(`
	Dear %s,

	We hope this email finds you well. We're excited to inform you that you have received a new mentor request for your project "%s".

	Request Details:
	- Student: %s
	- Project: %s

	To review and respond to this request, please click on the following link:
	%s

	If you have any questions or need assistance, please don't hesitate to reach out to our support team at support@mentorfinder.com.

	Best regards,
	The Mentor Finder Team

	Note: This is an automated message. Please do not reply directly to this email.
	`, mentorName, projectTitle, studentName, projectTitle, link)
}

func EmailtoStudent(studentName, mentorName, projectTitle, status string) string {
	var statusMessage, nextSteps string
	if status == "Accepted" {
		statusMessage = "We're excited to inform you that your mentor request has been accepted!"
		nextSteps = `Next Steps:
		1. Reach out to your mentor to schedule your first meeting.
		2. Prepare any questions or specific areas you'd like to focus on.
		3. Review your project goals and be ready to discuss them with your mentor.`
	} else {
		statusMessage = "We regret to inform you that your mentor request has not been accepted at this time."
		nextSteps = `Next Steps:
		1. Don't be discouraged! This is a common part of the mentorship process.
		2. Review other available mentors who might be a good fit for your project.
		3. Feel free to reach out to our support team if you need assistance finding a suitable mentor.`
	}

	return fmt.Sprintf(`
		Dear %s,

		%s

		Request Details:
		- Mentor: %s
		- Project: %s
		- Status: %s

		%s

		If you have any questions or need further assistance, please don't hesitate to contact our support team at support@mentorfinder.com.

		Best regards,
		The Mentor Finder Team

		Note: This is an automated message. Please do not reply directly to this email.
		`, studentName, statusMessage, mentorName, projectTitle, status, nextSteps)
}
