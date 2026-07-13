package utils

import (
	"bytes"
	"encoding/json"

	"github.com/nyaruka/gocommon/httpx"
)

// ExtractJSON makes a best effort to extract valid JSON from the given bytes
func ExtractJSON(body []byte) []byte {
	// we make a best effort to turn the body into JSON, so we strip out:
	//  1. any invalid UTF-8 sequences
	//  2. null chars
	//  3. escaped null chars (\u0000)
	cleaned := bytes.ToValidUTF8(body, nil)
	cleaned = bytes.ReplaceAll(cleaned, []byte{0}, nil)
	cleaned = []byte(httpx.ReplaceEscapedNulls(string(cleaned), ""))

	if json.Valid(cleaned) {
		return cleaned
	}
	return nil
}
