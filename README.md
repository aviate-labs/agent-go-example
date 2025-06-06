# Agent Go Example: Transactions

This is an example of the [Go Agent](https://github.com/aviate-labs/agent-go) in action.

In this example we will use the agent go fetch transactions from the ICP main network. We will not be using the Rosetta
API in this example, but instead will be using the agent to fetch transactions directly. We will explore the simple CLI
that comes with the agent to generate a client.

This repository contains the end result that is described in the README. If you would like to follow along, you can
clone this repository and checkout the `start` tag.

## Setting Up Go

You can find more info on how to set up Go [here](https://golang.org/doc/install).

```shell
go mod init github.com/aviate-labs/agent-go-example
```

<details>
<summary>Tree</summary>

```text
.
├── LICENSE
├── README.md
└── go.mod
```

</details>

## Installing the Agent CLI

We can install the CLI by running the following command:

```shell
go install github.com/aviate-labs/agent-go/cmd/goic@v0.7.3
```

More info on the CLI can be found [here](https://github.com/aviate-labs/agent-go/tree/main/cmd/goic).

We can validate that the CLI is installed by running:

```shell
goic version
```

## Generating a Client

The CLI provides two ways of generating clients, either by providing a configuration file or by providing the canister
id.

We will be using the candid files that can be
found [here](https://github.com/dfinity/ic/tree/release-2023-11-01_23-01/rs/rosetta-api/icp_ledger).

### Getting the Candid Files

We can fetch the candid interfaces by running:

```shell
goic fetch ryjl3-tyaaa-aaaaa-aaaba-cai --output=ledger.did
goic fetch qjdve-lqaaa-aaaaa-aaaeq-cai --output=ledger_archive.did
```

This should have generated two files, `ledger.did` and `ledger_archive.did`.

<details>
<summary>Tree</summary>

```text
.
├── LICENSE
├── README.md
├── go.mod
├── ledger.did
└── ledger_archive.did
```

</details>

### Generating the Client

We can generate the client by running:

```shell
mkdir ledger
goic generate did ledger.did ledger --output=ledger/agent.go --packageName=ledger

mkdir archive
goic generate did ledger_archive.did archive --output=archive/agent.go --packageName=archive
```

<details>
<summary>Tree</summary>

```text
.
├── LICENSE
├── README.md
├── archive
│   └── agent.go
├── go.mod
├── ledger
│   └── agent.go
├── ledger.did
└── ledger_archive.did
```

</details>

If you open `ledger.go` and `archive.go` you should see the generated code. You will also notice that some dependencies
are missing. We can fetch these dependencies by running:

```shell
go get github.com/aviate-labs/agent-go@v0.7.3
go mod tidy
```

```text
.
├── LICENSE
├── README.md
├── archive
│   └── agent.go
├── go.mod
├── go.sum
├── ledger
│   └── agent.go
├── ledger.did
└── ledger_archive.did
```

## Fetching Transactions

Now that we have generated the client, we can use it to fetch transactions from the ICP main network.

### Creating `main.go`

Create a file called `main.go` and add the following code:

```go
package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}

```

<details>
<summary>Tree</summary>

```text
.
├── LICENSE
├── README.md
├── archive
│   └── agent.go
├── go.mod
├── go.sum
├── ledger
│   └── agent.go
├── ledger.did
├── ledger_archive.did
└── main.go
```

</details>

### Setting Up the Client

We can now set up the client by adding the following code to `main.go`:

```go
package main

import (
	"log"

	"github.com/aviate-labs/agent-go"
	"github.com/aviate-labs/agent-go-example/ledger"
	"github.com/aviate-labs/agent-go/principal"
)

var LEDGER_PRINCIPAL = principal.MustDecode("ryjl3-tyaaa-aaaaa-aaaba-cai")

func main() {
	// The default configuration is fine for most use cases, it uses the anonymous identity to create requests.
	ledgerAgent, err := ledger.NewAgent(LEDGER_PRINCIPAL, agent.DefaultConfig)
	if err != nil {
		log.Fatal(err)
	}

	_ = ledgerAgent
}

```

### Fetching The Block Height

The next step is to fetch the block height. We do not know yet what the last block is, so we query the block height
first.

```go
package main

import (
	"fmt"
	"log"

	"github.com/aviate-labs/agent-go"
	"github.com/aviate-labs/agent-go-example/ledger"
	"github.com/aviate-labs/agent-go/principal"
)

var LEDGER_PRINCIPAL = principal.MustDecode("ryjl3-tyaaa-aaaaa-aaaba-cai")

func main() {
	ledgerAgent, err := ledger.NewAgent(LEDGER_PRINCIPAL, agent.DefaultConfig)
	if err != nil {
		log.Fatal(err)
	}

	// We do not know yet what the last block is, so we query the block height first.
	blockHeight, err := ledgerAgent.QueryBlocks(ledger.GetBlocksArgs{})
	if err != nil {
		log.Fatal(err)
	}

	lastBlock := blockHeight.ChainLength // The last block on the ledger.
	fmt.Println(lastBlock)
}

```

### Fetching Last 10 Transactions

Now that we know the last block, we can fetch the last 10 blocks. Based on the operation type we can determine what kind
of transaction it is. We will only be looking at `Transfer`, `Burn` and `Mint` operations.

```go
package main

import (
	"fmt"
	"log"

	"github.com/aviate-labs/agent-go"
	"github.com/aviate-labs/agent-go-example/ledger"
	"github.com/aviate-labs/agent-go/principal"
)

var LEDGER_PRINCIPAL = principal.MustDecode("ryjl3-tyaaa-aaaaa-aaaba-cai")

func main() {
	ledgerAgent, err := ledger.NewAgent(LEDGER_PRINCIPAL, agent.DefaultConfig)
	if err != nil {
		log.Fatal(err)
	}

	// We do not know yet what the last block is, so we query the block height first.
	blockHeight, err := ledgerAgent.QueryBlocks(ledger.GetBlocksArgs{})
	if err != nil {
		log.Fatal(err)
	}

	lastBlock := blockHeight.ChainLength // The last block on the ledger.

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
			copy(from[:], transfer.From[4:])

			var to principal.AccountIdentifier
			copy(to[:], transfer.To[4:])

			fmt.Printf("Block %d: %s -> %s: %.2f ICP.\n", int(lastBlock)+i, from, to, float64(transfer.Amount.E8s)/1e8)
		} else if burn := operation.Burn; burn != nil {
			var from principal.AccountIdentifier
			copy(from[:], burn.From[4:])

			fmt.Printf("Block %d: %s: %.2f ICP burned.\n", int(lastBlock)+i, from, float64(burn.Amount.E8s)/1e8)
		} else if mint := operation.Mint; mint != nil {
			var to principal.AccountIdentifier
			copy(to[:], mint.To[4:])

			fmt.Printf("Block %d: %s: %.2f ICP minted.\n", int(lastBlock)+i, to, float64(mint.Amount.E8s)/1e8)
		}
	}
}

```

### Other Examples

<details>
<summary>ICRC-1</summary>

```
go run icrc1/main.go
```


```go
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
	copy(to[:], genesisBlock.To)

	fmt.Printf("Block %d: %s: %.2f ICP minted.\n", 0, to, float64(genesisBlock.Amount.E8s)/1e8)
}

```
</details>