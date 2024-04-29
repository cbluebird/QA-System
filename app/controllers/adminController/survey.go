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
	user, err := sessionService.GetUserSession(c)
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
	err = adminService.CreateSurvey(user.ID,data.Title, data.Desc, data.Img, data.Questions, data.Status, time)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

// 修改问卷状态
type UpdateSurveyStatusData struct {
	ID     int `json:"id" binding:"required"`
	Status int `json:"status" binding:"required,oneof=1 2"`
}

func UpdateSurveyStatus(c *gin.Context) {
	var data UpdateSurveyStatusData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}
	//鉴权
	user, err := sessionService.GetUserSession(c)
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
	//判断权限
	if (user.AdminType !=2)&&(user.AdminType !=1||survey.UserID != user.ID) {
		utils.JsonErrorResponse(c, apiException.NoPermission)
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

type UpdateSurveyData struct {
	ID        int                     `json:"id" binding:"required"`
	Title     string                  `json:"title"`
	Desc      string                  `json:"desc" `
	Img       string                  `json:"img" `
	Time      string                  `json:"time"`
	Questions []adminService.Question `json:"questions"`
}

func UpdateSurvey(c *gin.Context) {
	var data UpdateSurveyData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}
	//鉴权
	user, err := sessionService.GetUserSession(c)
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
	//判断权限
	if (user.AdminType !=2)&&(user.AdminType !=1||survey.UserID != user.ID) {
		utils.JsonErrorResponse(c, apiException.NoPermission)
		return
	}
	if !adminService.UserInManage(user.ID,survey.ID){
		utils.JsonErrorResponse(c, apiException.NoPermission)
		return
	}
	//判断问卷状态
	if survey.Status !=1 {
		utils.JsonErrorResponse(c, apiException.StatusRepeatError)
		return
	}
	// 判断问卷的填写数量是否为零
	if survey.Num != 0 {
		utils.JsonErrorResponse(c, apiException.SurveyNumError)
		return
	}
	//解析时间转换为中国时间(UTC+8)
	time, err := time.Parse(time.RFC3339, data.Time)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//修改问卷
	err = adminService.UpdateSurvey(data.ID, data.Title, data.Desc, data.Img, data.Questions, time)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

// 删除问卷
type DeleteSurveyData struct {
	ID int `form:"id" binding:"required"`
}

func DeleteSurvey(c *gin.Context) {
	var data DeleteSurveyData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}
	//鉴权
	user, err := sessionService.GetUserSession(c)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.NotLogin)
		return
	}
	// 获取问卷
	survey, err := adminService.GetSurveyByID(data.ID)
	if err == gorm.ErrRecordNotFound {
		utils.JsonErrorResponse(c, apiException.SurveyNotExist)
		return
	} else if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//判断权限
	if (user.AdminType != 2) && (user.AdminType != 1 || survey.UserID != user.ID) && !adminService.UserInManage(user.ID, survey.ID) {
		utils.JsonErrorResponse(c, apiException.NoPermission)
		return
	}
	//删除问卷
	err = adminService.DeleteSurvey(data.ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

// 获取问卷收集数据
type GetSurveyAnswersData struct {
	ID       int `form:"id" binding:"required"`
	PageNum  int `form:"page_num" binding:"required"`
	PageSize int `form:"page_size" binding:"required"`
}


func GetSurveyAnswers(c *gin.Context) {
	var data GetSurveyAnswersData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}
	//鉴权
	user, err := sessionService.GetUserSession(c)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.NotLogin)
		return
	}
	// 获取问卷
	survey, err := adminService.GetSurveyByID(data.ID)
	if err == gorm.ErrRecordNotFound {
		utils.JsonErrorResponse(c, apiException.SurveyNotExist)
		return
	} else if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//判断权限
	if (user.AdminType != 2) && (user.AdminType != 1 || survey.UserID != user.ID) && !adminService.UserInManage(user.ID, survey.ID) {
		utils.JsonErrorResponse(c, apiException.NoPermission)
		return
	}
	//获取问卷收集数据
	var num *int64
	answers,num, err := adminService.GetSurveyAnswers(data.ID, data.PageNum, data.PageSize)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, gin.H{
		"data": answers,
		"total_page_num": math.Ceil(float64(*num) / float64(data.PageSize)),
	})
}

func GetAllSurvey(c *gin.Context) {
	user, err := sessionService.GetUserSession(c)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.NotLogin)
		return
	}
	// 获取问卷
	response := make([]interface{}, 0)
	if user.AdminType == 2 {
		response, err = adminService.GetAllSurvey()
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
	} else {
		response, err = adminService.GetAllSurveyByUserID(user.ID)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		managedSurveys, err := adminService.GetManageredSurveyByUserID(user.ID)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		for _, manage := range managedSurveys {
			managedSurvey, err := adminService.GetSurveyByID(manage.SurveyID)
			if err != nil {
				utils.JsonErrorResponse(c, apiException.ServerError)
				return
			}
			managedSurveyResponse := map[string]interface{}{
				"id":     managedSurvey.ID,
				"title":  managedSurvey.Title,
				"status": managedSurvey.Status,
				"num":    managedSurvey.Num,
			}
			response = append(response, managedSurveyResponse)
		}
	}

	utils.JsonSuccessResponse(c, response)
}

type GetSurveyData struct {
	ID int `form:"id" binding:"required"`
}

type SurveyData struct {
	ID        int                    `json:"id"`
	Time      string                 `json:"time"`
	Desc      string                 `json:"desc"`
	Img       string                 `json:"img"`
	Questions []userService.Question `json:"questions"`
}

// 管理员获取问卷题面
func GetSurvey(c *gin.Context) {
	var data GetSurveyData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}
	user, err := sessionService.GetUserSession(c)
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
	//判断权限
	if (user.AdminType != 2) && (user.AdminType != 1 || survey.UserID != user.ID) {
		utils.JsonErrorResponse(c, apiException.NoPermission)
		return
	}
	if !adminService.UserInManage(user.ID, survey.ID) {
		utils.JsonErrorResponse(c, apiException.NoPermission)
		return
	}
	// 获取相应的问题
	questions, err := userService.GetQuestionsBySurveyID(survey.ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	// 构建问卷响应
	questionsResponse := make([]map[string]interface{}, 0)
	for _, question := range questions {
		options, err := userService.GetOptionsByQuestionID(question.ID)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		optionsResponse := make([]map[string]interface{}, 0)
		for _, option := range options {
			optionResponse := map[string]interface{}{
				"option_type": option.OptionType,
				"content":     option.Content,
				"serial_num":  option.SerialNum,
			}
			optionsResponse = append(optionsResponse, optionResponse)
		}
		questionMap := map[string]interface{}{
			"id":            question.ID,
			"serial_num":    question.SerialNum,
			"subject":       question.Subject,
			"describe":      question.Description,
			"required":      question.Required,
			"unique":        question.Unique,
			"img":           question.Img,
			"question_type": question.QuestionType,
			"reg":           question.Reg,
			"options":       optionsResponse,
		}
		questionsResponse = append(questionsResponse, questionMap)
	}
	response := map[string]interface{}{
		"id":        survey.ID,
		"title":     survey.Title,
		"time":      survey.Deadline.Format("2006-01-02 15:04:05"),
		"desc":      survey.Desc,
		"img":       survey.Img,
		"questions": questionsResponse,
	}

	utils.JsonSuccessResponse(c, response)
}

