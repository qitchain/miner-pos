package entity

import (
	"bytes"
)

type MiningInfo struct {
	Height         uint32
	Challenge      []byte
	Difficulty     uint64
	Epoch          int64
	ScanIterations int64
	ReceiveTime    int64
	FilterBits     int
	ServerTime     int64
	BestQuality    uint64
}

func (m *MiningInfo) IsSame(miningInfo *MiningInfo) bool {
	return bytes.Equal(m.Challenge, miningInfo.Challenge) && m.ScanIterations == miningInfo.ScanIterations
}
