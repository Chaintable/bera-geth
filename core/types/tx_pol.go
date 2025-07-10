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

// PoLTx represents an BRIP-0004 transaction.
type PoLTx struct {
	ChainID     *big.Int
	Pubkey      *Pubkey
	To          *common.Address // address of the PoL Distributor contract.
	BlockNumber uint64
}

// ===== TxData interface implementation =====

func (*PoLTx) txType() byte { return PoLTxType }

func (tx *PoLTx) copy() TxData {
	return &PoLTx{
		ChainID:     tx.ChainID,
		Pubkey:      tx.Pubkey,
		To:          tx.To,
		BlockNumber: tx.BlockNumber,
	}
}

func (*PoLTx) chainID() *big.Int      { return big.NewInt(0) }
func (*PoLTx) accessList() AccessList { return nil }
func (tx *PoLTx) data() []byte        { return mustGetDistributeForData(tx.Pubkey) }
func (*PoLTx) gas() uint64            { return 0 }
func (*PoLTx) gasPrice() *big.Int     { return new(big.Int) }
func (*PoLTx) gasTipCap() *big.Int    { return new(big.Int) }
func (*PoLTx) gasFeeCap() *big.Int    { return new(big.Int) }
func (*PoLTx) value() *big.Int        { return new(big.Int) }
func (tx *PoLTx) nonce() uint64       { return tx.BlockNumber }
func (tx *PoLTx) to() *common.Address { return tx.To }

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

// TODO: make like DynamicFeeTx.sigHash
func (tx *PoLTx) sigHash(chainID *big.Int) common.Hash {
	return common.Hash{}
}

// mustGetDistributeForData returns the tx data for the distributeFor method.
func mustGetDistributeForData(pubkey *Pubkey) []byte {
	bytesT, err := abi.NewType("bytes", "", nil)
	if err != nil {
		panic(err)
	}

	distributeForMethod := abi.NewMethod(
		"distributeFor", "distributeFor", abi.Function, "nonpayable", false, false, []abi.Argument{
			{Name: "pubkey", Type: bytesT, Indexed: false},
		}, nil,
	)

	data, err := distributeForMethod.Inputs.Pack(pubkey)
	if err != nil {
		panic(err)
	}
	return data
}
