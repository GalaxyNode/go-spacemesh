// Code generated by MockGen. DO NOT EDIT.
// Source: ./layers.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	network "github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-core/peer"
	protocol "github.com/libp2p/go-libp2p-core/protocol"
	types "github.com/spacemeshos/go-spacemesh/common/types"
)

// MockatxHandler is a mock of atxHandler interface.
type MockatxHandler struct {
	ctrl     *gomock.Controller
	recorder *MockatxHandlerMockRecorder
}

// MockatxHandlerMockRecorder is the mock recorder for MockatxHandler.
type MockatxHandlerMockRecorder struct {
	mock *MockatxHandler
}

// NewMockatxHandler creates a new mock instance.
func NewMockatxHandler(ctrl *gomock.Controller) *MockatxHandler {
	mock := &MockatxHandler{ctrl: ctrl}
	mock.recorder = &MockatxHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockatxHandler) EXPECT() *MockatxHandlerMockRecorder {
	return m.recorder
}

// HandleAtxData mocks base method.
func (m *MockatxHandler) HandleAtxData(arg0 context.Context, arg1 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleAtxData", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleAtxData indicates an expected call of HandleAtxData.
func (mr *MockatxHandlerMockRecorder) HandleAtxData(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleAtxData", reflect.TypeOf((*MockatxHandler)(nil).HandleAtxData), arg0, arg1)
}

// MockblockHandler is a mock of blockHandler interface.
type MockblockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockblockHandlerMockRecorder
}

// MockblockHandlerMockRecorder is the mock recorder for MockblockHandler.
type MockblockHandlerMockRecorder struct {
	mock *MockblockHandler
}

// NewMockblockHandler creates a new mock instance.
func NewMockblockHandler(ctrl *gomock.Controller) *MockblockHandler {
	mock := &MockblockHandler{ctrl: ctrl}
	mock.recorder = &MockblockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockblockHandler) EXPECT() *MockblockHandlerMockRecorder {
	return m.recorder
}

