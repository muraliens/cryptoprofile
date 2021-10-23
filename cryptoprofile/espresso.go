package main

// #include <stdio.h>
// #include <memory.h>
// #include <stdlib.h>
/*

typedef unsigned char* BytePtr;

typedef struct {
  unsigned char key[16];  //128 bit key
  unsigned char iv[12];   //96 bit IV
  unsigned char ls[256]; //Shift register
  unsigned int ctr;
} espresso_ctx;

int update_ls(espresso_ctx *ctx, unsigned char init) {
  unsigned char  out;
  unsigned char *ls = ctx->ls;
  unsigned char n[256];

  memset(n, 0, 256);

  // Save new variables and output
  out  = ls[80] ^ ls[99] ^ ls[137] ^ ls[227] ^ ls[222] ^ ls[187] ^ \
         ls[243]&ls[217] ^ ls[247]&ls[231] ^ ls[213]&ls[235] ^ \
         ls[255]&ls[251] ^ ls[181]&ls[239] ^ ls[174]&ls[44]  ^  \
         ls[164]&ls[29]  ^ ls[255]&ls[247]&ls[243]&ls[213]&ls[181]&ls[174];
  n[255] = ls[0] ^ ls[41]&ls[70];
  n[251] = ls[42]&ls[83]  ^ ls[8];
  n[247] = ls[44]&ls[102] ^ ls[40];
  n[243] = ls[43]&ls[118] ^ ls[103];
  n[239] = ls[46]&ls[141] ^ ls[117];
  n[235] = ls[67]&ls[90]&ls[110]&ls[137];
  n[231] = ls[50]&ls[159] ^ ls[189];
  n[217] = ls[3]&ls[32];
  n[213] = ls[4]&ls[45];
  n[209] = ls[6]&ls[64];
  n[205] = ls[5]&ls[80];
  n[201] = ls[8]&ls[103];
  n[197] = ls[29]&ls[52]&ls[72]&ls[99];
  n[193] = ls[12]&ls[121];
  if (init) {
	  n[255] ^= out;
	  n[217] ^= out;
  }

  for (int i = 0; i < 255; i++)
  {
    ls[i] = ls[i+1] ^ n[i];
  }

  ls[255] = n[255];

  return out;
}


int init_ls(espresso_ctx *ctx, int numRounds) {
  unsigned int i,j;
  unsigned char *ls = ctx->ls;

  // Load key and IV
  for (i=0;i<16;i++) for (j=0;j<8;j++) ls[8*i + j] = (((ctx->key[i])>>j)&1);
  for (i=0;i<12;i++) for (j=0;j<8;j++) ls[128 + 8*i + j] = (((ctx->iv[i])>>j)&1);
  for (i=0;i<31;i++) ls[128+96+i] = 1;
  ls[255] = 0;
  ctx->ctr=0;
  for (i=0;i<numRounds;i++) update_ls(ctx,1);
  return 0;
}
unsigned char *espresso(int numRounds, unsigned char *key, unsigned char *iv, int numStream) {
	int i, j;
	espresso_ctx ctx;
	unsigned char *keystream;
	keystream = (unsigned char *)malloc(numStream);
	memset(keystream, 0, numStream);
	for (i=0;i<16;++i) (&ctx)->key[i] = (unsigned char) key[i];
  for (i=0;i<12;++i) (&ctx)->iv[i] = (unsigned char) iv[i];

	// Initiate cipher
  	init_ls(&ctx, numRounds);

	for (i=0;i<numStream;++i) for (j=0;j<8;++j) keystream[i] ^= (update_ls(&ctx,0)<<j);
	return keystream;
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

func (h *Handle) EspressoStream(numRounds int) cryptoprofile.BitStream {
	if len(h.key) == 0 {
		h.key = make([]byte, 16)
		rand.Read(h.key)
	}
	if len(h.iv) == 0 {
		h.iv = make([]byte, 12)
		rand.Read(h.iv)
	}

	key := C.CBytes(h.key)
	iv := C.CBytes(h.iv)

	numberStream := h.streamLength / 8

	if numberStream*8 != h.streamLength {
		numberStream++
	}

	if numRounds == 0 || numRounds > 256 {
		numRounds = 256
	}

	keystreams := C.espresso(C.int(numRounds), C.BytePtr(key), C.BytePtr(iv), C.int(numberStream))

	streams := C.GoBytes(unsafe.Pointer(keystreams), C.int(numberStream))

	defer func() {
		C.free(unsafe.Pointer(key))
		C.free(unsafe.Pointer(iv))
		C.free(unsafe.Pointer(keystreams))
	}()

	return cryptoprofile.ParseBytes(h.streamLength, streams)
}
