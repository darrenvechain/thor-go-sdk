package transactions

import (
	"fmt"

	"github.com/darrenvechain/thor-go-sdk/client"
	"github.com/darrenvechain/thor-go-sdk/crypto/transaction"
	"github.com/ethereum/go-ethereum/common"
)

type Transactor struct {
	client   *client.Client
	clauses  []*transaction.Clause
	builder  *transaction.Builder
	caller   common.Address
	gasPayer *common.Address
}

func NewTransactor(client *client.Client, clauses []*transaction.Clause, caller common.Address) *Transactor {
	builder := new(transaction.Builder)
	return &Transactor{
		client:  client,
		clauses: clauses,
		caller:  caller,
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
func (t *Transactor) Simulate() (Simulation, error) {
	request := client.InspectRequest{
		Clauses: t.clauses,
		Caller:  &t.caller,
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

	return Simulation{
		consumedGas:  consumedGas,
		vmError:      lastResult.VmError,
		reverted:     lastResult.Reverted,
		outputs:      response,
		intrinsicGas: intrinsicGas,
	}, nil
}

// Build constructs the transaction, applying defaults where necessary.
func (t *Transactor) Build() (*transaction.Transaction, error) {
	unsanitized := t.builder.Build()
	chainTag := t.client.ChainTag()

	builder := new(transaction.Builder).
		GasPriceCoef(unsanitized.GasPriceCoef()).
		ChainTag(chainTag).
		Features(unsanitized.Features()).
		DependsOn(unsanitized.DependsOn())

	for _, clause := range t.clauses {
		builder.Clause(clause)
	}

	// Set gas
	if unsanitized.Gas() == 0 {
		simulation, err := t.Simulate()
		if err != nil {
			return nil, err
		}
		builder.Gas(simulation.TotalGas())
	} else {
		builder.Gas(unsanitized.Gas())
	}

	// Set block reference
	if unsanitized.BlockRef().Number() == 0 {
		best, err := t.client.BestBlock()
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
	SignTransaction(tx *transaction.Transaction) ([]byte, error)
}

func (t *Transactor) Send(signer TxSigner) (*Visitor, error) {
	tx, err := t.Build()
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
