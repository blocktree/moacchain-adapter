package moacchain_txsigner

import (
	"encoding/hex"
	"fmt"

	"github.com/blocktree/go-owcdrivers/moacchainTransaction"
)

var Default = &TransactionSigner{}

type TransactionSigner struct {
}

// SignTransactionHash 交易哈希签名算法
// required
func (singer *TransactionSigner) SignTransactionHash(msg []byte, privateKey []byte, eccType uint32) ([]byte, error) {
	signature, err := moacchainTransaction.SignRawTransaction(hex.EncodeToString(msg), privateKey)
	if err != nil {
		return nil, fmt.Errorf("ECC sign hash failed")
	}
	return signature, nil
}
