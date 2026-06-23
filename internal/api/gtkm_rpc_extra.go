package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
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

// ListRecentTransactions scans the last `blocks` blocks (including latest) and returns a
// slice of Transactions found in those blocks. This is a best-effort, potentially heavy
// operation and should be used with care.
func (c *GTKMClient) ListRecentTransactions(blocks int) ([]Transaction, error) {
	if blocks <= 0 {
		blocks = 50
	}
	latest, err := c.GetBlockNumber()
	if err != nil {
		return nil, err
	}

	var out []Transaction
	for i := 0; i < blocks; i++ {
		if latest < uint64(i) {
			break
		}
		n := latest - uint64(i)
		// Fetch block with full transactions
		var blk map[string]interface{}
		numHex := fmt.Sprintf("0x%x", n)
		if err := c.client.Call(&blk, "eth_getBlockByNumber", numHex, true); err != nil {
			// continue on error to try other blocks
			continue
		}
		if blk == nil {
			continue
		}
		// extract transactions
		if txs, ok := blk["transactions"].([]interface{}); ok {
			for _, raw := range txs {
				if m, ok := raw.(map[string]interface{}); ok {
					var t Transaction
					if h, ok := m["hash"].(string); ok {
						t.Hash = common.HexToHash(h)
					}
					if f, ok := m["from"].(string); ok {
						t.From = common.HexToAddress(f)
					}
					if to, ok := m["to"].(string); ok && to != "" {
						t.To = common.HexToAddress(to)
					}
					if v, ok := m["value"].(string); ok {
						t.Value = new(big.Int)
						if strings.HasPrefix(v, "0x") {
							t.Value.SetString(v[2:], 16)
						} else {
							t.Value.SetString(v, 10)
						}
					}
					if g, ok := m["gas"].(string); ok {
						fmt.Sscanf(g, "0x%x", &t.Gas)
					}
					if gp, ok := m["gasPrice"].(string); ok {
						t.GasPrice = new(big.Int)
						if strings.HasPrefix(gp, "0x") {
							t.GasPrice.SetString(gp[2:], 16)
						} else {
							t.GasPrice.SetString(gp, 10)
						}
					}
					if nstr, ok := m["nonce"].(string); ok {
						fmt.Sscanf(nstr, "0x%x", &t.Nonce)
					}
					if inp, ok := m["input"].(string); ok && inp != "0x" {
						if b, err := hex.DecodeString(strings.TrimPrefix(inp, "0x")); err == nil {
							t.Data = b
						}
					}
					if bn, ok := m["blockNumber"].(string); ok && bn != "" {
						fmt.Sscanf(bn, "0x%x", &t.BlockNum)
					}
					// Get receipt status if blockNumber present
					if t.BlockNum > 0 {
						if receipt, err := c.GetTransactionReceipt(t.Hash); err == nil && receipt != nil {
							if status, ok := receipt["status"].(string); ok {
								fmt.Sscanf(status, "0x%x", &t.Status)
							}
						}
					}
					out = append(out, t)
				}
			}
		}
	}

	return out, nil
}
