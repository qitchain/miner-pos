package miner

import (
	"bytes"
	"chia-miner/miner/entity"
	"chia-miner/pkg/config"
	"chia-miner/utils"
	"encoding/hex"
	"github.com/sirupsen/logrus"
	"time"
)

var singleMiner *Miner

func init() {
	singleMiner = &Miner{}
}

func GetMiner() *Miner {
	return singleMiner
}

type Miner struct {
	config         *config.Config
	spaces         []*Space
	miningInfo     *entity.MiningInfo
	scanIterations int64
	scanTime       int64
}

func (m *Miner) Start(config *config.Config) {
	m.config = config
	InitJsonRpc(config)
	for _, filepath := range m.config.Path {
		m.spaces = append(m.spaces, NewSpace(filepath, config))
	}

	utils.StartTime(m.onTimer, 1000)
}
func (m *Miner) onTimer() {
	miningInfo, err := GetJsonRpc().GetMiningInfo()
	if err != nil {
		logrus.Errorf("error getting mining info, please check server config %v", err)
		return
	}

	needScan := false
	if m.miningInfo == nil || !m.miningInfo.IsSame(miningInfo) {
		m.scanTime = time.Now().Unix()
		needScan = true
		if m.miningInfo != nil && !bytes.Equal(m.miningInfo.Challenge, miningInfo.Challenge) {
			m.scanIterations = 0
		} else {
			m.scanIterations++
		}
		m.miningInfo = miningInfo
		m.miningInfo.ScanIterations = m.scanIterations
	}

	if needScan {
		for _, space := range m.spaces {
			space.requestScan(m.miningInfo)
		}
		logrus.Infof("new block: height%v difficulty[%v] challenge[%v] scanIterations[%v] ",
			m.miningInfo.Height, m.miningInfo.Difficulty, hex.EncodeToString(m.miningInfo.Challenge),
			m.miningInfo.ScanIterations)
	}
}
