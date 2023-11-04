package main

import (
	"fmt"
	"github.com/aviate-labs/agent-go"
	"github.com/aviate-labs/agent-go-example/ledger"
	"github.com/aviate-labs/agent-go/ic"
	"github.com/aviate-labs/agent-go/principal"
	"log"
)

func main() {
	ledgerAgent, err := ledger.NewAgent(ic.LEDGER_PRINCIPAL, agent.DefaultConfig)
	if err != nil {
		log.Fatal(err)
	}

	// We do not know yet what the last block is, so we query the block height first.
	blockHeight, err := ledgerAgent.QueryBlocks(ledger.GetBlocksArgs{})
	if err != nil {
		log.Fatal(err)
	}
	// oldestBlock = blockHeight.FirstBlockIndex // The first block that we can query the ledger.
	lastBlock := blockHeight.ChainLength // The last block that we can query the ledger.

	// Query the last 10 blocks.
	response, err := ledgerAgent.QueryBlocks(ledger.GetBlocksArgs{
		Start:  lastBlock - 10,
		Length: 10,
	})
	if err != nil {
		log.Fatal(err)
	}
	for i, block := range response.Blocks {
		operation := block.Transaction.Operation
		if transfer := operation.Transfer; transfer != nil {
			var from principal.AccountIdentifier
			copy(from[:], transfer.From)

			var to principal.AccountIdentifier
			copy(to[:], transfer.To)

			fmt.Printf("Block %d: %s -> %s: %.2f ICP.\n", int(lastBlock)+i, from, to, float64(transfer.Amount.E8s)/1e8)
		} else if burn := operation.Burn; burn != nil {
			var from principal.AccountIdentifier
			copy(from[:], burn.From)

			fmt.Printf("Block %d: %s: %.2f ICP burned.\n", int(lastBlock)+i, from, float64(burn.Amount.E8s)/1e8)
		} else if mint := operation.Mint; mint != nil {
			var to principal.AccountIdentifier
			copy(to[:], mint.To)

			fmt.Printf("Block %d: %s: %.2f ICP minted.\n", int(lastBlock)+i, to, float64(mint.Amount.E8s)/1e8)
		}
	}
}
