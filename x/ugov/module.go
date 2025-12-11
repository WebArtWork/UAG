package ugov

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	abci "github.com/cometbft/cometbft/abci/types"

	"uagd/x/ugov/keeper"
	"uagd/x/ugov/types"
)

type AppModuleBasic struct{}

func (AppModuleBasic) Name() string { return types.ModuleName }
func (AppModuleBasic) RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {}
func (AppModuleBasic) RegisterInterfaces(_ codec.Types)              {}

func (AppModuleBasic) DefaultGenesis(_ codec.JSONCodec) json.RawMessage {
	bz, _ := json.Marshal(types.DefaultGenesis())
	return bz
}

func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var gs types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &gs); err != nil {
		return err
	}
	return gs.Params.Validate()
}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(_ client.Context, _ *runtime.ServeMux) {}
func (AppModuleBasic) GetTxCmd() *client.Command    { return nil }
func (AppModuleBasic) GetQueryCmd() *client.Command { return nil }

type AppModule struct {
	AppModuleBasic
	k keeper.Keeper
}

func NewAppModule(k keeper.Keeper) AppModule { return AppModule{k: k} }

func (am AppModule) RegisterServices(_ module.Configurator) {
	// Once proto exists:
	// types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.k))
	// types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(am.k))
}

func (am AppModule) InitGenesis(ctx module.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var gs types.GenesisState
	cdc.MustUnmarshalJSON(data, &gs)
	InitGenesis(ctx, am.k, gs)
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx module.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := ExportGenesis(ctx, am.k)
	return cdc.MustMarshalJSON(gs)
}
