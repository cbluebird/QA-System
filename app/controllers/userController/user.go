package userController

import (
	"QA-System/app/apiException"
	"QA-System/app/services/adminService"
	"QA-System/app/services/userService"
	"QA-System/app/utils"
	"fmt"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

type SubmitServeyData struct {
	ID            int                         `json:"id" binding:"required"`
	QuestionsList []userService.QuestionsList `json:"questions_list"`
}

func SubmitSurvey(c *gin.Context) {
	var data SubmitServeyData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}
	// 判断问卷问题和答卷问题数目是否一致
	survey, err := userService.GetSurveyByID(data.ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	questions, err := userService.GetQuestionsBySurveyID(survey.ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	if len(questions) != len(data.QuestionsList) {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	// 判断填写时间是否在问卷有效期内
	if !survey.Deadline.IsZero() && survey.Deadline.Before(time.Now()) {
		utils.JsonErrorResponse(c, apiException.TimeBeyondError)
		return
	}
	// 逐个判断问题答案
	for _, q := range data.QuestionsList {
		question, err := userService.GetQuestionByID(q.QuestionID)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		// 判断必填字段是否为空
		if question.Required && q.Answer == "" {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		// 判断正则是否匹配
		if question.Reg != "" {
			match, err := regexp.MatchString(question.Reg, q.Answer)
			if err != nil {
				utils.JsonErrorResponse(c, apiException.ServerError)
				return
			}
			if !match {
				utils.JsonErrorResponse(c, apiException.RegError)
				return
			}
		}
		// 判断唯一字段是否唯一
		if question.Unique {
			fmt.Println(1)
			unique, err := userService.CheckUnique(data.ID, q.QuestionID, question.SerialNum, q.Answer)
			if err != nil {
				fmt.Println(2)
				utils.JsonErrorResponse(c, apiException.ServerError)
				return
			}
			if !unique {
				utils.JsonErrorResponse(c, apiException.UniqueError)
				return
			}

		}
	}
	fmt.Println(2)
	// 提交问卷
	err = userService.SubmitSurvey(data.ID, data.QuestionsList)
	if err != nil {
		fmt.Println(5)
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

type GetSurveyData struct {
	ID int `form:"id" binding:"required"`
}

type SurveyData struct {
	ID        int            `json:"id"`
	Time      string         `json:"time"`
	Desc      string         `json:"desc"`
	Img       string         `json:"img"`
	Questions []QuestionData `json:"questions"`
}

type QuestionData struct {
	ID           int          `json:"id"`
	SerialNum    int          `json:"serial_num"`
	Subject      string       `json:"subject"`
	Description  string       `json:"describe"`
	Required     bool         `json:"required"`
	Unique       bool         `json:"unique"`
	Img          string       `json:"img"`
	QuestionType int          `json:"question_type"`
	Reg          string       `json:"reg"`
	Options      []OptionData `json:"options"`
}

type OptionData struct {
	OptionType int    `json:"option_type"`
	Content    string `json:"content"`
	SerialNum  int    `json:"serial_num"`
}

// 用户获取问卷
func GetSurvey(c *gin.Context) {
	var data GetSurveyData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}
	// 获取问卷
	survey, err := adminService.GetSurveyByID(data.ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
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
		"time":      survey.Deadline.Format("2006-01-02 15:04:05"),
		"desc":      survey.Desc,
		"img":       survey.Img,
		"questions": questionsResponse,
	}

	utils.JsonSuccessResponse(c, response)
}
