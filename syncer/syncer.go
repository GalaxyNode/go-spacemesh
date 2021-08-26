package syncer

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spacemeshos/go-spacemesh/common/util"
	"github.com/spacemeshos/go-spacemesh/events"
	"github.com/spacemeshos/go-spacemesh/layerfetcher"
	"github.com/spacemeshos/go-spacemesh/log"
	"github.com/spacemeshos/go-spacemesh/mesh"
	"github.com/spacemeshos/go-spacemesh/p2p/peers"
)

type layerTicker interface {
	GetCurrentLayer() types.LayerID
	LayerToTime(types.LayerID) time.Time
}

type layerFetcher interface {
	PollLayerHash(ctx context.Context, layerID types.LayerID) chan layerfetcher.LayerHashResult
	PollLayerBlocks(ctx context.Context, layerID types.LayerID, hashes map[types.Hash32][]peers.Peer) chan layerfetcher.LayerPromiseResult
	GetEpochATXs(ctx context.Context, id types.EpochID) error
	GetTortoiseBeacon(ctx context.Context, id types.EpochID) error // TODO(nkryuchkov): remove
	SetTortoiseBeacon(ctx context.Context, epochID types.EpochID, beacon types.Hash32) error
}

// Configuration is the config params for syncer
type Configuration struct {
	SyncInterval time.Duration
	// the sync process will try to validate the current layer if ValidationDelta has elapsed.
	ValidationDelta     time.Duration
	TBCalculationLayers uint32
	AlwaysListen        bool
}

const (
	outOfSyncThreshold  uint32 = 3 // see notSynced
	numGossipSyncLayers uint32 = 2 // see gossipSync
)

type syncState uint32

const (
	// notSynced is the state where the node is outOfSyncThreshold layers or more behind the current layer.
	notSynced syncState = iota
	// gossipSync is the state in which a node listens to at least one full layer of gossip before participating
	// in the protocol. this is to protect the node from participating in the consensus without full information.
	// for example, when a node wakes up in the middle of layer N, since it didn't receive all relevant messages and
	// blocks of layer N, it shouldn't vote or produce blocks in layer N+1. it instead listens to gossip for all
	// through layer N+1 and starts producing blocks and participates in hare committee in layer N+2
	gossipSync
	// synced is the state where the node is in sync with its peers.
	synced
)

func (s syncState) String() string {
	switch s {
	case notSynced:
		return "notSynced"
	case gossipSync:
		return "gossipSync"
	case synced:
		return "synced"
	default:
		return "unknown"
	}
}

// Syncer is responsible to keep the node in sync with the network.
type Syncer struct {
	logger log.Log

	conf     Configuration
	ticker   layerTicker
	mesh     *mesh.Mesh
	fetcher  layerFetcher
	syncOnce sync.Once
	// access via atomic.[Load|Store]Uint32
	syncState syncState
	// access via atomic.[Load|Store]Uint32
	isBusy    uint32
	syncTimer *time.Ticker
	// targetSyncedLayer is used to signal at which layer we can set this node to synced state
	targetSyncedLayer unsafe.Pointer

	tbCount         map[string]int
	seenTB          map[string]struct{}
	mostUsedTBCount int
	mostUsedTB      []byte

	// awaitSyncedCh is the list of subscribers' channels to notify when this node enters synced state
	awaitSyncedCh []chan struct{}
	awaitSyncedMu sync.Mutex

	shutdownCtx context.Context
	cancelFunc  context.CancelFunc

	// recording the run # since started. for logging/debugging only.
	run uint64
}

// NewSyncer creates a new Syncer instance.
func NewSyncer(ctx context.Context, conf Configuration, ticker layerTicker, mesh *mesh.Mesh, fetcher layerFetcher, logger log.Log) *Syncer {
	shutdownCtx, cancel := context.WithCancel(ctx)
	return &Syncer{
		logger:            logger,
		conf:              conf,
		ticker:            ticker,
		mesh:              mesh,
		fetcher:           fetcher,
		syncState:         notSynced,
		syncTimer:         time.NewTicker(conf.SyncInterval),
		targetSyncedLayer: unsafe.Pointer(&types.LayerID{}),
		tbCount:           make(map[string]int),
		seenTB:            make(map[string]struct{}),
		awaitSyncedCh:     make([]chan struct{}, 0),
		shutdownCtx:       shutdownCtx,
		cancelFunc:        cancel,
	}
}

