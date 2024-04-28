package router

import (
	"QA-System/app/controllers/adminController"
	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {

	const pre = "/api"

	api := r.Group(pre)
	{
		api.POST("/admin/login", adminController.Login)
		admin := api.Group("/admin")
		{
			admin.POST("/create", adminController.CreateSurvey)
			admin.PUT("/update/status", adminController.UpdateSurveyStatus)
		}
	}
}
