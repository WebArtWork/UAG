package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	circuitkeeper "cosmossdk.io/x/circuit/keeper"
	circuittypes "cosmossdk.io/x/circuit/types"
	consensuskeeper "cosmossdk.io/x/consensus/keeper"
	consensustypes "cosmossdk.io/x/consensus/types"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"

	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramsproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	// IBC v8
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	ibctypes "github.com/cosmos/ibc-go/v8/modules/core/types"

	// wasm
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

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
	txConfig          client.TxConfig
	cdc               *codec.LegacyAmino
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
	CrisisKeeper          *crisiskeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	AuthzKeeper           authzkeeper.Keeper
	FeegrantKeeper        feegrantkeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper
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
		&app.CrisisKeeper,
		&app.UpgradeKeeper,
		&app.AuthzKeeper,
		&app.FeegrantKeeper,
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
	enabledSignModes := append(authtx.DefaultSignModes, authtypes.SignMode_SIGN_MODE_TEXTUAL)
	txConfigOpts := authtx.ConfigOptions{
		EnabledSignModes:           enabledSignModes,
		TextualCoinMetadataQueryFn: authtx.NewGRPCCoinMetadataQueryFn(app.BaseApp.GRPCQueryRouter(), app.BankKeeper),
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
			HandlerOptions: ante.HandlerOptions{
				AccountKeeper:   app.AccountKeeper,
				BankKeeper:      app.BankKeeper,
				FeegrantKeeper:  app.FeegrantKeeper,
				SignModeHandler: txConfig.SignModeHandler(),
				SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
			},
			IBCKeeper:         app.IBCKeeper,
			StakingKeeper:     app.StakingKeeper,
			CircuitKeeper:     &app.CircuitBreakerKeeper,
			WasmConfig:        &wasmtypes.WasmConfig{},
			WasmKeeper:        app.WasmKeeper,
			TXCounterStoreKey: app.keys[wasmtypes.StoreKey],
		},
	)
	if err != nil {
		panic(err)
	}
	app.SetAnteHandler(anteHandler)

	// register invariants
	app.CrisisKeeper.RegisterRoute(module.NewManager().Modules()...)

	// register routes
	if err := app.moduleManager.RegisterRoutes(app.BaseApp.MsgServiceRouter(), app.BaseApp.GRPCQueryRouter()); err != nil {
		panic(err)
	}

	// register grpc services
	app.moduleManager.RegisterServices(app.configurator)

	// register snapshots
	if manager := app.SnapshotManager(); manager != nil {
		err := manager.RegisterExtensions()
		if err != nil {
			panic(fmt.Errorf("failed to register snapshot extensions: %w", err))
		}
	}

	// set the governance module account as the authority for the upgrade module
	app.UpgradeKeeper.SetAuthority(authtypes.NewModuleAddress(govtypes.ModuleName).String())

	// wasm: set x/wasm authority to ugov module account (custom governance)
	app.WasmKeeper.SetAuthority(authtypes.NewModuleAddress(ugovtypes.ModuleName).String())

	// load state streaming if enabled
	if err := app.LoadStreamingServices(appOpts, app.appCodec, app.keys); err != nil {
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
		AppConfig(),
		depinject.Provide(
			ProvideFundGrowthKeeper,
		),
	)
}

// Name returns the name of the App
func (app *App) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (app *App) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	return app.moduleManager.BeginBlock(ctx)
}

// EndBlocker application updates every end block
func (app *App) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	return app.moduleManager.EndBlock(ctx)
}

// InitChainer application update at chain initialization
func (app *App) InitChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	var genesisState map[string]json.RawMessage
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	return app.moduleManager.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads a particular height
func (app *App) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *App) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)

	accs := make(map[string]bool)
	for acc := range maccPerms {
		accs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return accs
}

// LegacyAmino returns the application's LegacyAmino codec.
func (app *App) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
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

// RegisterAPIRoutes registers all application module routes with the provided API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx

	// Register grpc-gateway routes for all modules
	app.moduleManager.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register swagger API if enabled
	if apiConfig.Swagger {
		RegisterSwaggerAPI(apiSvr.Router)
	}
}

// RegisterTxService implements the Application interface.
func (app *App) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application interface.
func (app *App) RegisterTendermintService(clientCtx client.Context) {
	cmtservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		app.Query,
	)
}

// RegisterNodeService registers the node gRPC service.
func (app *App) RegisterNodeService(clientCtx client.Context, cfg config.Config) {
	nodeservice.RegisterNodeService(clientCtx, app.BaseApp.GRPCQueryRouter(), cfg)
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

// SimulationManager implements the SimulationApp interface
func (app *App) SimulationManager() *module.SimulationManager {
	return app.sm
}

// setupUpgradeHandlers sets up the upgrade handlers
func (app *App) setupUpgradeHandlers() {
	// Add your upgrade handlers here
}

// setupHooks sets up the hooks for the app
func (app *App) setupHooks() {
	// Add your hooks here
}
