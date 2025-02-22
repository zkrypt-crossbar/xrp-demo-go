package rpc

import (
	"encoding/hex"
	"fmt"

	bip32 "github.com/tyler-smith/go-bip32"
	bip39 "github.com/tyler-smith/go-bip39"
	"github.com/zkrypt-crossbar/ripple_skd/crypto"
	"github.com/zkrypt-crossbar/ripple_skd/data"
)

// Derivation path for XRP Ledger (BIP-44 standard)
const (
	XRP_BIP44_PATH = "m/44'/144'/0'/0/0"
	HardenedOffset = 0x80000000
)

// GenerateAddress generates a new XRP address from a mnemonic.
func GenerateAddress(mnemonic string) (privKey, pubKey, address string, err error) {
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return "", "", "", err
	}

	hdKey, err := deriveHDKey(masterKey)
	if err != nil {
		return "", "", "", err
	}

	privateKeyBytes := hdKey.Key[:16]
	key, err := crypto.NewECDSAKey(privateKeyBytes)
	if err != nil {
		return
	}

	addressHash, err := crypto.AccountId(key, nil)
	if err != nil {
		return
	}

	address = addressHash.String()
	privKey = hex.EncodeToString(key.Private(nil))
	pubKey = hex.EncodeToString(key.Public(nil))

	return
}

// deriveHDKey derives the HD key using the BIP-44 path.
func deriveHDKey(masterKey *bip32.Key) (*bip32.Key, error) {
	path := []uint32{
		HardenedOffset + 44,  // Purpose
		HardenedOffset + 144, // Coin (XRP)
		HardenedOffset + 0,   // Account 0
		0,                    // External chain
		0,                    // First address
	}

	hdKey := masterKey
	for _, childIndex := range path {
		var err error
		hdKey, err = hdKey.NewChildKey(childIndex)
		if err != nil {
			return nil, err
		}
	}
	return hdKey, nil
}

// CreateECDSATx creates and signs an ECDSA transaction.
func CreateECDSATx(from, to, currency, value, fee, privateKey, pubKey string, accountSequence, lastLedgerSequence uint32) (string, error) {
	fromAccount, err := data.NewAccountFromAddress(from)
	if err != nil {
		return "", err
	}

	toAccount, err := data.NewAccountFromAddress(to)
	if err != nil {
		return "", err
	}

	transferAmount, err := createTransferAmount(value, currency)
	if err != nil {
		return "", err
	}

	feeAmount, err := data.NewValue(fee, true)
	if err != nil {
		return "", err
	}

	txnBase := data.TxBase{
		TransactionType:    data.PAYMENT,
		Account:            *fromAccount,
		Sequence:           accountSequence,
		Fee:                *feeAmount,
		LastLedgerSequence: &lastLedgerSequence,
	}

	payment := &data.Payment{
		TxBase:      txnBase,
		Destination: *toAccount,
		Amount:      *transferAmount,
	}

	// txBlob, err := SignOffline(payment, privateKey)
	txBlob, err := CustomSignOffline(payment, privateKey, pubKey)
	if err != nil {
		return "", err
	}
	return txBlob, nil
}

// createTransferAmount creates a transfer amount based on the currency.
func createTransferAmount(value, currency string) (*data.Amount, error) {
	tmpAmount := value
	if currency != "" {
		tmpAmount += "/" + currency
	}
	return data.NewAmount(tmpAmount)
}

// SignOffline signs the payment transaction offline.
func SignOffline(payment *data.Payment, privateKey string) (string, error) {
	privateKeyBytes, _ := hex.DecodeString(privateKey)
	key := crypto.LoadECDSAKey(privateKeyBytes)

	err := data.Sign(payment, key, nil)
	if err != nil {
		return "", err
	}

	return MakeTxBlob(payment)
}

// Custom SignOffline
func CustomSignOffline(payment *data.Payment, privateKey string, publicKey string) (string, error) {
	// Convert public key to bytes
	publicKeyBytes, _ := hex.DecodeString(publicKey)

	// Get message hash
	hash, msg, err := data.GetMessageHash(payment, publicKeyBytes, nil)
	if err != nil {
		return "", err
	}
	fmt.Println("HASHED MESSAGE: ", hex.EncodeToString(hash.Bytes()))
	fmt.Println("RAW MESSAGE: ", hex.EncodeToString(msg))

	// Custom sign message
	privateKeyBytes, _ := hex.DecodeString(privateKey)
	key := crypto.LoadECDSAKey(privateKeyBytes)

	sig, err := crypto.SignECDSA(key.Private(nil), hash.Bytes())
	if err != nil {
		return "", err
	}

	*payment.GetSignature() = sig
	hash, _, err = data.Raw(payment)
	if err != nil {
		return "", err
	}
	copy(payment.GetHash().Bytes(), hash.Bytes())
	return MakeTxBlob(payment)
}

// MakeTxBlob creates a transaction blob from the payment.
func MakeTxBlob(payment *data.Payment) (string, error) {
	_, raw, err := data.Raw(data.Transaction(payment))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%X", raw), nil
}
