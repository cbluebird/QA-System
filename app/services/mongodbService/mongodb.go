package mongodbService

import (
	"QA-System/config/database"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type Answer struct {
	QuestionID int    `json:"question_id"` //问题ID
	SerialNum  int    `json:"serial_num"`  //问题序号
	Subject    string `json:"subject"`     //问题
	Content    string `json:"content"`     //回答内容
}

type AnswerSheet struct {
	SurveyID int      `json:"survey_id"` //问卷ID
    Time    string   `json:"time"`      //回答时间
	Answers  []Answer `json:"answers"`   //回答
}

func SaveAnswerSheet(answerSheet AnswerSheet) error {
	_, err := database.MDB.InsertOne(context.Background(), answerSheet)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func GetAnswerSheetBySurveyID(surveyID int) ([]AnswerSheet, error) {
    var answerSheets []AnswerSheet
    filter := bson.M{"surveyid": surveyID}
    cur, err := database.MDB.Find(context.Background(), filter)
    if err != nil {
        return nil, err
    }
    defer cur.Close(context.Background()) // 关闭游标
    for cur.Next(context.Background()) {
        var answerSheet AnswerSheet
        err := cur.Decode(&answerSheet)
        if err != nil {
            return nil, err
        }
        answerSheets = append(answerSheets, answerSheet)
    }
    if err := cur.Err(); err != nil {
        return nil, err
    }
    return answerSheets, nil
}

func DeleteAnswerSheetBySurveyID(surveyID int) error {
    filter := bson.M{"surveyid": surveyID}
    // 删除所有满足条件的文档
    _, err := database.MDB.DeleteMany(context.Background(), filter)
    if err != nil {
        return err
    }
    return nil
}