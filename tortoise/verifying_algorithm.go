package tortoise

import (
	"context"
	"sync"
	"time"

	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spacemeshos/go-spacemesh/log"
	"github.com/spacemeshos/go-spacemesh/mesh"
)

// ThreadSafeVerifyingTortoise is a thread safe verifying tortoise wrapper, it just locks all actions.
type ThreadSafeVerifyingTortoise struct {
	trtl          *turtle
	logger        log.Log
	rerunInterval time.Duration
	lastRerun     time.Time
	mutex         sync.RWMutex
}

// Config holds the arguments and dependencies to create a verifying tortoise instance.
type Config struct {
	LayerSize       int
	Database        blockDataProvider
	Hdist           int // hare lookback distance: the distance over which we use the input vector/hare results
	Zdist           int // hare result wait distance: the distance over which we're willing to wait for hare results
	ConfidenceParam int // confidence wait distance: how long we wait for global consensus to be established
	WindowSize      int // tortoise sliding window: how many layers we store data for
	Log             log.Log
	Recovered       bool
	RerunInterval   time.Duration // how often to rerun from genesis
}

// NewVerifyingTortoise creates a new verifying tortoise wrapper
func NewVerifyingTortoise(ctx context.Context, cfg Config) *ThreadSafeVerifyingTortoise {
	if cfg.Recovered {
		return recoveredVerifyingTortoise(cfg.Database, cfg.Log)
	}
	return verifyingTortoise(
		ctx,
		cfg.LayerSize,
		cfg.Database,
		cfg.Hdist,
		cfg.Zdist,
		cfg.ConfidenceParam,
		cfg.WindowSize,
		cfg.RerunInterval,
		cfg.Log)
}

// verifyingTortoise creates a new verifying tortoise wrapper
func verifyingTortoise(
	ctx context.Context,
	layerSize int,
	mdb blockDataProvider,
	hdist,
	zdist,
	confidenceParam,
	windowSize int,
	rerunInterval time.Duration,
	logger log.Log,
) *ThreadSafeVerifyingTortoise {
	if hdist < zdist {
		logger.With().Panic("hdist must be >= zdist", log.Int("hdist", hdist), log.Int("zdist", zdist))
	}
	alg := &ThreadSafeVerifyingTortoise{
		trtl: newTurtle(mdb, hdist, zdist, confidenceParam, windowSize, layerSize),
	}
	alg.logger = logger
	alg.rerunInterval = rerunInterval
	alg.lastRerun = time.Now()
	alg.trtl.SetLogger(logger.WithFields(log.String("rerun", "false")))
	alg.trtl.init(ctx, mesh.GenesisLayer())
	return alg
}

// NewRecoveredVerifyingTortoise recovers a previously persisted tortoise copy from mesh.DB
func recoveredVerifyingTortoise(mdb blockDataProvider, logger log.Log) *ThreadSafeVerifyingTortoise {
	tmp, err := RecoverVerifyingTortoise(mdb)
	if err != nil {
		logger.With().Panic("could not recover tortoise state from disk", log.Err(err))
	}

	trtl := tmp.(*turtle)

	logger.Info("recovered tortoise from disk")
	trtl.bdp = mdb
	trtl.logger = logger

	return &ThreadSafeVerifyingTortoise{trtl: trtl}
}

// LatestComplete returns the latest verified layer. TODO: rename?
func (trtl *ThreadSafeVerifyingTortoise) LatestComplete() types.LayerID {
	trtl.mutex.RLock()
	verified := trtl.trtl.Verified
	trtl.mutex.RUnlock()
	return verified
}

// BaseBlock chooses a base block and creates a differences list. needs the hare results for latest layers.
func (trtl *ThreadSafeVerifyingTortoise) BaseBlock(ctx context.Context) (types.BlockID, [][]types.BlockID, error) {
	trtl.mutex.Lock()
	block, diffs, err := trtl.trtl.BaseBlock(ctx)
	trtl.mutex.Unlock()
	if err != nil {
		return types.BlockID{}, nil, err
	}
	return block, diffs, err
}

