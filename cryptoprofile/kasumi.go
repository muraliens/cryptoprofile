package main

// #include <stdio.h>
// #include <memory.h>
// #include <stdlib.h>
// #include <stdint.h>
/*
typedef unsigned char* BytePtr;
typedef uint8_t u8;
typedef uint16_t  u16;

typedef uint32_t u32;


typedef union {
	u32 b32;
	u16 b16[2];
	u8  b8[4];
} REGISTER32;

typedef union {
	u16 b16;
	u8  b8[2];
} REGISTER16;

typedef union {
	u32 b32[2];
	u16 b16[4];
	u8  b8[8];
} REGISTER64;

#define ROL16(a,b) (u16)((a<<b)|(a>>(16-b)))

static u16 KLi1[8], KLi2[8];
static u16 KOi1[8], KOi2[8], KOi3[8];
static u16 KIi1[8], KIi2[8], KIi3[8];


static u16 FI( u16 in, u16 subkey )
{
	u16 nine, seven;
	static u16 S7[] = {
		54, 50, 62, 56, 22, 34, 94, 96, 38, 6, 63, 93, 2, 18,123, 33,
		55,113, 39,114, 21, 67, 65, 12, 47, 73, 46, 27, 25,111,124, 81,
		53, 9,121, 79, 52, 60, 58, 48,101,127, 40,120,104, 70, 71, 43,
		20,122, 72, 61, 23,109, 13,100, 77, 1, 16, 7, 82, 10,105, 98,
		117,116, 76, 11, 89,106, 0,125,118, 99, 86, 69, 30, 57,126, 87,
		112, 51, 17, 5, 95, 14, 90, 84, 91, 8, 35,103, 32, 97, 28, 66,
		102, 31, 26, 45, 75, 4, 85, 92, 37, 74, 80, 49, 68, 29,115, 44,
		64,107,108, 24,110, 83, 36, 78, 42, 19, 15, 41, 88,119, 59, 3};
	static u16 S9[] = {
		167,239,161,379,391,334,  9,338, 38,226, 48,358,452,385, 90,397,
		183,253,147,331,415,340, 51,362,306,500,262, 82,216,159,356,177,
		175,241,489, 37,206, 17,  0,333, 44,254,378, 58,143,220, 81,400,
		 95,  3,315,245, 54,235,218,405,472,264,172,494,371,290,399, 76,
		165,197,395,121,257,480,423,212,240, 28,462,176,406,507,288,223,
		501,407,249,265, 89,186,221,428,164, 74,440,196,458,421,350,163,
		232,158,134,354, 13,250,491,142,191, 69,193,425,152,227,366,135,
		344,300,276,242,437,320,113,278, 11,243, 87,317, 36, 93,496, 27,
		487,446,482, 41, 68,156,457,131,326,403,339, 20, 39,115,442,124,
		475,384,508, 53,112,170,479,151,126,169, 73,268,279,321,168,364,
		363,292, 46,499,393,327,324, 24,456,267,157,460,488,426,309,229,
		439,506,208,271,349,401,434,236, 16,209,359, 52, 56,120,199,277,
		465,416,252,287,246,  6, 83,305,420,345,153,502, 65, 61,244,282,
		173,222,418, 67,386,368,261,101,476,291,195,430, 49, 79,166,330,
		280,383,373,128,382,408,155,495,367,388,274,107,459,417, 62,454,
		132,225,203,316,234, 14,301, 91,503,286,424,211,347,307,140,374,
		 35,103,125,427, 19,214,453,146,498,314,444,230,256,329,198,285,
		 50,116, 78,410, 10,205,510,171,231, 45,139,467, 29, 86,505, 32,
		 72, 26,342,150,313,490,431,238,411,325,149,473, 40,119,174,355,
		185,233,389, 71,448,273,372, 55,110,178,322, 12,469,392,369,190,
		  1,109,375,137,181, 88, 75,308,260,484, 98,272,370,275,412,111,
		336,318,  4,504,492,259,304, 77,337,435, 21,357,303,332,483, 18,
		 47, 85, 25,497,474,289,100,269,296,478,270,106, 31,104,433, 84,
		414,486,394, 96, 99,154,511,148,413,361,409,255,162,215,302,201,
		266,351,343,144,441,365,108,298,251, 34,182,509,138,210,335,133,
		311,352,328,141,396,346,123,319,450,281,429,228,443,481, 92,404,
		485,422,248,297, 23,213,130,466, 22,217,283, 70,294,360,419,127,
		312,377,  7,468,194,  2,117,295,463,258,224,447,247,187, 80,398,
		284,353,105,390,299,471,470,184, 57,200,348, 63,204,188, 33,451,
		 97, 30,310,219, 94,160,129,493, 64,179,263,102,189,207,114,402,
		438,477,387,122,192, 42,381,  5,145,118,180,449,293,323,136,380,
		 43, 66, 60,455,341,445,202,432, 8,237, 15,376,436,464, 59,461};


	nine  = (u16)(in>>7);
	seven = (u16)(in&0x7F);


	nine  = (u16)(S9[nine]  ^ seven);
	seven = (u16)(S7[seven] ^ (nine & 0x7F));

	seven ^= (subkey>>9);
	nine  ^= (subkey&0x1FF);

	nine  = (u16)(S9[nine]  ^ seven);
	seven = (u16)(S7[seven] ^ (nine & 0x7F));

	in = (u16)((seven<<9) + nine);

	return( in );
}


static u32 FO( u32 in, int index )
{
	u16 left, right;


	left  = (u16)(in>>16);
	right = (u16) in;


	left ^= KOi1[index];
	left  = FI( left, KIi1[index] );
	left ^= right;

	right ^= KOi2[index];
	right  = FI( right, KIi2[index] );
	right ^= left;

	left ^= KOi3[index];
	left  = FI( left, KIi3[index] );
	left ^= right;

	in = (((u32)right)<<16)+left;

	return( in );
}

static u32 FL( u32 in, int index )
{
	u16 l, r, a, b;

	l = (u16)(in>>16);
	r = (u16)(in);


	a  = (u16) (l & KLi1[index]);
	r ^= ROL16(a,1);

	b  = (u16)(r | KLi2[index]);
	l ^= ROL16(b,1);


	in = (((u32)l)<<16) + r;

	return( in );
}


void kasumi( u8 *data )
{
	u32 left, right, temp;
	REGISTER32 *d;
	int n;

	d = (REGISTER32*)data;
	left  = (((u32)d[0].b8[0])<<24)+(((u32)d[0].b8[1])<<16)
            +(d[0].b8[2]<<8)+(d[0].b8[3]);
	right = (((u32)d[1].b8[0])<<24)+(((u32)d[1].b8[1])<<16)
            +(d[1].b8[2]<<8)+(d[1].b8[3]);
	n = 0;
	do {
	    temp = FL( left, n   );
		temp = FO( temp,  n++ );
		right ^= temp;
		temp = FO( right, n   );
		temp = FL( temp,   n++ );
		left ^= temp;
	} while( n<=7 );

	d[0].b8[0] = (u8)(left>>24);		d[1].b8[0] = (u8)(right>>24);
	d[0].b8[1] = (u8)(left>>16);		d[1].b8[1] = (u8)(right>>16);
	d[0].b8[2] = (u8)(left>>8);		    d[1].b8[2] = (u8)(right>>8);
	d[0].b8[3] = (u8)(left);			d[1].b8[3] = (u8)(right);

}

void kasumi_key_schedule( u8 *k )
{
	static u16 C[] = {
		0x0123,0x4567,0x89AB,0xCDEF, 0xFEDC,0xBA98,0x7654,0x3210 };
	u16 key[8], Kprime[8];
	REGISTER16 *k16;
	int n;

	k16 = (REGISTER16 *)k;
	for( n=0; n<8; ++n )
		key[n] = (u16)((k16[n].b8[0]<<8) + (k16[n].b8[1]));


	for( n=0; n<8; ++n )
		Kprime[n] = (u16)(key[n] ^ C[n]);


	for( n=0; n<8; ++n )
	{
		KLi1[n] = ROL16(key[n],1);
		KLi2[n] = Kprime[(n+2)&0x7];
		KOi1[n] = ROL16(key[(n+1)&0x7],5);
		KOi2[n] = ROL16(key[(n+5)&0x7],8);
		KOi3[n] = ROL16(key[(n+6)&0x7],13);
		KIi1[n] = Kprime[(n+4)&0x7];
		KIi2[n] = Kprime[(n+3)&0x7];
		KIi3[n] = Kprime[(n+7)&0x7];
	}
}

void kasumi_f8(u8 *key, u32 count, u32 bearer, u32 dir, u8 *data, int length)
{
	REGISTER64 A;
	REGISTER64 temp;
	int i, n;
	int lastbits = (8-(length%8)) % 8;
	u8  ModKey[16];
	u16 blkcnt;


	temp.b32[0]  = temp.b32[1]  = 0;
	A.b32[0]     = A.b32[1]     = 0;

	A.b8[0]  = (u8) (count>>24);
	A.b8[1]  = (u8) (count>>16);
	A.b8[2]  = (u8) (count>>8);
	A.b8[3]  = (u8) (count);
	A.b8[4]  = (u8) (bearer<<3);
	A.b8[4] |= (u8) (dir<<2);

	for( n=0; n<16; ++n )
		ModKey[n] = (u8)(key[n] ^ 0x55);
	kasumi_key_schedule( ModKey );

	kasumi( A.b8 );

	blkcnt = 0;
	kasumi_key_schedule( key );


	while( length > 0 )
	{

		temp.b32[0] ^= A.b32[0];
		temp.b32[1] ^= A.b32[1];
		temp.b8[7] ^= (u8)  blkcnt;
		temp.b8[6] ^= (u8) (blkcnt>>8);


		kasumi( temp.b8 );


		if( length >= 64 )
			n = 8;
		else
			n = (length+7)/8;


		for( i=0; i<n; ++i )
			*data++ ^= temp.b8[i];

		length -= 64;
		++blkcnt;
	}


#if 0
	if (lastbits)
		*data-- ;
#else
	if (lastbits)
		data-- ;
#endif
		*data &= 256 - (1<<lastbits) ;
}

u8 *kasumi_f9(u8 *key, u32 count, u32 fresh, u32 dir, u8 *data, int length)
{
	REGISTER64 A;
	REGISTER64 B;
	u8  FinalBit[8] = {0x80, 0x40, 0x20, 0x10, 8,4,2,1};
	u8  ModKey[16];
	static u8 mac_i[4];
	int i, n;


	kasumi_key_schedule( key );


	for( n=0; n<4; ++n )
	{
		A.b8[n]   = (u8)(count>>(24-(n*8)));
		A.b8[n+4] = (u8)(fresh>>(24-(n*8)));
	}
	kasumi( A.b8 );
	B.b32[0] = A.b32[0];
	B.b32[1] = A.b32[1];


	while( length >= 64 )
	{
		for( n=0; n<8; ++n )
			A.b8[n] ^= *data++;
		kasumi( A.b8 );
		length -= 64;
		B.b32[0] ^= A.b32[0];
		B.b32[1] ^= A.b32[1];
	}


	n = 0;
	while( length >=8 )
	{
		A.b8[n++] ^= *data++;
		length -= 8;
	}


	if( length )
	{
		i = *data;
		if( dir )
			i |= FinalBit[length];
	}
	else
		i = dir ? 0x80 : 0;

	A.b8[n++] ^= (u8)i;


	if( (length==7) && (n==8) )
	{
		kasumi( A.b8 );
		B.b32[0] ^= A.b32[0];
		B.b32[1] ^= A.b32[1];

		A.b8[0] ^= 0x80;
		i = 0x80;
		n = 1;
	}
	else
	{
		if( length == 7 )
			A.b8[n] ^= 0x80;
		else
			A.b8[n-1] ^= FinalBit[length+1];
	}

	kasumi( A.b8 );
	B.b32[0] ^= A.b32[0];
	B.b32[1] ^= A.b32[1];


	for( n=0; n<16; ++n )
		ModKey[n] = (u8)*key++ ^ 0xAA;
	kasumi_key_schedule( ModKey );
	kasumi( B.b8 );

	for( n=0; n<4; ++n )
		mac_i[n] = B.b8[n];

	return( mac_i );
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

func (h *Handle) KasumiStream() cryptoprofile.BitStream {
	if len(h.key) == 0 {
		h.key = make([]byte, 16)
		rand.Read(h.key)
	}
	if len(h.iv) == 0 {
		h.iv = make([]byte, 12)
		rand.Read(h.iv)
	}

	key := C.CBytes(h.key)

	numberStream := h.streamLength / 8

	if numberStream*8 != h.streamLength {
		numberStream++
	}

	keystreams := (C.BytePtr)(C.malloc(C.size_t(numberStream)))

	C.memset(unsafe.Pointer(keystreams), 0, C.size_t(numberStream))

	C.kasumi_f8(C.BytePtr(key), 0x12345678, 0x56273, 0, C.BytePtr(keystreams), C.int(numberStream))

	streams := C.GoBytes(unsafe.Pointer(keystreams), C.int(numberStream))

	defer func() {
		C.free(unsafe.Pointer(key))
		C.free(unsafe.Pointer(keystreams))
	}()

	return cryptoprofile.ParseBytes(h.streamLength, streams)
}
