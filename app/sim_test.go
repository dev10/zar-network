package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authsimops "github.com/cosmos/cosmos-sdk/x/auth/simulation/operations"
	banksimops "github.com/cosmos/cosmos-sdk/x/bank/simulation/operations"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrsimops "github.com/cosmos/cosmos-sdk/x/distribution/simulation/operations"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govsimops "github.com/cosmos/cosmos-sdk/x/gov/simulation/operations"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsimops "github.com/cosmos/cosmos-sdk/x/params/simulation/operations"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingsimops "github.com/cosmos/cosmos-sdk/x/slashing/simulation/operations"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingsimops "github.com/cosmos/cosmos-sdk/x/staking/simulation/operations"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

func init() {
	simapp.GetSimulatorFlags()
}

func testAndRunTxs(app *ZarApp, config simulation.Config) []simulation.WeightedOperation {
	ap := make(simulation.AppParams)

	paramChanges := app.sm.GenerateParamChanges(config.Seed)

	if config.ParamsFile != "" {
		bz, err := ioutil.ReadFile(config.ParamsFile)
		if err != nil {
			panic(err)
		}

		app.cdc.MustUnmarshalJSON(bz, &ap)
	}

	// nolint: govet
	return []simulation.WeightedOperation{
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(app.cdc, OpWeightDeductFee, &v, nil,
					func(_ *rand.Rand) {
						v = 5
					})
				return v
			}(nil),
			authsimops.SimulateDeductFee(app.accountKeeper, app.supplyKeeper),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(app.cdc, OpWeightMsgSend, &v, nil,
					func(_ *rand.Rand) {
						v = 100
					})
				return v
			}(nil),
			banksimops.SimulateMsgSend(app.accountKeeper, app.bankKeeper),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(app.cdc, OpWeightSingleInputMsgMultiSend, &v, nil,
					func(_ *rand.Rand) {
						v = 10
					})
				return v
			}(nil),
			banksimops.SimulateSingleInputMsgMultiSend(app.accountKeeper, app.bankKeeper),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(app.cdc, OpWeightMsgSetWithdrawAddress, &v, nil,
					func(_ *rand.Rand) {
						v = 50
					})
				return v
			}(nil),
			distrsimops.SimulateMsgSetWithdrawAddress(app.distrKeeper),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(app.cdc, OpWeightMsgWithdrawDelegationReward, &v, nil,
					func(_ *rand.Rand) {
						v = 50
					})
				return v
			}(nil),
			distrsimops.SimulateMsgWithdrawDelegatorReward(app.distrKeeper),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(app.cdc, OpWeightMsgWithdrawValidatorCommission, &v, nil,
					func(_ *rand.Rand) {
						v = 50
					})
				return v
			}(nil),
			distrsimops.SimulateMsgWithdrawValidatorCommission(app.distrKeeper),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(app.cdc, OpWeightSubmitVotingSlashingTextProposal, &v, nil,
					func(_ *rand.Rand) {
						v = 5
					})
				return v
			}(nil),
			govsimops.SimulateSubmittingVotingAndSlashingForProposal(app.govKeeper, govsimops.SimulateTextProposalContent),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(app.cdc, OpWeightSubmitVotingSlashingCommunitySpendProposal, &v, nil,
					func(_ *rand.Rand) {
						v = 5
					})
				return v
			}(nil),
			govsimops.SimulateSubmittingVotingAndSlashingForProposal(app.govKeeper, distrsimops.SimulateCommunityPoolSpendProposalContent(app.distrKeeper)),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(app.cdc, OpWeightSubmitVotingSlashingParamChangeProposal, &v, nil,
					func(_ *rand.Rand) {
						v = 5
					})
				return v
			}(nil),
			govsimops.SimulateSubmittingVotingAndSlashingForProposal(app.govKeeper, paramsimops.SimulateParamChangeProposalContent(paramChanges)),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(app.cdc, OpWeightMsgDeposit, &v, nil,
					func(_ *rand.Rand) {
						v = 100
					})
				return v
			}(nil),
			govsimops.SimulateMsgDeposit(app.govKeeper),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(app.cdc, OpWeightMsgCreateValidator, &v, nil,
					func(_ *rand.Rand) {
						v = 100
					})
				return v
			}(nil),
			stakingsimops.SimulateMsgCreateValidator(app.accountKeeper, app.stakingKeeper),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(app.cdc, OpWeightMsgEditValidator, &v, nil,
					func(_ *rand.Rand) {
						v = 5
					})
				return v
			}(nil),
			stakingsimops.SimulateMsgEditValidator(app.stakingKeeper),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(app.cdc, OpWeightMsgDelegate, &v, nil,
					func(_ *rand.Rand) {
						v = 100
					})
				return v
			}(nil),
			stakingsimops.SimulateMsgDelegate(app.accountKeeper, app.stakingKeeper),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(app.cdc, OpWeightMsgUndelegate, &v, nil,
					func(_ *rand.Rand) {
						v = 100
					})
				return v
			}(nil),
			stakingsimops.SimulateMsgUndelegate(app.accountKeeper, app.stakingKeeper),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(app.cdc, OpWeightMsgBeginRedelegate, &v, nil,
					func(_ *rand.Rand) {
						v = 100
					})
				return v
			}(nil),
			stakingsimops.SimulateMsgBeginRedelegate(app.accountKeeper, app.stakingKeeper),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(app.cdc, OpWeightMsgUnjail, &v, nil,
					func(_ *rand.Rand) {
						v = 100
					})
				return v
			}(nil),
			slashingsimops.SimulateMsgUnjail(app.slashingKeeper),
		},
	}
}

