package router

import (
	"QA-System/app/controllers/adminController"
	"QA-System/app/controllers/userController"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {

	const pre = "/api"

	api := r.Group(pre)
	{
		api.POST("/admin/login", adminController.Login)
		user := api.Group("/user")
		{
			user.POST("/submit", userController.SubmitSurvey)
		}
		admin := api.Group("/admin")
		{
			admin.POST("/create", adminController.CreateSurvey)
			admin.PUT("/update/status", adminController.UpdateSurveyStatus)
			admin.PUT("/update/questions", adminController.UpdateSurvey)
		}
	}
}
