package main

import (
	"dynamixel-controller/pkg/motor"
	"fmt"
	"log"
)

func main() {
	// MotorController 초기화 (포트 이름과 보드레이트 설정)
	motor1, err := motor.NewMotorController("/dev/ttyUSB0", 1000000, 1) // 포트 이름, 보드레이트, 모터 ID
	if err != nil {
		log.Fatalf("Error initializing motor 1: %v", err)
	}
	defer motor1.Close()

	// Motor 1의 위치 설정
	err = motor1.SetPosition(512) // 목표 위치 설정
	if err != nil {
		log.Fatalf("Error setting position for motor 1: %v", err)
	}

	// Motor 1의 현재 위치 읽기
	position, err := motor1.GetPosition()
	if err != nil {
		log.Fatalf("Error reading position for motor 1: %v", err)
	}
	fmt.Printf("Motor 1 current position: %d\n", position)

	// Ping 명령어로 Motor 1의 연결 확인
	connected, err := motor1.Ping()
	if err != nil {
		log.Fatalf("Error pinging motor 1: %v", err)
	}
	if connected {
		fmt.Println("Motor 1 is connected successfully!")
	} else {
		log.Println("Motor 1 is not connected.")
	}

	// 여러 모터를 동시에 제어하는 예시 (SyncWrite)
	positions := map[byte]int{
		1: 512,  // Motor 1 목표 위치
		2: 1024, // Motor 2 목표 위치
	}
	err = motor1.SyncWrite(positions)
	if err != nil {
		log.Fatalf("Error syncing write: %v", err)
	}
	fmt.Println("Synchronized write completed for multiple motors.")
}
