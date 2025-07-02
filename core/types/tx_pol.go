// PoLTx represents an BRIP-0004 transaction.
package types

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

// PoLTx represents an BRIP-0004 transaction.
type PoLTx struct {
	ChainID *big.Int
	Nonce   uint64

	// To is the address of the PoL Distributor contract.
	To common.Address

	// ABI encoding of "distributeFor(bytes calldata pubkey)"
	Data []byte
}

// ===== TxData interface implementation =====

func (*PoLTx) txType() byte { return PoLTxType }

func (tx *PoLTx) copy() TxData {
	return &PoLTx{
		ChainID: tx.ChainID,
		Nonce:   tx.Nonce,
		To:      tx.To,
		Data:    common.CopyBytes(tx.Data),
	}
}

func (*PoLTx) chainID() *big.Int      { return big.NewInt(0) }
func (*PoLTx) accessList() AccessList { return nil }
func (tx *PoLTx) data() []byte        { return tx.Data }
func (*PoLTx) gas() uint64            { return 0 }
func (*PoLTx) gasPrice() *big.Int     { return new(big.Int) }
func (*PoLTx) gasTipCap() *big.Int    { return new(big.Int) }
func (*PoLTx) gasFeeCap() *big.Int    { return new(big.Int) }
func (*PoLTx) value() *big.Int        { return new(big.Int) }
func (tx *PoLTx) nonce() uint64       { return tx.Nonce }
func (tx *PoLTx) to() *common.Address { return &tx.To }

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
