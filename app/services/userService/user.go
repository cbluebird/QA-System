package userService

import (
	"QA-System/app/models"
	"QA-System/app/services/mongodbService"
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
	SerialNum    int      `json:"serial_num"`    //题目序号
	Subject      string   `json:"subject"`       //问题
	Description  string   `json:"description"`   //问题描述
	Img          string   `json:"img"`           //图片
	Required     bool     `json:"required"`      //是否必填
	Unique       bool     `json:"unique"`        //是否唯一
	QuestionType int      `json:"question_type"` //问题类型 1单选2多选3填空4简答5图片
	Reg          string   `json:"reg"`           //正则表达式
	Options      []Option `json:"options"`       //选项
}

type QuestionsList struct {
	QuestionID int    `json:"question_id" binding:"required"`
	SerialNum  int    `json:"serial_num"`
	Answer     string `json:"answer"`
}

func GetSurveyByID(id int) (models.Survey, error) {
	var survey models.Survey
	err := database.DB.Where("id = ?", id).First(&survey).Error
	return survey, err
}

func GetQuestionsBySurveyID(id int) ([]models.Question, error) {
	var questions []models.Question
	err := database.DB.Where("survey_id = ?", id).Find(&questions).Error
	return questions, err
}

func GetOptionsByQuestionID(questionId int) ([]models.Option, error) {
	var options []models.Option
	err := database.DB.Where("question_id = ?", questionId).Find(&options).Error
	return options, err
}

func GetQuestionByID(id int) (models.Question, error) {
	var question models.Question
	err := database.DB.Where("id = ?", id).First(&question).Error
	return question, err
}

func CheckUnique(sid int, qid int, serial_num int, content string) (bool, error) {
	var answerSheets []mongodbService.AnswerSheet
	answerSheets,_, err := mongodbService.GetAnswerSheetBySurveyID(sid,0,0)
	if err != nil {
		return false, err
	}

	for _, answerSheet := range answerSheets {
		for _, answer := range answerSheet.Answers {
			if answer.QuestionID == qid && answer.SerialNum == serial_num && answer.Content == content {
				return false, nil
			}
		}
	}
	return true, nil
}

func SubmitSurvey(sid int, data []QuestionsList) error {
	var answerSheet mongodbService.AnswerSheet
	answerSheet.SurveyID = sid
	answerSheet.Time = time.Now().Format("2006-01-02 15:04:05")
	for _, q := range data {
		var answer mongodbService.Answer
		answer.QuestionID = q.QuestionID
		answer.SerialNum = q.SerialNum
		answer.Content = q.Answer
		answerSheet.Answers = append(answerSheet.Answers, answer)
	}
	err := mongodbService.SaveAnswerSheet(answerSheet)
	if err != nil {
		return err
	}
	return nil
}
