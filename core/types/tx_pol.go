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
	"github.com/ethereum/go-ethereum/rlp"
)

// PoLTx implements the TxData interface.
var _ TxData = (*PoLTx)(nil)

// PoLTx represents an BRIP-0004 transaction. No gas is consumed for execution.
type PoLTx struct {
	ChainID            *big.Int
	From               *common.Address // should be the system address --> sender
	DistributorAddress *common.Address // address of the PoL Distributor contract --> to
	Pubkey             *Pubkey         // pubkey distributing for --> data
	BlockNumber        uint64          // block number distributing for --> nonce
	GasLimit           uint64          // artificial gas limit for the PoL tx --> gas
}

// NewPoLTx creates a new PoL transaction.
func NewPoLTx(
	chainID *big.Int,
	from common.Address,
	distributorAddress common.Address,
	pubkey *Pubkey,
	blockNumber *big.Int,
	gasLimit uint64,
) *Transaction {
	return NewTx(&PoLTx{
		ChainID:            chainID,
		From:               &from,
		DistributorAddress: &distributorAddress,
		Pubkey:             pubkey,
		BlockNumber:        blockNumber.Uint64(),
		GasLimit:           gasLimit,
	})
}

func (*PoLTx) txType() byte { return PoLTxType }

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *PoLTx) copy() TxData {
	cpy := &PoLTx{
		ChainID:            new(big.Int),
		From:               copyAddressPtr(tx.From),
		DistributorAddress: copyAddressPtr(tx.DistributorAddress),
		Pubkey:             copyPubkeyPtr(tx.Pubkey),
		BlockNumber:        tx.BlockNumber,
		GasLimit:           tx.GasLimit,
	}
	if tx.ChainID != nil {
		cpy.ChainID.Set(tx.ChainID)
	}
	return cpy
}

func (tx *PoLTx) chainID() *big.Int   { return tx.ChainID }
func (*PoLTx) accessList() AccessList { return nil }
func (tx *PoLTx) data() []byte        { return mustGetDistributeForData(tx.Pubkey) }
func (tx *PoLTx) gas() uint64         { return tx.GasLimit }
func (*PoLTx) gasPrice() *big.Int     { return new(big.Int) }
func (*PoLTx) gasTipCap() *big.Int    { return new(big.Int) }
func (*PoLTx) gasFeeCap() *big.Int    { return new(big.Int) }
func (*PoLTx) value() *big.Int        { return new(big.Int) }
func (tx *PoLTx) nonce() uint64       { return tx.BlockNumber }
func (tx *PoLTx) to() *common.Address { return tx.DistributorAddress }

// No-op: PoLTx is system-signed and carries no signature.
func (*PoLTx) rawSignatureValues() (v, r, s *big.Int) {
	return nil, nil, nil
}

func (*PoLTx) setSignatureValues(chainID, v, r, s *big.Int) {
	// No-op: PoLTx is system-signed and carries no signature.
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
			chainID,               // chainID: EIP-155 chain ID
			tx.From,               // from = system address
			tx.DistributorAddress, // to = address of the PoL Distributor contract
			tx.Pubkey,             // data ~= pubkey distributing for
			tx.BlockNumber,        // nonce = block number distributing for
			tx.GasLimit,           // gasLimit = artificial gas limit for execution
		})
}

// mustGetDistributeForData returns the tx data for the distributeFor method.
func mustGetDistributeForData(pubkey *Pubkey) []byte {
	bytesT, err := abi.NewType("bytes", "", nil)
	if err != nil {
		// NOTE: this should never happen.
		panic(err)
	}

	distributeForMethod := abi.NewMethod(
		"distributeFor", "distributeFor", abi.Function, "nonpayable", false, false, []abi.Argument{
			{Name: "pubkey", Type: bytesT, Indexed: false},
		}, nil,
	)

	data, err := distributeForMethod.Inputs.Pack(pubkey)
	if err != nil {
		// NOTE: this should never happen.
		panic(err)
	}
	return data
}
