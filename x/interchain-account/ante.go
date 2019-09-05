package interchain_account

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// IBC Ante handler checks that transaction has right permissions.
// With looping msgs in tx, check that msg's signer accounts are interchain account that the sending chain has made.
func NewIBCAnteHandler(keeper Keeper) IBCAnteHandler {
	return func(ctx sdk.Context, tx ChainAccountTx) sdk.Error {
		seen := map[string]bool{}
		var signers []sdk.AccAddress
		for _, msg := range tx.Msgs {
			for _, addr := range msg.GetSigners() {
				if !seen[addr.String()] {
					signers = append(signers, addr)
					seen[addr.String()] = true
				}
			}
		}

		path := tx.SourcePort() + "/" + tx.SourceChannel()

		for _, signer := range signers {
			if keeper.hasPrivilege(ctx, path, signer) == false {
				return sdk.ErrUnauthorized("Unauthorized")
			}
		}

		return nil
	}
}
