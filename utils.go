package main

import (
	//"log"
	"bytes"
	"crypto/rand"
	xrand "math/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"
)

var prime *big.Int

/**
 * Returns a random number from the range (0, prime-1) inclusive
**/
func random() (*big.Int, error) {
	max := big.NewInt(0).Set(prime) // max = prime
	max = max.Sub(max, big.NewInt(1)) // max = max -1
	return rand.Int(rand.Reader, max)
}

/**
 * Converts an array of big.Ints to the original byte array, removing any
 * least significant nulls
**/
func mergeIntToByte(secret []*big.Int) []byte {
	var final []byte
	for _, elem := range secret {
		byteString := elem.Bytes()
		final = append(final, byteString...)
	}

	final = bytes.TrimRight(final, "\x00")

	return final
}

/**
 * inNumbers(array, value) returns boolean whether or not value is in array
**/
func inNumbers(numbers []*big.Int, value *big.Int) bool {
	for n := range numbers {
		if numbers[n].Cmp(value) == 0 {
			return true
		}
	}

	return false
}

/**
 * Returns the big.Int number base10 in base64 representation; note: this is
 * not a string representation; the base64 output is exactly 256 bits long
**/
func toBase64(number *big.Int) string {
	byteString := number.Bytes()
	return base64.URLEncoding.EncodeToString(byteString)
}

/**
 * Returns the number base64 in base 10 big.Int representation; note: this is
 * not coming from a string representation; the base64 input is exactly 256
 * bits long, and the output is an arbitrary size base 10 integer.
 *
 * Returns -1 on failure
**/
func fromBase64(number string) (*big.Int, error) {
	z := new(big.Int)

	bytedata, err := base64.URLEncoding.DecodeString(number)
	if err != nil {
		return z, err
	}

	return z.SetBytes(bytedata), nil
}

// DO NOT EDIT THE BELOW WITHOUT KNOWING WHAT YOU ARE DOING
/**
 * Computes the multiplicative inverse of the number on the field prime; more
 * specifically, number * inverse == 1; Note: number should never be zero
**/
func modInverse(number *big.Int) (*big.Int, error) {
	if number == big.NewInt(0) {
		return big.NewInt(-1), fmt.Errorf("Number is zero")
	}
	copy := big.NewInt(0).Set(number)
	copy = copy.Mod(copy, prime)
	pcopy := big.NewInt(0).Set(prime)
	x := big.NewInt(0)
	y := big.NewInt(0)

	copy.GCD(x, y, pcopy, copy)

	result := big.NewInt(0).Set(prime)

	result = result.Add(result, y)
	result = result.Mod(result, prime)
	return result, nil
}

/**
 * Converts a byte array into an a 256-bit big.Int, arraied based upon size of
 * the input byte; all values are right-padded to length 256, even if the most
 * significant bit is zero.
**/
func splitByteToInt(secret []byte) []*big.Int {
	hex_data := hex.EncodeToString(secret)
	count := int(math.Ceil(float64(len(hex_data)) / 64.0))

	result := make([]*big.Int, count)

	for i := 0; i < count; i++ {
		if (i+1)*64 < len(hex_data) {
			result[i], _ = big.NewInt(0).SetString(hex_data[i*64:(i+1)*64], 16)
		} else {
			data := strings.Join([]string{hex_data[i*64:], strings.Repeat("0", 64-(len(hex_data)-i*64))}, "")
			result[i], _ = big.NewInt(0).SetString(data, 16)
		}
	}

	return result
}

/**
 * Evauluates a polynomial with coefficients specified in reverse order:
 * evaluatePolynomial([a, b, c, d], x):
 * 		returns a + bx + cx^2 + dx^3
**/
func evaluatePolynomial(polynomial []*big.Int, value *big.Int) *big.Int {
	last := len(polynomial) - 1
	var result *big.Int = big.NewInt(0).Set(polynomial[last])

	for s := last - 1; s >= 0; s-- {
		result = result.Mul(result, value)
		result = result.Add(result, polynomial[s])
		result = result.Mod(result, prime)
	}

	return result
}

func GetRandomString(n int) string {
	// random string implementation courtesy: icza
	// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
	const (
		letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	var src = xrand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
