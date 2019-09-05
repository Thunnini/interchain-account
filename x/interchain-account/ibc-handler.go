package interchain_account

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Mock up for recv packet handler
func (keeper Keeper) OnRecvPacket(ctx sdk.Context, packet Packet) {
	switch packet := packet.(type) {
	case RegisterIBCAccountPacketData :
		addr, err := keeper.RegisterIBCAccount(ctx, packet)
		if err != nil {
			SendPacket(ResultRegisterPacketData{
				Address: []byte{},
				Success: false,
			})
		} else {
			SendPacket(ResultRegisterPacketData{
				Address: addr,
				Success: true,
			})
		}
	case RunTxPacketData:
		tx := ChainAccountTx{}
		// Unmarshal data to chain account tx.
		err := keeper.cdc.UnmarshalBinaryBare(packet.GetData(), &tx)
		if err != nil {
			// If unmashal fail, return failed data.
			SendPacket(ResultRunTxPacketData{
				Hash: packet.Hash(),
				Code: 1,
				Data: []byte("Fail to unmarshal"),
			})
			return
		}

		// Run interchain account tx.
		err = keeper.RunTxIBCTx(ctx, tx)
		if err != nil {
			// If tx fail, return failed data.
			SendPacket(ResultRunTxPacketData{
				Hash: packet.Hash(),
				Code: 1,
			})
			return
		}

		// Transaction succeed.
		SendPacket(ResultRunTxPacketData{
			Hash: packet.Hash(),
			Code: 0,
		})
	}
}
