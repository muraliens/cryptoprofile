package main

// #include <stdio.h>
// #include <memory.h>
// #include <stdlib.h>
/*

typedef unsigned char* BytePtr;
typedef struct {
  unsigned char key[16];  //128 bit key
  unsigned char iv[12];   //96 bit IV
  unsigned char ls[2000]; //Shift register
  unsigned int ctr;
} espresso_ctx;

int update_ls(espresso_ctx *ctx, unsigned char init) {
  unsigned char n255, n251, n247, n243, n239, n235, n231, n217, n213, n209, n205, n201, n197, n193, out;
  unsigned char *ls = ctx->ls;
  unsigned int *ctr = &(ctx->ctr);

  // Save new variables and output
  out  = ls[80+*ctr] ^ ls[99+*ctr] ^ ls[137+*ctr] ^ ls[227+*ctr] ^ ls[222+*ctr] ^ ls[187+*ctr] ^ \
         ls[243+*ctr]&ls[217+*ctr] ^ ls[247+*ctr]&ls[231+*ctr] ^ ls[213+*ctr]&ls[235+*ctr] ^ \
         ls[255+*ctr]&ls[251+*ctr] ^ ls[181+*ctr]&ls[239+*ctr] ^ ls[174+*ctr]&ls[44+*ctr]  ^  \
         ls[164+*ctr]&ls[29+*ctr]  ^ ls[255+*ctr]&ls[247+*ctr]&ls[243+*ctr]&ls[213+*ctr]&ls[181+*ctr]&ls[174+*ctr];
  n255 = ls[0+*ctr] ^ ls[41+*ctr]&ls[70+*ctr];
  n251 = ls[42+*ctr]&ls[83+*ctr]  ^ ls[8+*ctr];
  n247 = ls[44+*ctr]&ls[102+*ctr] ^ ls[40+*ctr];
  n243 = ls[43+*ctr]&ls[118+*ctr] ^ ls[103+*ctr];
  n239 = ls[46+*ctr]&ls[141+*ctr] ^ ls[117+*ctr];
  n235 = ls[67+*ctr]&ls[90+*ctr]&ls[110+*ctr]&ls[137+*ctr];
  n231 = ls[50+*ctr]&ls[159+*ctr] ^ ls[189+*ctr];
  n217 = ls[3+*ctr]&ls[32+*ctr];
  n213 = ls[4+*ctr]&ls[45+*ctr];
  n209 = ls[6+*ctr]&ls[64+*ctr];
  n205 = ls[5+*ctr]&ls[80+*ctr];
  n201 = ls[8+*ctr]&ls[103+*ctr];
  n197 = ls[29+*ctr]&ls[52+*ctr]&ls[72+*ctr]&ls[99+*ctr];
  n193 = ls[12+*ctr]&ls[121+*ctr];
  if (init) {
	  n255 ^= out;
	  n217 ^= out;
  }

  // Update state
  (ctx->ctr)++;
  for (int i = 254; i >=0; i--)
  {
    int cnt = *ctr-1;
    ls[i+*ctr] = ls[i+1+cnt];
  }
  ls[255+*ctr] = n255;
  ls[251+*ctr] ^= n251;
  ls[247+*ctr] ^= n247;
  ls[243+*ctr] ^= n243;
  ls[239+*ctr] ^= n239;
  ls[235+*ctr] ^= n235;
  ls[231+*ctr] ^= n231;
  ls[217+*ctr] ^= n217;
  ls[213+*ctr] ^= n213;
  ls[209+*ctr] ^= n209;
  ls[205+*ctr] ^= n205;
  ls[201+*ctr] ^= n201;
  ls[197+*ctr] ^= n197;
  ls[193+*ctr] ^= n193;

  if ((ctx->ctr) == 1700) {
	memcpy(ls, ls+1700, 256);
	(ctx->ctr) = 0;
  }

  return out;
}


int init_ls(espresso_ctx *ctx) {
  unsigned int i,j;
  unsigned char *ls = ctx->ls;

  // Load key and IV
  for (i=0;i<16;++i) for (j=0;j<8;++j) ls[8*i + j] = (((ctx->key[i])>>j)&1);
  for (i=0;i<12;++i) for (j=0;j<8;++j) ls[128 + 8*i + j] = (((ctx->iv[i])>>j)&1);
  for (i=0;i<31;++i) ls[128+96+i] = 1;
  ls[255] = 0;
  ctx->ctr=0;
  for (i=0;i<256;++i) update_ls(ctx,1);
  return 0;
}
unsigned char *espresso(unsigned char *key, unsigned char *iv, int numStream) {
	int i, j;
	espresso_ctx ctx;
	unsigned char *keystream;
	keystream = (unsigned char *)malloc(numStream);
	memset(keystream, 0, numStream);
	for (i=0;i<16;++i) (&ctx)->key[i] = (unsigned char) key[i];
  for (i=0;i<12;++i) (&ctx)->iv[i] = (unsigned char) iv[i];

	// Initiate cipher
  	init_ls(&ctx);

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

func (h *Handle) EspressoStream() cryptoprofile.BitStream {
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

	keystreams := C.espresso(C.BytePtr(key), C.BytePtr(iv), C.int(numberStream))

	streams := C.GoBytes(unsafe.Pointer(keystreams), C.int(numberStream))

	defer func() {
		C.free(unsafe.Pointer(key))
		C.free(unsafe.Pointer(iv))
		C.free(unsafe.Pointer(keystreams))
	}()

	return cryptoprofile.ParseBytes(h.streamLength, streams)
}