// Close stops the syncing process and the goroutines syncer spawns.
func (s *Syncer) Close() {
	// TODO: ensure goroutines are all terminated before shutting down
	s.cancelFunc()
}

// RegisterChForSynced registers ch for notification when the node enters synced state
func (s *Syncer) RegisterChForSynced(ctx context.Context, ch chan struct{}) {
	if s.IsSynced(ctx) {
		close(ch)
		return
	}
	s.awaitSyncedMu.Lock()
	defer s.awaitSyncedMu.Unlock()
	s.awaitSyncedCh = append(s.awaitSyncedCh, ch)
}

// ListenToGossip returns true if the node is listening to gossip for blocks/TXs/ATXs data
func (s *Syncer) ListenToGossip() bool {
	return s.conf.AlwaysListen || s.getSyncState() >= gossipSync
}

// IsSynced returns true if the node is in synced state
func (s *Syncer) IsSynced(ctx context.Context) bool {
	// TODO: at startup, ctx contains no sessionId here
	res := s.getSyncState() == synced
	s.logger.WithContext(ctx).With().Info("node sync state",
		log.Bool("synced", res),
		log.FieldNamed("current", s.ticker.GetCurrentLayer()),
		log.FieldNamed("latest", s.mesh.LatestLayer()),
		log.FieldNamed("processed", s.mesh.ProcessedLayer()))
	return res
}

// Start starts the main sync loop that tries to sync data for every SyncInterval
func (s *Syncer) Start(ctx context.Context) {
	s.syncOnce.Do(func() {
		if s.ticker.GetCurrentLayer().Uint32() <= 1 {
			s.setSyncState(ctx, synced)
		}
		for {
			select {
			case <-s.shutdownCtx.Done():
				s.logger.WithContext(ctx).Info("stopping sync to shutdown")
				return
			case <-s.syncTimer.C:
				s.synchronize(ctx)
			}
		}
	})
}

// ForceSync manually starts a sync process outside the main sync loop. If the node is already running a sync process,
// ForceSync will be ignored.
func (s *Syncer) ForceSync(ctx context.Context) {
	s.logger.WithContext(ctx).Debug("executing ForceSync")
	go s.synchronize(ctx)
}

func (s *Syncer) isClosed() bool {
	select {
	case <-s.shutdownCtx.Done():
		return true
	default:
		return false
	}
}

func (s *Syncer) getSyncState() syncState {
	return (syncState)(atomic.LoadUint32((*uint32)(&s.syncState)))
}

func (s *Syncer) setSyncState(ctx context.Context, newState syncState) {
	oldState := syncState(atomic.SwapUint32((*uint32)(&s.syncState), uint32(newState)))
	if oldState != newState {
		s.logger.WithContext(ctx).With().Info("sync state change",
			log.String("from_state", oldState.String()),
			log.String("to_state", newState.String()),
			log.FieldNamed("current", s.ticker.GetCurrentLayer()),
			log.FieldNamed("latest", s.mesh.LatestLayer()),
			log.FieldNamed("processed", s.mesh.ProcessedLayer()))
		events.ReportNodeStatusUpdate()
		if newState != synced {
			return
		}
		// notify subscribes
		s.awaitSyncedMu.Lock()
		defer s.awaitSyncedMu.Unlock()
		for _, ch := range s.awaitSyncedCh {
			close(ch)
		}
		s.awaitSyncedCh = make([]chan struct{}, 0)
	}
}

// setSyncerBusy returns false if the syncer is already running a sync process.
// otherwise it sets syncer to be busy and returns true.
func (s *Syncer) setSyncerBusy() bool {
	return atomic.CompareAndSwapUint32(&s.isBusy, 0, 1)
}

func (s *Syncer) setSyncerIdle() {
	atomic.StoreUint32(&s.isBusy, 0)
}

