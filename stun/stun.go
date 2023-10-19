package stun

import (
	"crypto/rand"
	"io"
)

func NewTxID() TxID {
	id := make(TxID, 12)
	io.ReadFull(rand.Reader, id)
	return id
}
