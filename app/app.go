package app

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	circuitkeeper "cosmossdk.io/x/circuit/keeper"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"

	abci "github.com/cometbft/cometbft/abci/types"
	dbm "github.com/cosmos/cosmos-db"

	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	// custom modules
	citizenkeeper "uag/x/citizen/keeper"
	growthkeeper "uag/x/growth/keeper"
)

const (
	AccountAddressPrefix = "uag"
	Name                 = "uag"
	ChainCoinType        = 118
)

var (
	DefaultNodeHome string
)

var maccPerms = map[string][]string{
	authtypes.FeeCollectorName:     nil,
	distrtypes.ModuleName:          nil,
	stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
	govtypes.ModuleName:            {authtypes.Burner},
}

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
	runtimeApp        *runtime.App

	// The module manager
	moduleManager *module.Manager
	// simulation manager
	sm *module.SimulationManager
	// module configurator
	configurator module.Configurator

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
	// custom keepers
	CitizenKeeper citizenkeeper.Keeper
	GrowthKeeper  growthkeeper.Keeper

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
			AppConfigOptions(),
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
		&app.GrowthKeeper,
		&app.keys,
		&app.tkeys,
		&app.memKeys,
	); err != nil {
		panic(err)
	}

	// setup baseapp
	builtApp := appBuilder.Build(db, traceStore, baseAppOptions...)
	app.BaseApp = builtApp.BaseApp
	app.runtimeApp = builtApp
	app.moduleManager = builtApp.ModuleManager

	app.setupUpgradeHandlers()

	app.setupHooks()

	// enable sign mode textual by overwriting the default TxConfig
	app.txConfig = authtx.NewTxConfig(app.appCodec, authtx.DefaultSignModes)

	// set the baseapp's parameter store
	app.BaseApp.SetParamStore(app.ConsensusParamsKeeper.ParamsStore)

	// set ante handler
	anteHandler, err := ante.NewAnteHandler(ante.HandlerOptions{
		AccountKeeper:   app.AccountKeeper,
		BankKeeper:      app.BankKeeper,
		FeegrantKeeper:  app.FeegrantKeeper,
		SignModeHandler: app.txConfig.SignModeHandler(),
	})
	if err != nil {
		panic(err)
	}
	app.SetAnteHandler(anteHandler)

	// load state streaming if enabled
	if err := app.BaseApp.RegisterStreamingServices(appOpts, app.keys); err != nil {
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
func AppConfigOptions() depinject.Config {
	return depinject.Configs(
		AppConfig,
	)
}

// Name returns the name of the App
func (app *App) Name() string { return app.BaseApp.Name() }

// ModuleManager accessor for tests/CLI wiring.
func (app *App) ModuleManager() *module.Manager {
	return app.moduleManager
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
	accs := make(map[string]bool)
	for acc := range maccPerms {
		accs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return accs
}

// BlockedAddresses returns the set of module account addresses that are blocked from receiving funds.
func BlockedAddresses() map[string]bool {
	accs := make(map[string]bool)
	for acc := range maccPerms {
		accs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return accs
}

// DefaultGenesis returns default genesis state.
func (app *App) DefaultGenesis() map[string]json.RawMessage {
	if app.runtimeApp != nil {
		return app.runtimeApp.DefaultGenesis()
	}
	return nil
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

// GetStoreKeys returns all store keys.
func (app *App) GetStoreKeys() []storetypes.StoreKey {
	out := make([]storetypes.StoreKey, 0, len(app.keys))
	for _, k := range app.keys {
		out = append(out, k)
	}
	return out
}

// RegisterAPIRoutes registers all application module routes with the provided API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	if app.runtimeApp != nil {
		app.runtimeApp.RegisterAPIRoutes(apiSvr, apiConfig)
		return
	}
}

// RegisterTxService implements the Application interface.
func (app *App) RegisterTxService(clientCtx client.Context) {
	if app.runtimeApp != nil {
		app.runtimeApp.RegisterTxService(clientCtx)
		return
	}
}

// RegisterTendermintService implements the Application interface.
func (app *App) RegisterTendermintService(clientCtx client.Context) {
	if app.runtimeApp != nil {
		app.runtimeApp.RegisterTendermintService(clientCtx)
		return
	}
}

// RegisterNodeService registers the node gRPC service.
func (app *App) RegisterNodeService(clientCtx client.Context, cfg config.Config) {
	if app.runtimeApp != nil {
		app.runtimeApp.RegisterNodeService(clientCtx, cfg)
		return
	}
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
