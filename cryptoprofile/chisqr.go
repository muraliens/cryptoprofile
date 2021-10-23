package main

/*
#include <stdio.h>
#include <math.h>
long double KM(long double S, long double Z)
{
	long double Sum = 1.0;
	long double Nom = 1.0;
	long double Denom = 1.0;

	for(int I = 0; I < 1000; I++) // Loops for 1000 iterations
	{
		Nom *= Z;
		S++;
		Denom *= S;
		Sum += (Nom / Denom);
	}

    return Sum;
}

long double log_igf(long double S, long double Z)
{
	if(Z < 0.0)
	{
		return 0.0;
	}
	long double Sc, K;
	Sc = (logl(Z) * S) - Z - logl(S);

    K = KM(S, Z);

    return logl(K) + Sc;
}


double igf(double S, double Z)
{
	if(Z < 0.0)
	{
		return 0.0;
	}
	long double Sc = (1.0 / S);
	Sc *= powl(Z, S);
	Sc *= expl(-Z);

	long double Sum = 1.0;
	long double Nom = 1.0;
	long double Denom = 1.0;

	for(int I = 0; I < 200; I++) // 200
	{
		Nom *= Z;
		S++;
		Denom *= S;
		Sum += (Nom / Denom);
	}

	return Sum * Sc;
}

#define A 15 // 15

double gamma(double N)
{

	//const long double SQRT2PI = sqrtl(atanl(1.0) * 8.0);
    const long double SQRT2PI = 2.5066282746310005024157652848110452530069867406099383;

    long double Z = (long double)N;
    long double Sc = powl((Z + A), (Z + 0.5));
	Sc *= expl(-1.0 * (Z + A));
    Sc /= Z;

	long double F = 1.0;
	long double Ck;
    long double Sum = SQRT2PI;


	for(int K = 1; K < A; K++)
	{
	    Z++;
		Ck = powl(A - K, K - 0.5);
		Ck *= expl(A - K);
		Ck /= F;

		Sum += (Ck / Z);

		F *= (-1.0 * K);
	}

	return (double)(Sum * Sc);
}

long double log_gamma(double N)
{


	//const long double SQRT2PI = sqrtl(atanl(1.0) * 8.0);
    const long double SQRT2PI = 2.5066282746310005024157652848110452530069867406099383;

    long double Z = (long double)N;
    long double Sc;

    Sc = (logl(Z + A) * (Z + 0.5)) - (Z + A) - logl(Z);

	long double F = 1.0;
	long double Ck;
    long double Sum = SQRT2PI;


	for(int K = 1; K < A; K++)
	{
	    Z++;
		Ck = powl(A - K, K - 0.5);
		Ck *= expl(A - K);
		Ck /= F;

		Sum += (Ck / Z);

		F *= (-1.0 * K);
	}

	return logl(Sum) + Sc;
}

double approx_gamma(double Z)
{
    const double RECIP_E = 0.36787944117144232159552377016147;  // RECIP_E = (E^-1) = (1.0 / E)
    const double TWOPI = 6.283185307179586476925286766559;  // TWOPI = 2.0 * PI

    double D = 1.0 / (10.0 * Z);
    D = 1.0 / ((12 * Z) - D);
    D = (D + Z) * RECIP_E;
    D = pow(D, Z);
    D *= sqrt(TWOPI / Z);

    return D;
}

long double approx_log_gamma(double N)
{
    const double LOGPIHALF = 0.24857493634706692717563414414545; // LOGPIHALF = (log10(PI) / 2.0)

    double D;

    D = 1.0 + (2.0 * N);
    D *= 4.0 * N;
    D += 1.0;
    D *= N;
    D = log10(D) * (1.0 / 6.0);
    D += N + (LOGPIHALF);
    D = (N * log(N)) - D;
    return D;

}


double chisqr(int Dof, double Cv)
{
    // printf("Dof:  %i\n", Dof);
    // printf("Cv:  %f\n", Cv);
    if(Cv < 0 || Dof < 1)
    {
        return 0.0;
    }
	double K = ((double)Dof) * 0.5;
	double X = Cv * 0.5;
	if(Dof == 2)
	{
		return exp(-1.0 * X);
	}
	long double PValue, Gam;
    long double ln_PV;
    ln_PV = log_igf(K, X);

    //Gam = approx_gamma(K);
    //Gam = lgammal(K);
    Gam = log_gamma(K);

    ln_PV -= Gam;
    PValue = 1.0 - expl(ln_PV);

	return (double)PValue;

}

*/
import (
	"C"
)
import (
	"fmt"
	"math"
	"os"

	"github.com/muraliens/cryptoprofile"
)

func ChiSqr(dv int, cv float64) float64 {
	p := C.chisqr(C.int(dv), C.double(cv))
	return float64(p)
}

func ChiSquareTest(filename string, crypto string, key []byte, iv []byte, rs cryptoprofile.BitStream, evps *cryptoprofile.EigenProfiles, evpsr *cryptoprofile.EigenProfiles) float64 {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Failed to create file")
		return 0
	}
	defer f.Close()
	E := make([]float64, len(evps.Profiles))
	O := make([]float64, len(evps.Profiles))
	C := make([]float64, len(evps.Profiles))
	CHS := float64(0)
	total := 0
	for i := 0; i < len(evps.Profiles); i++ {
		total = total + evps.Profiles[i].Count
	}
	totalr := 0
	for i := 0; i < len(evpsr.Profiles); i++ {
		totalr = totalr + evpsr.Profiles[i].Count
	}
	f.WriteString("-----------------------------------------------------------------------------\n")
	f.WriteString("   Bins     Expected Frequency (E)     Observed Frequency (O)     (O-E)^2/E\n")
	f.WriteString("-----------------------------------------------------------------------------\n")
	for i := 0; i < len(evps.Profiles); i++ {
		found := false
		index := 0
		for j := range evpsr.Profiles {
			if cryptoprofile.IsEigenProfileMatch(evpsr.Profiles[j].Profile, evps.Profiles[i].Profile) {
				found = true
				index = j
				break
			}
		}
		if found {
			E[i] = (float64(evps.Profiles[i].Count) / float64(total)) * float64(totalr)
			O[i] = float64(evpsr.Profiles[index].Count)
			C[i] = math.Pow((O[i]-E[i]), 2) / E[i]
			CHS = CHS + C[i]
		}
		str := fmt.Sprintf("%6d %20.06f    %20.06f    %20.06f\n", i+1, E[i], O[i], C[i])
		f.WriteString(str)
	}
	f.WriteString("-----------------------------------------------------------------------------\n")
	str := fmt.Sprintf("                                                Ch^2 Sum    %15.06f\n", CHS)
	f.WriteString(str)
	f.WriteString("-----------------------------------------------------------------------------\n")
	pvalue := ChiSqr(len(evps.Profiles)-1, CHS)
	str = fmt.Sprintf("\np-Value =   %15.06f\n", pvalue)
	f.WriteString(str)
	f.Close()
	fmt.Printf("P-Value =   %15.06f\n", pvalue)
	return pvalue
}
