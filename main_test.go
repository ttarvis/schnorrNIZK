package main

import(
	"testing"
	"math/big"
	"fmt"
)

func genKey() KeyPairs {
	kp := new(KeyPairs);
	kp.g = big.NewInt(0)
	kp.g.SetString(gHex, 16);
	
	kp.q = big.NewInt(0);
	kp.q.SetString(qHex, 16);

	kp.p = big.NewInt(0);
	kp.p.SetString(phex, 16);

	return kp;
}

func Testmain(t *testing.T) {
 
}

func TestSchnorrHash(t *testing.T) {
	// update this to take more inputs
	g := big.NewInt(0);
	V := big.NewInt(0);
	A := big.NewInt(0); 
	h := SchnorrHash(g, V, A, "test");
	fmt.Printf("%x", h);
}
