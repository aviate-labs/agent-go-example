// Package archive provides a client for the "archive" canister.
// Do NOT edit this file. It was automatically generated by https://github.com/aviate-labs/agent-go.
package archive

import (
	"github.com/aviate-labs/agent-go"
	"github.com/aviate-labs/agent-go/candid/idl"
	"github.com/aviate-labs/agent-go/principal"
)

type BlockIndex = uint64

type Memo = uint64

type AccountIdentifier = []byte

type Tokens = struct {
	E8s uint64 `ic:"e8s"`
}

type Timestamp = struct {
	TimestampNanos uint64 `ic:"timestamp_nanos"`
}

type Operation = struct {
	Mint *struct {
		To     AccountIdentifier `ic:"to"`
		Amount Tokens            `ic:"amount"`
	} `ic:"Mint,variant"`
	Burn *struct {
		From    AccountIdentifier  `ic:"from"`
		Spender *AccountIdentifier `ic:"spender,omitempty"`
		Amount  Tokens             `ic:"amount"`
	} `ic:"Burn,variant"`
	Transfer *struct {
		From    AccountIdentifier `ic:"from"`
		To      AccountIdentifier `ic:"to"`
		Amount  Tokens            `ic:"amount"`
		Fee     Tokens            `ic:"fee"`
		Spender *[]uint8          `ic:"spender,omitempty"`
	} `ic:"Transfer,variant"`
	Approve *struct {
		From              AccountIdentifier `ic:"from"`
		Spender           AccountIdentifier `ic:"spender"`
		AllowanceE8s      idl.Int           `ic:"allowance_e8s"`
		Allowance         Tokens            `ic:"allowance"`
		Fee               Tokens            `ic:"fee"`
		ExpiresAt         *Timestamp        `ic:"expires_at,omitempty"`
		ExpectedAllowance *Tokens           `ic:"expected_allowance,omitempty"`
	} `ic:"Approve,variant"`
}

type Transaction = struct {
	Memo          Memo       `ic:"memo"`
	Icrc1Memo     *[]byte    `ic:"icrc1_memo,omitempty"`
	Operation     *Operation `ic:"operation,omitempty"`
	CreatedAtTime Timestamp  `ic:"created_at_time"`
}

type Block = struct {
	ParentHash  *[]byte     `ic:"parent_hash,omitempty"`
	Transaction Transaction `ic:"transaction"`
	Timestamp   Timestamp   `ic:"timestamp"`
}

type GetBlocksArgs = struct {
	Start  BlockIndex `ic:"start"`
	Length uint64     `ic:"length"`
}

type BlockRange = struct {
	Blocks []Block `ic:"blocks"`
}

type GetBlocksError = struct {
	BadFirstBlockIndex *struct {
		RequestedIndex  BlockIndex `ic:"requested_index"`
		FirstValidIndex BlockIndex `ic:"first_valid_index"`
	} `ic:"BadFirstBlockIndex,variant"`
	Other *struct {
		ErrorCode    uint64 `ic:"error_code"`
		ErrorMessage string `ic:"error_message"`
	} `ic:"Other,variant"`
}

type GetBlocksResult = struct {
	Ok  *BlockRange     `ic:"Ok,variant"`
	Err *GetBlocksError `ic:"Err,variant"`
}

type GetEncodedBlocksResult = struct {
	Ok  *[][]byte       `ic:"Ok,variant"`
	Err *GetBlocksError `ic:"Err,variant"`
}

// Agent is a client for the "archive" canister.
type Agent struct {
	a          *agent.Agent
	canisterId principal.Principal
}

// NewAgent creates a new agent for the "archive" canister.
func NewAgent(canisterId principal.Principal, config agent.Config) (*Agent, error) {
	a, err := agent.New(config)
	if err != nil {
		return nil, err
	}
	return &Agent{
		a:          a,
		canisterId: canisterId,
	}, nil
}

// GetBlocks calls the "get_blocks" method on the "archive" canister.
func (a Agent) GetBlocks(arg0 GetBlocksArgs) (*GetBlocksResult, error) {
	var r0 GetBlocksResult
	if err := a.a.Query(
		a.canisterId,
		"get_blocks",
		[]any{arg0},
		[]any{&r0},
	); err != nil {
		return nil, err
	}
	return &r0, nil
}

// GetEncodedBlocks calls the "get_encoded_blocks" method on the "archive" canister.
func (a Agent) GetEncodedBlocks(arg0 GetBlocksArgs) (*GetEncodedBlocksResult, error) {
	var r0 GetEncodedBlocksResult
	if err := a.a.Query(
		a.canisterId,
		"get_encoded_blocks",
		[]any{arg0},
		[]any{&r0},
	); err != nil {
		return nil, err
	}
	return &r0, nil
}