// fauxMerkleModeOpt returns a BaseApp option to use a dbStoreAdapter instead of
// an IAVLStore for faster simulation speed.
func fauxMerkleModeOpt(bapp *baseapp.BaseApp) {
	bapp.SetFauxMerkleMode()
}

// interBlockCacheOpt returns a BaseApp option function that sets the persistent
// inter-block write-through cache.
func interBlockCacheOpt() func(*baseapp.BaseApp) {
	return baseapp.SetInterBlockCache(store.NewCommitKVStoreCacheManager())
}

// Profile with:
// /usr/local/go/bin/go test -benchmem -run=^$ github.com/cosmos/cosmos-sdk/GaiaApp -bench ^BenchmarkFullAppSimulation$ -Commit=true -cpuprofile cpu.out
func BenchmarkFullAppSimulation(b *testing.B) {
	logger := log.NewNopLogger()
	config := simapp.NewConfigFromFlags()

	var db dbm.DB
	dir, _ := ioutil.TempDir("", "goleveldb-app-sim")
	db, _ = sdk.NewLevelDB("Simulation", dir)
	defer func() {
		db.Close()
		_ = os.RemoveAll(dir)
	}()

	gapp := NewZarApp(logger, db, nil, true, simapp.FlagPeriodValue, interBlockCacheOpt())

	// Run randomized simulation
	// TODO: parameterize numbers, save for a later PR
	_, simParams, simErr := simulation.SimulateFromSeed(
		b, os.Stdout, gapp.BaseApp, simapp.AppStateFn(gapp.Codec(), gapp.sm),
		testAndRunTxs(gapp, config), gapp.ModuleAccountAddrs(), config,
	)

	// export state and params before the simulation error is checked
	if config.ExportStatePath != "" {
		if err := ExportStateToJSON(gapp, config.ExportStatePath); err != nil {
			fmt.Println(err)
			b.Fail()
		}
	}

	if config.ExportParamsPath != "" {
		if err := simapp.ExportParamsToJSON(simParams, config.ExportParamsPath); err != nil {
			fmt.Println(err)
			b.Fail()
		}
	}

	if simErr != nil {
		fmt.Println(simErr)
		b.FailNow()
	}

	if config.Commit {
		fmt.Println("\nGoLevelDB Stats")
		fmt.Println(db.Stats()["leveldb.stats"])
		fmt.Println("GoLevelDB cached block size", db.Stats()["leveldb.cachedblock"])
	}
}

