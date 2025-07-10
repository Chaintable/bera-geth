// Copyright 2025 Berachain Foundation
// This file is part of the bera-geth library.
//
// The bera-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The bera-geth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the bera-geth library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// PubkeySize represents the expected byte length of a BLS12-381 public key
// as used by the beacon chain.
const PubkeySize = 48

// Pubkey represents a fixed-length 48-byte BLS public key.
// JSON and text serialization use 0x-prefixed hex strings.
type Pubkey [PubkeySize]byte

// Bytes returns a copy of the underlying byte slice.
func (p Pubkey) Bytes() []byte {
	b := make([]byte, PubkeySize)
	copy(b, p[:])
	return b
}

// String returns the hex-encoded string representation of the pubkey.
func (p Pubkey) String() string {
	return hexutil.Encode(p[:])
}

// MarshalText encodes the pubkey as a 0x-prefixed hex string.
func (p Pubkey) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

// UnmarshalText decodes a 0x-prefixed hex string into the pubkey.
func (p *Pubkey) UnmarshalText(text []byte) error {
	// Use hexutil.Bytes helper for consistent behaviour.
	var b hexutil.Bytes
	if err := b.UnmarshalText(text); err != nil {
		return err
	}
	if len(b) != PubkeySize {
		return fmt.Errorf("invalid pubkey length: expected %d bytes, got %d", PubkeySize, len(b))
	}
	copy(p[:], b)
	return nil
}

// MarshalJSON encodes the pubkey as a JSON string.
func (p Pubkey) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

// UnmarshalJSON decodes a JSON string containing the 0x-prefixed hex pubkey.
func (p *Pubkey) UnmarshalJSON(input []byte) error {
	var s string
	if err := json.Unmarshal(input, &s); err != nil {
		return err
	}
	return p.UnmarshalText([]byte(s))
}

// copyPubkeyPtr copies a pubkey.
func copyPubkeyPtr(p *Pubkey) *Pubkey {
	if p == nil {
		return nil
	}
	cpy := *p
	return &cpy
}
