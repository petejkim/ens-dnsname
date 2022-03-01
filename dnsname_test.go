package dnsname_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	dnsname "github.com/petejkim/ens-dnsname"
)

func mustDecodeHex(hexString string) []byte {
	b, err := hex.DecodeString(hexString)
	if err != nil {
		panic(err)
	}
	return b
}

func shouldBeNil(t *testing.T, actual interface{}) {
	if actual != nil {
		t.Errorf("expected %v to be nil", actual)
	}
}

func shouldEqual(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func shouldEqualBytes(t *testing.T, actual, expected []byte) {
	if !bytes.Equal(actual, expected) {
		t.Errorf("expected %s, got %s", hex.EncodeToString(expected), hex.EncodeToString(actual))
	}
}

func shouldHaveError(t *testing.T, actual error, expected string) {
	if actual.Error() != expected {
		t.Errorf("expected %s, got %v", expected, actual)
	}
}

func TestDecodeHappyCases(t *testing.T) {
	for _, tc := range [][]string{
		{"04746573740365746800", "test.eth"},
		{"047065746504746573740365746800", "pete.test.eth"},
		{"047468697302697304746573740365746800", "this.is.test.eth"},
		{"076578616d706c650378797a00", "example.xyz"},
	} {
		d, err := dnsname.Decode(mustDecodeHex(tc[0]))
		shouldBeNil(t, err)
		shouldEqual(t, d, tc[1])
	}
}

func TestDecodeSadCases(t *testing.T) {
	for _, tc := range [][]string{
		{"047465737403657468", "out of bounds"},
		{"09746573740365746800", "out of bounds"},
		{"04746573740065746800", "unexpected terminator"},
		{"04740073740365746800", "unexpected null-zero"},
		{"40303030303030303030313030303030303030303230303030303030303033303030303030303030343030303030303030303530303030303030303036313233340365746800", "label too long"},
	} {
		d, err := dnsname.Decode(mustDecodeHex(tc[0]))
		shouldEqual(t, d, "")
		shouldHaveError(t, err, tc[1])
	}
}

func TestEncodeHappyCases(t *testing.T) {
	for _, name := range []string{"test.eth", ".test.eth", "test.eth.", "..test.eth..."} {
		e, err := dnsname.Encode(name)
		shouldBeNil(t, err)
		shouldEqualBytes(t, e, mustDecodeHex("04746573740365746800"))
	}

	for _, tc := range [][]string{
		{"pete.test.eth", "047065746504746573740365746800"},
		{"this.is.test.eth", "047468697302697304746573740365746800"},
		{"example.xyz", "076578616d706c650378797a00"},
		{"", "0000"},
	} {
		e, err := dnsname.Encode(tc[0])
		shouldBeNil(t, err)
		shouldEqualBytes(t, e, mustDecodeHex(tc[1]))
	}
}

func TestEncodeSadCases(t *testing.T) {
	e, err := dnsname.Encode("0000000001000000000200000000030000000004000000000500000000061234.eth")
	shouldEqualBytes(t, e, []byte{})
	shouldHaveError(t, err, "label too long")
}