// targetSyncedLayer is used to signal at which layer we can set this node to synced state
func (s *Syncer) setTargetSyncedLayer(ctx context.Context, layerID types.LayerID) {
	oldSyncLayer := *(*types.LayerID)(atomic.SwapPointer(&s.targetSyncedLayer, unsafe.Pointer(&layerID)))
	s.logger.WithContext(ctx).With().Info("target synced layer changed",
		log.Uint32("from_layer", oldSyncLayer.Uint32()),
		log.Uint32("to_layer", layerID.Uint32()),
		log.FieldNamed("current", s.ticker.GetCurrentLayer()),
		log.FieldNamed("latest", s.mesh.LatestLayer()),
		log.FieldNamed("processed", s.mesh.ProcessedLayer()))
}

func (s *Syncer) getTargetSyncedLayer() types.LayerID {
	return *(*types.LayerID)(atomic.LoadPointer(&s.targetSyncedLayer))
}

func (s *Syncer) synchronize(ctx context.Context) bool {
	logger := s.logger.WithContext(ctx)

	if s.isClosed() {
		logger.Warning("attempting to sync while shutting down")
		return false
	}

	// at most one synchronize process can run at any time
	if !s.setSyncerBusy() {
		logger.Info("sync is already running, giving up")
		return false
	}
	// no need to worry about race condition for s.run. only one instance of synchronize can run at a time
	s.run++
	logger.With().Info(fmt.Sprintf("starting sync run #%v", s.run),
		log.String("sync_state", s.getSyncState().String()),
		log.FieldNamed("current", s.ticker.GetCurrentLayer()),
		log.FieldNamed("latest", s.mesh.LatestLayer()),
		log.FieldNamed("processed", s.mesh.ProcessedLayer()))

	s.setStateBeforeSync(ctx)
	// start a dedicated process for validation.
	// do not use a unbuffered channel for vQueue. we don't want it to block if the receiver isn't ready. i.e.
	// if validation for the last layer is still running
	vQueue := make(chan *types.Layer, s.ticker.GetCurrentLayer().Uint32())
	vDone := make(chan struct{})
	go s.startValidating(ctx, s.run, vQueue, vDone)

	success := false
	defer func() {
		close(vQueue)
		<-vDone
		s.setStateAfterSync(ctx, success)
		logger.With().Info(fmt.Sprintf("finished sync run #%v", s.run),
			log.Bool("success", success),
			log.String("sync_state", s.getSyncState().String()),
			log.FieldNamed("current", s.ticker.GetCurrentLayer()),
			log.FieldNamed("latest", s.mesh.LatestLayer()),
			log.FieldNamed("processed", s.mesh.ProcessedLayer()))
		s.setSyncerIdle()
	}()

	// using ProcessedLayer() instead of LatestLayer() so we can validate layers on a best-efforts basis.
	// our clock starts ticking from 1 so it is safe to skip layer 0
	// always sync to currentLayer-1 to reduce race with gossip and hare/tortoise
	for layerID := s.mesh.ProcessedLayer().Add(1); layerID.Before(s.ticker.GetCurrentLayer()); layerID = layerID.Add(1) {
		if layerID.FirstInEpoch() {
			s.tbCount = make(map[string]int)
			s.seenTB = make(map[string]struct{})
			s.mostUsedTBCount = 0
			s.mostUsedTB = nil
		}

		layer, err := s.syncLayer(ctx, layerID)
		if err != nil {
			logger.With().Error("failed to sync to layer", layerID, log.Err(err))
			return false
		}

		if layerID.OrdinalInEpoch() < s.conf.TBCalculationLayers {
			s.processTortoiseBeacons(ctx, layer)
		} else if layerID.OrdinalInEpoch() == s.conf.TBCalculationLayers {
			epoch := layerID.GetEpoch()
			if err := s.fetcher.SetTortoiseBeacon(ctx, epoch, types.BytesToHash(s.mostUsedTB)); err != nil {
				logger.With().Error("failed to write synced tortoise beacon into DB",
					log.Uint32("epoch", uint32(epoch)),
					log.String("beacon", util.Bytes2Hex(s.mostUsedTB)))
			}
		}

		if len(layer.Blocks()) == 0 {
			logger.With().Info("setting layer to zero-block", layerID)
			if err := s.mesh.SetZeroBlockLayer(layerID); err != nil {
				logger.With().Error("failed to set zero-block for layer", layerID, log.Err(err))
			}
		}

		if s.shouldValidateLayer(layerID) {
			vQueue <- layer
		}
		logger.With().Debug("finished data sync", layerID)
	}

	logger.With().Debug("data is synced, waiting for validation",
		log.FieldNamed("current", s.ticker.GetCurrentLayer()),
		log.FieldNamed("latest", s.mesh.LatestLayer()),
		log.FieldNamed("processed", s.mesh.ProcessedLayer()))
	success = true

	return true
}

