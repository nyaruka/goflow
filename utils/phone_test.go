package utils_test

import (
	"testing"

	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestParsePhoneNumber(t *testing.T) {
	assert.Equal(t, "", utils.ParsePhoneNumber("", ""))
	assert.Equal(t, "", utils.ParsePhoneNumber("", "RW"))
	assert.Equal(t, "", utils.ParsePhoneNumber("oui", "RW"))
	assert.Equal(t, "", utils.ParsePhoneNumber("tel = +250788383383", "RW"))
	assert.Equal(t, "", utils.ParsePhoneNumber("Hi my number is +250788383383 thanks", "RW"))
	assert.Equal(t, "", utils.ParsePhoneNumber("Hi my phone is (202) 555-0110 thanks", "US"))
	assert.Equal(t, "", utils.ParsePhoneNumber("+810000000977123456", "CD"))

	assert.Equal(t, "+250788383383", utils.ParsePhoneNumber("+250788383383", ""))
	assert.Equal(t, "+250788383383", utils.ParsePhoneNumber("0788 383 383", "RW"))
	assert.Equal(t, "+12025550110", utils.ParsePhoneNumber("   (202) 555-0110   ", "US"))
	assert.Equal(t, "+12025550110", utils.ParsePhoneNumber("202.555.0110", "US"))
	assert.Equal(t, "+12025550110", utils.ParsePhoneNumber("202-555-0110", "US"))

	// country can be omitted for full numbers
	assert.Equal(t, "+250788383383", utils.ParsePhoneNumber("+250788383383", ""))
	assert.Equal(t, "", utils.ParsePhoneNumber("0788383383", ""))
	assert.Equal(t, "+12025550110", utils.ParsePhoneNumber("+12025550110", ""))
	assert.Equal(t, "", utils.ParsePhoneNumber("(202) 555-0110", ""))
}

func TestFindPhoneNumbers(t *testing.T) {
	assert.Equal(t, []string{}, utils.FindPhoneNumbers("", ""))
	assert.Equal(t, []string{}, utils.FindPhoneNumbers("", "RW"))

	// can match exact
	assert.Equal(t, []string{"+250788383383"}, utils.FindPhoneNumbers("0788 383 383", "RW"))
	assert.Equal(t, []string{"+12025550110"}, utils.FindPhoneNumbers("(202) 555-0110", "US"))

	// can find anywhere in string
	assert.Equal(t, []string{"+250788383383"}, utils.FindPhoneNumbers("tel = 0788383383", "RW"))
	assert.Equal(t, []string{"+250788383383"}, utils.FindPhoneNumbers("Hi my phone is +250788383383 thanks", "RW"))
	assert.Equal(t, []string{"+250788383383"}, utils.FindPhoneNumbers("Hi my phone is 0788383383 thanks", "RW"))
	assert.Equal(t, []string{"+12025550110"}, utils.FindPhoneNumbers("Hi my phone is +12025550110 thanks", "US"))
	assert.Equal(t, []string{"+12025550110"}, utils.FindPhoneNumbers("Hi my phone is (202) 555-0110 thanks", "US"))

	// returns all if more than one
	assert.Equal(t, []string{"+12025550110", "+12025550111"}, utils.FindPhoneNumbers("Mine is (202) 555-0110, his is (202) 555-0111 thanks", "US"))

	// country can be omitted for full numbers
	assert.Equal(t, []string{"+250788383383"}, utils.FindPhoneNumbers("Hi my phone is +250788383383 thanks", ""))
	assert.Equal(t, []string{}, utils.FindPhoneNumbers("Hi my phone is 0788383383 thanks", ""))
	assert.Equal(t, []string{"+12025550110"}, utils.FindPhoneNumbers("Hi my phone is +12025550110 thanks", ""))
	assert.Equal(t, []string{}, utils.FindPhoneNumbers("Hi my phone is (202) 555-0110 thanks", ""))
}
