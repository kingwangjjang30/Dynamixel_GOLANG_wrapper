package utils

import (
	"fmt"
	"log"
	"os"
)

// 로그 파일을 생성하거나 기존 파일에 추가하는 함수
func initLogger(logFile string) (*log.Logger, *os.File, error) {
	// 로그 파일 열기 (없으면 생성)
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open log file: %v", err)
	}

	// 로그 포맷 설정: 날짜, 시간, 로그 레벨, 메시지
	logger := log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	return logger, file, nil
}

// 로그 파일을 생성하고 메시지를 기록하는 함수
func LogMessage(logFile string, message string) error {
	logger, file, err := initLogger(logFile)
	if err != nil {
		return fmt.Errorf("error initializing logger: %v", err)
	}
	defer file.Close()

	// 로그 메시지 기록
	logger.Println(message)
	return nil
}

// 로그 레벨에 따라 다른 로그를 기록하는 함수
func LogError(logFile string, message string) error {
	logger, file, err := initLogger(logFile)
	if err != nil {
		return fmt.Errorf("error initializing logger: %v", err)
	}
	defer file.Close()

	// 에러 로그 기록
	logger.SetPrefix("ERROR: ")
	logger.Println(message)
	return nil
}

// Debug 수준의 로그를 기록하는 함수
func LogDebug(logFile string, message string) error {
	logger, file, err := initLogger(logFile)
	if err != nil {
		return fmt.Errorf("error initializing logger: %v", err)
	}
	defer file.Close()

	// 디버그 로그 기록
	logger.SetPrefix("DEBUG: ")
	logger.Println(message)
	return nil
}
