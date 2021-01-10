// Package constants contains commonly used Nano constants
package constants

const (
	// SeedLength is the length of a Nano seed in bytes
	SeedLength = 32

	// KeyLength is the length of a Nano key in bytes
	KeyLength = 32

	// KeyNumBits is the number of bits in a Nano key
	KeyNumBits = 256

	// POWSize is the number of bytes in a proof of work
	POWSize = 8

	// ZeroKey is a zeroed out key
	ZeroKey = "0000000000000000000000000000000000000000000000000000000000000000"

	// DecimalPlaces is the number of fractionals for a single Nano
	DecimalPlaces = 30

	// CheckSumLength is the number of bytes for the checksum of an address
	CheckSumLength = 5
)
