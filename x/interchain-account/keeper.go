package interchain_account

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/tendermint/tendermint/crypto/tmhash"
)

type Keeper struct {
	cdc *codec.Codec
	key sdk.StoreKey
	router sdk.Router
	ante IBCAnteHandler
	accountKeeper auth.AccountKeeper
}

func NewKeepr(cdc *codec.Codec, key sdk.StoreKey, router sdk.Router, ante IBCAnteHandler, accountKeeper auth.AccountKeeper) Keeper {
	return Keeper {
		cdc: cdc,
		key: key,
		router: router,
		ante: ante,
		accountKeeper: accountKeeper,
	}
}

func (keeper Keeper) RegisterIBCAccount(ctx sdk.Context, packet RegisterIBCAccountPacketData ) ([]byte, sdk.Error) {
	path := packet.SourcePort() + "/" + packet.SourceChannel()
	// Address is sha256(channel + salt)
	address := tmhash.Sum(append([]byte(path), packet.Salt...))

	acc := keeper.accountKeeper.GetAccount(ctx, address)
	if acc == nil {
		acc = keeper.accountKeeper.NewAccountWithAddress(ctx, address)
	} else {
		// Only fail if exising account is interchain account
		if acc.GetSequence() != 1 || acc.GetPubKey() != nil {
			return []byte{}, sdk.ErrUnauthorized("Interchain account already exists")
		}
	}

	// Don't block even if normal account with same address already exists.
	// Otherwise, attackers can disrupt creating interchain account by sending small asset to expected address in advance.
	// Because IBC is exeucted asynchronously, attackers can know that sending chain wants to create interchain account,
	// and send asset to expected address faster than relayer.

	// Set interchain account's sequence 1
	// Normal account can't be sequence 1 without publish pub key.
	// So, pubkey == nil && sequence == 1 is mark for interchain account.
	err := acc.SetSequence(1)
	if err != nil {
		return []byte{}, sdk.ErrInternal(err.Error())
	}
	keeper.accountKeeper.SetAccount(ctx, acc)

	// Set interchain account's permission with channel data.
	store := ctx.KVStore(keeper.key)
	store.Set(address, []byte(path))

	return address, nil
}

func (keeper Keeper) RunTxIBCTx(ctx sdk.Context, packet ChainAccountTx) sdk.Error {
	err := keeper.ante(ctx, packet)
	if err != nil {
		return err
	}

	for _, msg := range packet.Msgs {
		msgRoute := msg.Route()
		// Get handler from router.
		handler := keeper.router.Route(msgRoute)
		if handler == nil {
			return sdk.ErrUnknownRequest("Unrecognized Msg type: " + msgRoute)
		}

		result := handler(ctx, msg)
		if result.IsOK() == false {
			return sdk.ErrInternal(fmt.Sprintf("Fail to run msg: %s", msg.Type()))
		}
	}

	return nil
}

func (keeper Keeper) hasPrivilege(ctx sdk.Context, path string, account sdk.AccAddress) bool {
	acc := keeper.accountKeeper.GetAccount(ctx, []byte(account))
	if acc == nil {
		return false
	} else {
		if acc.GetSequence() != 1 || acc.GetPubKey() != nil {
			return false
		}
	}

	// Check that account has right permission for channel.
	store := ctx.KVStore(keeper.key)
	accountPath := store.Get([]byte(account))
	if string(accountPath) != path {
		return false
	}

	return true
}
