package adminController

import (
	"QA-System/app/apiException"
	"QA-System/app/services/adminService"
	"QA-System/app/services/sessionService"
	"QA-System/app/utils"

	"github.com/gin-gonic/gin"
)

type CreatePermissionData struct {
	UserName string `json:"username"`
	SurveyID int    `json:"survey_id"`
}

func CreatrPermission(c *gin.Context) {
	var data CreatePermissionData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}
	//鉴权
	admin, err := sessionService.GetUserSession(c)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.NotLogin)
		return
	}
	if admin.AdminType != 2 {
		utils.JsonErrorResponse(c, apiException.NoPermission)
		return
	}
	user ,err:= adminService.GetUserByName(data.UserName)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	_, err = adminService.GetSurveyByID(data.SurveyID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//创建权限
	err = adminService.CreatePermission(user.ID, data.SurveyID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

type DeletePermissionData struct {
	UserName string `form:"username"`
	SurveyID int    `form:"survey_id"`
}

func DeletePermission(c *gin.Context) {
	var data DeletePermissionData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}
	//鉴权
	admin, err := sessionService.GetUserSession(c)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.NotLogin)
		return
	}
	if admin.AdminType != 2 {
		utils.JsonErrorResponse(c, apiException.NoPermission)
		return
	}
	user ,err:= adminService.GetUserByName(data.UserName)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	_, err = adminService.GetSurveyByID(data.SurveyID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//删除权限
	err = adminService.DeletePermission(user.ID, data.SurveyID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}