package lift_base_board

import (
	"encoding/binary"
	"errors"
	"math/rand"
	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"reflect"
	"time"
)

var (
	ErrRespondNotMatch = errors.New("respond not match")
	ErrSerialNotFound  = errors.New("serial not found")
)

type BaseBoard struct {
	dev *i2c.Dev
}

func (bb BaseBoard) Ping() error {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, rand.Uint32())

	err := bb.sendSetCommand(CommandPing, b)
	if err != nil {
		return err
	}

	return nil
}

func (bb BaseBoard) GetVersion() (uint16, error) {
	res, err := bb.sendCommand(CommandGetVersion, make([]byte, 2))
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(res), nil
}

func (bb BaseBoard) GetTemp() (float32, float32, error) {
	res, err := bb.sendCommand(CommandGetTemp, make([]byte, 4))
	if err != nil {
		return 0, 0, err
	}
	return float32(binary.BigEndian.Uint16(res[0:2])) / 10, float32(binary.BigEndian.Uint16(res[2:4])) / 10, nil
}

func (bb BaseBoard) SetWatchdog(sec uint16) error {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, sec)

	err := bb.sendSetCommand(CommandSetWatchdog, b)
	if err != nil {
		return err
	}

	return nil
}

func (bb BaseBoard) GetTime() (time.Time, error) {
	res, err := bb.sendCommand(CommandGetTime, make([]byte, 4))
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(0, int64(binary.BigEndian.Uint32(res))*int64(time.Millisecond)), nil
}

func (bb BaseBoard) SetTime(t time.Time) error {
	ms := uint32(t.UnixNano() / int64(time.Millisecond))
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, ms)

	err := bb.sendSetCommand(CommandSetTime, b)
	if err != nil {
		return err
	}

	return nil
}

func (bb BaseBoard) SetLed(b byte) error {
	err := bb.sendSetCommand(CommandSetLed, []byte{b})
	if err != nil {
		return err
	}

	return nil
}

// GetButtons TODO return button state
func (bb BaseBoard) GetButtons() error {
	err := bb.sendSetCommand(CommandReadButtons, make([]byte, 1))
	if err != nil {
		return err
	}

	return nil
}

func (bb BaseBoard) GetBatteryVoltage() (float32, error) {
	res, err := bb.sendCommand(CommandReadBatLevel, make([]byte, 2))
	if err != nil {
		return 0, err
	}

	return float32(binary.BigEndian.Uint16(res)) / 10, nil
}

func (bb BaseBoard) GetPowerSupplyVoltage() (float32, error) {
	res, err := bb.sendCommand(CommandReadVinLevel, make([]byte, 2))
	if err != nil {
		return 0, err
	}

	return float32(binary.BigEndian.Uint16(res)) / 10, nil
}

func (bb BaseBoard) GetDipSwitchState() ([]byte, error) {
	res, err := bb.sendCommand(CommandReactDipsSwitch, make([]byte, 3))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (bb BaseBoard) GetSerialBufferCount() ([]byte, error) {
	res, err := bb.sendCommand(CommandBufferCount, make([]byte, 4))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (bb BaseBoard) GetSerialBufferData(id int, count byte) ([]byte, error) {
	var cmd Command
	switch id {
	case 1:
		cmd = CommandReadBuffer1
	case 2:
		cmd = CommandReadBuffer2
	case 3:
		cmd = CommandReadBuffer3
	case 4:
		cmd = CommandReadBuffer4
	default:
		return nil, ErrSerialNotFound
	}

	res, err := bb.sendCommand(cmd, make([]byte, count))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (bb BaseBoard) WriteSerialBufferData(id int, count byte) ([]byte, error) {
	var cmd Command
	switch id {
	case 1:
		cmd = CommandWriteBuffer1
	case 2:
		cmd = CommandWriteBuffer2
	case 3:
		cmd = CommandWriteBuffer3
	case 4:
		cmd = CommandWriteBuffer4
	default:
		return nil, ErrSerialNotFound
	}

	res, err := bb.sendCommand(cmd, make([]byte, count))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (bb BaseBoard) sendSetCommand(cmd Command, data []byte) error {
	res, err := bb.sendCommand(cmd, data)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(data, res) {
		return ErrRespondNotMatch
	}

	return nil
}

func (bb BaseBoard) sendCommand(cmd Command, data []byte) ([]byte, error) {
	res := make([]byte, len(data)+1)

	err := bb.dev.Tx(append([]byte{cmd}, data...), res)
	if err != nil {
		return nil, err
	}

	return res[1:], nil
}

func NewBaseBoard(dev string, addr uint16) (*BaseBoard, error) {
	if _, err := driverreg.Init(); err != nil {
		return nil, err
	}

	b, err := i2creg.Open(dev)
	if err != nil {
		return nil, err
	}

	d := &i2c.Dev{Addr: addr, Bus: b}

	return &BaseBoard{
		dev: d,
	}, nil
}
