package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	"cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramsproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"

	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"

	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/comet"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	circuitkeeper "cosmossdk.io/x/circuit/keeper"
	circuittypes "cosmossdk.io/x/circuit/types"
	consensuskeeper "cosmossdk.io/x/consensus/keeper"
	consensustypes "cosmossdk.io/x/consensus/types"
	minttypes "cosmossdk.io/x/mint/types"

	// ibc v10
	ibcclienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v10/modules/core/03-connection/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v10/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"
	ibctypes "github.com/cosmos/ibc-go/v10/modules/core/types"
	// wasm
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	// custom modules
	citizenkeeper "uagd/x/citizen/keeper"
	fundkeeper "uagd/x/fund/keeper"
	growthkeeper "uagd/x/growth/keeper"
	uagdkeeper "uagd/x/uagd/keeper"
	ugovkeeper "uagd/x/ugov/keeper"
	ugovtypes "uagd/x/ugov/types"
)

const (
	AccountAddressPrefix = "uag"
	Name                 = "uagd"
)

var (
	DefaultNodeHome string
)

var (
	_ servertypes.Application = (*App)(nil)
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, "."+Name)
}

// App extends an ABCI application, but with keepers, managers and other useful components.
type App struct {
	*baseapp.BaseApp
	appCodec          codec.Codec
	cdc               *codec.ProtoCodec
	interfaceRegistry codectypes.InterfaceRegistry

	// The module manager
	moduleManager *module.Manager
	// simulation manager
	sm *module.SimulationManager
	// module configurator
	configurator module.Configurator

	// Default chain address codecs
	ac address.Codec
	vc address.ValidatorCodec

	// keepers
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bankkeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	GovKeeper             *govkeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	AuthzKeeper           authzkeeper.Keeper
	ConsensusParamsKeeper consensuskeeper.Keeper
	CircuitBreakerKeeper  circuitkeeper.Keeper
	ParamsKeeper          paramskeeper.Keeper
	IBCKeeper             *ibckeeper.Keeper
	WasmKeeper            *wasmkeeper.Keeper
	CapabilityKeeper      *capabilitykeeper.Keeper
	ScopedIBCKeeper       capabilitykeeper.ScopedKeeper
	ScopedWasmKeeper      capabilitykeeper.ScopedKeeper

	// custom keepers
	CitizenKeeper citizenkeeper.Keeper
	FundKeeper    fundkeeper.Keeper
	GrowthKeeper  growthkeeper.Keeper
	UAGDKeeper    uagdkeeper.Keeper
	UGovKeeper    ugovkeeper.Keeper

	// the module store keys
	keys map[string]*storetypes.KVStoreKey

	// the transient store keys
	tkeys map[string]*storetypes.TransientStoreKey

	// the memory store keys
	memKeys map[string]*storetypes.MemoryStoreKey
}

// New creates a new App object
func New(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {
	var (
		app        = &App{}
		appBuilder *runtime.AppBuilder

		appConfig = depinject.Configs(
			AppConfig(),
			depinject.Supply(
				logger,
				appOpts,
			),
		)
	)

	if err := depinject.Inject(appConfig,
		&appBuilder,
		&app.appCodec,
		&app.cdc,
		&app.interfaceRegistry,
		&app.moduleManager,
		&app.sm,
		&app.configurator,
		&app.AccountKeeper,
		&app.BankKeeper,
		&app.StakingKeeper,
		&app.SlashingKeeper,
		&app.DistrKeeper,
		&app.GovKeeper,
		&app.UpgradeKeeper,
		&app.AuthzKeeper,
		&app.ConsensusParamsKeeper,
		&app.CircuitBreakerKeeper,
		&app.ParamsKeeper,
		&app.CitizenKeeper,
		&app.FundKeeper,
		&app.GrowthKeeper,
		&app.UAGDKeeper,
		&app.UGovKeeper,
		&app.IBCKeeper,
		&app.WasmKeeper,
		&app.CapabilityKeeper,
		&app.ScopedIBCKeeper,
		&app.ScopedWasmKeeper,
		&app.keys,
		&app.tkeys,
		&app.memKeys,
	); err != nil {
		panic(err)
	}

	// setup baseapp
	app.BaseApp = appBuilder.Build(db, traceStore, baseAppOptions...)

	app.setupUpgradeHandlers()

	app.setupHooks()

	// enable sign mode textual by overwriting the default TxConfig
	enabledSignModes := append(authtx.DefaultSignModes, sigtypes.SignMode_SIGN_MODE_TEXTUAL)
	txConfigOpts := authtx.ConfigOptions{
		EnabledSignModes:           enabledSignModes,
		TextualCoinMetadataQueryFn: authtx.NewGRPCCoinMetadataQueryFn(app.BankKeeper),
	}
	txConfig, err := authtx.NewTxConfigWithOptions(
		app.appCodec,
		txConfigOpts,
	)
	if err != nil {
		panic(err)
	}
	app.txConfig = txConfig

	// set the baseapp's parameter store
	app.BaseApp.SetParamStore(app.ConsensusParamsKeeper.ParamsStore)

	// set ante handler
	anteHandler, err := NewAnteHandler(
		app,
		AnteHandlerOptions{
			HandlerOptions: authante.HandlerOptions{
				AccountKeeper:   app.AccountKeeper,
				BankKeeper:      app.BankKeeper,
				FeegrantKeeper:  app.FeegrantKeeper,
				SignModeHandler: txConfig.SignModeHandler(),
				SigGasConsumer:  authante.DefaultSigVerificationGasConsumer,
			},
			IBCKeeper:         app.IBCKeeper,
			StakingKeeper:     app.StakingKeeper,
			CircuitKeeper:     app.CircuitBreakerKeeper,
			WasmConfig:        wasmkeeper.DefaultWasmConfig(),
			WasmKeeper:        app.WasmKeeper,
			TXCounterStoreKey: app.keys[wasmtypes.StoreKey],
		},
	)
	if err != nil {
		panic(err)
	}
	app.SetAnteHandler(anteHandler)

	// register invariants
	app.moduleManager.RegisterInvariants(app.CrisisKeeper)

	// register routes
	app.moduleManager.RegisterRoutes(app.Router(), app.QueryRouter(), app.GrpcQueryRouter())

	// register grpc services
	app.moduleManager.RegisterServices(app.configurator)

	// register snapshots
	if app.SnapshotManager() != nil {
		err := app.moduleManager.RegisterSnapshotExtensions(app.SnapshotManager())
		if err != nil {
			panic(fmt.Errorf("failed to register snapshot extensions: %w", err))
		}
	}

	// set the governance module account as the authority for the upgrade module
	app.UpgradeKeeper.SetAuthority(authtypes.NewModuleAddress(govtypes.ModuleName).String())

	// wasm: set x/wasm authority to ugov module account (custom governance)
	app.WasmKeeper.SetAuthority(authtypes.NewModuleAddress(ugovtypes.ModuleName).String())

	// load state streaming if enabled
	if err := app.LoadStreamingServices(appOpts, app.cdc); err != nil {
		panic(err)
	}

	// load latest version
	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			panic(err)
		}
	}

	return app
}

