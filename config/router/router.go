package router

import (
	"QA-System/app/controllers/adminController"
	"QA-System/app/controllers/userController"
	"QA-System/app/midwares"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {

	const pre = "/api"

	api := r.Group(pre)
	{
		api.POST("/admin/reg", adminController.Register)
		api.POST("/admin/login", adminController.Login)
		user := api.Group("/user")
		{
			user.POST("/submit", userController.SubmitSurvey)
			user.GET("/get", userController.GetSurvey)
			user.POST("/upload", userController.UploadImg)
		}
		admin := api.Group("/admin", midwares.CheckLogin)
		{
			admin.POST("/create", adminController.CreateSurvey)
			admin.PUT("/update/status", adminController.UpdateSurveyStatus)
			admin.PUT("/update/questions", adminController.UpdateSurvey)
			admin.GET("/list/answers", adminController.GetSurveyAnswers)
			admin.DELETE("/delete", adminController.DeleteSurvey)

			admin.POST("/permission/create", adminController.CreatrPermission)
			admin.DELETE("/permission/delete", adminController.DeletePermission)

			admin.GET("/list/questions", adminController.GetAllSurvey)
			admin.GET("/single/question", adminController.GetSurvey)
			admin.GET("/download", adminController.DownloadFile)

			admin.GET("/log", adminController.GetLogMsg)

		}
	}
}
