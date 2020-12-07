package nackdetector

import (
	"context"
	"sync"
	"time"

	envoy_service_discovery_v3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/solo-io/go-utils/contextutils"
	xds "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/server"
)

type State int

const (
	New State = iota
	InSync
	OutOfSync
	OutOfSyncNack
	Gone
)

func (s State) String() string {
	switch s {
	case New:
		return "New"
	case InSync:
		return "InSync"
	case OutOfSync:
		return "OutOfSync"
	case OutOfSyncNack:
		return "OutOfSyncNack"
	case Gone:
		return "Gone"
	default:
		return "unknown"
	}
}

type DiscoveryServiceId struct {
	GrpcStreamId int64
	TypeUrl      string
}

type TimeProvider interface {
	Now() time.Time
	After(time.Duration) <-chan time.Time
}

type stdTime struct{}

func (s stdTime) Now() time.Time {
	return time.Now()
}
func (s stdTime) After(t time.Duration) <-chan time.Time {
	return time.After(t)
}

type EnvoyStatusId struct {
	NodeId   string
	StreamId DiscoveryServiceId
}

type StateChangedHandler interface {
	StateChanged(ctx context.Context, eid EnvoyStatusId, oldState, newState State)
}

type StateChangedCallback func(ctx context.Context, eid EnvoyStatusId, oldState, newState State)

func (f StateChangedCallback) StateChanged(ctx context.Context, eid EnvoyStatusId, oldState, newState State) {
	f(ctx, eid, oldState, newState)
}

type EnvoyStatus struct {
	EnvoyStatusId EnvoyStatusId
	State         State
	LastModified  time.Time
}

type EnvoyState struct {
	ServerVersion string
	ServerNonce   string

	EnvoyStatus EnvoyStatus
}

type EnvoysState struct {
	values  map[int64]map[string]*EnvoyState
	maplock sync.RWMutex

	stateChangedHandler StateChangedHandler

	changes chan EnvoyStatus

	WaitTimeForSync time.Duration

	TimeProvider TimeProvider
}

func (ess *EnvoysState) Get(id DiscoveryServiceId) *EnvoyState {
	var stcopy *EnvoyState

	ess.maplock.RLock()
	submap := ess.values[id.GrpcStreamId]
	if submap != nil {
		if st, ok := submap[id.TypeUrl]; ok {
			clone := *st
			stcopy = &clone
		}
	}
	ess.maplock.RUnlock()
	if stcopy == nil {
		stcopy = &EnvoyState{
			EnvoyStatus: EnvoyStatus{
				EnvoyStatusId: EnvoyStatusId{
					StreamId: id,
				},
			},
		}
	}

	return stcopy
}

func (ess *EnvoysState) Delete(ctx context.Context, id int64) {
	ess.maplock.Lock()
	deadEnvoy := ess.values[id]
	delete(ess.values, id)
	ess.maplock.Unlock()

	// TODO:
	// envoy removed! notify about envoy out of sync?
	// we provide the stream id to make sure we can deal with a race with envoy reconnecting.
	for _, v := range deadEnvoy {
		v.EnvoyStatus.State = Gone
		ess.processStateChange(ctx, v)
	}

}
func (ess *EnvoysState) Set(id DiscoveryServiceId, es *EnvoyState) {
	copyEs := *es
	ess.maplock.Lock()
	if submap, ok := ess.values[id.GrpcStreamId]; ok {
		submap[id.TypeUrl] = &copyEs
	} else {
		ess.values[id.GrpcStreamId] = map[string]*EnvoyState{
			id.TypeUrl: &copyEs,
		}
	}
	ess.maplock.Unlock()
}

