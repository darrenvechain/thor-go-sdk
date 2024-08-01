package hex

import (
	"encoding/json"
	"math/big"
)

type Int struct {
	*big.Int
}

// MarshalJSON converts the Int to a JSON-encoded hex string with 0x prefix.
func (a *Int) MarshalJSON() ([]byte, error) {
	if a == nil || a.Int == nil {
		return json.Marshal(nil)
	}

	// Convert big.Int to a hex string with 0x prefix
	hexString := "0x" + a.Int.Text(16)
	return json.Marshal(hexString)
}

// UnmarshalJSON parses a JSON-encoded hex string with 0x prefix into an Int.
func (a *Int) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	// Parse the JSON string
	var hexString string
	if err := json.Unmarshal(data, &hexString); err != nil {
		return err
	}

	// Remove the 0x prefix if present
	if len(hexString) > 2 && hexString[:2] == "0x" {
		hexString = hexString[2:]
	}

	// Parse the hex string into big.Int
	parsedInt := new(big.Int)
	_, success := parsedInt.SetString(hexString, 16)
	if !success {
		return json.Unmarshal(data, &a.Int) // Fallback in case of invalid hex string
	}

	a.Int = parsedInt
	return nil
}
