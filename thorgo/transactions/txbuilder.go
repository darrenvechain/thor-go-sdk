package transactions

import (
	"fmt"

	"github.com/darrenvechain/thor-go-sdk/client"
	"github.com/darrenvechain/thor-go-sdk/transaction"
	"github.com/ethereum/go-ethereum/common"
)

type Builder struct {
	client  *client.Client
	clauses []*transaction.Clause
	builder *transaction.Builder
	caller  common.Address
}

func NewBuilder(client *client.Client, clauses []*transaction.Clause, caller common.Address) *Builder {
	builder := new(transaction.Builder)
	return &Builder{
		client:  client,
		clauses: clauses,
		caller:  caller,
		builder: builder,
	}
}

// Gas sets the gas provision for the transaction. If not set, it will be estimated.
func (b *Builder) Gas(gas uint64) *Builder {
	b.builder.Gas(gas)
	return b
}

// GasPriceCoef sets the gas price coefficient. Defaults to 0 if not set.
func (b *Builder) GasPriceCoef(coef uint8) *Builder {
	b.builder.GasPriceCoef(coef)
	return b
}

// Expiration sets the expiration block count. Defaults to 30 blocks (5 minutes) if not set.
func (b *Builder) Expiration(exp uint32) *Builder {
	b.builder.Expiration(exp)
	return b
}

// Nonce sets the nonce. Defaults to a random value if not set.
func (b *Builder) Nonce(nonce uint64) *Builder {
	b.builder.Nonce(nonce)
	return b
}

// BlockRef sets the block reference. Defaults to the "best" block reference if not set.
func (b *Builder) BlockRef(br transaction.BlockRef) *Builder {
	b.builder.BlockRef(br)
	return b
}

// DependsOn sets the dependent transaction ID. Defaults to nil if not set.
func (b *Builder) DependsOn(txID *common.Hash) *Builder {
	b.builder.DependsOn(txID)
	return b
}

// Delegated enables transaction delegation. If not set, it will be nil.
func (b *Builder) Delegate() *Builder {
	b.builder.Features(transaction.DelegationFeature)
	return b
}

// Simulate estimates the gas usage and checks for errors or reversion in the transaction.
func (b *Builder) Simulate() (Simulation, error) {
	request := client.InspectRequest{
		Clauses: b.clauses,
		Caller:  &b.caller,
	}

	response, err := b.client.Inspect(request)
	if err != nil {
		return Simulation{}, err
	}

	lastResult := response[len(response)-1]

	var consumedGas uint64
	for _, res := range response {
		consumedGas += res.GasUsed
	}

	intrinsicGas, err := transaction.IntrinsicGas(b.clauses...)
	if err != nil {
		return Simulation{}, err
	}

	return Simulation{
		consumedGas:  consumedGas,
		vmError:      lastResult.VmError,
		reverted:     lastResult.Reverted,
		outputs:      response,
		intrinsicGas: intrinsicGas,
	}, nil
}

// Build constructs the transaction, applying defaults where necessary.
func (b *Builder) Build() (*transaction.Transaction, error) {
	unsanitized := b.builder.Build()
	chainTag := b.client.ChainTag()

	builder := new(transaction.Builder).
		GasPriceCoef(unsanitized.GasPriceCoef()).
		ChainTag(chainTag).
		Features(unsanitized.Features()).
		DependsOn(unsanitized.DependsOn())

	for _, clause := range b.clauses {
		builder.Clause(clause)
	}

	// Set gas
	if unsanitized.Gas() == 0 {
		simulation, err := b.Simulate()
		if err != nil {
			return nil, err
		}
		builder.Gas(simulation.TotalGas())
	} else {
		builder.Gas(unsanitized.Gas())
	}

	// Set block reference
	if unsanitized.BlockRef().Number() == 0 {
		best, err := b.client.BestBlock()
		if err != nil {
			return nil, err
		}
		builder.BlockRef(best.BlockRef())
	} else {
		builder.BlockRef(unsanitized.BlockRef())
	}

	// Set expiration
	if unsanitized.Expiration() == 0 {
		builder.Expiration(30)
	} else {
		builder.Expiration(unsanitized.Expiration())
	}

	// Set nonce
	if unsanitized.Nonce() == 0 {
		builder.Nonce(transaction.Nonce())
	} else {
		builder.Nonce(unsanitized.Nonce())
	}

	return builder.Build(), nil
}

type TxSigner interface {
	SignTransaction(tx *transaction.Transaction) (*transaction.Transaction, error)
}

func (b *Builder) Send(signer TxSigner) (*Visitor, error) {
	tx, err := b.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build transaction: %w", err)
	}

	tx, err = signer.SignTransaction(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	res, err := b.client.SendTransaction(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	return New(b.client, res.ID), nil
}
