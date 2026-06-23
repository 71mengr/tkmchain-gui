package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
)

// Below additions extend GTKMClient with utility RPC methods used by the UI.

// GetGasPrice returns the current gas price from the node (eth_gasPrice)
func (c *GTKMClient) GetGasPrice() (*big.Int, error) {
	var result string
	if err := c.client.Call(&result, "eth_gasPrice"); err != nil {
		return nil, err
	}
	gp := new(big.Int)
	if strings.HasPrefix(result, "0x") {
		gp.SetString(result[2:], 16)
	} else {
		gp.SetString(result, 10)
	}
	return gp, nil
}

// EstimateGas calls eth_estimateGas with a transaction object and returns the gas estimate.
func (c *GTKMClient) EstimateGas(tx map[string]interface{}) (uint64, error) {
	var result string
	if err := c.client.Call(&result, "eth_estimateGas", tx); err != nil {
		return 0, err
	}
	var gas uint64
	if strings.HasPrefix(result, "0x") {
		fmt.Sscanf(result[2:], "%x", &gas)
	} else {
		fmt.Sscanf(result, "%x", &gas)
	}
	return gas, nil
}

// SendRawTransaction sends a signed, RLP-encoded transaction via eth_sendRawTransaction.
func (c *GTKMClient) SendRawTransaction(rawTxHex string) (string, error) {
	// rawTxHex may include 0x prefix or not
	if strings.HasPrefix(rawTxHex, "0x") == false {
		rawTxHex = "0x" + rawTxHex
	}
	var result string
	if err := c.client.Call(&result, "eth_sendRawTransaction", rawTxHex); err != nil {
		return "", err
	}
	return result, nil
}

// RK* methods map to rotating-king specific RPC methods used by the gtkm fork.

func (c *GTKMClient) RKAdd(params map[string]interface{}) (interface{}, error) {
	var result interface{}
	if err := c.client.Call(&result, "rk_add", params); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *GTKMClient) RKList() (interface{}, error) {
	var result interface{}
	if err := c.client.Call(&result, "rk_list"); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *GTKMClient) RKStats() (interface{}, error) {
	var result interface{}
	if err := c.client.Call(&result, "rk_stats"); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *GTKMClient) RKStatus() (interface{}, error) {
	var result interface{}
	if err := c.client.Call(&result, "rk_status"); err != nil {
		return nil, err
	}
	return result, nil
}

// Helper: pretty-print JSON for UI
func toJSON(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(b)
}
