package zabbix

type (
	AvailableType int
	StatusType    int
)

const (
	Available   AvailableType = 1
	Unavailable AvailableType = 2

	Monitored   StatusType = 0
	Unmonitored StatusType = 1
)
