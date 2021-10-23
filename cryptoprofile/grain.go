package main

// #include <stdio.h>
// #include <memory.h>
// #include <stdlib.h>
// #include <stdint.h>
/*
typedef uint8_t u8;
typedef uint32_t u32;
typedef unsigned char* BytePtr;

#define GRAIN_MAXKEYSIZE 128
#define GRAIN_KEYSIZE(i) (128 + (i))

#define GRAIN_MAXIVSIZE 96
#define GRAIN_IVSIZE(i) (96 + (i))

typedef struct {
    u32 LFSR[128];
    u32 NFSR[128];
    const u8* p_key;
    u32 keysize;
    u32 ivsize;

} GRAIN_ctx;

GRAIN_ctx _ctx;

#define INITCLOCKS 160
#define N(i) (ctx->NFSR[80-i])
#define L(i) (ctx->LFSR[80-i])
#define X0 (ctx->LFSR[3])
#define X1 (ctx->LFSR[25])
#define X2 (ctx->LFSR[46])
#define X3 (ctx->LFSR[64])
#define X4 (ctx->NFSR[63])

static const u8 NFTable[1024]= {0,1,0,0,1,0,1,1,1,0,1,1,0,1,0,0,1,0,1,1,0,1,0,0,0,1,0,0,0,1,0,1,
								1,0,1,1,0,1,0,0,0,1,0,0,1,0,1,1,1,0,1,1,0,1,0,0,0,1,0,0,0,1,0,1,
								1,0,1,1,0,1,0,0,0,1,0,0,1,0,1,1,0,1,0,0,1,0,1,1,1,0,1,1,1,0,1,0,
								0,1,0,0,1,0,1,1,1,0,1,1,0,1,0,0,0,1,0,0,1,0,1,1,1,0,1,1,1,0,1,0,
								1,0,1,1,0,1,0,0,0,1,0,0,1,0,1,1,0,1,0,0,1,0,1,1,1,0,1,1,1,0,1,0,
								0,1,0,0,1,0,1,1,1,0,1,1,0,1,0,0,0,1,0,0,1,0,1,1,1,0,1,1,1,0,1,0,
								0,1,0,0,1,0,1,1,1,0,1,1,0,1,0,0,1,0,1,1,0,1,0,0,0,1,0,0,0,1,0,1,
								1,0,1,1,0,1,0,0,0,1,0,0,1,0,1,1,1,0,1,1,0,1,0,0,0,1,0,0,1,0,1,0,
								1,0,1,1,0,1,0,0,0,1,0,0,1,0,1,1,0,1,0,0,1,0,1,1,1,0,1,1,1,0,1,0,
								0,1,0,0,1,0,1,1,1,0,1,1,0,1,0,0,0,1,0,0,1,0,1,1,1,0,1,1,1,0,1,0,
								0,1,0,0,1,0,1,1,1,0,1,1,0,1,0,0,1,0,1,1,0,1,0,0,0,1,0,0,0,1,0,1,
								1,0,1,1,0,1,0,0,0,1,0,0,1,0,1,1,1,0,1,1,0,1,0,0,0,1,0,0,0,1,0,1,
								0,1,0,0,1,0,1,1,1,0,1,1,0,1,0,0,1,0,1,1,0,1,0,0,0,1,0,0,0,1,0,1,
								1,0,1,1,0,1,0,0,0,1,0,0,1,0,1,1,0,1,0,0,1,0,1,1,1,0,1,1,1,0,1,0,
								0,1,0,0,1,0,1,1,1,0,1,1,0,1,0,0,1,0,1,1,0,1,0,0,0,1,0,0,0,1,0,1,
								1,0,1,1,0,1,0,0,0,1,0,0,1,0,1,1,0,1,0,0,1,0,1,1,1,0,1,1,0,1,0,1,
								0,1,0,0,1,0,1,1,1,0,1,1,0,1,0,0,1,0,1,1,0,1,0,0,0,1,0,0,0,1,0,1,
								1,0,1,1,0,1,0,0,0,1,0,0,1,0,1,1,1,0,1,1,0,1,0,0,0,1,0,0,0,1,0,1,
								1,0,1,1,0,1,0,0,0,0,0,1,1,1,1,0,0,1,0,0,1,0,1,1,1,1,1,0,1,1,1,1,
								0,1,0,0,1,0,1,1,1,1,1,0,0,0,0,1,0,1,0,0,1,0,1,1,1,1,1,0,1,1,1,1,
								1,0,1,1,0,1,0,0,0,1,0,0,1,0,1,1,0,1,0,0,1,0,1,1,1,0,1,1,1,0,1,0,
								0,1,0,0,1,0,1,1,1,0,1,1,0,1,0,0,0,1,0,0,1,0,1,1,1,0,1,1,1,0,1,0,
								0,1,0,0,1,0,1,1,1,1,1,0,0,0,0,1,1,0,1,1,0,1,0,0,0,0,0,1,0,0,0,0,
								1,0,1,1,0,1,0,0,0,0,0,1,1,1,1,0,1,0,1,1,0,1,0,0,0,0,0,1,1,1,1,1,
								0,1,0,0,1,0,0,0,1,0,1,1,0,1,1,1,1,0,1,1,0,1,1,1,0,1,0,0,0,1,1,0,
								1,0,1,1,0,1,1,1,0,1,0,0,1,0,0,0,1,0,1,1,0,1,1,1,0,1,0,0,0,1,1,0,
								1,0,1,1,0,1,1,1,0,0,0,1,1,1,0,1,0,1,0,0,1,0,0,0,1,1,1,0,1,1,0,0,
								0,1,0,0,1,0,0,0,1,1,1,0,0,0,1,0,0,1,0,0,1,0,0,0,1,1,1,0,1,1,0,0,
								1,0,1,1,0,1,1,1,0,1,0,0,1,0,0,0,0,1,0,0,1,0,0,0,1,0,1,1,1,0,0,1,
								0,1,0,0,1,0,0,0,1,0,1,1,0,1,1,1,1,0,1,1,0,1,1,1,0,1,0,0,0,1,1,0,
								1,0,1,1,0,1,1,1,0,0,0,1,1,1,0,1,0,1,0,0,1,0,0,0,1,1,1,0,1,1,0,0,
								1,0,1,1,0,1,1,1,0,0,0,1,1,1,0,1,0,1,0,0,1,0,0,0,1,1,1,0,0,0,1,1};



static const u8 boolTable[32] = {0,0,1,1,0,0,1,0,0,1,1,0,1,1,0,1,1,1,0,0,1,0,1,1,0,1,1,0,0,1,0,0};

u8 grain_keystream(GRAIN_ctx* ctx, int numRounds) {

    u8 i, NBit, LBit, outbit;


    outbit = ctx->NFSR[2]; //^ctx->NFSR[15]^ctx->NFSR[36]^ctx->NFSR[45]^ctx->NFSR[64]^ctx->NFSR[73]^ctx->NFSR[89]^ctx->LFSR[93]^(ctx->NFSR[12]&ctx->LFSR[8])^(ctx->LFSR[13]&ctx->LFSR[20])^(ctx->NFSR[95]&ctx->LFSR[42])^(ctx->LFSR[60]&ctx->LFSR[79])^(ctx->NFSR[12]&ctx->NFSR[95]&ctx->LFSR[95]);
    NBit = ctx->LFSR[0];  //^ctx->NFSR[0]^ctx->NFSR[26]^ctx->NFSR[56]^ctx->NFSR[91]^ctx->NFSR[96]^(ctx->NFSR[3]&ctx->NFSR[67])^(ctx->NFSR[11]&ctx->NFSR[13])^(ctx->NFSR[17]&ctx->NFSR[18])^(ctx->NFSR[27]&ctx->NFSR[59])^(ctx->NFSR[40]&ctx->NFSR[48])^(ctx->NFSR[61]&ctx->NFSR[65])^(ctx->NFSR[68]&ctx->NFSR[84]);
    LBit = ctx->LFSR[0]; //^ctx->LFSR[7]^ctx->LFSR[38]^ctx->LFSR[70]^ctx->LFSR[81]^ctx->LFSR[96];


    if (numRounds > 1) {
        outbit ^= ctx->NFSR[15]; //^ctx->NFSR[36]^ctx->NFSR[45]^ctx->NFSR[64]^ctx->NFSR[73]^ctx->NFSR[89]^ctx->LFSR[93]^(ctx->NFSR[12]&ctx->LFSR[8])^(ctx->LFSR[13]&ctx->LFSR[20])^(ctx->NFSR[95]&ctx->LFSR[42])^(ctx->LFSR[60]&ctx->LFSR[79])^(ctx->NFSR[12]&ctx->NFSR[95]&ctx->LFSR[95]);
        NBit ^= ctx->NFSR[0]; //^ctx->NFSR[26]^ctx->NFSR[56]^ctx->NFSR[91]^ctx->NFSR[96]^(ctx->NFSR[3]&ctx->NFSR[67])^(ctx->NFSR[11]&ctx->NFSR[13])^(ctx->NFSR[17]&ctx->NFSR[18])^(ctx->NFSR[27]&ctx->NFSR[59])^(ctx->NFSR[40]&ctx->NFSR[48])^(ctx->NFSR[61]&ctx->NFSR[65])^(ctx->NFSR[68]&ctx->NFSR[84]);
    }
    if (numRounds > 2) {
        outbit ^= ctx->NFSR[36]; //^ctx->NFSR[45]^ctx->NFSR[64]^ctx->NFSR[73]^ctx->NFSR[89]^ctx->LFSR[93]^(ctx->NFSR[12]&ctx->LFSR[8])^(ctx->LFSR[13]&ctx->LFSR[20])^(ctx->NFSR[95]&ctx->LFSR[42])^(ctx->LFSR[60]&ctx->LFSR[79])^(ctx->NFSR[12]&ctx->NFSR[95]&ctx->LFSR[95]);
        NBit ^= ctx->NFSR[26]; //^ctx->NFSR[56]^ctx->NFSR[91]^ctx->NFSR[96]^(ctx->NFSR[3]&ctx->NFSR[67])^(ctx->NFSR[11]&ctx->NFSR[13])^(ctx->NFSR[17]&ctx->NFSR[18])^(ctx->NFSR[27]&ctx->NFSR[59])^(ctx->NFSR[40]&ctx->NFSR[48])^(ctx->NFSR[61]&ctx->NFSR[65])^(ctx->NFSR[68]&ctx->NFSR[84]);
        LBit ^= ctx->LFSR[7]; //^ctx->LFSR[38]^ctx->LFSR[70]^ctx->LFSR[81]^ctx->LFSR[96];
    }
    if (numRounds > 3) {
        outbit ^= ctx->NFSR[45]; //^ctx->NFSR[64]^ctx->NFSR[73]^ctx->NFSR[89]^ctx->LFSR[93]^(ctx->NFSR[12]&ctx->LFSR[8])^(ctx->LFSR[13]&ctx->LFSR[20])^(ctx->NFSR[95]&ctx->LFSR[42])^(ctx->LFSR[60]&ctx->LFSR[79])^(ctx->NFSR[12]&ctx->NFSR[95]&ctx->LFSR[95]);
        NBit ^= ctx->NFSR[56]; //^ctx->NFSR[91]^ctx->NFSR[96]^(ctx->NFSR[3]&ctx->NFSR[67])^(ctx->NFSR[11]&ctx->NFSR[13])^(ctx->NFSR[17]&ctx->NFSR[18])^(ctx->NFSR[27]&ctx->NFSR[59])^(ctx->NFSR[40]&ctx->NFSR[48])^(ctx->NFSR[61]&ctx->NFSR[65])^(ctx->NFSR[68]&ctx->NFSR[84]);
    }
    if (numRounds > 4) {
        outbit ^= ctx->NFSR[64]; //^ctx->NFSR[73]^ctx->NFSR[89]^ctx->LFSR[93]^(ctx->NFSR[12]&ctx->LFSR[8])^(ctx->LFSR[13]&ctx->LFSR[20])^(ctx->NFSR[95]&ctx->LFSR[42])^(ctx->LFSR[60]&ctx->LFSR[79])^(ctx->NFSR[12]&ctx->NFSR[95]&ctx->LFSR[95]);
        NBit ^= ctx->NFSR[91];  //^ctx->NFSR[96]^(ctx->NFSR[3]&ctx->NFSR[67])^(ctx->NFSR[11]&ctx->NFSR[13])^(ctx->NFSR[17]&ctx->NFSR[18])^(ctx->NFSR[27]&ctx->NFSR[59])^(ctx->NFSR[40]&ctx->NFSR[48])^(ctx->NFSR[61]&ctx->NFSR[65])^(ctx->NFSR[68]&ctx->NFSR[84]);
        LBit ^= ctx->LFSR[38]; //^ctx->LFSR[70]^ctx->LFSR[81]^ctx->LFSR[96];
    }
    if (numRounds > 5) {
        outbit ^= ctx->NFSR[73]; //^ctx->NFSR[89]^ctx->LFSR[93]^(ctx->NFSR[12]&ctx->LFSR[8])^(ctx->LFSR[13]&ctx->LFSR[20])^(ctx->NFSR[95]&ctx->LFSR[42])^(ctx->LFSR[60]&ctx->LFSR[79])^(ctx->NFSR[12]&ctx->NFSR[95]&ctx->LFSR[95]);
        NBit ^= ctx->NFSR[96]; //^(ctx->NFSR[3]&ctx->NFSR[67])^(ctx->NFSR[11]&ctx->NFSR[13])^(ctx->NFSR[17]&ctx->NFSR[18])^(ctx->NFSR[27]&ctx->NFSR[59])^(ctx->NFSR[40]&ctx->NFSR[48])^(ctx->NFSR[61]&ctx->NFSR[65])^(ctx->NFSR[68]&ctx->NFSR[84]);
    }
    if (numRounds > 6) {
        outbit ^= ctx->NFSR[89]; //^ctx->LFSR[93]^(ctx->NFSR[12]&ctx->LFSR[8])^(ctx->LFSR[13]&ctx->LFSR[20])^(ctx->NFSR[95]&ctx->LFSR[42])^(ctx->LFSR[60]&ctx->LFSR[79])^(ctx->NFSR[12]&ctx->NFSR[95]&ctx->LFSR[95]);
        NBit ^= (ctx->NFSR[3] & ctx->NFSR[67]); //^(ctx->NFSR[11]&ctx->NFSR[13])^(ctx->NFSR[17]&ctx->NFSR[18])^(ctx->NFSR[27]&ctx->NFSR[59])^(ctx->NFSR[40]&ctx->NFSR[48])^(ctx->NFSR[61]&ctx->NFSR[65])^(ctx->NFSR[68]&ctx->NFSR[84]);
        LBit ^= ctx->LFSR[70];  //^ctx->LFSR[81]^ctx->LFSR[96];
    }
    if (numRounds > 7) {
        outbit ^= ctx->LFSR[93]; //^(ctx->NFSR[12]&ctx->LFSR[8])^(ctx->LFSR[13]&ctx->LFSR[20])^(ctx->NFSR[95]&ctx->LFSR[42])^(ctx->LFSR[60]&ctx->LFSR[79])^(ctx->NFSR[12]&ctx->NFSR[95]&ctx->LFSR[95]);
        NBit ^= (ctx->NFSR[11] & ctx->NFSR[13]); //^(ctx->NFSR[17]&ctx->NFSR[18])^(ctx->NFSR[27]&ctx->NFSR[59])^(ctx->NFSR[40]&ctx->NFSR[48])^(ctx->NFSR[61]&ctx->NFSR[65])^(ctx->NFSR[68]&ctx->NFSR[84]);
    }
    if (numRounds > 8) {
        outbit ^=(ctx->NFSR[12] & ctx->LFSR[8]); //^(ctx->LFSR[13]&ctx->LFSR[20])^(ctx->NFSR[95]&ctx->LFSR[42])^(ctx->LFSR[60]&ctx->LFSR[79])^(ctx->NFSR[12]&ctx->NFSR[95]&ctx->LFSR[95]);
        NBit ^= (ctx->NFSR[17] & ctx->NFSR[18]); //^(ctx->NFSR[27]&ctx->NFSR[59])^(ctx->NFSR[40]&ctx->NFSR[48])^(ctx->NFSR[61]&ctx->NFSR[65])^(ctx->NFSR[68]&ctx->NFSR[84]);
    }
    if (numRounds > 9) {
        outbit ^= (ctx->LFSR[13] & ctx->LFSR[20]); //^(ctx->NFSR[95]&ctx->LFSR[42])^(ctx->LFSR[60]&ctx->LFSR[79])^(ctx->NFSR[12]&ctx->NFSR[95]&ctx->LFSR[95]);
        NBit ^= (ctx->NFSR[27] & ctx->NFSR[59]); //^(ctx->NFSR[40]&ctx->NFSR[48])^(ctx->NFSR[61]&ctx->NFSR[65])^(ctx->NFSR[68]&ctx->NFSR[84]);
    }
    if (numRounds > 10) {
        outbit ^= (ctx->NFSR[95] & ctx->LFSR[42]); //^(ctx->LFSR[60]&ctx->LFSR[79])^(ctx->NFSR[12]&ctx->NFSR[95]&ctx->LFSR[95]);
        NBit ^= (ctx->NFSR[40] & ctx->NFSR[48]); //^(ctx->NFSR[61]&ctx->NFSR[65])^(ctx->NFSR[68]&ctx->NFSR[84]);
        LBit ^= ctx->LFSR[81];   //^ctx->LFSR[96];
    }
    if (numRounds > 11) {
        outbit ^= (ctx->LFSR[60] & ctx->LFSR[79]); //^(ctx->NFSR[12]&ctx->NFSR[95]&ctx->LFSR[95]);
        NBit ^= (ctx->NFSR[61] & ctx->NFSR[65]);   //^(ctx->NFSR[68]&ctx->NFSR[84]);
    }
    if (numRounds > 12) {
        outbit ^= (ctx->NFSR[12] & ctx->NFSR[95] & ctx->LFSR[95]);
        NBit ^= (ctx->NFSR[68] & ctx->NFSR[84]);
        LBit ^= ctx->LFSR[96];
    }

    for (i = 1; i < (ctx->keysize); ++i) {
        ctx->NFSR[i - 1] = ctx->NFSR[i];
        ctx->LFSR[i - 1] = ctx->LFSR[i];
    }
    ctx->NFSR[(ctx->keysize) - 1] = NBit;
    ctx->LFSR[(ctx->keysize) - 1] = LBit;
    return outbit;
}

void ECRYPT_keysetup(const u8* key,
                                   u32 keysize,
                                   u32 ivsize)
{
    GRAIN_ctx* ctx = &_ctx;
    ctx->p_key = key;
    ctx->keysize = keysize;
    ctx->ivsize = ivsize;
}

void ECRYPT_ivsetup(const u8* iv) {
    GRAIN_ctx* ctx = &_ctx;
    u32 i, j;
    u8 outbit;

    for (i = 0; i < (ctx->ivsize) / 8; ++i) {
        for (j = 0; j < 8; ++j) {
            ctx->NFSR[i * 8 + j] = ((ctx->p_key[i] >> j) & 1);
            ctx->LFSR[i * 8 + j] = ((iv[i] >> j) & 1);
        }
    }
    for (i = (ctx->ivsize) / 8; i < (ctx->keysize) / 8; ++i) {
        for (j = 0; j < 8; ++j) {
            ctx->NFSR[i * 8 + j] = ((ctx->p_key[i] >> j) & 1);
            ctx->LFSR[i * 8 + j] = 1;
        }
    }

    for (i = 0; i < 256; ++i) {
        outbit = grain_keystream(ctx, 13);
        ctx->LFSR[127] ^= outbit;
        ctx->NFSR[127] ^= outbit;
    }
}

void GRAIN_keystream_bytes(u8* keystream, u32 msglen) {
    u32 i, j;
    GRAIN_ctx* ctx = &_ctx;
	for (i = 0; i < msglen; ++i) {
        keystream[i] = 0;
        for (j = 0; j < 8; ++j) {
            keystream[i] |= (grain_keystream(ctx, 13) << j);
        }
    }
}

*/
import (
	"C"
)
import (
	"crypto/rand"
	"unsafe"

	"github.com/muraliens/cryptoprofile"
)

func (h *Handle) GrainStream() cryptoprofile.BitStream {
	if len(h.key) == 0 {
		h.key = make([]byte, 10)
		rand.Read(h.key)
	}
	if len(h.iv) == 0 {
		h.iv = make([]byte, 8)
		rand.Read(h.iv)
	}

	key := C.CBytes(h.key)
	iv := C.CBytes(h.iv)

	numberStream := h.streamLength / 8

	if numberStream*8 != h.streamLength {
		numberStream++
	}

	keystreams := (C.BytePtr)(C.malloc(C.size_t(numberStream)))
	C.ECRYPT_keysetup(C.BytePtr(key), 80, 64)
	C.ECRYPT_ivsetup(C.BytePtr(iv))

	C.GRAIN_keystream_bytes(C.BytePtr(keystreams), C.uint32_t(numberStream))

	streams := C.GoBytes(unsafe.Pointer(keystreams), C.int(numberStream))

	defer func() {
		C.free(unsafe.Pointer(key))
		C.free(unsafe.Pointer(iv))
		C.free(unsafe.Pointer(keystreams))
	}()

	return cryptoprofile.ParseBytes(h.streamLength, streams)
}
