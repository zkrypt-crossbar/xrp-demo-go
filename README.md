# XRP Demo Go

A simple Go application for interacting with the XRP Ledger. This project demonstrates how to generate XRP addresses, create transactions, and interact with the XRP network using the Ripple API.

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [Commands](#commands)
- [Contributing](#contributing)
- [License](#license)

## Features

- Generate XRP addresses from mnemonics.
- Create and sign transactions.
- Fetch account information and server state.
- Submit transactions to the XRP Ledger.

## Requirements

- Go 1.16 or later
- Access to the XRP Testnet or Mainnet

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/xrp-demo-go.git
   cd xrp-demo-go
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

## Usage

To run the application, use the following command:

```bash
go run main.go <command>
```

Replace `<command>` with one of the following:

- `genkey`: Generate a new wallet (private key, public key and address) from a mnemonic
- `transfer`: Transfer XRP between accounts
- `account`: Query account information

To faucet XRP, visit the following link:

- [XRP Faucet](https://xrpl.org/resources/dev-tools/xrp-faucets)

To prepare the fund for the transaction, you can use the following link:

- [XRP Sender](https://xrpl.org/resources/dev-tools/tx-sender)

Check the balance of the account:

- [XRP Explorer](https://testnet.xrpl.org/)

Check the raw transaction:

- [XRP Decoder](https://fluxw42.github.io/ripple-tx-decoder/)

## Acknowledgements

- [github.com/rubblelabs/ripple/crypto](https://github.com/rubblelabs/ripple/crypto)
- [github.com/rubblelabs/ripple/data](https://github.com/rubblelabs/ripple/data)
- [github.com/rubblelabs/ripple/rpc](https://github.com/rubblelabs/ripple/rpc)
