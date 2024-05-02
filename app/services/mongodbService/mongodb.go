package mongodbService

import (
	"QA-System/config/database"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Answer struct {
	QuestionID int    `json:"question_id"` //问题ID
	SerialNum  int    `json:"serial_num"`  //问题序号
	Subject    string `json:"subject"`     //问题
	Content    string `json:"content"`     //回答内容
}

type AnswerSheet struct {
	SurveyID int      `json:"survey_id"` //问卷ID
	Time     string   `json:"time"`      //回答时间
	Answers  []Answer `json:"answers"`   //回答
}

func SaveAnswerSheet(answerSheet AnswerSheet) error {
	_, err := database.MDB.InsertOne(context.Background(), answerSheet)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func GetAnswerSheetBySurveyID(surveyID int, pageNum int, pageSize int) ([]AnswerSheet, *int64, error) {
	var answerSheets []AnswerSheet
	filter := bson.M{"surveyid": surveyID}

	// 设置总记录数查询过滤条件
	countFilter := bson.M{"surveyid": surveyID}

	// 设置总记录数查询选项
	countOpts := options.Count()

	// 执行总记录数查询
	total, err := database.MDB.CountDocuments(context.Background(), countFilter, countOpts)
	if err != nil {
		return nil, nil, err
	}

	// 设置分页查询选项
	opts := options.Find()
	if pageNum != 0 && pageSize != 0 {
		opts.SetSkip(int64((pageNum - 1) * pageSize)) // 计算要跳过的文档数
		opts.SetLimit(int64(pageSize))             // 设置返回的文档数
	}
	// 执行分页查询
	cur, err := database.MDB.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, nil, err
	}
	defer cur.Close(context.Background())

	// 迭代查询结果
	for cur.Next(context.Background()) {
		var answerSheet AnswerSheet
		if err := cur.Decode(&answerSheet); err != nil {
			return nil, nil, err
		}
		answerSheets = append(answerSheets, answerSheet)
	}
	if err := cur.Err(); err != nil {
		return nil, nil, err
	}

	// 返回分页数据和总记录数
	return answerSheets, &total, nil
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
