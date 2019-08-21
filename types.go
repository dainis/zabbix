package zabbix

type (
	AvailableType int
	StatusType    int
	TimestampType uint64
)

const (
	Available   AvailableType = 1
	Unavailable AvailableType = 2

	Monitored   StatusType = 0
	Unmonitored StatusType = 1

	ActiveProxy  StatusType = 5
	PassiveProxy StatusType = 6
)

const (
	ZbxApiErrorParameters int = -32602
	ZbxApiErrorInternal   int = -32500
)
