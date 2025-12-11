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
		TaxIndex:     tax.String(),
		GdpIndex:     gdp.String(),
		ExportsIndex: exports.String(),
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
	decs := []string{msg.TaxIndex, msg.GdpIndex, msg.ExportsIndex}
	for _, d := range decs {
		v, err := sdkmath.LegacyNewDecFromStr(d)
		if err != nil {
			return err
		}
		if v.IsNegative() {
			return ErrInvalidMetric
		}
	}
	return nil
}

func (msg *MsgSetRegionMetric) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}
