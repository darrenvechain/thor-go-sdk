package transactions

import (
	"github.com/darrenvechain/thor-go-sdk/client"
)

type Simulation struct {
	consumedGas  uint64
	reverted     bool
	outputs      []client.InspectResponse
	vmError      string
	intrinsicGas uint64
}

func (s *Simulation) TotalGas() uint64 {
	return s.consumedGas + s.intrinsicGas
}

func (s *Simulation) IsSuccess() bool {
	return !s.reverted && s.vmError == ""
}

func (s *Simulation) ConsumedGas() uint64 {
	return s.consumedGas
}

func (s *Simulation) Reverted() bool {
	return s.reverted
}

func (s *Simulation) Outputs() []client.InspectResponse {
	return s.outputs
}

func (s *Simulation) VMError() string {
	return s.vmError
}

func (s *Simulation) IntrinsicGas() uint64 {
	return s.intrinsicGas
}
