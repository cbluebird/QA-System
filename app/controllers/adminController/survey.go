package adminController

import (
	"QA-System/app/apiException"
	"QA-System/app/services/adminService"
	"QA-System/app/services/sessionService"
	"QA-System/app/utils"

	"time"

	"github.com/gin-gonic/gin"
)

// 新建问卷
type CreateSurveyData struct {
	Title     string                  `json:"title"`
	Desc      string                  `json:"desc" `
	Img       string                  `json:"img" `
	Status    int                     `json:"status" `
	Time      string                  `json:"time"`
	Questions []adminService.Question `json:"questions"`
}

func CreateSurvey(c *gin.Context) {
	var data CreateSurveyData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}
	//鉴权
	_, err = sessionService.GetUserSession(c)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.NotLogin)
		return
	}
	//解析时间转换为中国时间(UTC+8)
	time, err := time.Parse(time.RFC3339, data.Time)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//创建问卷
	err = adminService.CreateSurvey(data.Title, data.Desc, data.Img, data.Questions, data.Status, time)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

//修改问卷状态
type UpdateSurveyStatusData struct {
	ID     int `json:"id" binding:"required"`
	Status int `json:"status" binding:"required"`
}

func UpdateSurveyStatus(c *gin.Context) {
	var data UpdateSurveyStatusData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}
	//鉴权
	_, err = sessionService.GetUserSession(c)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.NotLogin)
		return
	}
	// 获取问卷
	survey, err := adminService.GetSurveyByID(data.ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//判断问卷状态
	if survey.Status == data.Status {
		utils.JsonErrorResponse(c, apiException.StatusRepeatError)
		return
	}
	//修改问卷状态
	err = adminService.UpdateSurveyStatus(data.ID, data.Status)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

