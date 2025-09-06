package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"math/big"
	"strconv"
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

	return makeKeyPairFromRsaKey(keyPair)
}

func makeKeyPairFromRsaKey(rsaKey *rsa.PrivateKey) (*KeyPair, error) {
	primes := rsaKey.Primes
	n := rsaKey.N

	lambdaN := lambda(primes)

	e, d, err := getEAndD(n, lambdaN)
	if err != nil {
		return nil, err
	}

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

func getEAndD(n *big.Int, lambdaN *big.Int) (*big.Int, *big.Int, error) {
	const maxAttempts = 100
	for i := 0; i < maxAttempts; i++ {
		e, err := randFromRange(big.NewInt(minExponent), n)
		if err != nil {
			return nil, nil, err
		}

		d := new(big.Int).ModInverse(e, lambdaN)
		if d != nil && d.Cmp(big.NewInt(minExponent)) > 0 {
			return e, d, nil
		}
	}

	return nil, nil, errors.New("failed to generate valid e and d after " + strconv.Itoa(maxAttempts) + " attempts")
}

func randFromRange(min, max *big.Int) (*big.Int, error) {
	e, err := rand.Int(rand.Reader, new(big.Int).Sub(max, min))
	if err != nil {
		return nil, err
	}

	e.Add(e, min)
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

func ValidatePublicKey(publicKey *PublicKey) error {
	return validateKey(publicKey.N, publicKey.E)
}

func ValidatePrivateKey(privateKey *PrivateKey) error {
	return validateKey(privateKey.N, privateKey.D)
}

func validateKey(n *big.Int, exponent *big.Int) error {
	nErr := validateN(n)
	exponentErr := validateExponent(exponent)
	
	result := ""
	if nErr != nil {
		result += nErr.Error()
	}
	if exponentErr != nil {
		if result != "" {
			result += "; "
		}
		result += exponentErr.Error()
	}

	if result != "" {
		return errors.New(result)
	}

	return nil
}

func validateExponent(exponent *big.Int) error {
	if exponent == nil {
		return errors.New("exponent is nil")
	}
	
	if exponent.Sign() <= 0 {
		return errors.New("exponent is not positive")
	}

	if exponent.Cmp(big.NewInt(minExponent)) < 0 {
		return fmt.Errorf("exponent is too small: %s (exponent) < %d (min exponent)", exponent.String(), minExponent)
	}

	return nil
}

func validateN(n *big.Int) error {
	if n == nil {
		return errors.New("'N' is nil")
	}

	if n.Sign() <= 0 {
		return errors.New("'N' is not positive")
	}

	if n.BitLen() < keySize {
		return fmt.Errorf("'N' is too small: %d bits (N) < %d bits (key size)", n.BitLen(), keySize)
	}

	return nil
}