func TestFullAppSimulation(t *testing.T) {
	if !simapp.FlagEnabledValue {
		t.Skip("skipping application simulation")
	}

	var logger log.Logger
	config := simapp.NewConfigFromFlags()

	if simapp.FlagVerboseValue {
		logger = log.TestingLogger()
	} else {
		logger = log.NewNopLogger()
	}

	var db dbm.DB
	dir, _ := ioutil.TempDir("", "goleveldb-app-sim")
	db, _ = sdk.NewLevelDB("Simulation", dir)

	defer func() {
		db.Close()
		_ = os.RemoveAll(dir)
	}()

	gapp := NewZarApp(logger, db, nil, true, simapp.FlagPeriodValue, fauxMerkleModeOpt)
	require.Equal(t, "GaiaApp", gapp.Name())

	// Run randomized simulation
	_, simParams, simErr := simulation.SimulateFromSeed(
		t, os.Stdout, gapp.BaseApp, simapp.AppStateFn(gapp.Codec(), gapp.sm),
		testAndRunTxs(gapp, config), gapp.ModuleAccountAddrs(), config,
	)

	// export state and params before the simulation error is checked
	if config.ExportStatePath != "" {
		err := ExportStateToJSON(gapp, config.ExportStatePath)
		require.NoError(t, err)
	}

	if config.ExportParamsPath != "" {
		err := simapp.ExportParamsToJSON(simParams, config.ExportParamsPath)
		require.NoError(t, err)
	}

	require.NoError(t, simErr)

	if config.Commit {
		// for memdb:
		// fmt.Println("Database Size", db.Stats()["database.size"])
		fmt.Println("\nGoLevelDB Stats")
		fmt.Println(db.Stats()["leveldb.stats"])
		fmt.Println("GoLevelDB cached block size", db.Stats()["leveldb.cachedblock"])
	}
}

func TestAppImportExport(t *testing.T) {
	if !simapp.FlagEnabledValue {
		t.Skip("skipping application import/export simulation")
	}

	var logger log.Logger
	config := simapp.NewConfigFromFlags()

	if simapp.FlagVerboseValue {
		logger = log.TestingLogger()
	} else {
		logger = log.NewNopLogger()
	}

	var db dbm.DB
	dir, _ := ioutil.TempDir("", "goleveldb-app-sim")
	db, _ = sdk.NewLevelDB("Simulation", dir)

	defer func() {
		db.Close()
		_ = os.RemoveAll(dir)
	}()

	app := NewZarApp(logger, db, nil, true, simapp.FlagPeriodValue, fauxMerkleModeOpt)
	require.Equal(t, "SimApp", app.Name())

	// Run randomized simulation
	_, simParams, simErr := simulation.SimulateFromSeed(
		t, os.Stdout, app.BaseApp, simapp.AppStateFn(app.Codec(), app.sm),
		testAndRunTxs(app, config), app.ModuleAccountAddrs(), config,
	)

	// export state and simParams before the simulation error is checked
	if config.ExportStatePath != "" {
		err := ExportStateToJSON(app, config.ExportStatePath)
		require.NoError(t, err)
	}

	if config.ExportParamsPath != "" {
		err := simapp.ExportParamsToJSON(simParams, config.ExportParamsPath)
		require.NoError(t, err)
	}

	require.NoError(t, simErr)

	if config.Commit {
		// for memdb:
		// fmt.Println("Database Size", db.Stats()["database.size"])
		fmt.Println("\nGoLevelDB Stats")
		fmt.Println(db.Stats()["leveldb.stats"])
		fmt.Println("GoLevelDB cached block size", db.Stats()["leveldb.cachedblock"])
	}

	fmt.Printf("exporting genesis...\n")

	appState, _, err := app.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err)
	fmt.Printf("importing genesis...\n")

	newDir, _ := ioutil.TempDir("", "goleveldb-app-sim-2")
	newDB, _ := sdk.NewLevelDB("Simulation-2", dir)

	defer func() {
		newDB.Close()
		_ = os.RemoveAll(newDir)
	}()

	newApp := NewZarApp(log.NewNopLogger(), newDB, nil, true, simapp.FlagPeriodValue, fauxMerkleModeOpt)
	require.Equal(t, "SimApp", newApp.Name())

	var genesisState simapp.GenesisState
	err = app.cdc.UnmarshalJSON(appState, &genesisState)
	require.NoError(t, err)

	ctxB := newApp.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
	newApp.mm.InitGenesis(ctxB, genesisState)

	fmt.Printf("comparing stores...\n")
	ctxA := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})

	type StoreKeysPrefixes struct {
		A        sdk.StoreKey
		B        sdk.StoreKey
		Prefixes [][]byte
	}

	storeKeysPrefixes := []StoreKeysPrefixes{
		{app.keys[baseapp.MainStoreKey], newApp.keys[baseapp.MainStoreKey], [][]byte{}},
		{app.keys[auth.StoreKey], newApp.keys[auth.StoreKey], [][]byte{}},
		{app.keys[staking.StoreKey], newApp.keys[staking.StoreKey],
			[][]byte{
				staking.UnbondingQueueKey, staking.RedelegationQueueKey, staking.ValidatorQueueKey,
			}}, // ordering may change but it doesn't matter
		{app.keys[slashing.StoreKey], newApp.keys[slashing.StoreKey], [][]byte{}},
		{app.keys[mint.StoreKey], newApp.keys[mint.StoreKey], [][]byte{}},
		{app.keys[distr.StoreKey], newApp.keys[distr.StoreKey], [][]byte{}},
		{app.keys[supply.StoreKey], newApp.keys[supply.StoreKey], [][]byte{}},
		{app.keys[params.StoreKey], newApp.keys[params.StoreKey], [][]byte{}},
		{app.keys[gov.StoreKey], newApp.keys[gov.StoreKey], [][]byte{}},
	}

	for _, storeKeysPrefix := range storeKeysPrefixes {
		storeKeyA := storeKeysPrefix.A
		storeKeyB := storeKeysPrefix.B
		prefixes := storeKeysPrefix.Prefixes

		storeA := ctxA.KVStore(storeKeyA)
		storeB := ctxB.KVStore(storeKeyB)

		failedKVAs, failedKVBs := sdk.DiffKVStores(storeA, storeB, prefixes)
		require.Equal(t, len(failedKVAs), len(failedKVBs), "unequal sets of key-values to compare")

		fmt.Printf("compared %d key/value pairs between %s and %s\n", len(failedKVAs), storeKeyA, storeKeyB)
		require.Len(t, failedKVAs, 0, simapp.GetSimulationLog(storeKeyA.Name(), app.sm.StoreDecoders, app.cdc, failedKVAs, failedKVBs))
	}
}

