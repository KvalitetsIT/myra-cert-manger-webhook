package testutil

import (
	"crypto/rand"
	"math/big"
)

func Get_random_port() (uint16, error) {
	min := 15534
	max := 65535
	randPort, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		return 0, err
	}
	return uint16(randPort.Int64() + int64(min)), nil
}
