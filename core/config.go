package core

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/hyperledger/burrow/config"
	"github.com/hyperledger/burrow/consensus/abci"
	"github.com/hyperledger/burrow/consensus/tendermint"
	"github.com/hyperledger/burrow/execution"
	"github.com/hyperledger/burrow/execution/registry"
	"github.com/hyperledger/burrow/keys"
	"github.com/hyperledger/burrow/logging/logconfig"
	"github.com/hyperledger/burrow/logging/structure"
	"github.com/hyperledger/burrow/project"
	tmConfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/node"
	tmTypes "github.com/tendermint/tendermint/types"
)

// LoadKeysFromConfig sets the keyClient & keyStore based on the given config
func (kern *Kernel) LoadKeysFromConfig(conf *keys.KeysConfig) (err error) {
	kern.keyStore = keys.NewFilesystemKeyStore(conf.KeysDirectory, conf.AllowBadFilePermissions)
	if conf.RemoteAddress != "" {
		kern.keyClient, err = keys.NewRemoteKeyClient(conf.RemoteAddress, kern.Logger)
		if err != nil {
			return err
		}
	} else {
		kern.keyClient = keys.NewLocalKeyClient(kern.keyStore, kern.Logger)
	}
	return nil
}

// LoadLoggerFromConfig adds a logging configuration to the kernel
func (kern *Kernel) LoadLoggerFromConfig(conf *logconfig.LoggingConfig) error {
	logger, err := conf.NewLogger()
	kern.SetLogger(logger)
	return err
}

// LoadExecutionOptionsFromConfig builds the execution options for the kernel
func (kern *Kernel) LoadExecutionOptionsFromConfig(conf *execution.ExecutionConfig) error {
	if conf != nil {
		exeOptions, err := conf.ExecutionOptions()
		if err != nil {
			return err
		}
		kern.exeOptions = exeOptions
		kern.timeoutFactor = conf.TimeoutFactor
	}
	return nil
}

// LoadTendermintFromConfig loads our consensus engine into the kernel
func (kern *Kernel) LoadTendermintFromConfig(conf *config.BurrowConfig, privVal tmTypes.PrivValidator) (err error) {
	if conf.Tendermint == nil || !conf.Tendermint.Enabled {
		return nil
	}

	authorizedPeersProvider := conf.Tendermint.DefaultAuthorizedPeersProvider()
	if conf.Tendermint.IdentifyPeers {
		authorizedPeersProvider = registry.NewNodeFilter(kern.State)
	}

	kern.database.Stats()

	kern.info = fmt.Sprintf("Burrow_%s_%s_ValidatorID:%X", project.History.CurrentVersion().String(),
		kern.Blockchain.ChainID(), privVal.GetPubKey().Address())

	app := abci.NewApp(kern.info, kern.Blockchain, kern.State, kern.checker, kern.committer, kern.txCodec,
		authorizedPeersProvider, kern.Panic, kern.Logger)

	// We could use this to provide/register our own metrics (though this will register them with us). Unfortunately
	// Tendermint currently ignores the metrics passed unless its own server is turned on.
	metricsProvider := node.DefaultMetricsProvider(&tmConfig.InstrumentationConfig{
		Prometheus:           false,
		PrometheusListenAddr: "",
	})

	genesisDoc := kern.Blockchain.GenesisDoc()

	tmGenesisDoc := tendermint.DeriveGenesisDoc(&genesisDoc, kern.Blockchain.AppHashAfterLastBlock())
	heightValuer := log.Valuer(func() interface{} { return kern.Blockchain.LastBlockHeight() })
	tmLogger := kern.Logger.With(structure.CallerKey, log.Caller(LoggingCallerDepth+1)).With("height", heightValuer)
	tmConf, err := conf.TendermintConfig()
	if err != nil {
		return fmt.Errorf("could not build Tendermint config: %v", err)
	}
	kern.Node, err = tendermint.NewNode(tmConf, privVal, tmGenesisDoc, app, metricsProvider, tmLogger)
	return err
}

// LoadKernelFromConfig builds and returns a Kernel based solely on the supplied configuration
func LoadKernelFromConfig(conf *config.BurrowConfig) (*Kernel, error) {
	kern, err := NewKernel(conf.BurrowDir)
	if err != nil {
		return nil, fmt.Errorf("could not create initial kernel: %v", err)
	}

	if err = kern.LoadLoggerFromConfig(conf.Logging); err != nil {
		return nil, fmt.Errorf("could not configure logger: %v", err)
	}

	err = kern.LoadKeysFromConfig(conf.Keys)
	if err != nil {
		return nil, fmt.Errorf("could not configure keys: %v", err)
	}

	err = kern.LoadExecutionOptionsFromConfig(conf.Execution)
	if err != nil {
		return nil, fmt.Errorf("could not add execution options: %v", err)
	}

	err = kern.LoadState(conf.GenesisDoc)
	if err != nil {
		return nil, fmt.Errorf("could not load state: %v", err)
	}

	if conf.ValidatorAddress == nil {
		return nil, fmt.Errorf("Address must be set")
	}

	privVal, err := kern.PrivValidator(*conf.ValidatorAddress)
	if err != nil {
		return nil, fmt.Errorf("could not form PrivValidator from Address: %v", err)
	}

	err = kern.LoadTendermintFromConfig(conf, privVal)
	if err != nil {
		return nil, fmt.Errorf("could not configure Tendermint: %v", err)
	}

	kern.AddProcesses(DefaultProcessLaunchers(kern, conf.RPC, conf.Keys)...)
	return kern, nil
}
