package adminService

import (
	"bufio"
	"encoding/json"
	"os"
)

func GetLastLinesFromLogFile(filePath string, numLines int, logType int) ([]map[string]interface{}, error) {
	// 打开日志文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

    levelMap := map[int]string{
        0: "",
        1: "error",
        2: "warn",
        3: "info",
        4: "debug",
    }
    level := levelMap[logType]

	// 用于存储解析后的日志内容
	var logs []map[string]interface{}

	// 逐行读取文件内容
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// 解析JSON字符串为map类型
		var logData map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &logData); err != nil {
			// 如果解析失败，跳过这行日志继续处理下一行
			continue
		}

		// 根据logType筛选日志
		if level != "" {
			if logData["level"] == level {
				logs = append(logs, logData)
			}
		} else {
			logs = append(logs, logData)
		}

	}

	// 检查是否发生了读取错误
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// 如果文件中的行数不足以满足需求，直接返回所有行
	if len(logs) <= numLines {
		return logs, nil
	}

	// 如果文件中的行数超过需求，提取最后几行并返回
	startIndex := len(logs) - numLines
	return logs[startIndex:], nil
}
