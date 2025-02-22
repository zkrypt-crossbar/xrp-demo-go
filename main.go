package main

import (
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/pkg/errors"
	rpc "github.com/zkrypt-crossbar/xrp-demo-go/rpc"
)

const (
	TestnetServer = "https://s.altnet.rippletest.net:51234"
	MainnetServer = "https://s2.ripple.com:51234"

	ApiURLTestNet = "https://testnet.data.api.ripple.com"
	ApiURLMainNet = "https://data.ripple.com"
)

type XRPLClient struct {
	ServerURL string
	client    *rpc.RPCClient
}

func NewXRPLClient(testnet bool) *XRPLClient {
	if testnet {
		return &XRPLClient{
			ServerURL: TestnetServer,
			client:    rpc.NewRPCClient(TestnetServer),
		}
	}
	return &XRPLClient{
		ServerURL: MainnetServer,
		client:    rpc.NewRPCClient(MainnetServer),
	}
}

func CreateTransaction(client *XRPLClient, fromPrivateKey, fromPublicKey, fromAddress, toAddress string, amountDrops *big.Int) (string, error) {
	currency := "XRP"
	amountXRP := new(big.Float).Quo(new(big.Float).SetInt(amountDrops), big.NewFloat(1000000))

	//fetch fee
	serverInfo, err := client.client.GetServerState()
	if err != nil {
		return "", err
	}
	latestSequence := serverInfo.Result.State.ValidatedLedger.Seq + 100

	fee := serverInfo.Result.State.ValidatedLedger.BaseFee
	feeString := fmt.Sprintf("%d", fee)
	fmt.Println("ESTIMATED FEE: ", feeString)

	account, err := client.client.GetAccountInfo(fromAddress)
	if err != nil {
		return "", errors.Wrap(err, "GetAccountInfo")
	}
	accountSequence := account.Result.AccountData.Sequence

	txBlob, err := rpc.CreateECDSATx(fromAddress, toAddress, currency, amountXRP.String(), feeString, fromPrivateKey, fromPublicKey, accountSequence, latestSequence)
	if err != nil {
		return "", errors.Wrap(err, "SignECDSA")
	}

	return txBlob, nil
}

const (
	mnemonic     = "ritual about elephant exotic melt tool emotion onion brother need bike coral"
	to           = "rBzjUggmJhcDjsu6trPEmiL6ng2CrYdYXf"
	from         = "r979UwkMamWKHAEPnkqU6zih1CdLmvfjsY"
	from_private = "0db0c389d75c29fa7a50105053886fb311341cdd5bce0d7a48bf13c989981e83"
	from_pubkey  = "036c6b7d11750deb8becbad7d1f2d56203f2a439d16839a3b99dd00e267e48aa34"
)

// Example main function to demonstrate usage
func main() {
	// Check if a command is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <command>")
		fmt.Println("Available commands: genkey, transfer, account")
		os.Exit(1)
	}

	xrpClient := NewXRPLClient(true)

	// Get the command from the command line
	cmd := os.Args[1]
	switch cmd {
	case "genkey":
		privKey, pubKey, address, err := rpc.GenerateAddress(mnemonic)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("privKey", privKey)
		fmt.Println("pubKey", pubKey)
		fmt.Println("address", address)
	case "transfer":
		txBlob, err := CreateTransaction(xrpClient, from_private, from_pubkey, from, to, big.NewInt(100))
		if err != nil {
			log.Fatal(err)
		}

		submitResult, err := xrpClient.client.SubmitTransaction(txBlob)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("SUBMITTED TRANSACTION: ", submitResult)
	case "account":
		accountInfo, err := xrpClient.client.GetAccountInfo("rP1k2vWcSA7N1sCiARfHNA1y9i3fPT3pRQ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("accountInfo", accountInfo.Result)
	}
}
