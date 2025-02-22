package rpc

import (
	"fmt"
	"testing"
)

func TestDeriveXRPAddress(t *testing.T) {

	// Example mnemonic (you can use a valid BIP39 mnemonic)
	mnemonic := "ritual about elephant exotic melt tool emotion onion brother need bike coral"

	privateKeyHex, publicKeyHex, xrpAddress, err := GenerateAddress(mnemonic)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	fmt.Println("privateKeyHex", privateKeyHex)
	fmt.Println("publicKeyHex", publicKeyHex)
	fmt.Println("xrpAddress", xrpAddress)
	// Check if the private key, public key, and XRP address are not empty
	if privateKeyHex == "" {
		t.Error("Expected non-empty private key")
	}
	if publicKeyHex == "" {
		t.Error("Expected non-empty public key")
	}
	if xrpAddress == "" {
		t.Error("Expected non-empty XRP address")
	}

	// Optionally, you can add more specific checks for the expected values
	// For example, you can check the length of the XRP address
	if len(xrpAddress) != 34 {
		t.Errorf("Expected XRP address length of 34, got %d", len(xrpAddress))
	}
	panic("test")
}
