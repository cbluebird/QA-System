package adminController

import (
	"QA-System/app/apiException"
	"QA-System/app/services/adminService"
	"QA-System/app/services/sessionService"
	"QA-System/app/services/userService"
	"QA-System/app/utils"
	"QA-System/config/config"
	"fmt"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"math"
	"os"
	"strconv"
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
	err = adminService.CreateSurvey(user.ID, data.Title, data.Desc, data.Img, data.Questions, data.Status, time)
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
	if (user.AdminType != 2) && (user.AdminType != 1 || survey.UserID != user.ID) && !adminService.UserInManage(user.ID, survey.ID) {
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
	if (user.AdminType != 2) && (user.AdminType != 1 || survey.UserID != user.ID) && !adminService.UserInManage(user.ID, survey.ID) {
		utils.JsonErrorResponse(c, apiException.NoPermission)
		return
	}
	//判断问卷状态
	if survey.Status != 1 {
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
	answers, num, err := adminService.GetSurveyAnswers(data.ID, data.PageNum, data.PageSize)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, gin.H{
		"answers_data":   answers,
		"total_page_num": math.Ceil(float64(*num) / float64(data.PageSize)),
	})
}

type GetAllSurveyData struct {
	PageNum  int    `form:"page_num" binding:"required"`
	PageSize int    `form:"page_size" binding:"required"`
	Title    string `form:"title"`
}

func GetAllSurvey(c *gin.Context) {
	var data GetAllSurveyData
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
	response := make([]interface{}, 0)
	var totalPageNum *int64
	if user.AdminType == 2 {
		response, totalPageNum = adminService.GetAllSurvey(data.PageNum, data.PageSize, data.Title)
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
		response, totalPageNum = adminService.ProcessResponse(response, data.PageNum, data.PageSize, data.Title)
	}

	utils.JsonSuccessResponse(c, gin.H{
		"survey_list":    response,
		"total_page_num": math.Ceil(float64(*totalPageNum) / float64(data.PageSize)),
	})
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
	if (user.AdminType != 2) && (user.AdminType != 1 || survey.UserID != user.ID) && !adminService.UserInManage(user.ID, survey.ID) {
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
			"id":            question.SerialNum,
			"subject":       question.Subject,
			"describe":      question.Description,
			"required":      question.Required,
			"unique":        question.Unique,
			"other_option":  question.OtherOption,
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

type DownloadFileData struct {
	ID int `form:"id" binding:"required"`
}

// 下载
func DownloadFile(c *gin.Context) {
	var data DownloadFileData
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
	// 判断权限
	if (user.AdminType != 2) && (user.AdminType != 1 || survey.UserID != user.ID) && !adminService.UserInManage(user.ID, survey.ID) {
		utils.JsonErrorResponse(c, apiException.NoPermission)
		return
	}
	// 获取数据
	answers, err := adminService.GetAllSurveyAnswers(data.ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	questionAnswers := answers.QuestionAnswers
	times := answers.Time
	// 创建一个新的Excel文件
	f := excelize.NewFile()
	streamWriter, err := f.NewStreamWriter("Sheet1")
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	// 设置字体样式
	styleID, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	// 计算每列的最大宽度
	maxWidths := make(map[int]int)
	maxWidths[0] = 7
	maxWidths[1] = 20
	for i, qa := range questionAnswers {
		maxWidths[i+2] = len(qa.Title)
		for _, answer := range qa.Answers {
			if len(answer) > maxWidths[i+2] {
				maxWidths[i+2] = len(answer)
			}
		}
	}
	// 设置列宽
	for colIndex, width := range maxWidths {
		if err := streamWriter.SetColWidth(colIndex+1, colIndex+1, float64(width)); err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
	}
	// 写入标题行
	rowData := make([]interface{}, 0)
	rowData = append(rowData, excelize.Cell{Value: "序号", StyleID: styleID}, excelize.Cell{Value: "提交时间", StyleID: styleID})
	for _, qa := range questionAnswers {
		rowData = append(rowData, excelize.Cell{Value: qa.Title, StyleID: styleID})
	}
	if err := streamWriter.SetRow("A1", rowData); err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	// 写入数据
	for i, time := range times {
		row := []interface{}{i + 1, time}
		for j, qa := range questionAnswers {
			if len(qa.Answers) <= i {
				continue
			}
			answer := qa.Answers[i]
			row = append(row, answer)
			colName, _ := excelize.ColumnNumberToName(j + 3)
			if err := f.SetCellValue("Sheet1", colName+strconv.Itoa(i+2), answer); err != nil {
				utils.JsonErrorResponse(c, apiException.ServerError)
				return
			}
		}
		if err := streamWriter.SetRow(fmt.Sprintf("A%d", i+2), row); err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
	}
	// 关闭
	if err := streamWriter.Flush(); err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	// 保存Excel文件
	fileName := survey.Title + ".xlsx"
	filePath := "./files/" + fileName
	if _, err := os.Stat("./files"); os.IsNotExist(err) {
		err := os.Mkdir("./files", 0755)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
	}
	// 删除旧文件
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
	}
	// 保存
	if err := f.SaveAs(filePath); err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, config.Config.GetString("url.host")+fileName)
}