func TestAppSimulationAfterImport(t *testing.T) {
	if !simapp.FlagEnabledValue {
		t.Skip("skipping application simulation after import")
	}

	var logger log.Logger
	config := simapp.NewConfigFromFlags()

	if simapp.FlagVerboseValue {
		logger = log.TestingLogger()
	} else {
		logger = log.NewNopLogger()
	}

	dir, _ := ioutil.TempDir("", "goleveldb-app-sim")
	db, _ := sdk.NewLevelDB("Simulation", dir)

	defer func() {
		db.Close()
		_ = os.RemoveAll(dir)
	}()

	gapp := NewZarApp(logger, db, nil, true, simapp.FlagPeriodValue, fauxMerkleModeOpt)
	require.Equal(t, "GaiaApp", gapp.Name())

	// Run randomized simulation
	// Run randomized simulation
	stopEarly, simParams, simErr := simulation.SimulateFromSeed(
		t, os.Stdout, gapp.BaseApp, simapp.AppStateFn(gapp.Codec(), gapp.sm),
		testAndRunTxs(gapp, config), gapp.ModuleAccountAddrs(), config,
	)

	// export state and params before the simulation error is checked
	if config.ExportStatePath != "" {
		err := ExportStateToJSON(gapp, config.ExportStatePath)
		require.NoError(t, err)
	}

	if config.ExportParamsPath != "" {
		err := simapp.ExportParamsToJSON(simParams, config.ExportParamsPath)
		require.NoError(t, err)
	}

	require.NoError(t, simErr)

	if config.Commit {
		// for memdb:
		// fmt.Println("Database Size", db.Stats()["database.size"])
		fmt.Println("\nGoLevelDB Stats")
		fmt.Println(db.Stats()["leveldb.stats"])
		fmt.Println("GoLevelDB cached block size", db.Stats()["leveldb.cachedblock"])
	}

	if stopEarly {
		// we can't export or import a zero-validator genesis
		fmt.Printf("We can't export or import a zero-validator genesis, exiting test...\n")
		return
	}

	fmt.Printf("Exporting genesis...\n")

	appState, _, err := gapp.ExportAppStateAndValidators(true, []string{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Importing genesis...\n")

	newDir, _ := ioutil.TempDir("", "goleveldb-app-sim-2")
	newDB, _ := sdk.NewLevelDB("Simulation-2", dir)

	defer func() {
		newDB.Close()
		_ = os.RemoveAll(newDir)
	}()

	newApp := NewZarApp(log.NewNopLogger(), newDB, nil, true, 0, fauxMerkleModeOpt)
	require.Equal(t, "GaiaApp", newApp.Name())

	newApp.InitChain(abci.RequestInitChain{
		AppStateBytes: appState,
	})

	// Run randomized simulation on imported app
	_, _, err = simulation.SimulateFromSeed(
		t, os.Stdout, newApp.BaseApp, simapp.AppStateFn(gapp.Codec(), gapp.sm),
		testAndRunTxs(newApp, config), newApp.ModuleAccountAddrs(), config,
	)

	require.NoError(t, err)
}

// TODO: Make another test for the fuzzer itself, which just has noOp txs
// and doesn't depend on the application.
func TestAppStateDeterminism(t *testing.T) {
	if !simapp.FlagEnabledValue {
		t.Skip("skipping application simulation")
	}

	config := simapp.NewConfigFromFlags()
	config.InitialBlockHeight = 1
	config.ExportParamsPath = ""
	config.OnOperation = false
	config.AllInvariants = false

	numSeeds := 3
	numTimesToRunPerSeed := 5
	appHashList := make([]json.RawMessage, numTimesToRunPerSeed)

	for i := 0; i < numSeeds; i++ {
		config.Seed = rand.Int63()

		for j := 0; j < numTimesToRunPerSeed; j++ {
			logger := log.NewNopLogger()
			db := dbm.NewMemDB()
			app := NewZarApp(logger, db, nil, true, simapp.FlagPeriodValue, interBlockCacheOpt())

			fmt.Printf(
				"running non-determinism simulation; seed %d: %d/%d, attempt: %d/%d\n",
				config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
			)

			_, _, err := simulation.SimulateFromSeed(
				t, os.Stdout, app.BaseApp, simapp.AppStateFn(app.Codec(), app.sm),
				testAndRunTxs(app, config), app.ModuleAccountAddrs(), config,
			)
			require.NoError(t, err)

			appHash := app.LastCommitID().Hash
			appHashList[j] = appHash

			if j != 0 {
				require.Equal(
					t, appHashList[0], appHashList[j],
					"non-determinism in seed %d: %d/%d, attempt: %d/%d\n", config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
				)
			}
		}
	}
}

func BenchmarkInvariants(b *testing.B) {
	logger := log.NewNopLogger()

	config := simapp.NewConfigFromFlags()
	config.AllInvariants = false

	dir, _ := ioutil.TempDir("", "goleveldb-app-invariant-bench")
	db, _ := sdk.NewLevelDB("simulation", dir)

	defer func() {
		db.Close()
		os.RemoveAll(dir)
	}()

	gapp := NewZarApp(logger, db, nil, true, simapp.FlagPeriodValue, interBlockCacheOpt())

	// 2. Run parameterized simulation (w/o invariants)
	_, simParams, simErr := simulation.SimulateFromSeed(
		b, ioutil.Discard, gapp.BaseApp, simapp.AppStateFn(gapp.Codec(), gapp.sm),
		testAndRunTxs(gapp, config), gapp.ModuleAccountAddrs(), config,
	)

	// export state and params before the simulation error is checked
	if config.ExportStatePath != "" {
		if err := ExportStateToJSON(gapp, config.ExportStatePath); err != nil {
			fmt.Println(err)
			b.Fail()
		}
	}

	if config.ExportParamsPath != "" {
		if err := simapp.ExportParamsToJSON(simParams, config.ExportParamsPath); err != nil {
			fmt.Println(err)
			b.Fail()
		}
	}

	if simErr != nil {
		fmt.Println(simErr)
		b.FailNow()
	}

	ctx := gapp.NewContext(true, abci.Header{Height: gapp.LastBlockHeight() + 1})

	// 3. Benchmark each invariant separately
	//
	// NOTE: We use the crisis keeper as it has all the invariants registered with
	// their respective metadata which makes it useful for testing/benchmarking.
	for _, cr := range gapp.crisisKeeper.Routes() {
		b.Run(fmt.Sprintf("%s/%s", cr.ModuleName, cr.Route), func(b *testing.B) {
			if res, stop := cr.Invar(ctx); stop {
				fmt.Printf("broken invariant at block %d of %d\n%s", ctx.BlockHeight()-1, config.NumBlocks, res)
				b.FailNow()
			}
		})
	}
}
