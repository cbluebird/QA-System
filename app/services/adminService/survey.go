package adminService

import (
	"QA-System/app/models"
	"QA-System/config/database"
	"time"
)

type Option struct {
	SerialNum  int    `json:"serial_num"`  //选项序号
	Content    string `json:"content"`     //选项内容
	OptionType int    `json:"option_type"` //选项类型 1文字2图片
}

type Question struct {
	ID           int      `json:"id"`
	Subject      string   `json:"subject"`       //问题
	Description  string   `json:"description"`   //问题描述
	Img          string   `json:"img"`           //图片
	Required     bool     `json:"required"`      //是否必填
	QuestionType int      `json:"question_type"` //问题类型 1单选2多选3填空4简答5图片
	Reg          string   `json:"reg"`           //正则表达式
	Options      []Option `json:"options"`       //选项
}

func CreateSurvey(title string, desc string, img string, questions []Question, status int, time time.Time) error {
	var survey models.Survey
	survey.Title = title
	survey.Desc = desc
	survey.Img = img
	survey.Status = status
	survey.Deadline = time

	err := database.DB.Create(&survey).Error
	if err != nil {
		return err
	}
	for _, question := range questions {
		var q models.Question
		q.SurveyID = survey.ID
		q.Subject = question.Subject
		q.Description = question.Description
		q.Img = question.Img
		q.Required = question.Required
		q.QuestionType = question.QuestionType
		err := database.DB.Create(&q).Error
		if err != nil {
			return err
		}
		for _, option := range question.Options {
			var o models.Option
			o.QuestionID = q.ID
			o.Content = option.Content
			o.SerialNum = option.SerialNum
			o.OptionType = option.OptionType
			err := database.DB.Create(&o).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}
