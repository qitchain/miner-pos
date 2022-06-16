#ifndef CHIA_H
#define CHIA_H

#include <stddef.h>
#include <stdint.h>

typedef struct DiskProver DiskProver;
typedef struct Verifier Verifier;

#ifdef _WIN32
#define PORT  __declspec(dllexport)
#else
#define PORT
#endif
struct Qualities {
    int nLen;
    unsigned char** qualities;
};

#if defined(__cplusplus)
extern "C"  {
#endif

// DiskProver
PORT DiskProver* CreateDiskProver(const char* filename, char *msg);
PORT void GetMemo(DiskProver* p, char* pMemo);
PORT unsigned int GetMemoSize(DiskProver* p);
PORT void GetId(DiskProver* p, char* pId);
PORT const char * GetFilename(DiskProver* p);
PORT unsigned char GetSize(DiskProver* p);
PORT struct Qualities* GetQualitiesForChallenge(DiskProver* p, const char*challenge, int* success, int pf);
PORT void getQualities(struct Qualities* p, int nIndex, char* buf, int nLen);
PORT int getQualitiesCount(struct Qualities* p);
PORT int GetFullProof(DiskProver* p, const char*challenge, unsigned int index, char *pProof, int pf);

// Verifier
// PORT unsigned char* ValidateProof(const char* id, unsigned char k, const char* challenge, unsigned char* proof_bytes, unsigned short proof_size);
PORT int ValidateProof(const char* id, unsigned char k, const char* challenge, unsigned char* proof_bytes, unsigned short proof_size, char *quality_buf);

PORT void releaseQualities(struct Qualities* p);

PORT void setMaxCache(unsigned int max);
#if defined(__cplusplus)
}
#endif

#endif