package moacchain

import (
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

const (
	testNodeAPI = "http://"
)

func Test_getBlockHeight(t *testing.T) {

	c := NewClient(testNodeAPI, true)

	r, err := c.getBlockHeight()

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("height:", r)
	}

}

func Test_getBlockByHeight(t *testing.T) {

	c := NewClient(testNodeAPI, true)
	r, err := c.getBlockByHeight(2758235)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
}
func Test_getBlockByHash(t *testing.T) {
	hash := "3Uvb87ukKKwVeU6BFsZ21hy9sSbSd3Rd5QZTWbNop1d3TaY9ZzceJAT54vuY8XXQmw6nDx8ZViPV3cVznAHTtiVE"

	c := NewClient(testNodeAPI, true)

	r, err := c.Call("blocks/signature/"+hash, nil)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
}

func Test_getBlockHash(t *testing.T) {

	c := NewClient(testNodeAPI, true)

	height := uint64(2758235)

	r, err := c.getBlockHash(height)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}

}

func Test_getBalance(t *testing.T) {

	c := NewClient(testNodeAPI, true)

	address := "0x25ff183be76c9583db211e66dbd481a923f40635"

	r, err := c.getBalance(address)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}

}

func Test_getGasPrice(t *testing.T) {

	c := NewClient(testNodeAPI, true)

	r, err := c.getGasPrice()

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
}

func Test_getGasEstimated(t *testing.T) {
	c := NewClient(testNodeAPI, true)
	from := "0x25ff183be76c9583db211e66dbd481a923f40635"
	to := "0x39aa046cf77c2877cc5f42e0224c27969b3ad8d9"
	gasLimit := big.NewInt(10000)
	gasPrice := big.NewInt(200000000000000)
	amount := big.NewInt(2000000000000000000)
	r, err := c.getGasEstimated(from, to, gasLimit, gasPrice, amount)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
}

func Test_getTransaction(t *testing.T) {

	c := NewClient(testNodeAPI, true)
	txid := "0x448c135168ed2b2c387ecde805fe3beef8b4a906fd727152dbf6f01d91dd1760"
	r, err := c.getTransaction(txid)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}

	// txid = "1b0147a6b5660215e9a37ca34fe9a6298988e45f7eefbd8c4b98993f4e762c3e" //"9KBoALfTjvZLJ6CAuJCGyzRA1aWduiNFMvbqTchfBVpF"

	// r, err = c.getTransaction(txid)

	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(r)
	// }

	// txid = "3ca0888b232df90d910a921d2f4004bb61a80bbbe27caee7107de282576e38a0" //"9KBoALfTjvZLJ6CAuJCGyzRA1aWduiNFMvbqTchfBVpF"

	// r, err = c.getTransaction(txid)

	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(r)
	// }
}

func Test_convert(t *testing.T) {

	amount := uint64(5000000001)

	amountStr := fmt.Sprintf("%d", amount)

	fmt.Println(amountStr)

	d, _ := decimal.NewFromString(amountStr)

	w, _ := decimal.NewFromString("100000000")

	d = d.Div(w)

	fmt.Println(d.String())

	d = d.Mul(w)

	fmt.Println(d.String())

	r, _ := strconv.ParseInt(d.String(), 10, 64)

	fmt.Println(r)

	fmt.Println(time.Now().UnixNano())
}

func Test_getTransactionByAddresses(t *testing.T) {
	addrs := "ARAA8AnUYa4kWwWkiZTTyztG5C6S9MFTx11"

	c := NewClient(testNodeAPI, true)
	result, err := c.getMultiAddrTransactions(0, -1, addrs)

	if err != nil {
		t.Error("get transactions failed!")
	} else {
		for _, tx := range result {
			fmt.Println(tx.TxID)
		}
	}
}

func Test_tmp(t *testing.T) {

	c := NewClient(testNodeAPI, true)

	block, err := c.getNonce("0x3d5ecba0f712179c29254928002ee58cca2b09a3")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(block)
}
