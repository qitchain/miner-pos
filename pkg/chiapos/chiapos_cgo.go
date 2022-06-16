package chiapos

//#cgo CFLAGS: -I./ -DBLAKE3_NO_AVX512=1 -DBLAKE3_NO_SSE41=1 -DBLAKE3_NO_SSE2=1
//#cgo LDFLAGS: -lchiapos -lfse -luint128 -lstdc++
//#cgo linux,amd64 LDFLAGS:-L./c-bindings/libs/linux -lm -lpthread -static -static-libgcc -static-libstdc++
//#cgo linux,arm LDFLAGS:-L./c-bindings/libs/linux-arm -lm -lpthread -static -static-libgcc -static-libstdc++
//#cgo linux,arm64 LDFLAGS:-L./c-bindings/libs/linux-arm64 -lm -lpthread -static -static-libgcc -static-libstdc++
//#cgo darwin,amd64 LDFLAGS:-L./c-bindings/libs/drawin -lm -stdlib=libc++
//#cgo windows,amd64 LDFLAGS:-L./c-bindings/libs/windows -static -static-libgcc -static-libstdc++
/*
#include "./c-bindings/include/chiapos.h"

char *getQualitiesIndex(struct Qualities* p, int nIndex){
	return p->qualities[nIndex];
}
*/
import "C"
import (
	"errors"
	"unsafe"
)

func byteToString(data []byte) string {
	for i, b := range data {
		if b == 0 {
			return string(data[:i])
		}
	}
	return ""
}
func getFile(dp uintptr) *C.DiskProver {
	return (*C.DiskProver)(unsafe.Pointer(dp))
}

func CreateDiskProver(fileName string) (uintptr, error) {
	message := make([]byte, 1024)
	var dp uintptr = uintptr(unsafe.Pointer(C.CreateDiskProver(C.CString(fileName), (*C.char)(unsafe.Pointer(&message[0])))))
	if dp == 0 {
		return 0, errors.New(byteToString(message))
	} else {
		return dp, nil
	}
}

func GetQualitiesForChallenge(dp uintptr, challenge []byte, fp int) ([][]byte, int) {
	var success int = 0
	arrChallenge := make([][]byte, 0)
	p := C.GetQualitiesForChallenge(getFile(dp),
		(*C.char)(unsafe.Pointer(&challenge[0])),
		(*C.int)(unsafe.Pointer(&success)),
		C.int(fp))
	if p != nil {
		for i := 0; i < int(p.nLen); i++ {
			buf := make([]byte, 32, 32)
			ch := C.getQualitiesIndex(p, C.int(i))
			copy(buf, (*[32]byte)(unsafe.Pointer(ch))[:32])
			arrChallenge = append(arrChallenge, buf)
		}
		C.releaseQualities(p)
	}
	return arrChallenge, success
}

func GetId(dp uintptr) []byte {
	id := make([]byte, IdLen)
	C.GetId(getFile(dp), (*C.char)(unsafe.Pointer(&id[0])))
	return id
}

func GetMemo(dp uintptr) []byte {
	nSize := C.GetMemoSize(getFile(dp))
	memo := make([]byte, nSize)
	C.GetMemo(getFile(dp), (*C.char)(unsafe.Pointer(&memo[0])))
	return memo
}

func GetFullProof(dp uintptr, challenge []byte, index int, fp int) ([]byte, bool) {
	proof := make([]byte, ByteAlign(GetSize(dp)*64)/8)
	proofSize := C.GetFullProof(getFile(dp),
		(*C.char)(unsafe.Pointer(&challenge[0])),
		C.uint(index), (*C.char)(unsafe.Pointer(&proof[0])),
		C.int(fp))
	return proof, proofSize != 0
}

func GetSize(dp uintptr) uint32 {
	size := C.GetSize(getFile(dp))
	return uint32(size)
}

func ValidateProofStatic(id, challenge, proofBytes []byte, k uint32) ([]byte, bool) {
	quality := make([]byte, 32)
	ok := C.ValidateProof(
		(*C.char)(unsafe.Pointer(&id[0])),
		C.uchar(k),
		(*C.char)(unsafe.Pointer(&challenge[0])),
		(*C.uchar)(unsafe.Pointer(&proofBytes[0])),
		C.ushort(len(proofBytes)),
		(*C.char)(unsafe.Pointer(&quality[0])))
	return quality, bool(ok == 1)
}
func SetMaxCache(size uint32) {
	C.setMaxCache(C.uint(size))
}
