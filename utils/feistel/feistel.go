package feistel

import (
	"errors"
	"strings"
)

const (
	// Safe uppercase alphabet, base-32 (no I, L, O, 0, 1)
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

func Encode(id int64, keys []int64) (string, error) {
	if id <= 0 || id > maxID {
		return "", errors.New("encode requires id between 1 and 9,999,999,999")
	}

	var obfuscated int64
	var codeLength int

	if id < threshold6Chars {
		obfuscated = feistel30Encrypt(id, keys)
		codeLength = 6
	} else {
		normalizedID := id - threshold6Chars
		obfuscated = feistel34Encrypt(normalizedID, keys)
		codeLength = 7
	}

	chars := make([]byte, codeLength)
	for i := codeLength - 1; i >= 0; i-- {
		chars[i] = alphabet[obfuscated%alphabetBase]
		obfuscated /= alphabetBase
	}
	return string(chars), nil
}

func Decode(code string, keys []int64) (int64, error) {
	code = strings.ToUpper(code)
	var num int64
	for _, c := range code {
		idx := strings.IndexRune(alphabet, c)
		if idx == -1 {
			return 0, errors.New("invalid character in code")
		}
		num = num*alphabetBase + int64(idx)
	}

	switch len(code) {
	case 6:
		return feistel30Decrypt(num, keys), nil
	case 7:
		normalizedID := feistel34Decrypt(num, keys)
		return normalizedID + threshold6Chars, nil
	default:
		return 0, errors.New("code must be 6 or 7 characters")
	}
}

func feistel30Encrypt(n int64, keys []int64) int64 {
	left := (n >> f30HalfSize) & f30Mask
	right := n & f30Mask

	for _, k := range keys {
		newLeft := right
		right = left ^ (feistelRound(right, k, f30Mask) & f30Mask)
		left = newLeft
	}
	return (left << f30HalfSize) | right
}

func feistel30Decrypt(n int64, keys []int64) int64 {
	left := (n >> f30HalfSize) & f30Mask
	right := n & f30Mask

	for i := len(keys) - 1; i >= 0; i-- {
		newRight := left
		left = right ^ (feistelRound(left, keys[i], f30Mask) & f30Mask)
		right = newRight
	}
	return (left << f30HalfSize) | right
}

func feistel34Encrypt(n int64, keys []int64) int64 {
	left := (n >> f34HalfSize) & f34Mask
	right := n & f34Mask

	for _, k := range keys {
		newLeft := right
		right = left ^ (feistelRound(right, k, f34Mask) & f34Mask)
		left = newLeft
	}
	return (left << f34HalfSize) | right
}

func feistel34Decrypt(n int64, keys []int64) int64 {
	left := (n >> f34HalfSize) & f34Mask
	right := n & f34Mask

	for i := len(keys) - 1; i >= 0; i-- {
		newRight := left
		left = right ^ (feistelRound(left, keys[i], f34Mask) & f34Mask)
		right = newRight
	}
	return (left << f34HalfSize) | right
}

func feistelRound(r, k, mask int64) int64 {
	r = (r ^ k) * 0x45D9F3B
	r ^= r >> 16
	return r & mask
}
