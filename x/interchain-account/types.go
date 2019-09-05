package interchain_account

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/tendermint/crypto/tmhash"
)

// This type is for transaction for interchain account module.
// Interchain tx consists of minimal data set without signatures...
// Signature is not neccessary because receiving chain can check that signers has right permission in application logic.
type ChainAccountTx struct {
	Msgs []sdk.Msg
}

func (ChainAccountTx) SourcePort() string {
	return "not implemented"
}

func (ChainAccountTx) SourceChannel() string {
	return "not implemented"
}

// IBC Ante handler checks that transaction has right permissions.
type IBCAnteHandler func(ctx sdk.Context, packet ChainAccountTx) sdk.Error

// Just mock up for IBC packet.
type Packet interface {
	GetData() []byte
	SourcePort() string
	SourceChannel() string
}

// Make an interchain account for sending chain.
// Address is determined by sha256(channel + salt).
type RegisterIBCAccountPacketData struct {
	Salt []byte
}

var _ Packet = RegisterIBCAccountPacketData{}

func (RegisterIBCAccountPacketData) GetData() []byte {
	return []byte{}
}

func (RegisterIBCAccountPacketData) SourcePort() string {
	return "not implemented"
}

func (RegisterIBCAccountPacketData) SourceChannel() string {
	return "not implemented"
}

// Return result data to sending chain.
type ResultRegisterPacketData struct {
	Address sdk.AccAddress
	Success bool
}

var _ Packet = ResultRegisterPacketData{}

func (ResultRegisterPacketData) GetData() []byte{
	return []byte{}
}

func (ResultRegisterPacketData) SourcePort() string {
	return "not implemented"
}

func (ResultRegisterPacketData) SourceChannel() string {
	return "not implemented"
}

// When IBC account module receive this packet,
// first, decode(unmarshal) tx bytes by their codec.
// In this example, it will be decodec to `ChainAccountTx`
// Then, check the permission by ibc ante handler.
// And, execute msgs.
// Finally, return result to sending chain.
type RunTxPacketData struct {
	TxBytes []byte
}

var _ Packet = RunTxPacketData{}

func (RunTxPacketData) GetData() []byte {
	return []byte{}
}

func (RunTxPacketData) SourcePort() string {
	return "not implemented"
}

func (RunTxPacketData) SourceChannel() string {
	return "not implemented"
}

func (packet RunTxPacketData) Hash() []byte {
	return tmhash.Sum(packet.TxBytes)
}

// Result of interchain account tx.
type ResultRunTxPacketData struct {
	Hash []byte // hash of tx
	Code sdk.CodeType // response code
	Data []byte // any data returned from the app
}

var _ Packet = ResultRunTxPacketData{}

func (ResultRunTxPacketData) GetData() []byte {
	return []byte{}
}

func (ResultRunTxPacketData) SourcePort() string {
	return "not implemented"
}

func (ResultRunTxPacketData) SourceChannel() string {
	return "not implemented"
}
