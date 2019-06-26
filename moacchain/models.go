/*
 * Copyright 2018 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package moacchain

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/blocktree/openwallet/crypto"
	"github.com/blocktree/openwallet/openwallet"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tidwall/gjson"
)

// type Vin struct {
// 	Coinbase string
// 	TxID     string
// 	Vout     uint64
// 	N        uint64
// 	Addr     string
// 	Value    string
// }

// type Vout struct {
// 	N            uint64
// 	Addr         string
// 	Value        string
// 	ScriptPubKey string
// 	Type         string
// }

type Block struct {
	Hash                  string // actually block signature in MOAC chain
	PrevBlockHash         string // actually block signature in MOAC chain
	TransactionMerkleRoot string
	Timestamp             uint64
	Height                uint64
	Transactions          []string
}

type Transaction struct {
	IsCoinBase  bool
	TxID        string
	Fee         *big.Int
	From        string
	To          string
	Amount      *big.Int
	BlockHeight uint64
	BlockHash   string
	Status      string
}

func (c *Client) NewTransaction(json *gjson.Result) *Transaction {
	obj := &Transaction{}
	if gjson.Get(json.Raw, "from").String() == "0x0000000000000000000000000000000000000064" ||
		gjson.Get(json.Raw, "to").String() == "0x0000000000000000000000000000000000000065" ||
		(gjson.Get(json.Raw, "v").String() == "0x0" && gjson.Get(json.Raw, "r").String() == "0x0" && gjson.Get(json.Raw, "s").String() == "0x0") {
		obj.IsCoinBase = true
		return obj
	}

	obj.TxID = gjson.Get(json.Raw, "hash").String()
	gasUsed, _ := c.getGasUsed(obj.TxID)
	gasPrice, _ := big.NewInt(0).SetString(gjson.Get(json.Raw, "gasPrice").String()[2:], 16)
	obj.Fee = gasUsed.Mul(gasUsed, gasPrice)
	obj.Amount, _ = big.NewInt(0).SetString(gjson.Get(json.Raw, "value").String()[2:], 16)
	obj.From = gjson.Get(json.Raw, "from").String()
	obj.To = gjson.Get(json.Raw, "to").String()
	blockHeight, _ := strconv.ParseUint(gjson.Get(json.Raw, "blockNumber").String()[2:], 16, 64)
	obj.BlockHeight = blockHeight
	obj.BlockHash = gjson.Get(json.Raw, "blockHash").String()

	return obj
}

func NewBlock(json *gjson.Result) *Block {
	obj := &Block{}
	// 解析
	obj.Hash = gjson.Get(json.Raw, "hash").String()
	obj.PrevBlockHash = gjson.Get(json.Raw, "parentHash").String()
	obj.TransactionMerkleRoot = gjson.Get(json.Raw, "transactionsRoot").String()
	timestamp, _ := strconv.ParseUint(gjson.Get(json.Raw, "timestamp").String()[2:], 16, 64)
	obj.Timestamp = timestamp
	height, _ := strconv.ParseUint(gjson.Get(json.Raw, "number").String()[2:], 16, 64)
	obj.Height = height

	for _, tx := range gjson.Get(json.Raw, "transactions").Array() {
		obj.Transactions = append(obj.Transactions, tx.String())
	}

	return obj
}

//BlockHeader 区块链头
func (b *Block) BlockHeader() *openwallet.BlockHeader {

	obj := openwallet.BlockHeader{}
	//解析json
	obj.Hash = b.Hash
	//obj.Confirmations = b.Confirmations
	obj.Merkleroot = b.TransactionMerkleRoot
	obj.Previousblockhash = b.PrevBlockHash
	obj.Height = b.Height
	//obj.Version = uint64(b.Version)
	obj.Time = b.Timestamp
	obj.Symbol = Symbol

	return &obj
}

//UnscanRecords 扫描失败的区块及交易
type UnscanRecord struct {
	ID          string `storm:"id"` // primary key
	BlockHeight uint64
	TxID        string
	Reason      string
}

func NewUnscanRecord(height uint64, txID, reason string) *UnscanRecord {
	obj := UnscanRecord{}
	obj.BlockHeight = height
	obj.TxID = txID
	obj.Reason = reason
	obj.ID = common.Bytes2Hex(crypto.SHA256([]byte(fmt.Sprintf("%d_%s", height, txID))))
	return &obj
}
