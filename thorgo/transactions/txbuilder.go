package transactions

import (
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

	for _, clause := range clauses {
		builder.Clause(clause)
	}

	return &Builder{
		client:  client,
		clauses: clauses,
		builder: builder,
		caller:  caller,
	}
}

// Gas sets the gas provision for the transaction. If not set, it will be estimated.
func (tx *Builder) Gas(gas uint64) *Builder {
	tx.builder.Gas(gas)
	return tx
}

// GasPriceCoef sets the gas price coefficient. Defaults to 0 if not set.
func (tx *Builder) GasPriceCoef(coef uint8) *Builder {
	tx.builder.GasPriceCoef(coef)
	return tx
}

// Expiration sets the expiration block count. Defaults to 30 blocks (5 minutes) if not set.
func (tx *Builder) Expiration(exp uint32) *Builder {
	tx.builder.Expiration(exp)
	return tx
}

// Nonce sets the nonce. Defaults to a random value if not set.
func (tx *Builder) Nonce(nonce uint64) *Builder {
	tx.builder.Nonce(nonce)
	return tx
}

// BlockRef sets the block reference. Defaults to the "best" block reference if not set.
func (tx *Builder) BlockRef(br transaction.BlockRef) *Builder {
	tx.builder.BlockRef(br)
	return tx
}

// DependsOn sets the dependent transaction ID. Defaults to nil if not set.
func (tx *Builder) DependsOn(txID *common.Hash) *Builder {
	tx.builder.DependsOn(txID)
	return tx
}

// Delegated enables transaction delegation. If not set, it will be nil.
func (tx *Builder) Delegated() *Builder {
	tx.builder.Features(transaction.DelegationFeature)
	return tx
}

// Simulate estimates the gas usage and checks for errors or reversion in the transaction.
func (tx *Builder) Simulate() (Simulation, error) {
	request := client.InspectRequest{
		Clauses: tx.clauses,
		Caller:  &tx.caller,
	}

	response, err := tx.client.Inspect(request)
	if err != nil {
		return Simulation{}, err
	}

	inspection := *response
	lastResult := inspection[len(inspection)-1]

	var consumedGas uint64
	for _, res := range inspection {
		consumedGas += res.GasUsed
	}

	intrinsicGas, err := transaction.IntrinsicGas(tx.clauses...)
	if err != nil {
		return Simulation{}, err
	}

	return Simulation{
		consumedGas:  consumedGas,
		vmError:      lastResult.VmError,
		reverted:     lastResult.Reverted,
		outputs:      inspection,
		intrinsicGas: intrinsicGas,
	}, nil
}

// Build constructs the transaction, applying defaults where necessary.
func (tx *Builder) Build() (*transaction.Transaction, error) {
	unsanitized := tx.builder.Build()
	chainTag := tx.client.ChainTag()

	builder := new(transaction.Builder).
		GasPriceCoef(unsanitized.GasPriceCoef()).
		ChainTag(chainTag).
		Features(unsanitized.Features()).
		DependsOn(unsanitized.DependsOn())

	for _, clause := range unsanitized.Clauses() {
		builder.Clause(clause)
	}

	// Set gas
	if unsanitized.Gas() == 0 {
		simulation, err := tx.Simulate()
		if err != nil {
			return nil, err
		}
		builder.Gas(simulation.TotalGas())
	} else {
		builder.Gas(unsanitized.Gas())
	}

	// Set block reference
	if unsanitized.BlockRef().Number() == 0 {
		best, err := tx.client.BestBlock()
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