func (s *Syncer) setStateBeforeSync(ctx context.Context) {
	current := s.ticker.GetCurrentLayer()
	if current.Uint32() <= 1 {
		s.setSyncState(ctx, synced)
		return
	}
	latest := s.mesh.LatestLayer()
	if current.After(latest) && current.Difference(latest) >= outOfSyncThreshold {
		s.logger.WithContext(ctx).With().Info("node is too far behind",
			log.FieldNamed("current", current),
			log.FieldNamed("latest", latest),
			log.Uint32("behind_threshold", outOfSyncThreshold))
		s.setSyncState(ctx, notSynced)
	}
}

func (s *Syncer) setStateAfterSync(ctx context.Context, success bool) {
	if !success {
		s.setSyncState(ctx, notSynced)
		return
	}
	currSyncState := s.getSyncState()
	current := s.ticker.GetCurrentLayer()
	// if we have gossip-synced to the target synced layer, we are ready to participate in consensus
	if currSyncState == gossipSync && !s.getTargetSyncedLayer().After(current) {
		s.setSyncState(ctx, synced)
	} else if currSyncState == notSynced {
		// wait till s.ticker.GetCurrentLayer() + numGossipSyncLayers to participate in consensus
		s.setSyncState(ctx, gossipSync)
		s.setTargetSyncedLayer(ctx, current.Add(numGossipSyncLayers))
	}
}

func (s *Syncer) syncLayer(ctx context.Context, layerID types.LayerID) (*types.Layer, error) {
	if s.isClosed() {
		return nil, errors.New("shutdown")
	}

	layer, err := s.getLayerFromPeers(ctx, layerID)
	if err != nil {
		return nil, err
	}

	if err := s.getATXs(ctx, layerID); err != nil {
		return nil, err
	}

	// TODO(nkryuchkov): remove
	//if err := s.getTortoiseBeacon(ctx, layerID); err != nil {
	//	return nil, err
	//}

	return layer, nil
}

func (s *Syncer) getLayerFromPeers(ctx context.Context, layerID types.LayerID) (*types.Layer, error) {
	ch := s.fetcher.PollLayerHash(ctx, layerID)
	hashRes := <-ch
	if hashRes.Err != nil {
		return nil, fmt.Errorf("PollLayerHash: %w", hashRes.Err)
	}
	// TODO: resolve hash with peers
	hashes := make(map[types.Hash32][]peers.Peer)
	for lyrHash, peers := range hashRes.Hashes {
		hashes[lyrHash.Hash] = peers
	}
	bch := s.fetcher.PollLayerBlocks(ctx, layerID, hashes)
	res := <-bch
	if res.Err != nil {
		if res.Err == layerfetcher.ErrZeroLayer {
			return types.NewLayer(layerID), nil
		}
		return nil, fmt.Errorf("PollLayerBlocks: %w", res.Err)
	}

	layer, err := s.mesh.GetLayer(layerID)
	if err != nil {
		return nil, fmt.Errorf("GetLayer: %w", err)
	}

	return layer, nil
}

func (s *Syncer) getATXs(ctx context.Context, layerID types.LayerID) error {
	if layerID.GetEpoch() == 0 {
		s.logger.WithContext(ctx).Info("skip getting atx in epoch 0")
		return nil
	}
	epoch := layerID.GetEpoch()
	atCurrentEpoch := epoch == s.ticker.GetCurrentLayer().GetEpoch()
	atLastLayerOfEpoch := layerID == (epoch + 1).FirstLayer().Sub(1)
	// only get ATXs if
	// - layerID is in the current epoch
	// - layerID is the last layer of a previous epoch
	// i.e. for older epochs we sync ATXs once per epoch. for current epoch we sync ATXs in every layer
	if atCurrentEpoch || atLastLayerOfEpoch {
		s.logger.WithContext(ctx).With().Debug("getting atxs", epoch, layerID)
		ctx = log.WithNewRequestID(ctx, layerID.GetEpoch())
		if err := s.fetcher.GetEpochATXs(ctx, epoch); err != nil {
			// dont fail sync if we cannot fetch atxs for the current epoch before the last layer
			if !atCurrentEpoch || atLastLayerOfEpoch {
				s.logger.WithContext(ctx).With().Error("failed to fetch epoch atxs", layerID, epoch, log.Err(err))
				return err
			}
			s.logger.WithContext(ctx).With().Warning("failed to fetch epoch atxs", layerID, epoch, log.Err(err))
		}
	}
	return nil
}

