package api

import (
	"context"
	"encoding/hex"
//	"encoding/json"
	"fmt"
	"math/big"
	"strings"
//	"time"

	"github.com/ethereum/go-ethereum/common"
//	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

const Version = "1.0.0"

type GTKMClient struct {
	client   *rpc.Client
	url      string
	ctx      context.Context
	connected bool
	lastError error
}

// GTKM specific types
type KingInfo struct {
	CurrentKing     common.Address         `json:"current_king"`
	MainKing        common.Address         `json:"main_king"`
	NextKing        common.Address         `json:"next_king"`
	AllKings        []common.Address       `json:"all_kings"`
	RotationInfo    RotationInfo           `json:"rotation_info"`
	MonitoringTasks []MonitoringCategory   `json:"monitoring_tasks"`
	IsKing          bool                   `json:"is_king"`
}

type RotationInfo struct {
	CurrentHeight    uint64   `json:"current_height"`
	CurrentIndex     uint64   `json:"current_index"`
	NextIndex        uint64   `json:"next_index"`
	TotalKings       int      `json:"total_kings"`
	RotationInterval uint64   `json:"rotation_interval"`
	BlocksUntilNext  uint64   `json:"blocks_until_next"`
	NextKing         common.Address `json:"next_king"`
}

type MonitoringCategory struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tasks       []string `json:"tasks"`
	Priority    string   `json:"priority"`
}

type Transaction struct {
	Hash     common.Hash    `json:"hash"`
	From     common.Address `json:"from"`
	To       common.Address `json:"to"`
	Value    *big.Int       `json:"value"`
	Gas      uint64         `json:"gas"`
	GasPrice *big.Int       `json:"gasPrice"`
	Nonce    uint64         `json:"nonce"`
	Data     []byte         `json:"data"`
	BlockNum uint64         `json:"blockNum"`
	Status   uint64         `json:"status"` // 0=pending, 1=success, 2=failed
}

type BlockInfo struct {
	Number       uint64         `json:"number"`
	Hash         common.Hash    `json:"hash"`
	ParentHash   common.Hash    `json:"parentHash"`
	Timestamp    uint64         `json:"timestamp"`
	Difficulty   *big.Int       `json:"difficulty"`
	GasUsed      uint64         `json:"gasUsed"`
	GasLimit     uint64         `json:"gasLimit"`
	Coinbase     common.Address `json:"coinbase"`
	TxCount      int            `json:"txCount"`
	Miner        common.Address `json:"miner"`
}

type NodeInfo struct {
	Version    string   `json:"version"`
	NetworkID  uint64   `json:"networkId"`
	PeerCount  int      `json:"peerCount"`
	IsMining   bool     `json:"isMining"`
	Hashrate   float64  `json:"hashrate"`
	Syncing    bool     `json:"syncing"`
}

func NewGTKMClient(url string) (*GTKMClient, error) {
	client, err := rpc.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gtkm node: %w", err)
	}

	return &GTKMClient{
		client:    client,
		url:       url,
		ctx:       context.Background(),
		connected: true,
	}, nil
}

// CheckConnection verifies the connection to gtkm
func (c *GTKMClient) CheckConnection() error {
	var result string
	err := c.client.Call(&result, "eth_blockNumber")
	if err != nil {
		c.connected = false
		c.lastError = err
		return err
	}
	c.connected = true
	c.lastError = nil
	return nil
}

// ==================== ROTATING KING RPC METHODS ====================

func (c *GTKMClient) GetCurrentKing() (common.Address, error) {
	var result common.Address
	err := c.client.Call(&result, "rotatingking_getCurrentKing")
	return result, err
}

func (c *GTKMClient) GetMainKing() (common.Address, error) {
	var result common.Address
	err := c.client.Call(&result, "rotatingking_getMainKing")
	return result, err
}

func (c *GTKMClient) GetNextKing() (common.Address, error) {
	var result common.Address
	err := c.client.Call(&result, "rotatingking_getNextKing")
	return result, err
}

func (c *GTKMClient) GetKingAddresses() ([]common.Address, error) {
	var result []common.Address
	err := c.client.Call(&result, "rotatingking_getKingAddresses")
	return result, err
}

func (c *GTKMClient) GetRotationInfo(height uint64) (RotationInfo, error) {
	var result RotationInfo
	err := c.client.Call(&result, "rotatingking_getRotationInfo", height)
	return result, err
}

