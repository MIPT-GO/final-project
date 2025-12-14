package reservation

import "time"

type Reserve struct {
	Id     uint64
	Email  string
	Start  time.Time
	End    time.Time
	Hotel  string
	Number uint64
}
