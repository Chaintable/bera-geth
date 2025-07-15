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
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
)

// PoLTx implements the TxData interface.
var _ TxData = (*PoLTx)(nil)

// PoLTx represents an BRIP-0004 transaction. No gas is consumed for execution.
type PoLTx struct {
	ChainID  *big.Int
	To       common.Address // address of the PoL Distributor contract
	Nonce    uint64         // block number distributing for
	GasLimit uint64         // artificial gas limit for the PoL tx, not consumed against the block gas limit
	Data     []byte         // encodes the pubkey distributing for
}

// NewPoLTx creates a new PoL transaction.
func NewPoLTx(
	chainID *big.Int,
	distributorAddress common.Address,
	blockNumber *big.Int,
	gasLimit uint64,
	pubkey *Pubkey,
) (*Transaction, error) {
	data, err := getDistributeForData(pubkey)
	if err != nil {
		return nil, err
	}
	return NewTx(&PoLTx{
		ChainID:  chainID,
		To:       distributorAddress,
		Nonce:    blockNumber.Uint64(),
		GasLimit: gasLimit,
		Data:     data,
	}), nil
}

func (*PoLTx) txType() byte { return PoLTxType }

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *PoLTx) copy() TxData {
	cpy := &PoLTx{
		ChainID:  new(big.Int),
		To:       tx.To,
		Nonce:    tx.Nonce,
		GasLimit: tx.GasLimit,
		Data:     common.CopyBytes(tx.Data),
	}
	if tx.ChainID != nil {
		cpy.ChainID.Set(tx.ChainID)
	}
	return cpy
}

func (tx *PoLTx) chainID() *big.Int   { return tx.ChainID }
func (*PoLTx) accessList() AccessList { return nil }
func (tx *PoLTx) data() []byte        { return tx.Data }
func (tx *PoLTx) gas() uint64         { return tx.GasLimit }
func (*PoLTx) gasPrice() *big.Int     { return new(big.Int) }
func (*PoLTx) gasTipCap() *big.Int    { return new(big.Int) }
func (*PoLTx) gasFeeCap() *big.Int    { return new(big.Int) }
func (*PoLTx) value() *big.Int        { return new(big.Int) }
func (tx *PoLTx) nonce() uint64       { return tx.Nonce }
func (tx *PoLTx) to() *common.Address { return &tx.To }

// No-op: PoLTx is originated from the system address and carries no signature.
func (*PoLTx) rawSignatureValues() (v, r, s *big.Int) {
	return nil, nil, nil
}

func (*PoLTx) setSignatureValues(chainID, v, r, s *big.Int) {
	// No-op: PoLTx is originated from the system address and carries no signature.
}

// effectiveGasPrice is a no-op for PoLTx. PoLTx does not pay for gas.
func (*PoLTx) effectiveGasPrice(dst *big.Int, baseFee *big.Int) *big.Int {
	return dst.SetUint64(0)
}

func (tx *PoLTx) encode(b *bytes.Buffer) error {
	return rlp.Encode(b, tx)
}

func (tx *PoLTx) decode(input []byte) error {
	return rlp.DecodeBytes(input, tx)
}

func (tx *PoLTx) sigHash(chainID *big.Int) common.Hash {
	return prefixedRlpHash(
		PoLTxType, // tx type: 0x7D
		[]any{
			chainID,              // chainID: EIP-155 chain ID
			params.SystemAddress, // from = system address
			tx.To,                // to = address of the PoL Distributor contract
			tx.Nonce,             // nonce = block number distributing for
			tx.GasLimit,          // gasLimit = artificial gas limit for execution
			tx.Data,              // data ~= pubkey distributing for
		})
}

var (
	bytesType, _        = abi.NewType("bytes", "", nil)
	distributeForMethod = abi.NewMethod(
		"distributeFor", "distributeFor", abi.Function, "nonpayable", false, false, []abi.Argument{
			{Name: "pubkey", Type: bytesType, Indexed: false},
		}, nil,
	)
)

// getDistributeForData returns the tx data for the `distributeFor(bytes pubkey)` method.
func getDistributeForData(pubkey *Pubkey) ([]byte, error) {
	arguments, err := distributeForMethod.Inputs.Pack(pubkey.Bytes())
	if err != nil {
		return nil, err
	}
	return append(distributeForMethod.ID, arguments...), nil
}

// isDistributeForCall returns true if the provided calldata corresponds to a
// call to the `distributeFor(bytes pubkey)` method defined in BRIP-0004.
//
// The function checks that the first four bytes (the function selector) match
// the ID of the `distributeFor` ABI method declared in tx_pol.go.
func isDistributeForCall(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	return bytes.Equal(data[:4], distributeForMethod.ID)
}

// IsPoLDistribution returns true if the transaction is a PoL distribution.
func IsPoLDistribution(to *common.Address, data []byte, distributorAddress common.Address) bool {
	// Txs that call the `distributeFor(bytes pubkey)` method on the PoL Distributor
	// contract are also consideredPoL txs.
	return to != nil && *to == distributorAddress && isDistributeForCall(data)
}
