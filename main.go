/*
 Schnorr NIZK package.  Based on https://tools.ietf.org/html/rfc8235
*/

package main

import(
	"math/big"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
)

// it makes sense to just implement constant time comparison locally

// finite field signature
// names use capital for export
// but really they are V, A, c, r
// todo: check if there is any reason any of these need to be exported
type SchnorrSigFF struct {
	V, A, C, R *big.Int
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
func SignFF(p, q, g, A, a *big.Int) (*SchnorrSigFF, error) {
	// V is g ^ v
	// c is a bigInt and c = H(g || V || A || etc) after setting bytes
	// r = v - a*c mod q
	var V, c, r, ac *big.Int
	// v is random number in [0, q)
	v, err := rand.Int(rand.Reader, q) 
	if err != nil {
		return nil, fmt.Errorf("in SignFF, %v\n", err);
	}

	// V = g^v mod p
	// modular exponentiation may leak information
	V = big.NewInt(0);
	V.Exp(g, v, p);

	// c is the challenge
	// it has been redefined in the non-interactive version
	// c = H(g || V || A || UserID || OtherInfo)
	// https://tools.ietf.org/html/rfc8235#section-2.3
	// todo: change this out
	UserID := "test"
	c = SchnorrHash(g, V, A, UserID)

	// r = v-a*c mod q
	r = big.NewInt(0)

	// a*c mod q
	ac = big.NewInt(0)
	ac.Mul(a, c)
	ac.Mod(ac, q)

	// v - (a*c) mod q
	r.Sub(v, ac)
	r.Mod(r, q)

	// must send A, g, V, c, r		
	// hopefully the public key A, g is sent already
	// verifier must verify that 1) A is a valid public key
	// and that 2) V = g^r * A^c

	SchnorrSig := new(SchnorrSigFF);
	SchnorrSig.V = V // g^v we send
	SchnorrSig.A = A // A is public key
	SchnorrSig.C = c // c is hash challenge
	SchnorrSig.R = r // r = v-a*c mod q we send

	return SchnorrSig, nil;
}

// verifies a Schnorr Signature over a finite field
// need A, g, V, c, r, p
// check that V = g^r * A^c mod p
func SchnorrVerifyFF(V, g, r, A, p *big.Int) bool {
	var c, Ac *big.Int
	// todo verify that A is within [1, p) and
	// A^q = 1 mod p

	// A^c mod p
	UserID := "test"
	c = SchnorrHash(g, V, A, UserID);
	Ac = big.NewInt(0)
	Ac.Exp(A, c, p)


	// g^r mod p
	gr := big.NewInt(0)
	gr.Exp(g, r, p)

	// g^r * A^c mod p
	V1 := big.NewInt(0)
	V1.Mul(gr, Ac)
	V1.Mod(V1, p);
	// computed V1 = g^r * A^c mod p

	// replace with own methods
	isVerified := subtle.ConstantTimeCompare(V.Bytes(), V1.Bytes())

	if isVerified == 1 {
		return true
	} else {	
		return false
	}
}

func main() {

}
