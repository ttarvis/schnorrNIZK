package main

import(
	"testing"
	"math/big"
	"crypto/dsa"
	"crypto/rand"
	"fmt"
	"os"
)

func genKeys() *dsa.PrivateKey {
	params := new(dsa.Parameters);
	privKey := new(dsa.PrivateKey);

	// is this the right way to refer to parameterSizes?
	err := dsa.GenerateParameters(params, rand.Reader, dsa.L2048N256);
	if err != nil {
		fmt.Printf("%v, \n", err);
		os.Exit(1)
	}
	privKey.PublicKey.Parameters = *params;
	
	
	err = dsa.GenerateKey(privKey, rand.Reader)
	if err != nil {
		fmt.Printf("%v, \n", err);
		os.Exit(1);
	}
	return privKey;
}

func TestMain(t *testing.T) {

}

func TestSchnorrHash(t *testing.T) {
	// update this to take more inputs
	g := big.NewInt(0);
	V := big.NewInt(0);
	A := big.NewInt(0); 
	h := SchnorrHash(g, V, A, "test");
	fmt.Printf("%x\n", h);

}

func TestSignFF(t *testing.T) {
	privKey := genKeys();

	_, err := SignFF(privKey.P, privKey.Q, privKey.G, privKey.X, privKey.Y)
	if err != nil {
		t.Errorf("Failed! %v", err);
	}	

}

func TestSchnorrVerifyFF(t *testing.T) {
	var isVerified bool

	privKey := genKeys();

	sig := new(SchnorrSigFF);
	sig, err := SignFF(privKey.P, privKey.Q, privKey.G, privKey.Y, privKey.X)
	if err != nil {
		t.Errorf("Failure: %v\n", err);
	}

	isVerified = SchnorrVerifyFF(sig.V, privKey.G, sig.R, sig.A, privKey.P, privKey.Q);

	// this one should be legit
	if !isVerified {
		t.Errorf("Failure: %v\n", sig.V);
	}
}
