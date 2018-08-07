/*
 Schnorr NIZK package.  Based on https://tools.ietf.org/html/rfc8235
*/

package main

import(
	"math/big"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

type KeyPairs struct{
	g,p,q *big.Int
}

// for finite field, hashes returns c = H(g || V || A || etc)
// returns the data set inside a big Int
func SchnorrHash(g, V, A *big.Int, id string) *big.Int {
	var c *big.Int
	var buf []byte
	
	buf = append(buf, g.Bytes()...)
	buf = append(buf, V.Bytes()...)
	buf = append(buf, A.Bytes()...)
	buf = append(buf, []byte(id)...)

	c = big.NewInt(0);	
	digest := hashBuffer(buf);
	c.SetBytes(digest);
	return c;
}

// this is to make it easier to drop in a different hash function
func hashBuffer(buf []byte) []byte {
	sum := sha256.Sum256(buf);
	// returning a slice
	// this converts the array in to a slice to that array
	return sum[:];
}

// SignFF is a Schnorr signature over a Finite Field
// it proves knowledge of private key 'a' where A = g ^ a
// A is the public exponent
// g is a generator of G
func SignFF(p, q, g, A, a *big.Int) error {
	// V is g ^ v
	// c is a bigInt and c = H(g || V || A || etc) after setting bytes
	// r = v - a*c mod q
	var V, c, r *big.Int
	// v is random number in [0, q)
	v, err := rand.Int(rand.Reader, q) 
	if err != nil {
		return fmt.Errorf("in SignFF, %v\n", err);
	}

	// V = g^v mod p
	// modular exponentiation may leak information
	V = big.NewInt(0);
	V.Exp(g, v, p);

	// c is the challenge
	// it has been redefined in the non-interactive version
	// c = H(g || V || A || UserID || OtherInfo)
	// https://tools.ietf.org/html/rfc8235#section-2.3
	UserID := "test"
	c = SchnorrHash(g, V, A, UserID)

	// r = v-a*c mod q
	r = big.NewInt(0)
	// a*c mod q
	r.Mul(a, c)
	r.Mod(r, q)

	// v - (a*c mod q) mod q
	r.Sub(v, r)
	r.Mod(r, q)

	// must send A, g, V, c, r		
	// hopefully the public key A, g is sent already
	// verifier must verify that 1) A is a valid public key
	// and that 2) V = g^r * A^c
	return nil;
}

GenerateKeyPairs() {

}

func main() {

}
