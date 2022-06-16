// BLS C interface for c-binding
//
#ifndef BLS_CBINDINGS_H_
#define BLS_CBINDINGS_H_

#include <stddef.h>
#include <stdint.h>

// types
struct BNWrapperT;
typedef struct BNWrapperT BNWrapperT;

// private key
struct PrivateKeyT;
typedef struct PrivateKeyT PrivateKeyT;

// public key
struct G1ElementT;
typedef struct G1ElementT G1ElementT;

// signature
struct G2ElementT;
typedef struct G2ElementT G2ElementT;

// PrivateKeyT

// gen key from seed
PrivateKeyT *key_gen(const uint8_t *seed, size_t seed_len);

PrivateKeyT *key_from_bytes(const uint8_t *data, size_t data_len);

PrivateKeyT *key_aggregate(const PrivateKeyT **keyTs, size_t keyT_count);

PrivateKeyT *key_derive_child(const PrivateKeyT *keyT, uint32_t index);

void key_destroy(PrivateKeyT *keyT);

size_t key_bytes(const PrivateKeyT *keyT, uint8_t *buff, size_t buf_len);

G1ElementT *key_g1(const PrivateKeyT *keyT);

G2ElementT *key_g2(const PrivateKeyT *keyT);

G2ElementT *key_sign(const PrivateKeyT *keyT, const uint8_t *message, size_t message_len, const G1ElementT *prepend_pkT);

// G1ElementT
G1ElementT *g1_from_bytes(const uint8_t *data, size_t data_len);

G1ElementT *g1_add2(const G1ElementT *g1T1, const G1ElementT *g1T2);

void g1_destroy(G1ElementT *g1T);

size_t g1_bytes(const G1ElementT *g1T, uint8_t *buff, size_t buf_len);

int32_t g1_verify(const G1ElementT *pkT, const uint8_t *message, size_t message_len, const G2ElementT *signatureT);

// G2ElementT
G2ElementT *g2_from_bytes(const uint8_t *data, size_t data_len);

void g2_destroy(G2ElementT *g2T);

size_t g2_bytes(const G2ElementT *g2T, uint8_t *buff, size_t buf_len);

#endif // BLS_CBINDINGS_H_