package main

import (
	"net/http"

	"github.com/Denuwan-Kalubowila/mentor-finder/auth"
	"github.com/Denuwan-Kalubowila/mentor-finder/config"
	"github.com/Denuwan-Kalubowila/mentor-finder/handlers"
	"github.com/Denuwan-Kalubowila/mentor-finder/utils"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	DB := config.InitDataBase()
	defer DB.Close()

	data := &handlers.Handlers{DB: DB}
	authData := &handlers.AuthHandler{Db: DB}
	reqData := &handlers.RequestHandeler{DB: DB, Es: &utils.EmailService{}}
	//routes
	e.GET("/", func(c echo.Context) error {
		return c.String((http.StatusOK), "welocome to the Mentor-Finder!")
	})

	e.GET("/students", data.GetStudents)
	e.POST("/register/student", data.RegisterStudent)
	e.POST("/request_mentor", reqData.RequestToMentor)

	e.POST("/login", authData.Login)
	e.POST("/refresh", authData.RefreshToken)

	e.GET("/mentors", data.GetMentors)
	e.POST("/register/mentor", data.RegisterMentor)
	e.POST("/mentor/role", data.GetMentorByRole)
	e.PUT("/mentor/profile", data.UpdateMentor, auth.JWTMiddleware)

	e.GET("/mentorequest/:id", reqData.GetRequestDataById)
	e.PUT("/mentorequest/:id", reqData.AcceptStudentRequest)

	e.GET("/protected", handlers.ProtectedHandler, auth.JWTMiddleware)

	e.Logger.Fatal(e.Start(":5000"))
}
