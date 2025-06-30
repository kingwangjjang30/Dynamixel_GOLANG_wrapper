package motor

import (
	"fmt"
	"protocol" // protocol 패키지에서 Dynamixel을 임포트
)

type MotorController struct {
	dxl *protocol.Dynamixel // Dynamixel 구조체를 포함하여 시리얼 포트 관리
	id  byte                // 모터 ID
}

// NewMotorController는 새로운 MotorController를 초기화합니다.
func NewMotorController(portName string, baudRate int, motorID byte) (*MotorController, error) {
	// Dynamixel 인스턴스를 생성하여 연결
	dxl, err := protocol.NewDynamixel(portName, baudRate)
	if err != nil {
		return nil, err
	}
	return &MotorController{dxl: dxl, id: motorID}, nil
}

// SetPosition은 모터의 목표 위치를 설정합니다.
func (mc *MotorController) SetPosition(position int) error {
	// 목표 위치를 설정하기 위해 WritePosition 호출
	err := mc.dxl.WritePosition(mc.id, position)
	if err != nil {
		return fmt.Errorf("failed to set position for motor %d: %v", mc.id, err)
	}
	return nil
}

// GetPosition은 모터의 현재 위치를 반환합니다.
func (mc *MotorController) GetPosition() (int, error) {
	// 현재 위치를 읽어오기 위해 ReadData 호출
	positionData, err := mc.dxl.ReadData(mc.id, 0x24, 2) // 0x24 주소에서 2 바이트 길이로 데이터 읽기 (현재 위치)
	if err != nil {
		return 0, fmt.Errorf("failed to get position for motor %d: %v", mc.id, err)
	}
	// 2바이트로 된 위치 값을 합쳐서 반환
	return int(positionData[0]) | int(positionData[1])<<8, nil
}

// Ping은 모터가 정상적으로 연결되었는지 확인합니다.
func (mc *MotorController) Ping() (bool, error) {
	// Ping 명령어로 모터와 연결된 상태를 확인
	connected, err := mc.dxl.Ping(mc.id)
	if err != nil {
		return false, fmt.Errorf("ping failed for motor %d: %v", mc.id, err)
	}
	return connected, nil
}

// SyncWrite는 여러 모터의 목표 위치를 동시에 설정합니다.
func (mc *MotorController) SyncWrite(positions map[byte]int) error {
	var data [][]byte
	// 여러 모터에 대해 목표 위치 데이터 생성
	for id, pos := range positions {
		posL := byte(pos & 0xFF)
		posH := byte((pos >> 8) & 0xFF)
		data = append(data, []byte{id, posL, posH})
	}

	// SyncWrite로 모든 모터에 목표 위치 전송
	err := mc.dxl.SyncWrite(0x30, byte(len(data[0])), data)
	if err != nil {
		return fmt.Errorf("failed to sync write: %v", err)
	}
	return nil
}

// Close는 시리얼 포트를 닫습니다.
func (mc *MotorController) Close() {
	mc.dxl.Port.Close()
}
