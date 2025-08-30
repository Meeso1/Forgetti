package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"math/big"
)

const keySize int = 2048
const minExponent int64 = 256

type KeyPair struct {
	VerificationKey *PrivateKey
	BroadcastKey *PublicKey
}

func GenerateKeyPair() (*KeyPair, error) {
	keyPair, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, err
	}

	return MakeKeyPairFromRsaKey(keyPair)
}

func MakeKeyPairFromRsaKey(rsaKey *rsa.PrivateKey) (*KeyPair, error) {
	primes := rsaKey.Primes
	n := rsaKey.N

	lambdaN := lambda(primes)

	e, err := randomExponent(n)
	if err != nil {
		return nil, err
	}

	d := new(big.Int).ModInverse(e, lambdaN)

	return &KeyPair{
		VerificationKey: &PrivateKey{
			N: n,
			D: d,
		},
		BroadcastKey: &PublicKey{
			N: n,
			E: e,
		},
	}, nil
}

func randomExponent(n *big.Int) (*big.Int, error) {
	e, err := rand.Int(rand.Reader, new(big.Int).Sub(n, big.NewInt(minExponent)))
	if err != nil {
		return nil, err
	}

	e.Add(e, big.NewInt(minExponent))
	return e, nil
}

func lambda(primes []*big.Int) *big.Int {
	lambdas := make([]*big.Int, len(primes))
	for i := range primes {
		lambdas[i] = new(big.Int).Sub(primes[i], big.NewInt(1))
	}

	result := new(big.Int).Set(lambdas[0])
	for i := 1; i < len(lambdas); i++ {
		result = lcm(result, lambdas[i])
	}

	return result
}

func lcm(a, b *big.Int) *big.Int {
	return new(big.Int).Div(new(big.Int).Mul(a, b), new(big.Int).GCD(nil, nil, a, b))
}
