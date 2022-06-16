package chiapos

import "C"
import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
)

var (
	emptyQuality = make([]byte, 32)
)

const (
	IdLen = 32
)

func _Open(fileName string) (*File, error) {
	file := &File{
		fileName: fileName,
	}

	filePoint, err := CreateDiskProver(fileName)
	if filePoint == 0 {
		return nil, err
	}
	file.filePoint = filePoint
	// memo
	file.memo = file.GetMemo()
	file.fileName = fileName
	return file, nil
}

func Open(fileName string) (f *File, errRet error) {
	defer func() {
		if err := recover(); err != nil {
			f = nil
			errRet = errors.New(fmt.Sprintf("%v", err))
		}
	}()
	return _Open(fileName)
}

type ScannerInfo struct {
	ScannerTime     int64
	K               uint32
	PlotId          uint64
	PoolPublicKey   uint64
	FarmerPublicKey uint64
	Lucky           int32
}

type File struct {
	filePoint uintptr
	fileName  string
	memo      []byte
}

func (f *File) GetId() []byte {
	return GetId(f.filePoint)
}

func (f *File) GetMemo() []byte {
	return GetMemo(f.filePoint)
}

func (f *File) GetPoolPublicKey() (string, error) {
	if len(f.memo) == 128 {
		return hex.EncodeToString(f.memo[:48]), nil
	} else {
		return hex.EncodeToString(f.memo[:32]), nil
	}
}

func (f *File) GetPoolPublicKeyBinary() []byte {
	if len(f.memo) == 128 {
		return f.memo[:48]
	} else {
		return f.memo[:32]
	}
}

func (f *File) GetFarmerPublicKeyBinary() []byte {
	if len(f.memo) == 128 {
		return f.memo[48:96]
	} else {
		return f.memo[32:80]
	}
}

func (f *File) GetFarmerPublicKey() (string, error) {
	return hex.EncodeToString(f.GetFarmerPublicKeyBinary()), nil
}

func (f *File) GetSecurityKey() (string, error) {
	return hex.EncodeToString(f.GetSecurityKeyBinary()), nil
}

func (f *File) GetSecurityKeyBinary() []byte {
	if len(f.memo) == 128 {
		return f.memo[96:]
	} else {
		return f.memo[80:]
	}
}

func (f *File) GetFilename() string {
	return f.fileName
}

func (f *File) GetSize() uint32 {
	return GetSize(f.filePoint)
}

func (f *File) GetQualitiesForChallenge(challenge []byte) ([][]byte, int) {
	qualities, ok := GetQualitiesForChallenge(f.filePoint, challenge, f.getFileHandle())
	if len(qualities) > 0 {
		for i := 0; i < len(qualities); i++ {
			if bytes.Equal(qualities[i], emptyQuality) {
				continue
			}
		}
	}
	return qualities, ok
}

func (f *File) ByteAlign(numBits uint32) uint32 {
	return numBits + (8-((numBits)%8))%8
}

func (f *File) GetFullProof(challenge []byte, index int) ([]byte, bool) {
	return GetFullProof(f.filePoint, challenge, index, f.getFileHandle())
}

func (f *File) getFileHandle() int {
	return int(f.filePoint)
}
