package motor

import (
	"errors"
	"time"

	"github.com/tarm/serial"
)

type Dynamixel struct {
	Port *serial.Port
}

// NewDynamixel은 시리얼 포트를 열고 새로운 Dynamixel 인스턴스를 초기화합니다.
func NewDynamixel(portName string, baudRate int) (*Dynamixel, error) {
	c := &serial.Config{Name: portName, Baud: baudRate, ReadTimeout: time.Millisecond * 100}
	p, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}
	return &Dynamixel{Port: p}, nil
}

// WritePosition은 다이나믹셀 모터의 목표 위치를 설정합니다.
func (dxl *Dynamixel) WritePosition(id byte, position int) error {
	address := byte(30)                  // 목표 위치를 설정하는 주소
	posL := byte(position & 0xFF)        // 하위 바이트
	posH := byte((position >> 8) & 0xFF) // 상위 바이트
	params := []byte{address, posL, posH}
	packet := makePacket(id, 0x03, params) // 명령: 0x03 (write)
	_, err := dxl.Port.Write(packet)
	return err
}

// Ping은 다이나믹셀 모터가 정상적으로 연결되었는지 확인하는 명령입니다.
func (dxl *Dynamixel) Ping(id byte) (bool, error) {
	packet := makePacket(id, 0x01, nil) // 명령: 0x01 (ping)
	_, err := dxl.Port.Write(packet)
	if err != nil {
		return false, err
	}
	resp := make([]byte, 6) // 응답을 받을 버퍼
	n, err := dxl.Port.Read(resp)
	if err != nil || n < 6 {
		return false, err
	}
	// 응답의 유효성 검사 (첫 번째 두 바이트는 0xFF, 0xFF여야 하며 ID와 오류 코드 확인)
	return resp[0] == 0xFF && resp[1] == 0xFF && resp[2] == id && resp[4] == 0x00, nil
}

// ReadData는 다이나믹셀 모터의 특정 데이터를 읽습니다.
func (dxl *Dynamixel) ReadData(id byte, address byte, length byte) ([]byte, error) {
	params := []byte{address, length}
	packet := makePacket(id, 0x02, params) // 명령: 0x02 (read)
	_, err := dxl.Port.Write(packet)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 6+length) // FF FF ID LEN ERR DATA1..N CHKSUM
	n, err := dxl.Port.Read(buf)
	if err != nil || n < int(6+length) {
		return nil, errors.New("read timeout or incomplete")
	}
	if !validateChecksum(buf) {
		return nil, errors.New("invalid checksum")
	}
	return buf[5 : 5+length], nil
}

// SyncWrite는 여러 모터에 동시에 데이터를 쓰는 명령입니다.
func (dxl *Dynamixel) SyncWrite(address byte, dataLength byte, data [][]byte) error {
	params := []byte{address, dataLength}
	for _, d := range data {
		params = append(params, d...)
	}
	packet := makePacket(0xFE, 0x83, params) // 0xFE: broadcast ID
	_, err := dxl.Port.Write(packet)
	return err
}

// makePacket은 다이나믹셀 모터와의 통신을 위한 패킷을 생성합니다.
func makePacket(id byte, instruction byte, params []byte) []byte {
	length := byte(len(params) + 2) // params 길이에 2를 더한 길이
	packet := []byte{0xFF, 0xFF, id, length, instruction}
	packet = append(packet, params...)
	checksum := calcChecksum(packet[2:]) // 체크섬 계산
	packet = append(packet, checksum)
	return packet
}

// calcChecksum은 패킷의 체크섬을 계산합니다.
func calcChecksum(data []byte) byte {
	sum := 0
	for _, b := range data {
		sum += int(b)
	}
	return byte(^sum & 0xFF)
}

// validateChecksum은 패킷의 체크섬을 검증합니다.
func validateChecksum(packet []byte) bool {
	data := packet[2 : len(packet)-1]
	chk := packet[len(packet)-1]
	return calcChecksum(data) == chk
}
