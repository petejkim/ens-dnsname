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
	d, err := dnsname.Decode(mustDecodeHex("04746573740365746800"))
	shouldBeNil(t, err)
	shouldEqual(t, d, "test.eth")

	d, err = dnsname.Decode(mustDecodeHex("047065746504746573740365746800"))
	shouldBeNil(t, err)
	shouldEqual(t, d, "pete.test.eth")

	d, err = dnsname.Decode(mustDecodeHex("047468697302697304746573740365746800"))
	shouldBeNil(t, err)
	shouldEqual(t, d, "this.is.test.eth")

	d, err = dnsname.Decode(mustDecodeHex("076578616d706c650378797a00"))
	shouldBeNil(t, err)
	shouldEqual(t, d, "example.xyz")
}

func TestDecodeSadCases(t *testing.T) {
	d, err := dnsname.Decode(mustDecodeHex("047465737403657468"))
	shouldEqual(t, d, "")
	shouldHaveError(t, err, "out of bounds")

	d, err = dnsname.Decode(mustDecodeHex("09746573740365746800"))
	shouldEqual(t, d, "")
	shouldHaveError(t, err, "out of bounds")

	d, err = dnsname.Decode(mustDecodeHex("04746573740065746800"))
	shouldEqual(t, d, "")
	shouldHaveError(t, err, "unexpected terminator")

	d, err = dnsname.Decode(mustDecodeHex("04740073740365746800"))
	shouldEqual(t, d, "")
	shouldHaveError(t, err, "unexpected null-zero")

	d, err = dnsname.Decode(mustDecodeHex("40303030303030303030313030303030303030303230303030303030303033303030303030303030343030303030303030303530303030303030303036313233340365746800"))
	shouldEqual(t, d, "")
	shouldHaveError(t, err, "label too long")
}

func TestEncodeHappyCases(t *testing.T) {
	e, err := dnsname.Encode("test.eth")
	shouldBeNil(t, err)
	shouldEqualBytes(t, e, mustDecodeHex("04746573740365746800"))

	e, err = dnsname.Encode("pete.test.eth")
	shouldBeNil(t, err)
	shouldEqualBytes(t, e, mustDecodeHex("047065746504746573740365746800"))

	e, err = dnsname.Encode("this.is.test.eth")
	shouldBeNil(t, err)
	shouldEqualBytes(t, e, mustDecodeHex("047468697302697304746573740365746800"))

	e, err = dnsname.Encode("example.xyz")
	shouldBeNil(t, err)
	shouldEqualBytes(t, e, mustDecodeHex("076578616d706c650378797a00"))

	e, err = dnsname.Encode("")
	shouldBeNil(t, err)
	shouldEqualBytes(t, e, mustDecodeHex("0000"))
}

func TestEncodeSadCases(t *testing.T) {
	e, err := dnsname.Encode("0000000001000000000200000000030000000004000000000500000000061234.eth")
	shouldEqualBytes(t, e, []byte{})
	shouldHaveError(t, err, "label too long")
}
