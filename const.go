package lift_base_board

type Command = byte

const (
	CommandPing       Command = 0x01
	CommandGetVersion Command = 0x02

	CommandGetTemp Command = 0x03

	CommandSetWatchdog Command = 0x04

	CommandGetTime Command = 0x05
	CommandSetTime Command = 0x06

	CommandSetLed          Command = 0x07
	CommandReadButtons     Command = 0x08
	CommandReadBatLevel    Command = 0x09
	CommandReadVinLevel    Command = 0x10
	CommandReactDipsSwitch Command = 0x11

	CommandBufferCount  Command = 0x12
	CommandReadBuffer1  Command = 0x13
	CommandReadBuffer2  Command = 0x14
	CommandReadBuffer3  Command = 0x15
	CommandReadBuffer4  Command = 0x16
	CommandWriteBuffer1 Command = 0x17
	CommandWriteBuffer2 Command = 0x18
	CommandWriteBuffer3 Command = 0x19
	CommandWriteBuffer4 Command = 0x20
)
