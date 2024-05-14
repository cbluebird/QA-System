package adminController

import (
	"QA-System/app/apiException"
	"QA-System/app/services/adminService"
	"QA-System/app/services/sessionService"
	"QA-System/app/utils"

	"github.com/gin-gonic/gin"
)

type LogData struct {
	Num     int `form:"num" json:"num"`
	LogType int `form:"log_type" binding:"required,oneof=0 1 2 3 4"`
}

func GetLogMsg(c *gin.Context){
	var data LogData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		c.Error(&gin.Error{Err: err, Type: gin.ErrorTypeBind})
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}
	//鉴权
	_, err = sessionService.GetUserSession(c)
	if err != nil {
		c.Error(&gin.Error{Err: err, Type: gin.ErrorTypePublic})
		utils.JsonErrorResponse(c, apiException.NotLogin)
		return
	}
	response, err := adminService.GetLastLinesFromLogFile("app.log", data.Num, data.LogType)
	if err != nil {
		c.Error(&gin.Error{Err: err, Type: gin.ErrorTypePublic})
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, response)
}