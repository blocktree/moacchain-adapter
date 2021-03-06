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
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/blocktree/openwallet/log"
	"github.com/imroc/req"
	"github.com/tidwall/gjson"
)

type ClientInterface interface {
	Call(path string, request []interface{}) (*gjson.Result, error)
}

// A Client is a Elastos RPC client. It performs RPCs over HTTP using JSON
// request and responses. A Client must be configured with a secret token
// to authenticate with other Cores on the network.
type Client struct {
	BaseURL     string
	AccessToken string
	Debug       bool
	client      *req.Req
	//Client *req.Req
}

type Response struct {
	Code    int         `json:"code,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Message string      `json:"message,omitempty"`
	Id      string      `json:"id,omitempty"`
}

func NewClient(url string /*token string,*/, debug bool) *Client {
	c := Client{
		BaseURL: url,
		//	AccessToken: token,
		Debug: debug,
	}

	api := req.New()
	//trans, _ := api.Client().Transport.(*http.Transport)
	//trans.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c.client = api

	return &c
}

// Call calls a remote procedure on another node, specified by the path.
func (c *Client) Call(path string, request []interface{}) (*gjson.Result, error) {

	var (
		body = make(map[string]interface{}, 0)
	)

	if c.client == nil {
		return nil, errors.New("API url is not setup. ")
	}

	authHeader := req.Header{
		"Accept":        "application/json",
		"Authorization": "Basic " + c.AccessToken,
	}

	//json-rpc
	body["jsonrpc"] = "2.0"
	body["id"] = "101"
	body["method"] = path
	body["params"] = request

	if c.Debug {
		log.Std.Info("Start Request API...")
	}

	r, err := c.client.Post(c.BaseURL, req.BodyJSON(&body), authHeader)

	if c.Debug {
		log.Std.Info("Request API Completed")
	}

	if c.Debug {
		log.Std.Info("%+v", r)
	}

	if err != nil {
		return nil, err
	}

	resp := gjson.ParseBytes(r.Bytes())
	err = isError(&resp)
	if err != nil {
		return nil, err
	}

	result := resp.Get("result")

	return &result, nil
}

// See 2 (end of page 4) http://www.ietf.org/rfc/rfc2617.txt
// "To receive authorization, the client sends the userid and password,
// separated by a single colon (":") character, within a base64
// encoded string in the credentials."
// It is not meant to be urlencoded.
func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))

	//return username + ":" + password
}

//isError 是否报错
func isError(result *gjson.Result) error {
	var (
		err error
	)

	/*
		//failed 返回错误
		{
			"result": null,
			"error": {
				"code": -8,
				"message": "Block height out of range"
			},
			"id": "foo"
		}
	*/

	if !result.Get("error").IsObject() {

		if !result.Get("result").Exists() {
			return errors.New("Response is empty! ")
		}

		return nil
	}

	errInfo := fmt.Sprintf("[%d]%s",
		result.Get("error.code").Int(),
		result.Get("error.message").String())
	err = errors.New(errInfo)

	return err
}

// 获取当前区块高度
func (c *Client) getBlockHeight() (uint64, error) {

	request := []interface{}{}

	resp, err := c.Call("mc_blockNumber", request)
	if err != nil {
		return 0, err
	}
	height, err := strconv.ParseUint(resp.String()[2:], 16, 64)
	if err != nil {
		return 0, err
	}
	return height, nil
}

// 通过高度获取区块哈希
func (c *Client) getBlockHash(height uint64) (string, error) {
	request := []interface{}{
		uint64ToHexString(height),
		false,
	}
	resp, err := c.Call("mc_getBlockByNumber", request)

	if err != nil {
		return "", err
	}

	return resp.Get("hash").String(), nil
}

func (c *Client) getNonce(address string) (uint64, error) {
	request := []interface{}{
		address,
		"pending",
	}

	r, err := c.Call("mc_getTransactionCount", request)

	if err != nil {
		return 0, err
	}

	nonce, _ := strconv.ParseUint(r.String()[2:], 16, 64)

	return nonce, nil
}

func (c *Client) getGasEstimated(from, to string, gasLimit, gasPrice, amount *big.Int) (*big.Int, error) {
	if amount == nil {
		request := []interface{}{
			map[string]interface{}{
				"from":     from,
				"to":       to,
				"gas":      "0x" + gasLimit.Text(16),
				"gasPrice": "0x" + gasPrice.Text(16),
			},
		}

		r, err := c.Call("mc_estimateGas", request)
		if err != nil {
			return nil, err
		}

		gasUsed, _ := big.NewInt(0).SetString(r.String()[2:], 16)
		return gasUsed.Sub(gasUsed, big.NewInt(1)), nil
	}
	request := []interface{}{
		map[string]interface{}{
			"from":     from,
			"to":       to,
			"gas":      "0x" + gasLimit.Text(16),
			"gasPrice": "0x" + gasPrice.Text(16),
			"value":    "0x" + amount.Text(16),
		},
	}

	r, err := c.Call("mc_estimateGas", request)
	if err != nil {
		return nil, err
	}

	gasUsed, _ := big.NewInt(0).SetString(r.String()[2:], 16)
	return gasUsed.Sub(gasUsed, big.NewInt(1)), nil
}

// 获取地址余额
func (c *Client) getBalance(address string) (*AddrBalance, error) {
	request := []interface{}{
		address,
		"latest",
	}

	r, err := c.Call("mc_getBalance", request)

	if err != nil {
		return nil, err
	}
	balance, pass := big.NewInt(0).SetString(r.String()[2:], 16)
	if !pass {
		return nil, errors.New("Failed to get balance of :" + address)
	}
	return &AddrBalance{Address: address, Balance: balance}, nil
}

// 获取区块信息
func (c *Client) getBlock(hash string) (*Block, error) {
	request := []interface{}{
		hash,
		false,
	}
	resp, err := c.Call("mc_getBlockByHash", request)

	if err != nil {
		return nil, err
	}
	return NewBlock(resp), nil
}

func (c *Client) getGasPrice() (*big.Int, error) {
	request := []interface{}{}
	resp, err := c.Call("mc_gasPrice", request)
	if err != nil {
		return nil, err
	}
	gasPrice, _ := big.NewInt(0).SetString(resp.String()[2:], 16)
	return gasPrice, nil
}

func uint64ToHexString(value uint64) string {
	return "0x" + strconv.FormatUint(value, 16)
}
func (c *Client) getBlockByHeight(height uint64) (*Block, error) {
	request := []interface{}{
		uint64ToHexString(height),
		false,
	}
	resp, err := c.Call("mc_getBlockByNumber", request)

	if err != nil {
		return nil, err
	}
	return NewBlock(resp), nil
}

func (c *Client) getTransaction(txid string) (*Transaction, error) {
	request := []interface{}{
		txid,
	}
	resp, err := c.Call("mc_getTransactionByHash", request)
	if err != nil {
		return nil, err
	}

	if resp.Raw == "null" {
		return nil, errors.New("Transaction does not exist!")
	}
	return c.NewTransaction(resp), nil
}

func (c *Client) getGasUsed(txid string) (*big.Int, error) {
	request := []interface{}{
		txid,
	}
	resp, err := c.Call("mc_getTransactionReceipt", request)
	if err != nil {
		return nil, err
	}

	if resp.Raw == "null" {
		return nil, errors.New("Transaction does not exist!")
	}

	gasUsed, _ := big.NewInt(0).SetString(gjson.Get(resp.Raw, "gasUsed").String()[2:], 16)

	return gasUsed, nil
}

func (c *Client) sendTransaction(rawTx string) (string, error) {
	request := []interface{}{
		rawTx,
	}

	resp, err := c.Call("mc_sendRawTransaction", request)

	if err != nil {
		return "", err
	}

	return resp.String(), nil
}