// simple wrapper for thread safety and reading old and new pbase values
func (trtl *ThreadSafeVerifyingTortoise) runAndReportLayers(ctx context.Context, fn func()) (types.LayerID, types.LayerID) {
	trtl.mutex.Lock()
	defer trtl.mutex.Unlock()
	oldVerified := trtl.trtl.Verified
	fn()
	newVerified := trtl.trtl.Verified
	return oldVerified, newVerified
}

// HandleLateBlocks processes votes and goodness for late blocks (for late block definition see white paper)
// returns the old verified layer and new verified layer after taking into account the blocks votes
func (trtl *ThreadSafeVerifyingTortoise) HandleLateBlocks(ctx context.Context, blocks []*types.Block) (types.LayerID, types.LayerID) {
	return trtl.runAndReportLayers(ctx, func() {
		if err := trtl.trtl.ProcessNewBlocks(ctx, blocks); err != nil {
			trtl.logger.WithContext(ctx).With().Error("tortoise errored handling late blocks", log.Err(err))
		}
	})
}

// HandleIncomingLayer processes all layer block votes
// returns the old verified layer and new verified layer after taking into account the blocks votes
func (trtl *ThreadSafeVerifyingTortoise) HandleIncomingLayer(ctx context.Context, layerID types.LayerID) (types.LayerID, types.LayerID) {
	return trtl.runAndReportLayers(ctx, func() {
		if err := trtl.trtl.HandleIncomingLayer(ctx, layerID); err != nil {
			trtl.logger.WithContext(ctx).With().Error("tortoise errored handling incoming layer", log.Err(err))
		}

		// rerun if needed
		trtl.rerunIfNeeded(ctx)
	})
}

// trigger a rerun from genesis once in a while
func (trtl *ThreadSafeVerifyingTortoise) rerunIfNeeded(ctx context.Context) {
	// TODO: in future we can do something more sophisticated, using accounting to determine when enough changes to old
	//   layers have accumulated (in terms of block weight) that our opinion could actually change. For now, we do the
	//   Simplest Possible Thing (TM) and just rerun from genesis once in a while. This requires a different instance of
	//   tortoise since we don't want to mess with the state of the main tortoise. We re-stream layer data from genesis
	//   using the sliding window, simulating a full resync.
	// TODO: should this happen "in the background" in a separate goroutine? Should it hold the mutex?
	if time.Now().Sub(trtl.lastRerun) > trtl.rerunInterval {
		logger := trtl.logger.WithContext(ctx)
		logger.With().Info("triggering tortoise full rerun from genesis",
			log.Duration("rerun_interval", trtl.rerunInterval),
			log.Time("last_rerun", trtl.lastRerun))

		// start from scratch with a new tortoise instance for each rerun
		trtlForRerun := trtl.trtl.cloneTurtle()
		trtlForRerun.SetLogger(logger.WithFields(log.String("rerun", "true")))
		trtlForRerun.init(ctx, mesh.GenesisLayer())

		for layerID := types.GetEffectiveGenesis(); layerID < trtl.trtl.Last; layerID++ {
			logger.With().Debug("rerunning tortoise for layer", layerID)
			if err := trtlForRerun.HandleIncomingLayer(ctx, layerID); err != nil {
				logger.With().Error("tortoise rerun errored", log.Err(err))
				break
			}
		}

		trtl.lastRerun = time.Now()
	}
}

// Persist saves a copy of the current tortoise state to the database
func (trtl *ThreadSafeVerifyingTortoise) Persist(ctx context.Context) error {
	trtl.mutex.Lock()
	defer trtl.mutex.Unlock()
	trtl.trtl.logger.WithContext(ctx).Info("persist tortoise")
	return trtl.trtl.persist()
}
