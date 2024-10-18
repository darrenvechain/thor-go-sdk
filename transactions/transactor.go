package transactions

import (
	"fmt"
	"math"

	"github.com/darrenvechain/thorgo/client"
	"github.com/darrenvechain/thorgo/crypto/transaction"
	"github.com/ethereum/go-ethereum/common"
)

// Transactor is a transaction builder that can be used to simulate, build and send transactions.
type Transactor struct {
	client   *client.Client
	clauses  []*transaction.Clause
	builder  *transaction.Builder
	gasPayer *common.Address
}

func NewTransactor(client *client.Client, clauses []*transaction.Clause) *Transactor {
	builder := new(transaction.Builder)
	return &Transactor{
		client:  client,
		clauses: clauses,
		builder: builder,
	}
}

// GasPayer sets the gas payer for the transaction. This is used to simulate the transaction.
func (t *Transactor) GasPayer(payer common.Address) *Transactor {
	t.gasPayer = &payer
	return t
}

// Gas sets the gas provision for the transaction. If not set, it will be estimated.
func (t *Transactor) Gas(gas uint64) *Transactor {
	t.builder.Gas(gas)
	return t
}

// GasPriceCoef sets the gas price coefficient. Defaults to 0 if not set.
func (t *Transactor) GasPriceCoef(coef uint8) *Transactor {
	t.builder.GasPriceCoef(coef)
	return t
}

// Expiration sets the expiration block count. Defaults to 30 blocks (5 minutes) if not set.
func (t *Transactor) Expiration(exp uint32) *Transactor {
	t.builder.Expiration(exp)
	return t
}

// Nonce sets the nonce. Defaults to a random value if not set.
func (t *Transactor) Nonce(nonce uint64) *Transactor {
	t.builder.Nonce(nonce)
	return t
}

// BlockRef sets the block reference. Defaults to the "best" block reference if not set.
func (t *Transactor) BlockRef(br transaction.BlockRef) *Transactor {
	t.builder.BlockRef(br)
	return t
}

// DependsOn sets the dependent transaction ID. Defaults to nil if not set.
func (t *Transactor) DependsOn(txID *common.Hash) *Transactor {
	t.builder.DependsOn(txID)
	return t
}

// Delegate enables transaction delegation. If not set, it will be nil.
func (t *Transactor) Delegate() *Transactor {
	t.builder.Features(transaction.DelegationFeature)
	return t
}

// Simulate estimates the gas usage and checks for errors or reversion in the transaction.
func (t *Transactor) Simulate(caller common.Address) (Simulation, error) {
	request := client.InspectRequest{
		Clauses: t.clauses,
		Caller:  &caller,
	}

	if t.gasPayer != nil {
		request.GasPayer = t.gasPayer
	}

	response, err := t.client.Inspect(request)
	if err != nil {
		return Simulation{}, err
	}

	lastResult := response[len(response)-1]

	var consumedGas uint64
	for _, res := range response {
		consumedGas += res.GasUsed
	}

	intrinsicGas, err := transaction.IntrinsicGas(t.clauses...)
	if err != nil {
		return Simulation{}, err
	}

	if intrinsicGas > math.MaxInt64 {
		return Simulation{}, fmt.Errorf("intrinsic gas exceeds maximum int64")
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
func (t *Transactor) Build(caller common.Address) (*transaction.Transaction, error) {
	initial := t.builder.Build()
	chainTag := t.client.ChainTag()

	builder := new(transaction.Builder).
		GasPriceCoef(initial.GasPriceCoef()).
		ChainTag(chainTag).
		Features(initial.Features()).
		DependsOn(initial.DependsOn()).
		Gas(initial.Gas()).
		BlockRef(initial.BlockRef()).
		Expiration(initial.Expiration()).
		Nonce(initial.Nonce())

	for _, clause := range t.clauses {
		builder.Clause(clause)
	}

	// Check if gas is set
	if initial.Gas() == 0 {
		simulation, err := t.Simulate(caller)
		if err != nil {
			return nil, err
		}
		builder.Gas(simulation.TotalGas())
	}

	// Check if block reference is set
	if initial.BlockRef().Number() == 0 {
		best, err := t.client.BestBlock()
		if err != nil {
			return nil, err
		}
		builder.BlockRef(best.BlockRef())
	}

	// Set expiration
	if initial.Expiration() == 0 {
		builder.Expiration(30)
	}

	// Set nonce
	if initial.Nonce() == 0 {
		builder.Nonce(transaction.Nonce())
	}

	return builder.Build(), nil
}

type Signer interface {
	SignTransaction(tx *transaction.Transaction) ([]byte, error)
	Address() common.Address
}

// Send will submit the transaction to the network.
func (t *Transactor) Send(signer Signer) (*Visitor, error) {
	tx, err := t.Build(signer.Address())
	if err != nil {
		return nil, fmt.Errorf("failed to build transaction: %w", err)
	}

	signature, err := signer.SignTransaction(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}
	tx = tx.WithSignature(signature)

	res, err := t.client.SendTransaction(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	return New(t.client, res.ID), nil
}
