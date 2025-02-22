package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/zkrypt-crossbar/xrp-demo-go/models"
)

// RPCRequest defines the JSON-RPC request format
type RPCRequest struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params,omitempty"`
}

// RPCResponse defines the JSON-RPC response format
type RPCResponse struct {
	Result json.RawMessage `json:"result"`
	Error  *RPCError       `json:"error,omitempty"`
}

// RPCError defines the JSON-RPC error format
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// RPCClient represents an instance of an RPC client
type RPCClient struct {
	URL string
}

// NewRPCClient initializes a new RPC client
func NewRPCClient(url string) *RPCClient {
	return &RPCClient{
		URL: url,
	}
}

func PrintPrettyJSON(data []byte) {
	// Convert byte array to JSON map
	var resultData map[string]interface{}
	err := json.Unmarshal(data, &resultData)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
	} else {
		// Pretty-print the JSON response
		prettyJSON, _ := json.MarshalIndent(resultData, "", "  ")
		fmt.Println("XRP Ledger Info:\n", string(prettyJSON))
	}
}

// Call sends an RPC request and returns the response
func (c *RPCClient) Call(method string, params interface{}) ([]byte, error) {
	reqData := RPCRequest{
		Method: method,
		Params: []interface{}{params},
	}

	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	return respBody, nil
}

func (c *RPCClient) GetAccountInfo(address string) (*models.AccountInfoResp, error) {
	params := map[string]interface{}{
		"account":      address,
		"ledger_index": "validated",
	}
	accountInfo, err := c.Call("account_info", params)
	if err != nil {
		return nil, err
	}

	result := &models.AccountInfoResp{}
	err = json.Unmarshal(accountInfo, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RPCClient) GetServerState() (*models.ServerInfoResp, error) {
	// call server_state
	serverState, err := c.Call("server_state", nil)
	if err != nil {
		return nil, err
	}
	// unmarshal serverState.Result to ServerInfoResp
	result := &models.ServerInfoResp{}
	err = json.Unmarshal(serverState, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RPCClient) SubmitTransaction(txBlob string) (*models.SubmitResult, error) {
	params := map[string]interface{}{
		"tx_blob": txBlob,
	}
	submitTransaction, err := c.Call("submit", params)
	if err != nil {
		return nil, err
	}
	result := &models.SubmitResp{}
	err = json.Unmarshal(submitTransaction, result)
	if err != nil {
		return nil, err
	}
	return result.Result, nil
}
