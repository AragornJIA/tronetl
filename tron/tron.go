package tron

import (
	"bytes"
	"encoding/json"
	"io"
	"math/big"
	"math/rand"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type TronClient struct {
	httpURI string
	jsonURI string
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func NewTronClient(providerURL string) *TronClient {
	return &TronClient{
		httpURI: providerURL + ":8090",
		jsonURI: providerURL + ":50545/jsonrpc",
	}
}

func (c *TronClient) GetTxInfosByNumber(number uint64) []TxInfo {
	url := c.httpURI + "/wallet/gettransactioninfobyblocknum"
	payload, err := json.Marshal(map[string]any{
		"num": number,
	})
	chk(err)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	chk(err)
	body, err := io.ReadAll(resp.Body)
	chk(err)

	var txInfos []TxInfo
	err = json.Unmarshal(body, &txInfos)
	chk(err)

	return txInfos
}

func (c *TronClient) GetBlockByNumber(number uint64) *Block {
	payload, err := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params": []any{
			toBlockNumArg(new(big.Int).SetUint64(number)), true,
		},
		"id": rand.Int(),
	})
	chk(err)
	resp, err := http.Post(c.jsonURI, "application/json", bytes.NewBuffer(payload))
	chk(err)
	body, err := io.ReadAll(resp.Body)
	chk(err)

	var rpcResp JSONResponse
	var block Block
	err = json.Unmarshal(body, &rpcResp)
	chk(err)
	err = json.Unmarshal(rpcResp.Result, &block)
	chk(err)

	return &block
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}