func (ess *EnvoysState) CheckIsSync(ctx context.Context, id DiscoveryServiceId, vi, rn string) {

	envoyState := ess.Get(id)
	if envoyState.ServerVersion == "" {
		if vi == "" {
			// no state, it means this is the first request so nothing else to do

		} else {
			// envoy has a state, but we dont. that probably because gloo restarted and envoy didnt.
			// assume that envoy state is good. gloo will correct it shortly if not.
			ess.inSync(ctx, id, envoyState)
		}
		return
	}
	contextutils.LoggerFrom(ctx).Debugf("nackdetector: CheckIsSync: versions: %s %s nonces: %s %s", envoyState.ServerVersion, vi, envoyState.ServerNonce, rn)

	if envoyState.ServerVersion != vi {
		// NACK maybe detected (if the nonce match)
		nack := envoyState.ServerNonce == rn
		ess.outOfSync(ctx, id, envoyState, nack)
	} else {
		// ACK detected
		ess.inSync(ctx, id, envoyState)
	}
}

func (ess *EnvoysState) outOfSync(ctx context.Context, id DiscoveryServiceId, es *EnvoyState, nack bool) {
	if nack {
		ess.preprocessStateChange(ctx, id, es, OutOfSyncNack)
	} else {
		ess.preprocessStateChange(ctx, id, es, OutOfSync)
	}
}

func (ess *EnvoysState) inSync(ctx context.Context, id DiscoveryServiceId, es *EnvoyState) {
	ess.preprocessStateChange(ctx, id, es, InSync)
}

func (ess *EnvoysState) newEnvoy(ctx context.Context, id DiscoveryServiceId, es *EnvoyState) {
	ess.preprocessStateChange(ctx, id, es, New)
}

func (ess *EnvoysState) preprocessStateChange(ctx context.Context, id DiscoveryServiceId, es *EnvoyState, st State) {
	es.EnvoyStatus.State = st
	es.EnvoyStatus.LastModified = ess.TimeProvider.Now()
	// save and notify
	ess.Set(id, es)
	ess.processStateChange(ctx, es)
}

func (ess *EnvoysState) processStateChange(ctx context.Context, es *EnvoyState) {
	var envoyStatus EnvoyStatus
	envoyStatus = es.EnvoyStatus
	select {
	case ess.changes <- envoyStatus:
	default:
		contextutils.LoggerFrom(ctx).Warn("changes channel full! this might mean that the system is loaded. consider running more gloos")
	}
}

func (ess *EnvoysState) runNotifications(ctx context.Context) {

	notifiedState := map[EnvoyStatusId]State{}
	deadlines := map[EnvoyStatusId]time.Time{}

	logger := contextutils.LoggerFrom(ctx)

	notifyIfNeeded := func(nodeID EnvoyStatusId, st State) {

		if oldst, ok := notifiedState[nodeID]; ok && oldst == st {
			return
		} else {
			logger.Debugw("notifyNewState", "nodeID", nodeID, "old-state", oldst, "state", st)
			defer logger.Debug("notifyNewState - done")
			notifiedState[nodeID] = st
			ess.stateChangedHandler.StateChanged(ctx, nodeID, oldst, st)
		}
	}

	for {
		select {
		case <-ctx.Done():
			return
		case es := <-ess.changes:

			logger.Debugw("received statchange event", "event", es)

			nodeID := es.EnvoyStatusId
			// remove it from the list of deadlines as we have an update.
			delete(deadlines, nodeID)

			switch es.State {
			case OutOfSync:
				if st, ok := notifiedState[nodeID]; ok && st == OutOfSync {
					continue
				}
				// envoy is out of sync. this happens as part of natural configuration delivery.
				// we should get another discovery-request soon signaling either a NACK or an ACK
				// we will check again after the deadline if envoy is still out of sync and notify about it.
				// this indicates a bug \ crashed envoy.
				deadlines[nodeID] = es.LastModified.Add(ess.WaitTimeForSync)

			// TODO(yuval-k): consider changing this to something more correct (i.e. wait for the next
			// minimal time period). this is good enough if WaitTimeForSync is small
			case Gone:
				notifyIfNeeded(nodeID, Gone)
				delete(notifiedState, nodeID)
			case New:
				fallthrough
			case OutOfSyncNack:
				// envoy nacked us notify immediatly
				fallthrough
			case InSync:
				// we are in sync all is good!
				notifyIfNeeded(nodeID, es.State)
			}

		case now := <-ess.TimeProvider.After(ess.WaitTimeForSync):
			for nodeID, t := range deadlines {
				if now.After(t) {
					logger.Debugw("notifying of out of sync envoy", "nodeID", nodeID)
					notifyIfNeeded(nodeID, OutOfSync)
				}
			}
		}
	}
}

