package miner

import (
	entity2 "chia-miner/miner/entity"
	chiapos2 "chia-miner/pkg/chiapos"
	"chia-miner/pkg/config"
	"chia-miner/utils"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"github.com/sirupsen/logrus"
	"sync/atomic"
)

const targetDeadline = int64(180)

type Space struct {
	filepath string
	queue    *utils.Queue
	files    []*chiapos2.File
	cfg      *config.Config
}

func NewSpace(filepath string, cfg *config.Config) *Space {
	space := &Space{
		filepath: filepath,
		cfg:      cfg,
	}
	space.queue = utils.NewQueue(1024, space.run)
	files := utils.GetFileList(filepath, ".plot")
	for _, fileInfo := range files {
		f, err := chiapos2.Open(fileInfo.FilePath)
		if err != nil {
			logrus.Errorf("Failed to load, error %v %v", fileInfo.FilePath, err)
			continue
		}
		logrus.Debugf("Load chia file %v", fileInfo.FilePath)
		space.files = append(space.files, f)
	}
	return space
}

func (s *Space) requestScan(miningInfo *entity2.MiningInfo) {
	s.queue.Push(miningInfo)
}

func (s *Space) run(v interface{}) {
	miningInfo := v.(*entity2.MiningInfo)
	utils.RunTimeout(func() {
		s.scan(miningInfo, miningInfo.ScanIterations)
	}, 150*1000)
}
func (s *Space) checkFilter(filterBits int, filterData []byte) bool {
	filterValue := binary.LittleEndian.Uint32(filterData)
	data := filterValue << (32 - filterBits)
	return data == 0
}
func (s *Space) scan(miningInfo *entity2.MiningInfo, scanIterations int64) {
	for _, f := range s.files {
		var b8 [8]byte
		binary.BigEndian.PutUint64(b8[:], uint64(scanIterations))
		challengeBytes := s.sha256s(miningInfo.Challenge, b8[:])

		chHash := s.sha256s(f.GetId(), challengeBytes)
		if !s.checkFilter(miningInfo.FilterBits, chHash) {
			continue
		}

		arrQualities, _ := f.GetQualitiesForChallenge(challengeBytes)

		for i, qualities := range arrQualities {
			requiredIters := chiapos2.CalculateIterationsQuality(qualities, int32(f.GetSize()), miningInfo.Difficulty, challengeBytes)
			inflate := 80 * 512 / (1 << miningInfo.FilterBits)
			SubDeadline := (requiredIters * uint64(inflate)) / 24433591728

			if int64(SubDeadline) < targetDeadline {
				proof, ok := f.GetFullProof(challengeBytes, i)
				if !ok {
					logrus.Error("Failed to read proof")
					continue
				}
				fPubKey, err := f.GetFarmerPublicKey()
				if err != nil {
					continue
				}
				privateKey, ok := s.cfg.FarmerKey[fPubKey]
				if !ok {
					logrus.Errorf("Chia farmer private key is not configured, farmer public key %v", fPubKey)
					continue
				}

				if requiredIters > atomic.LoadUint64(&miningInfo.BestQuality) {
					continue
				}
				atomic.StoreUint64(&miningInfo.BestQuality, requiredIters)

				submitProof := &entity2.SubmitProof{
					//Quality:         requiredIters,
					Height:           miningInfo.Height,
					ScanIterations:   miningInfo.ScanIterations,
					Challenge:        hex.EncodeToString(miningInfo.Challenge),
					QualityString:    hex.EncodeToString(qualities),
					PlotSize:         f.GetSize(),
					PlotId:           hex.EncodeToString(f.GetId()),
					PoolPublicKey:    hex.EncodeToString(f.GetPoolPublicKeyBinary()),
					FarmerPublicKey:  hex.EncodeToString(f.GetFarmerPublicKeyBinary()),
					FarmerPrivateKey: privateKey,
					SecurityKey:      hex.EncodeToString(f.GetSecurityKeyBinary()),
					ResponseNumber:   int32(i),
					ProofXs:          hex.EncodeToString(proof),
					RequiredIters:    requiredIters,
				}
				GetJsonRpc().Submit(submitProof)
			}
		}
	}
}
func (s *Space) sha256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}
func (s *Space) sha256s(params ...[]byte) []byte {
	hash := sha256.New()
	for _, data := range params {
		hash.Write(data)
	}
	return hash.Sum(nil)
}
