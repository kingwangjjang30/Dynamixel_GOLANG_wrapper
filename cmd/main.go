package main

import (
	"DXL_GOLANG_wrapper/pkg/motor"
	"fmt"
	"log"
	"time"
)

func main() {
	// MotorController 초기화 (포트 이름과 보드레이트 설정)
	motor1, err := motor.NewMotorController("/dev/ttyUSB0", 1000000, 1)
	if err != nil {
		log.Fatalf("Error initializing motor 1: %v", err)
	}
	defer motor1.Close()

	motor2, err := motor.NewMotorController("/dev/ttyUSB0", 1000000, 2)
	if err != nil {
		log.Fatalf("Error initializing motor 2: %v", err)
	}
	defer motor2.Close()

	// 모터 1의 위치 설정
	err = motor1.SetPosition(512)
	if err != nil {
		log.Fatalf("Error setting position for motor 1: %v", err)
	}

	// 모터 2의 위치 설정
	err = motor2.SetPosition(1024)
	if err != nil {
		log.Fatalf("Error setting position for motor 2: %v", err)
	}

	// 모터 1의 현재 위치 출력
	position, err := motor1.GetPosition()
	if err != nil {
		log.Fatalf("Error reading position for motor 1: %v", err)
	}
	fmt.Printf("Motor 1 current position: %d\n", position)

	// 모터 2의 현재 위치 출력
	position, err = motor2.GetPosition()
	if err != nil {
		log.Fatalf("Error reading position for motor 2: %v", err)
	}
	fmt.Printf("Motor 2 current position: %d\n", position)

	// 두 모터를 동시에 제어 (SyncWrite)
	positions := map[byte]int{
		1: 512,
		2: 1024,
	}
	err = motor1.SyncWrite(positions)
	if err != nil {
		log.Fatalf("Error syncing write: %v", err)
	}
	fmt.Println("Synchronized write completed")

	// 잠시 대기하여 모터가 이동할 시간을 줌
	time.Sleep(2 * time.Second)
}
