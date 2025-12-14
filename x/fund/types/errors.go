package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrFundNotFound       = errorsmod.Register(ModuleName, 1, "fund not found")
	ErrFundInactive       = errorsmod.Register(ModuleName, 2, "fund is inactive")
	ErrInvalidDenom       = errorsmod.Register(ModuleName, 3, "invalid denom")
	ErrDelegationLimit    = errorsmod.Register(ModuleName, 4, "delegation limit exceeded")
	ErrPayrollLimit       = errorsmod.Register(ModuleName, 5, "payroll limit exceeded")
	ErrUnauthorized       = errorsmod.Register(ModuleName, 6, "unauthorized")
	ErrDirectExecDisabled = errorsmod.Register(ModuleName, 7, "direct execution disabled; use governance-approved plan")
	ErrRegionLocked       = errorsmod.Register(ModuleName, 8, "region is locked due to occupation")
)
