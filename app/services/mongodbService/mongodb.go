package mongodbService

import (
	"QA-System/config/database"
	"context"
	"log"
)

type Answer struct {
	QuestionID int    `json:"question_id"` //问题ID
	SerialNum  int    `json:"serial_num"`  //问题序号
	Subject    string `json:"subject"`     //问题
	Content    string `json:"content"`     //回答内容
}

type AnswerSheet struct {
	SurveyID int      `json:"survey_id"` //问卷ID
	Answers  []Answer `json:"answers"`   //回答
}

func SaveAnswerSheet(answerSheet AnswerSheet) error {
	// Insert the answer sheet document
	_, err := database.MDB.InsertOne(context.Background(), answerSheet)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func GetAnswerSheetBySurveyID(surveyID int) ([]AnswerSheet, error) {
	var answerSheets []AnswerSheet
	cur, err := database.MDB.Find(context.Background(), AnswerSheet{SurveyID: surveyID})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var answerSheet AnswerSheet
		err := cur.Decode(&answerSheet)
		if err != nil {
			log.Fatal(err)
		}
		answerSheets = append(answerSheets, answerSheet)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	return answerSheets, nil
}