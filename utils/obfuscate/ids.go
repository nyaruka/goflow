package obfuscate

import (
	"errors"
	"strings"
)

const (
	// Safe uppercase alphabet, base-32 (no I, O, 0, 1)
	alphabet     = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	alphabetBase = int64(len(alphabet)) // 32

	// Threshold where we switch from 6 → 7 chars
	threshold6Chars = int64(1 << 30) // 32^6 = 1,073,741,824

	maxID = int64(9_999_999_999)

	f30BitSize  = 30
	f30HalfSize = f30BitSize / 2         // 15 bits
	f30Mask     = (1 << f30HalfSize) - 1 // 0x7FFF

	f34BitSize  = 34
	f34HalfSize = f34BitSize / 2         // 17 bits
	f34Mask     = (1 << f34HalfSize) - 1 // 0x1FFFF
)

// EncodeID encodes a numeric ID into an obfuscated string using a Feistel cipher with the provided key.
func EncodeID(id int64, key [4]int32) (string, error) {
	if id <= 0 || id > maxID {
		return "", errors.New("encode requires id between 1 and 9,999,999,999")
	}

	var obfuscated int64
	var codeLength int

	if id < threshold6Chars {
		obfuscated = feistel30Encrypt(id, key)
		codeLength = 6
	} else {
		normalizedID := id - threshold6Chars
		obfuscated = feistel34Encrypt(normalizedID, key)
		codeLength = 7
	}

	chars := make([]byte, codeLength)
	for i := codeLength - 1; i >= 0; i-- {
		chars[i] = alphabet[obfuscated%alphabetBase]
		obfuscated /= alphabetBase
	}
	return string(chars), nil
}

// DecodeID decodes an obfuscated string back into the original numeric ID using a Feistel cipher with the provided key.
func DecodeID(code string, key [4]int32) (int64, error) {
	if !WasID(code) {
		return 0, errors.New("code is not a valid obfuscated ID")
	}

	code = strings.ToUpper(code)
	var num int64
	for _, c := range code {
		idx := strings.IndexRune(alphabet, c)
		num = num*alphabetBase + int64(idx)
	}

	if len(code) == 6 {
		return feistel30Decrypt(num, key), nil
	}

	// must be 7 characters
	normalizedID := feistel34Decrypt(num, key)
	return normalizedID + threshold6Chars, nil
}

// WasID returns true if the provided code is a valid obfuscated ID (6 or 7 characters from the alphabet).
func WasID(code string) bool {
	code = strings.ToUpper(code)
	if len(code) != 6 && len(code) != 7 {
		return false
	}
	for _, c := range code {
		if !strings.ContainsRune(alphabet, c) {
			return false
		}
	}
	return true
}

func feistel30Encrypt(n int64, key [4]int32) int64 {
	left := (n >> f30HalfSize) & f30Mask
	right := n & f30Mask

	for _, k := range key {
		newLeft := right
		right = left ^ (feistelRound(right, int64(k), f30Mask) & f30Mask)
		left = newLeft
	}
	return (left << f30HalfSize) | right
}

func feistel30Decrypt(n int64, key [4]int32) int64 {
	left := (n >> f30HalfSize) & f30Mask
	right := n & f30Mask

	for i := len(key) - 1; i >= 0; i-- {
		newRight := left
		left = right ^ (feistelRound(left, int64(key[i]), f30Mask) & f30Mask)
		right = newRight
	}
	return (left << f30HalfSize) | right
}

func feistel34Encrypt(n int64, key [4]int32) int64 {
	left := (n >> f34HalfSize) & f34Mask
	right := n & f34Mask

	for _, k := range key {
		newLeft := right
		right = left ^ (feistelRound(right, int64(k), f34Mask) & f34Mask)
		left = newLeft
	}
	return (left << f34HalfSize) | right
}

func feistel34Decrypt(n int64, key [4]int32) int64 {
	left := (n >> f34HalfSize) & f34Mask
	right := n & f34Mask

	for i := len(key) - 1; i >= 0; i-- {
		newRight := left
		left = right ^ (feistelRound(left, int64(key[i]), f34Mask) & f34Mask)
		right = newRight
	}
	return (left << f34HalfSize) | right
}

func feistelRound(r, k, mask int64) int64 {
	r = (r ^ k) * 0x45D9F3B
	r ^= r >> 16
	return r & mask
}
