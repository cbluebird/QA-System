package userService

import (
	"QA-System/app/models"
	"QA-System/app/services/mongodbService"
	"QA-System/config/database"

)

type QuestionsList struct {
	QuestionID int    `json:"question_id"`
	SerialNum   int    `json:"serial_num"`
	Answer      string `json:"answer"`
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

func GetQuestionByID(id int) (models.Question, error) {
	var question models.Question
	err := database.DB.Where("id = ?", id).First(&question).Error
	return question, err
}

func CheckUnique(sid int, qid int, serial_num int,content string) (bool, error) {
	var answerSheets []mongodbService.AnswerSheet
	answerSheets, err := mongodbService.GetAnswerSheetBySurveyID(sid)
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