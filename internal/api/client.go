package api

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// Client wraps both the raw RPC client and the ethclient.
type Client struct {
	RPC *rpc.Client
	ETH *ethclient.Client
	// Verified chain ID (nil if verification not done)
	ChainID *big.Int
}

// DialAndVerify dials a node and verifies the connection by fetching the chain ID.
func DialAndVerify(ctx context.Context, rawURL string) (*Client, error) {
	// basic URL sanity check (allow http, https, ws, wss)
	if _, err := url.ParseRequestURI(rawURL); err != nil {
		return nil, fmt.Errorf("invalid node URL: %w", err)
	}

	// Prefer rpc.DialContext to allow more transport choices and to be able to close the rpc client cleanly.
	rpcClient, err := rpc.DialContext(ctx, rawURL)
	if err != nil {
		return nil, fmt.Errorf("rpc dial: %w", err)
	}

	ethClient := ethclient.NewClient(rpcClient)

	// verify by requesting chain id / network id
	chainCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	chainID, err := ethClient.ChainID(chainCtx)
	if err != nil {
		// close acquired clients on error
		_ = rpcClient.Close()
		return nil, fmt.Errorf("failed to verify chain id: %w", err)
	}

	return &Client{
		RPC:     rpcClient,
		ETH:     ethClient,
		ChainID: chainID,
	}, nil
}

// NewClientWithRetry tries to dial and verify with retries and interval.
func NewClientWithRetry(rawURL string, dialTimeout time.Duration, attempts int, interval time.Duration) (*Client, error) {
	var lastErr error
	deadline := time.Now().Add(dialTimeout * time.Duration(attempts))
	// If attempts <= 0, attempt once
	if attempts <= 0 {
		attempts = 1
	}
	for i := 0; i < attempts; i++ {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			remaining = dialTimeout
		}
		ctx, cancel := context.WithTimeout(context.Background(), remaining)
		client, err := DialAndVerify(ctx, rawURL)
		cancel()
		if err == nil {
			return client, nil
		}
		lastErr = err
		// wait before retrying, but stop early if next attempt would exceed deadline
		time.Sleep(interval)
	}
	if lastErr == nil {
		lastErr = errors.New("unknown error dialing node")
	}
	return nil, fmt.Errorf("all attempts failed: %w", lastErr)
}

// Close closes underlying clients. Safe to call multiple times.
func (c *Client) Close() {
	if c == nil {
		return
	}
	if c.ETH != nil {
		_ = c.ETH.Close()
	}
	if c.RPC != nil {
		_ = c.RPC.Close()
	}
}