func (c *GTKMClient) GetMonitoringResponsibilities() ([]MonitoringCategory, error) {
	var result []MonitoringCategory
	err := c.client.Call(&result, "rotatingking_getMonitoringResponsibilities")
	return result, err
}

func (c *GTKMClient) IsKing(address common.Address) (bool, error) {
	var result bool
	err := c.client.Call(&result, "rotatingking_isKing", address)
	return result, err
}

// ==================== ETH RPC METHODS ====================

func (c *GTKMClient) GetBalance(address common.Address) (*big.Int, error) {
	var result string
	err := c.client.Call(&result, "eth_getBalance", address.Hex(), "latest")
	if err != nil {
		return nil, err
	}
	
	balance := new(big.Int)
	if strings.HasPrefix(result, "0x") {
		balance.SetString(result[2:], 16)
	} else {
		balance.SetString(result, 16)
	}
	return balance, nil
}

func (c *GTKMClient) GetBlockNumber() (uint64, error) {
	var result string
	err := c.client.Call(&result, "eth_blockNumber")
	if err != nil {
		return 0, err
	}
	
	var blockNum uint64
	if strings.HasPrefix(result, "0x") {
		fmt.Sscanf(result[2:], "%x", &blockNum)
	} else {
		fmt.Sscanf(result, "%x", &blockNum)
	}
	return blockNum, nil
}

func (c *GTKMClient) GetBlockByNumber(num uint64) (*BlockInfo, error) {
	var result map[string]interface{}
	numHex := fmt.Sprintf("0x%x", num)
	err := c.client.Call(&result, "eth_getBlockByNumber", numHex, true)
	if err != nil {
		return nil, err
	}
	
	if result == nil {
		return nil, fmt.Errorf("block not found")
	}
	
	block := &BlockInfo{}
	
	// Parse block fields
	if v, ok := result["number"].(string); ok {
		fmt.Sscanf(v, "0x%x", &block.Number)
	}
	if v, ok := result["hash"].(string); ok {
		block.Hash = common.HexToHash(v)
	}
	if v, ok := result["parentHash"].(string); ok {
		block.ParentHash = common.HexToHash(v)
	}
	if v, ok := result["timestamp"].(string); ok {
		fmt.Sscanf(v, "0x%x", &block.Timestamp)
	}
	if v, ok := result["difficulty"].(string); ok {
		block.Difficulty = new(big.Int)
		block.Difficulty.SetString(v[2:], 16)
	}
	if v, ok := result["gasUsed"].(string); ok {
		fmt.Sscanf(v, "0x%x", &block.GasUsed)
	}
	if v, ok := result["gasLimit"].(string); ok {
		fmt.Sscanf(v, "0x%x", &block.GasLimit)
	}
	if v, ok := result["miner"].(string); ok {
		block.Coinbase = common.HexToAddress(v)
	}
	if v, ok := result["transactions"].([]interface{}); ok {
		block.TxCount = len(v)
	}
	
	block.Miner = block.Coinbase
	
	return block, nil
}

func (c *GTKMClient) GetNodeInfo() (*NodeInfo, error) {
	info := &NodeInfo{}
	
	// Get network ID
	var networkID string
	err := c.client.Call(&networkID, "net_version")
	if err == nil {
		fmt.Sscanf(networkID, "%d", &info.NetworkID)
	}
	
	// Get peer count
	var peerCount string
	err = c.client.Call(&peerCount, "net_peerCount")
	if err == nil {
		fmt.Sscanf(peerCount, "0x%x", &info.PeerCount)
	}
	
	// Check if mining
	var mining bool
	err = c.client.Call(&mining, "eth_mining")
	if err == nil {
		info.IsMining = mining
	}
	
	// Get hashrate
	var hashrate string
	err = c.client.Call(&hashrate, "eth_hashrate")
	if err == nil {
		if strings.HasPrefix(hashrate, "0x") {
			var hr uint64
			fmt.Sscanf(hashrate[2:], "%x", &hr)
			info.Hashrate = float64(hr)
		}
	}
	
	// Check syncing
	var syncing interface{}
	err = c.client.Call(&syncing, "eth_syncing")
	if err == nil {
		if syncing != false {
			info.Syncing = true
		}
	}
	
	return info, nil
}

