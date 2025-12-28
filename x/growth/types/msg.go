package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMsgSetRegionMetric(authority, regionID, period string, tax, gdp, exports sdkmath.LegacyDec) *MsgSetRegionMetric {
	return &MsgSetRegionMetric{
		Authority:    authority,
		RegionId:     regionID,
		Period:       period,
		TaxIndex:     tax,
		GdpIndex:     gdp,
		ExportsIndex: exports,
	}
}

func (msg *MsgSetRegionMetric) ValidateBasic() error {
	if msg == nil {
		return ErrInvalidMetric
	}
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return err
	}
	if msg.RegionId == "" || msg.Period == "" {
		return ErrInvalidMetric
	}
	decs := []sdkmath.LegacyDec{msg.TaxIndex, msg.GdpIndex, msg.ExportsIndex}
	for _, d := range decs {
		if d.IsNegative() {
			return ErrInvalidMetric
		}
	}
	return nil
}

func (msg *MsgSetRegionMetric) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}
