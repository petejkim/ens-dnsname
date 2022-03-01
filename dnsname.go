package dnsname

import (
	"errors"
	"strings"
)

func Decode(encoded []byte) (string, error) {
	size := len(encoded)
	offset := 0
	var decoded string

	for {
		// read length
		l := int(encoded[offset])

		// length must be less than 64
		if l > 63 {
			return "", errors.New("label too long")
		}

		offset++

		// if null-zero is found
		if l == 0 {
			if offset == size {
				// we're done
				break
			} else {
				return "", errors.New("unexpected terminator")
			}
		}

		// after reading this label, there should be at least one byte left
		// for the terminator
		if size-offset-l < 1 {
			return "", errors.New("out of bounds")
		}

		// read the label
		label := string(encoded[offset : offset+l])

		// the label should not contain a null-zero character
		if strings.ContainsRune(label, rune(0)) {
			return "", errors.New("unexpected null-zero")
		}
		decoded += label
		offset += l

		// if we are not at the end of the name, append a period
		if size-offset > 1 {
			decoded += "."
		}
	}

	return decoded, nil
}

func Encode(name string) ([]byte, error) {
	name = strings.Trim(name, ".")

	encoded := make([]byte, len(name)+2)
	offset := 0

	// split name into labels
	labels := strings.Split(name, ".")

	for _, label := range labels {
		l := len(label)

		// length must be less than 64
		if l > 63 {
			return nil, errors.New("label too long")
		}

		// write length
		encoded[offset] = byte(l)
		offset++

		// write label
		copy(encoded[offset:offset+l], []byte(label))
		offset += l
	}

	return encoded, nil
}