func NewEnvoysState(ctx context.Context, stateChangedHandler StateChangedHandler) *EnvoysState {
	es := &EnvoysState{

		changes: make(chan EnvoyStatus, 100),

		WaitTimeForSync: time.Second,

		values: map[int64]map[string]*EnvoyState{},

		stateChangedHandler: stateChangedHandler,
		TimeProvider:        stdTime{},
	}
	go es.runNotifications(ctx)
	return es
}

type nackDetector struct {
	states              *EnvoysState
	stateChangedHandler StateChangedHandler

	ctx context.Context
}

func NewNackDetector(ctx context.Context, stChangedCB StateChangedHandler) *nackDetector {
	return &nackDetector{
		states: NewEnvoysState(ctx, stChangedCB),
	}
}

func NewNackDetectorWithEnvoysState(envoysState *EnvoysState) *nackDetector {
	return &nackDetector{states: envoysState}
}

func (n *nackDetector) OnStreamOpen(ctx context.Context, id int64, url string) error {
	contextutils.LoggerFrom(ctx).Debugf("nackdetector: envoy stream open %d url %s", id, url)
	return nil
}

func (n *nackDetector) OnStreamRequest(id int64, req *envoy_service_discovery_v3.DiscoveryRequest) error {
	contextutils.LoggerFrom(n.ctx).Debugf("nackdetector: envoy requested %s %s nonce %s", req.VersionInfo, req.TypeUrl, req.ResponseNonce)
	var dsid DiscoveryServiceId
	dsid.GrpcStreamId = id
	dsid.TypeUrl = req.TypeUrl

	envoyState := n.states.Get(dsid)
	if envoyState.EnvoyStatus.EnvoyStatusId.NodeId == "" {
		// first time we see this envoy!
		// set the id and type so that the callback can see them too.
		envoyState.EnvoyStatus.EnvoyStatusId.NodeId = req.Node.Id
		n.states.Set(dsid, envoyState)
		n.states.newEnvoy(n.ctx, dsid, envoyState)
	}

	n.states.CheckIsSync(n.ctx, dsid, req.VersionInfo, req.ResponseNonce)
	return nil
}

func (n *nackDetector) OnStreamResponse(id int64, req *envoy_service_discovery_v3.DiscoveryRequest, resp *envoy_service_discovery_v3.DiscoveryResponse) {
	var dsid DiscoveryServiceId
	dsid.GrpcStreamId = id
	dsid.TypeUrl = req.TypeUrl
	envoyState := n.states.Get(dsid)
	envoyState.ServerVersion = resp.VersionInfo
	envoyState.ServerNonce = resp.Nonce
	n.states.Set(dsid, envoyState)

	//  This check will mark envoy as out of sync... which covers the case of envoy not
	// requesting new config
	contextutils.LoggerFrom(n.ctx).Debugf("nackdetector: delivering new version %s nonce  %s to %s %s", resp.VersionInfo, resp.Nonce, envoyState.EnvoyStatus.EnvoyStatusId.NodeId, envoyState.EnvoyStatus.EnvoyStatusId.StreamId.TypeUrl)

	n.states.CheckIsSync(n.ctx, dsid, req.VersionInfo, req.ResponseNonce)
}

func (n *nackDetector) OnFetchRequest(ctx context.Context, request *envoy_service_discovery_v3.DiscoveryRequest) error {
	// gloo uses streaming api, so do nothing here
	return nil
}

func (n *nackDetector) OnFetchResponse(request *envoy_service_discovery_v3.DiscoveryRequest, response *envoy_service_discovery_v3.DiscoveryResponse) {
	// gloo uses streaming api, so do nothing here
}

func (n *nackDetector) OnStreamClosed(id int64) {
	n.states.Delete(n.ctx, id)
}

var _ xds.Callbacks = (*nackDetector)(nil)

/*
TODO:
use case:restart gloo
envoy connects and provides existing version and nonce.
gloo is ok with it. nacks detector not so much

*/
