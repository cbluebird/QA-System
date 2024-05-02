package mongodbService

import (
	"QA-System/config/database"
	"testing"
	"time"
)

// 单元测试
// TestSaveAnswerSheet 函数的单元测试
func TestSaveAnswerSheet(t *testing.T) {
	tests := []struct {
		name        string
		answerSheet AnswerSheet
		expectError error
	}{
		{
			name: "有效答卷",
			answerSheet: AnswerSheet{
				SurveyID: 1,
				Time:     time.Now().Format("2006-01-02 15:04:05"),
				Answers: []Answer{
					{QuestionID: 1, SerialNum: 1, Subject: "subject", Content: "content"},
					{QuestionID: 2, SerialNum: 2, Subject: "subject", Content: "content"},
				},
			},
			expectError: nil,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 初始化 MongoDB
			database.MongodbInit()

			// 保存答卷
			err := SaveAnswerSheet(tt.answerSheet)

			if (err != nil) != (tt.expectError != nil) {
				t.Errorf("Test case %q failed, expected error: %v, got: %v", tt.name, tt.expectError, err)
			}
		})
	}
}

// TestGetAnswerSheetBySurveyID 函数的单元测试
func TestGetAnswerSheetBySurveyID(t *testing.T) {
	tests := []struct {
		name        string
		surveyID    int
		pageNum     int
		pageSize    int
		expectError error 
	}{
		{
			name:        "有效的调查",
			surveyID:    1,
			pageNum:     50,
			pageSize:    100,
			expectError: nil,
		},
		// 在这里添加更多的测试用例
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 初始化 MongoDB
			database.MongodbInit()

			// 调用 GetAnswerSheetBySurveyID 函数获取答卷表
			_, _, err := GetAnswerSheetBySurveyID(tt.surveyID, tt.pageNum, tt.pageSize)

			if (err != nil) != (tt.expectError != nil){
				t.Errorf("测试用例 %q 失败，期望错误：%v，实际得到：%v", tt.name, tt.expectError, err)
			}
		})
	}
}

// TestDeleteAnswerSheetBySurveyID 函数的单元测试
func TestDeleteAnswerSheetBySurveyID(t *testing.T) {
	tests := []struct {
		name       string
		surveyID   int
		expectError  error
	}{
		{
			name:      "有效的调查ID",
			surveyID:  1,
			expectError: nil,
		},
		// 在这里添加更多的测试用例
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 初始化 MongoDB
			database.MongodbInit()

			// 删除指定 surveyID 的答卷表
			err := DeleteAnswerSheetBySurveyID(tt.surveyID)

			if(err != nil) != (tt.expectError != nil) {
				t.Errorf("测试用例 %q 失败，期望错误：%v，实际得到：%v", tt.name, tt.expectError, err)
			}
		})
	}
}

// 基准测试
// BenchmarkSaveAnswerSheet 函数的并发基准测试
func BenchmarkSaveAnswerSheet(b *testing.B) {
	// 初始化 MongoDB
	database.MongodbInit()

	// 并行度设置为 10
	b.SetParallelism(100)

	// 并发测试
	b.RunParallel(func(pb *testing.PB) {
		// 每个并发测试独立地运行 b.N 次
		for pb.Next() {
			// 创建 AnswerSheet 实例并保存
			answerSheet := AnswerSheet{
				SurveyID: 1,
				Time:     time.Now().Format("2006-01-02 15:04:05"),
				Answers: []Answer{
					{QuestionID: 1, SerialNum: 1, Subject: "subject", Content: "content"},
					{QuestionID: 2, SerialNum: 2, Subject: "subject", Content: "content"},
				},
			}
			if err := SaveAnswerSheet(answerSheet); err != nil {
				b.Errorf("SaveAnswerSheet() error = %v", err)
			}
		}
	})
}

// BenchmarkGetAnswerSheetBySurveyID 函数的并发基准测试
func BenchmarkGetAnswerSheetBySurveyID(b *testing.B) {
	// 初始化 MongoDB
	database.MongodbInit()

	// 并行度设置为 10
	b.SetParallelism(100)

	// 并发测试
	b.RunParallel(func(pb *testing.PB) {
		// 每个并发测试独立地运行 b.N 次
		for pb.Next() {
			// 调用 GetAnswerSheetBySurveyID 函数获取答卷表
			_, _, err := GetAnswerSheetBySurveyID(1, 50, 100)
			if err != nil {
				b.Errorf("GetAnswerSheetBySurveyID() error = %v", err)
			}
		}
	})
}

// BenchmarkDeleteAnswerSheetBySurveyID 函数的并发基准测试
func BenchmarkDeleteAnswerSheetBySurveyID(b *testing.B) {
	// 初始化 MongoDB
	database.MongodbInit()

	// 并行度设置为 10
	b.SetParallelism(10)

	// 并发测试
	b.RunParallel(func(pb *testing.PB) {
		// 每个并发测试独立地运行 b.N 次
		for pb.Next() {
			// 删除指定 surveyID 的答卷表
			surveyID := 1 // 修改为实际的 surveyID
			if err := DeleteAnswerSheetBySurveyID(surveyID); err != nil {
				b.Errorf("DeleteAnswerSheetBySurveyID() error = %v", err)
			}
		}
	})
}

// 示例测试
// ExampleSaveAnswerSheet 函数的示例测试
func ExampleSaveAnswerSheet() {
	// 初始化 MongoDB
	database.MongodbInit()

	// 创建 AnswerSheet 实例并保存
	answerSheet := AnswerSheet{
		SurveyID: 1,
		Time:     time.Now().Format("2006-01-02 15:04:05"),
		Answers: []Answer{
			{QuestionID: 1, SerialNum: 1, Subject: "subject", Content: "content"},
			{QuestionID: 2, SerialNum: 2, Subject: "subject", Content: "content"},
		},
	}
	if err := SaveAnswerSheet(answerSheet); err != nil {
	}
}

// ExampleGetAnswerSheetBySurveyID 函数的示例测试
func ExampleGetAnswerSheetBySurveyID() {
	// 初始化 MongoDB
	database.MongodbInit()

	// 调用 GetAnswerSheetBySurveyID 函数获取答卷表
	_, _, err := GetAnswerSheetBySurveyID(1, 50, 100)
	if err != nil {
	}
}

// ExampleDeleteAnswerSheetBySurveyID 函数的示例测试
func ExampleDeleteAnswerSheetBySurveyID() {
	// 初始化 MongoDB
	database.MongodbInit()

	// 删除指定 surveyID 的答卷表
	surveyID := 1 // 修改为实际的 surveyID
	if err := DeleteAnswerSheetBySurveyID(surveyID); err != nil {
	}
}
