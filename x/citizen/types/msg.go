package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMsgSetCitizenRegion(authority, address, regionID string) *MsgSetCitizenRegion {
	return &MsgSetCitizenRegion{Authority: authority, Address: address, RegionId: regionID}
}

func (msg *MsgSetCitizenRegion) ValidateBasic() error {
	if msg == nil {
		return ErrInvalidRegion
	}
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return err
	}
	if _, err := sdk.AccAddressFromBech32(msg.Address); err != nil {
		return err
	}
	if msg.RegionId == "" {
		return ErrInvalidRegion
	}
	return nil
}

func (msg *MsgSetCitizenRegion) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

func NewMsgClearCitizenRegion(authority, address string) *MsgClearCitizenRegion {
	return &MsgClearCitizenRegion{Authority: authority, Address: address}
}

func (msg *MsgClearCitizenRegion) ValidateBasic() error {
	if msg == nil {
		return ErrInvalidRegion
	}
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return err
	}
	if _, err := sdk.AccAddressFromBech32(msg.Address); err != nil {
		return err
	}
	return nil
}

func (msg *MsgClearCitizenRegion) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}
