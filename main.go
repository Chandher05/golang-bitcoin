package main

import (
	"fmt"
	"log"
	"os"

	"github.com/btcsuite/btcd/btcutil"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/joho/godotenv"
)

// var CustomSignetParams = chaincfg.Params{
// 	Name:        "customsignet",
// 	Net:         0xabcdef12, // Replace with your custom Signet magic value
// 	DefaultPort: "38332",
// 	DNSSeeds:    []chaincfg.DNSSeed{},
// 	// Address encoding magics
// 	HDCoinType:     1,
// }

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Get the RPC connection details from environment variables
	rpcHost := os.Getenv("BITCOIN_IP")
	rpcPort := os.Getenv("BITCOIN_RPC_PORT")
	rpcUser := os.Getenv("BITCOIN_USERNAME")
	rpcPassword := os.Getenv("BITCOIN_PASSWORD")

	// Create the RPC connection configuration
	connCfg := &rpcclient.ConnConfig{
		Host:         fmt.Sprintf("%s:%s", rpcHost, rpcPort),
		User:         rpcUser,
		Pass:         rpcPassword,
		HTTPPostMode: true,  // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true,  // Bitcoin core does not provide TLS by default
		// Params: "", 
	}

	// Connect to the Bitcoin RPC server
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatalf("Error creating new Bitcoin RPC client: %v", err)
	}
	defer client.Shutdown()


	blockchainInfo, err := client.GetBlockChainInfo()
	if err != nil {
		log.Fatalf("Error getting blockchain info: %v", err)
	}

	// Get the latest block hash
	blockHash, err := client.GetBlockHash(int64(blockchainInfo.Blocks))
	if err != nil {
		log.Fatalf("Error getting block hash: %v", err)
	}

	// Get the latest block
	block, err := client.GetBlock(blockHash)
	if err != nil {
		log.Fatalf("Error getting block: %v", err)
	}

	// Print the latest block information
	fmt.Printf("Latest Block Hash: %s\n", blockHash)
	fmt.Printf("Block Height: %d\n", blockchainInfo.Blocks)
	fmt.Printf("Block Version: %d\n", block.Header.Version)
	fmt.Printf("Previous Block Hash: %s\n", block.Header.PrevBlock)
	fmt.Printf("Merkle Root: %s\n", block.Header.MerkleRoot)
	fmt.Printf("Timestamp: %s\n", block.Header.Timestamp)
	fmt.Printf("Bits: %d\n", block.Header.Bits)
	fmt.Printf("Nonce: %d\n", block.Header.Nonce)
	fmt.Printf("Transaction Count: %d\n", len(block.Transactions))

	txHashStr := "9d817a0bf45548dc4879d6dba0173d53617ad0d250e202b5160fd9fcbf330b67" // Replace with the actual transaction hash
	txHash, err := chainhash.NewHashFromStr(txHashStr)
	if err != nil {
		log.Fatalf("Invalid transaction hash: %v", err)
	}

	// Get the raw transaction details
	tx, err := client.GetRawTransactionVerbose(txHash)
	if err != nil {
		log.Fatalf("Error getting raw transaction: %v", err)
	}

	// Print transaction details
	fmt.Printf("Transaction Hash: %s\n", tx.Hash)
	fmt.Printf("Confirmations: %d\n", tx.Confirmations)
	fmt.Printf("Block Hash: %s\n", tx.BlockHash)
	fmt.Printf("Block Time: %d\n", tx.Blocktime)
	fmt.Printf("Transaction Time: %d\n", tx.Time)
	fmt.Printf("Inputs:\n")
	for _, vin := range tx.Vin {
		fmt.Printf("  TxID: %s, Vout: %d, ScriptSig: %s\n", vin.Txid, vin.Vout, vin.ScriptSig)
	}
	fmt.Printf("Outputs:\n")
	for _, vout := range tx.Vout {
		fmt.Printf("  Value: %f, N: %d, ScriptPubKey: %s\n", vout.Value, vout.N, vout.ScriptPubKey)
	}

	address := "tb1q50p0jp9hdyzt3k7ep3r50hxaaqucm4dpnys08t"
	addr, err := btcutil.DecodeAddress(address, &chaincfg.MainNetParams)
	if err != nil {
		log.Fatalf("Invalid address: %v", err)
	}

	// Get the unspent transaction outputs (UTXOs) for the address
	utxos, err := client.ListUnspentMinMaxAddresses(1, 9999999, []btcutil.Address{addr})
	if err != nil {
		log.Fatalf("Error listing unspent transactions: %v", err)
	}

	// Calculate the total balance
	var totalBalance btcutil.Amount
	for _, utxo := range utxos {
		totalBalance += btcutil.Amount(utxo.Amount * btcutil.SatoshiPerBitcoin)
	}

	// Print the balance
	fmt.Printf("Balance for address %s: %v BTC\n", address, totalBalance.ToBTC())
}