// AppConfig returns the default app config.
func AppConfig() depinject.Config {
	return depinject.Configs(
		baseAppConfig(),
		depinject.Provide(
			ProvideFundGrowthKeeper,
		),
		depinject.Supply(
			genutiltypes.NewModuleBasics(
				genutil.AppModuleBasic{},
			),
		),
	)
}

// Name returns the name of the App
func (app *App) Name() string { return app.BaseApp.Name() }

// getGovProposalHandlers returns the governance proposal handlers.
func getGovProposalHandlers() []govclient.ProposalHandler {
	return []govclient.ProposalHandler{
		paramsproposal.Handler,
	}
}

// RegisterAPIRoutes registers all application module routes with the provided API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	rpc.RegisterRoutes(clientCtx, apiSvr.Router)
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register custom tx routes.
	// This gRPC service is defined in protobuf.
	// Since its location is in the app package, we need to register it ourselves.
	uagdtxtypes.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register swagger API from github.com/cosmos/cosmos-sdk/client/docs/swagger.yaml
	if apiConfig.Swagger {
		RegisterSwaggerAPI(apiSvr.Router)
	}

	app.moduleManager.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
}

// RegisterTxService implements the Application interface.
func (app *App) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application interface.
func (app *App) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(clientCtx, app.BaseApp.GRPCQueryRouter(), app.interfaceRegistry, app.QueryRouter())
}

// RegisterNodeService registers the node gRPC service.
func (app *App) RegisterNodeService(clientCtx client.Context, cfg config.Config) {
	nodeservice.RegisterNodeService(clientCtx, app.BaseApp.GRPCQueryRouter())
}

// RegisterSwaggerAPI registers swagger route with API Server
func RegisterSwaggerAPI(rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

// DefaultGenesis returns a default genesis state
func DefaultGenesis() map[string]json.RawMessage {
	return ModuleBasics.DefaultGenesis(app.appCodec)
}

// BeginBlocker application updates every begin block
func (app *App) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	return app.moduleManager.BeginBlock(ctx)
}

// EndBlocker application updates every end block
func (app *App) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	return app.moduleManager.EndBlock(ctx)
}

// InitChainer application update at chain initialization
func (app *App) InitChainer(ctx sdk.Context, req abci.RequestInitChain) (abci.ResponseInitChain, error) {
	var genesisState map[string]json.RawMessage
	app.appCodec.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	response, err := app.moduleManager.InitGenesis(ctx, app.appCodec, genesisState)
	if err != nil {
		panic(err)
	}

	return response, nil
}

// LoadHeight loads a particular height
func (app *App) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *App) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// LegacyAmino returns the application's LegacyAmino codec.
func (app *App) LegacyAmino() *codec.LegacyAmino {
	return app.cdc.Amino
}

// AppCodec returns the application's appCodec.
func (app *App) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns the application's InterfaceRegistry.
func (app *App) InterfaceRegistry() codectypes.InterfaceRegistry {
	return app.interfaceRegistry
}

// TxConfig returns the application's TxConfig.
func (app *App) TxConfig() client.TxConfig {
	return app.txConfig
}

// GetKey returns the KVStoreKey for the provided store key.
func (app *App) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
func (app *App) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemoryStoreKey for the provided store key.
func (app *App) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// RegisterLegacyAminoCodec registers the app's module interfaces with the LegacyAmino codec.
func (app *App) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	ModuleBasics.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the app's module interfaces.
func (app *App) RegisterInterfaces(reg codectypes.InterfaceRegistry) {
	ModuleBasics.RegisterInterfaces(reg)
}

// SimulationManager implements the SimulationApp interface
func (app *App) SimulationManager() *module.SimulationManager {
	return app.sm
}
