package main

import (
	"fmt"
	"log"

	"github.com/aviate-labs/agent-go"
	"github.com/aviate-labs/agent-go-example/archive"
	"github.com/aviate-labs/agent-go-example/ledger"
	"github.com/aviate-labs/agent-go/principal"
)

var LEDGER_PRINCIPAL = principal.MustDecode("ryjl3-tyaaa-aaaaa-aaaba-cai")

func main() {
	ledgerAgent, err := ledger.NewAgent(LEDGER_PRINCIPAL, agent.DefaultConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Query the first block.
	args := ledger.GetBlocksArgs{
		Start:  0,
		Length: 1,
	}
	response, err := ledgerAgent.QueryBlocks(args)
	if err != nil {
		log.Fatal(err)
	}

	archivedBlock := response.ArchivedBlocks[0]

	// We can either re-use the ledger agent here, or create an actual (generated) archive agent.
	var result archive.GetBlocksResult
	if err := ledgerAgent.Query(
		archivedBlock.Callback.Method.Principal,
		archivedBlock.Callback.Method.Method,
		[]any{args},
		[]any{&result},
	); err != nil {
		log.Fatal(err)
	}

	genesisBlock := result.Ok.Blocks[0].Transaction.Operation.Mint

	var to principal.AccountIdentifier
	copy(to[:], genesisBlock.To[4:])

	fmt.Printf("Block %d: %s: %.2f ICP minted.\n", 0, to, float64(genesisBlock.Amount.E8s)/1e8)

	// NOTE that the blocks still use account identifier, not ICRC1 accounts.
	fmt.Println(ledgerAgent.AccountBalance(ledger.AccountIdentifierByteBuf{
		Account: to.Bytes(),
	}))
}