// HandleBlockData mocks base method.
func (m *MockblockHandler) HandleBlockData(arg0 context.Context, arg1 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleBlockData", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleBlockData indicates an expected call of HandleBlockData.
func (mr *MockblockHandlerMockRecorder) HandleBlockData(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleBlockData", reflect.TypeOf((*MockblockHandler)(nil).HandleBlockData), arg0, arg1)
}

// MockballotHandler is a mock of ballotHandler interface.
type MockballotHandler struct {
	ctrl     *gomock.Controller
	recorder *MockballotHandlerMockRecorder
}

// MockballotHandlerMockRecorder is the mock recorder for MockballotHandler.
type MockballotHandlerMockRecorder struct {
	mock *MockballotHandler
}

// NewMockballotHandler creates a new mock instance.
func NewMockballotHandler(ctrl *gomock.Controller) *MockballotHandler {
	mock := &MockballotHandler{ctrl: ctrl}
	mock.recorder = &MockballotHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockballotHandler) EXPECT() *MockballotHandlerMockRecorder {
	return m.recorder
}

// HandleBallotData mocks base method.
func (m *MockballotHandler) HandleBallotData(arg0 context.Context, arg1 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleBallotData", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleBallotData indicates an expected call of HandleBallotData.
func (mr *MockballotHandlerMockRecorder) HandleBallotData(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleBallotData", reflect.TypeOf((*MockballotHandler)(nil).HandleBallotData), arg0, arg1)
}

// MockproposalHandler is a mock of proposalHandler interface.
type MockproposalHandler struct {
	ctrl     *gomock.Controller
	recorder *MockproposalHandlerMockRecorder
}

// MockproposalHandlerMockRecorder is the mock recorder for MockproposalHandler.
type MockproposalHandlerMockRecorder struct {
	mock *MockproposalHandler
}

// NewMockproposalHandler creates a new mock instance.
func NewMockproposalHandler(ctrl *gomock.Controller) *MockproposalHandler {
	mock := &MockproposalHandler{ctrl: ctrl}
	mock.recorder = &MockproposalHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockproposalHandler) EXPECT() *MockproposalHandlerMockRecorder {
	return m.recorder
}

// HandleProposalData mocks base method.
func (m *MockproposalHandler) HandleProposalData(arg0 context.Context, arg1 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleProposalData", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleProposalData indicates an expected call of HandleProposalData.
func (mr *MockproposalHandlerMockRecorder) HandleProposalData(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleProposalData", reflect.TypeOf((*MockproposalHandler)(nil).HandleProposalData), arg0, arg1)
}

// MocktxHandler is a mock of txHandler interface.
type MocktxHandler struct {
	ctrl     *gomock.Controller
	recorder *MocktxHandlerMockRecorder
}

// MocktxHandlerMockRecorder is the mock recorder for MocktxHandler.
type MocktxHandlerMockRecorder struct {
	mock *MocktxHandler
}

// NewMocktxHandler creates a new mock instance.
func NewMocktxHandler(ctrl *gomock.Controller) *MocktxHandler {
	mock := &MocktxHandler{ctrl: ctrl}
	mock.recorder = &MocktxHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocktxHandler) EXPECT() *MocktxHandlerMockRecorder {
	return m.recorder
}

// HandleSyncTransaction mocks base method.
func (m *MocktxHandler) HandleSyncTransaction(arg0 context.Context, arg1 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleSyncTransaction", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleSyncTransaction indicates an expected call of HandleSyncTransaction.
func (mr *MocktxHandlerMockRecorder) HandleSyncTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleSyncTransaction", reflect.TypeOf((*MocktxHandler)(nil).HandleSyncTransaction), arg0, arg1)
}

// MocklayerDB is a mock of layerDB interface.
type MocklayerDB struct {
	ctrl     *gomock.Controller
	recorder *MocklayerDBMockRecorder
}

// MocklayerDBMockRecorder is the mock recorder for MocklayerDB.
type MocklayerDBMockRecorder struct {
	mock *MocklayerDB
}

// NewMocklayerDB creates a new mock instance.
func NewMocklayerDB(ctrl *gomock.Controller) *MocklayerDB {
	mock := &MocklayerDB{ctrl: ctrl}
	mock.recorder = &MocklayerDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocklayerDB) EXPECT() *MocklayerDBMockRecorder {
	return m.recorder
}

// GetAggregatedLayerHash mocks base method.
func (m *MocklayerDB) GetAggregatedLayerHash(arg0 types.LayerID) types.Hash32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAggregatedLayerHash", arg0)
	ret0, _ := ret[0].(types.Hash32)
	return ret0
}

// GetAggregatedLayerHash indicates an expected call of GetAggregatedLayerHash.
func (mr *MocklayerDBMockRecorder) GetAggregatedLayerHash(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAggregatedLayerHash", reflect.TypeOf((*MocklayerDB)(nil).GetAggregatedLayerHash), arg0)
}

// GetHareConsensusOutput mocks base method.
func (m *MocklayerDB) GetHareConsensusOutput(arg0 types.LayerID) (types.BlockID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHareConsensusOutput", arg0)
	ret0, _ := ret[0].(types.BlockID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHareConsensusOutput indicates an expected call of GetHareConsensusOutput.
func (mr *MocklayerDBMockRecorder) GetHareConsensusOutput(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHareConsensusOutput", reflect.TypeOf((*MocklayerDB)(nil).GetHareConsensusOutput), arg0)
}

// GetLayerHash mocks base method.
func (m *MocklayerDB) GetLayerHash(arg0 types.LayerID) types.Hash32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLayerHash", arg0)
	ret0, _ := ret[0].(types.Hash32)
	return ret0
}

// GetLayerHash indicates an expected call of GetLayerHash.
func (mr *MocklayerDBMockRecorder) GetLayerHash(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLayerHash", reflect.TypeOf((*MocklayerDB)(nil).GetLayerHash), arg0)
}

// LayerBallotIDs mocks base method.
func (m *MocklayerDB) LayerBallotIDs(arg0 types.LayerID) ([]types.BallotID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LayerBallotIDs", arg0)
	ret0, _ := ret[0].([]types.BallotID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LayerBallotIDs indicates an expected call of LayerBallotIDs.
func (mr *MocklayerDBMockRecorder) LayerBallotIDs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LayerBallotIDs", reflect.TypeOf((*MocklayerDB)(nil).LayerBallotIDs), arg0)
}

// LayerBlockIds mocks base method.
func (m *MocklayerDB) LayerBlockIds(arg0 types.LayerID) ([]types.BlockID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LayerBlockIds", arg0)
	ret0, _ := ret[0].([]types.BlockID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LayerBlockIds indicates an expected call of LayerBlockIds.
func (mr *MocklayerDBMockRecorder) LayerBlockIds(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LayerBlockIds", reflect.TypeOf((*MocklayerDB)(nil).LayerBlockIds), arg0)
}

// ProcessedLayer mocks base method.
func (m *MocklayerDB) ProcessedLayer() types.LayerID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessedLayer")
	ret0, _ := ret[0].(types.LayerID)
	return ret0
}

// ProcessedLayer indicates an expected call of ProcessedLayer.
func (mr *MocklayerDBMockRecorder) ProcessedLayer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessedLayer", reflect.TypeOf((*MocklayerDB)(nil).ProcessedLayer))
}

// SaveHareConsensusOutput mocks base method.
func (m *MocklayerDB) SaveHareConsensusOutput(arg0 context.Context, arg1 types.LayerID, arg2 types.BlockID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveHareConsensusOutput", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveHareConsensusOutput indicates an expected call of SaveHareConsensusOutput.
func (mr *MocklayerDBMockRecorder) SaveHareConsensusOutput(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveHareConsensusOutput", reflect.TypeOf((*MocklayerDB)(nil).SaveHareConsensusOutput), arg0, arg1, arg2)
}

// SetZeroBlockLayer mocks base method.
func (m *MocklayerDB) SetZeroBlockLayer(arg0 types.LayerID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetZeroBlockLayer", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetZeroBlockLayer indicates an expected call of SetZeroBlockLayer.
func (mr *MocklayerDBMockRecorder) SetZeroBlockLayer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetZeroBlockLayer", reflect.TypeOf((*MocklayerDB)(nil).SetZeroBlockLayer), arg0)
}

// MockatxIDsDB is a mock of atxIDsDB interface.
type MockatxIDsDB struct {
	ctrl     *gomock.Controller
	recorder *MockatxIDsDBMockRecorder
}

// MockatxIDsDBMockRecorder is the mock recorder for MockatxIDsDB.
type MockatxIDsDBMockRecorder struct {
	mock *MockatxIDsDB
}

// NewMockatxIDsDB creates a new mock instance.
func NewMockatxIDsDB(ctrl *gomock.Controller) *MockatxIDsDB {
	mock := &MockatxIDsDB{ctrl: ctrl}
	mock.recorder = &MockatxIDsDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockatxIDsDB) EXPECT() *MockatxIDsDBMockRecorder {
	return m.recorder
}

// GetEpochAtxs mocks base method.
func (m *MockatxIDsDB) GetEpochAtxs(epochID types.EpochID) ([]types.ATXID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEpochAtxs", epochID)
	ret0, _ := ret[0].([]types.ATXID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEpochAtxs indicates an expected call of GetEpochAtxs.
func (mr *MockatxIDsDBMockRecorder) GetEpochAtxs(epochID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEpochAtxs", reflect.TypeOf((*MockatxIDsDB)(nil).GetEpochAtxs), epochID)
}

// MockpoetDB is a mock of poetDB interface.
type MockpoetDB struct {
	ctrl     *gomock.Controller
	recorder *MockpoetDBMockRecorder
}

// MockpoetDBMockRecorder is the mock recorder for MockpoetDB.
type MockpoetDBMockRecorder struct {
	mock *MockpoetDB
}

// NewMockpoetDB creates a new mock instance.
func NewMockpoetDB(ctrl *gomock.Controller) *MockpoetDB {
	mock := &MockpoetDB{ctrl: ctrl}
	mock.recorder = &MockpoetDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockpoetDB) EXPECT() *MockpoetDBMockRecorder {
	return m.recorder
}

// ValidateAndStoreMsg mocks base method.
func (m *MockpoetDB) ValidateAndStoreMsg(data []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateAndStoreMsg", data)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateAndStoreMsg indicates an expected call of ValidateAndStoreMsg.
func (mr *MockpoetDBMockRecorder) ValidateAndStoreMsg(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateAndStoreMsg", reflect.TypeOf((*MockpoetDB)(nil).ValidateAndStoreMsg), data)
}

// Mocknetwork is a mock of network interface.
type Mocknetwork struct {
	ctrl     *gomock.Controller
	recorder *MocknetworkMockRecorder
}

// MocknetworkMockRecorder is the mock recorder for Mocknetwork.
type MocknetworkMockRecorder struct {
	mock *Mocknetwork
}

// NewMocknetwork creates a new mock instance.
func NewMocknetwork(ctrl *gomock.Controller) *Mocknetwork {
	mock := &Mocknetwork{ctrl: ctrl}
	mock.recorder = &MocknetworkMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mocknetwork) EXPECT() *MocknetworkMockRecorder {
	return m.recorder
}

// GetPeers mocks base method.
func (m *Mocknetwork) GetPeers() []peer.ID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPeers")
	ret0, _ := ret[0].([]peer.ID)
	return ret0
}

// GetPeers indicates an expected call of GetPeers.
func (mr *MocknetworkMockRecorder) GetPeers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPeers", reflect.TypeOf((*Mocknetwork)(nil).GetPeers))
}

// Network mocks base method.
func (m *Mocknetwork) Network() network.Network {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Network")
	ret0, _ := ret[0].(network.Network)
	return ret0
}

// Network indicates an expected call of Network.
func (mr *MocknetworkMockRecorder) Network() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Network", reflect.TypeOf((*Mocknetwork)(nil).Network))
}

// NewStream mocks base method.
func (m *Mocknetwork) NewStream(arg0 context.Context, arg1 peer.ID, arg2 ...protocol.ID) (network.Stream, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "NewStream", varargs...)
	ret0, _ := ret[0].(network.Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewStream indicates an expected call of NewStream.
func (mr *MocknetworkMockRecorder) NewStream(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewStream", reflect.TypeOf((*Mocknetwork)(nil).NewStream), varargs...)
}

// PeerCount mocks base method.
func (m *Mocknetwork) PeerCount() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PeerCount")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// PeerCount indicates an expected call of PeerCount.
func (mr *MocknetworkMockRecorder) PeerCount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PeerCount", reflect.TypeOf((*Mocknetwork)(nil).PeerCount))
}

// SetStreamHandler mocks base method.
func (m *Mocknetwork) SetStreamHandler(arg0 protocol.ID, arg1 network.StreamHandler) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetStreamHandler", arg0, arg1)
}

// SetStreamHandler indicates an expected call of SetStreamHandler.
func (mr *MocknetworkMockRecorder) SetStreamHandler(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStreamHandler", reflect.TypeOf((*Mocknetwork)(nil).SetStreamHandler), arg0, arg1)
}
