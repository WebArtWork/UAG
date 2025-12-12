package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrUnauthorized  = errorsmod.Wrap(sdkerrors.ErrUnauthorized, "unauthorized")
	ErrInvalidRegion = errorsmod.Register(ModuleName, 2, "invalid region mapping")
)