func (s *Syncer) processTortoiseBeacons(ctx context.Context, layer *types.Layer) {
	for _, b := range layer.Blocks() {
		if _, ok := s.seenTB[b.MinerID().String()]; ok {
			// TODO(nkryuchkov): invalid block, miner attempted to include beacon in two blocks within the same epoch
			// TODO(nkryuchkov): handle this case
			s.logger.Warning("miner attempted to include beacon in two blocks within the same epoch")
		}

		s.seenTB[b.MinerID().String()] = struct{}{}

		tb := b.EligibilityProof.TortoiseBeacon
		tbString := string(tb)

		s.tbCount[tbString]++
		if s.tbCount[tbString] > s.mostUsedTBCount {
			s.mostUsedTBCount = s.tbCount[tbString]
			s.mostUsedTB = tb
		}
	}
}

func (s *Syncer) getTortoiseBeacon(ctx context.Context, layerID types.LayerID) error {
	epoch := layerID.GetEpoch()
	if epoch.IsGenesis() {
		s.logger.WithContext(ctx).Info("skip getting tortoise beacons in genesis epoch")
		return nil
	}

	currentEpoch := s.ticker.GetCurrentLayer().GetEpoch()
	// only get tortoise beacon if
	// - layerID is in the current epoch
	// - layerID is the last layer of a previous epoch
	// i.e. for older epochs we sync tortoise beacons once per epoch. for current epoch we sync tortoise beacons in every layer
	if epoch == currentEpoch || ((epoch+1).FirstLayer().Value > 0 && layerID == (epoch+1).FirstLayer().Sub(1)) {
		s.logger.WithContext(ctx).With().Debug("getting tortoise beacons", epoch, layerID)
		ctx = log.WithNewRequestID(ctx, layerID.GetEpoch())
		if err := s.fetcher.GetTortoiseBeacon(ctx, epoch); err != nil {
			s.logger.WithContext(ctx).With().Error("failed to fetch epoch tortoise beacons",
				layerID,
				epoch,
				log.Err(err))
			return err
		}
	}
	return nil
}

// always returns true if layerID is an old layer.
// for current layer, only returns true if current layer already elapsed ValidationDelta
func (s *Syncer) shouldValidateLayer(layerID types.LayerID) bool {
	if layerID == types.NewLayerID(0) {
		return false
	}
	current := s.ticker.GetCurrentLayer()
	return layerID.Before(current) || time.Now().Sub(s.ticker.LayerToTime(current)) > s.conf.ValidationDelta
}

// start a dedicated process to validate layers one by one
func (s *Syncer) startValidating(ctx context.Context, run uint64, queue chan *types.Layer, done chan struct{}) {
	logger := s.logger.WithContext(ctx).WithName("validation")
	logger.Debug("validation started for run #%v", run)
	defer func() {
		logger.Debug("validation done for run #%v", run)
		close(done)
	}()
	for layer := range queue {
		if s.isClosed() {
			return
		}
		s.validateLayer(ctx, layer)
	}
}

func (s *Syncer) validateLayer(ctx context.Context, layer *types.Layer) {
	if s.isClosed() {
		s.logger.WithContext(ctx).Error("shutting down")
		return
	}

	s.logger.WithContext(ctx).With().Debug("validating layer",
		layer.Index(),
		log.String("blocks", fmt.Sprint(types.BlockIDs(layer.Blocks()))))

	// TODO: re-architect this so the syncer does not need to actually wait for tortoise to finish running.
	//   It should be sufficient to call GetLayer (above), and maybe, to queue a request to tortoise to analyze this
	//   layer (without waiting for this to finish -- it should be able to run async).
	//   See https://github.com/spacemeshos/go-spacemesh/issues/2415
	s.mesh.ValidateLayer(ctx, layer)
}