func (c *GTKMClient) GetTransactionCount(address common.Address) (uint64, error) {
	var result string
	err := c.client.Call(&result, "eth_getTransactionCount", address.Hex(), "pending")
	if err != nil {
		return 0, err
	}
	
	var nonce uint64
	if strings.HasPrefix(result, "0x") {
		fmt.Sscanf(result[2:], "%x", &nonce)
	} else {
		fmt.Sscanf(result, "%x", &nonce)
	}
	return nonce, nil
}

func (c *GTKMClient) SendTransaction(tx Transaction) (common.Hash, error) {
	// Build transaction object for RPC
	txMap := map[string]interface{}{
		"from":  tx.From.Hex(),
		"to":    tx.To.Hex(),
		"value": fmt.Sprintf("0x%x", tx.Value),
		"gas":   fmt.Sprintf("0x%x", tx.Gas),
		"gasPrice": fmt.Sprintf("0x%x", tx.GasPrice),
		"nonce": fmt.Sprintf("0x%x", tx.Nonce),
	}
	
	if len(tx.Data) > 0 {
		txMap["data"] = "0x" + hex.EncodeToString(tx.Data)
	}
	
	var result string
	err := c.client.Call(&result, "eth_sendTransaction", txMap)
	if err != nil {
		return common.Hash{}, err
	}
	
	return common.HexToHash(result), nil
}

func (c *GTKMClient) GetTransactionReceipt(txHash common.Hash) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := c.client.Call(&result, "eth_getTransactionReceipt", txHash.Hex())
	return result, err
}

func (c *GTKMClient) GetTransactionByHash(txHash common.Hash) (*Transaction, error) {
	var result map[string]interface{}
	err := c.client.Call(&result, "eth_getTransactionByHash", txHash.Hex())
	if err != nil {
		return nil, err
	}
	
	if result == nil {
		return nil, fmt.Errorf("transaction not found")
	}
	
	tx := &Transaction{}
	tx.Hash = txHash
	
	if v, ok := result["from"].(string); ok {
		tx.From = common.HexToAddress(v)
	}
	if v, ok := result["to"].(string); ok && v != "" {
		tx.To = common.HexToAddress(v)
	}
	if v, ok := result["value"].(string); ok {
		tx.Value = new(big.Int)
		tx.Value.SetString(v[2:], 16)
	}
	if v, ok := result["gas"].(string); ok {
		fmt.Sscanf(v, "0x%x", &tx.Gas)
	}
	if v, ok := result["gasPrice"].(string); ok {
		tx.GasPrice = new(big.Int)
		tx.GasPrice.SetString(v[2:], 16)
	}
	if v, ok := result["nonce"].(string); ok {
		fmt.Sscanf(v, "0x%x", &tx.Nonce)
	}
	if v, ok := result["input"].(string); ok && v != "0x" {
		tx.Data, _ = hex.DecodeString(v[2:])
	}
	if v, ok := result["blockNumber"].(string); ok && v != "" {
		fmt.Sscanf(v, "0x%x", &tx.BlockNum)
	}
	
	// Get receipt for status
	if tx.BlockNum > 0 {
		receipt, err := c.GetTransactionReceipt(txHash)
		if err == nil {
			if status, ok := receipt["status"].(string); ok {
				fmt.Sscanf(status, "0x%x", &tx.Status)
			}
		}
	}
	
	return tx, nil
}

// ==================== COMPOSITE METHODS ====================

func (c *GTKMClient) GetFullKingInfo() (*KingInfo, error) {
	info := &KingInfo{
		AllKings: []common.Address{},
	}
	
	// Get current king
	if current, err := c.GetCurrentKing(); err == nil {
		info.CurrentKing = current
	}
	
	// Get main king
	if main, err := c.GetMainKing(); err == nil {
		info.MainKing = main
	}
	
	// Get next king
	if next, err := c.GetNextKing(); err == nil {
		info.NextKing = next
	}
	
	// Get all kings
	if all, err := c.GetKingAddresses(); err == nil {
		info.AllKings = all
	}
	
	// Get rotation info
	if height, err := c.GetBlockNumber(); err == nil {
		if rotation, err := c.GetRotationInfo(height); err == nil {
			info.RotationInfo = rotation
		}
	}
	
	// Get monitoring responsibilities
	if monitoring, err := c.GetMonitoringResponsibilities(); err == nil {
		info.MonitoringTasks = monitoring
	}
	
	return info, nil
}
