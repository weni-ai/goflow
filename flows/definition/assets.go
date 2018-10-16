package definition

import (
	"sync"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

// implemention of FlowAssets which provides lazy loading and validation of flows
type flowAssets struct {
	source        assets.AssetSource
	sessionAssets flows.SessionAssets

	mutex             sync.Mutex
	byUUID            map[assets.FlowUUID]flows.Flow
	validationStarted map[assets.FlowUUID]bool
}

// NewFlowAssets creates a new flow assets
func NewFlowAssets(source assets.AssetSource, sessionAssets flows.SessionAssets) flows.FlowAssets {
	return &flowAssets{
		source:        source,
		sessionAssets: sessionAssets,

		byUUID:            make(map[assets.FlowUUID]flows.Flow),
		validationStarted: make(map[assets.FlowUUID]bool),
	}
}

// Get returns the flow with the given UUID
func (a *flowAssets) Get(uuid assets.FlowUUID) (flows.Flow, error) {
	flow, shouldValidate, err := a.get(uuid)
	if err != nil {
		return nil, err
	}

	if shouldValidate {
		if err := flow.Validate(a.sessionAssets); err != nil {
			return nil, err
		}
	}

	return flow, nil
}

// gets the flow and determines whether the caller should validate it. Validation doesn't occur
// in this method because it
func (a *flowAssets) get(uuid assets.FlowUUID) (flows.Flow, bool, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	flow := a.byUUID[uuid]
	if flow != nil {
		return flow, false, nil
	}

	asset, err := a.source.Flow(uuid)
	if err != nil {
		return nil, false, err
	}

	flow, err = ReadFlow(asset.Definition())
	if err != nil {
		return nil, false, err
	}

	a.byUUID[flow.UUID()] = flow

	shouldValidate := !a.validationStarted[flow.UUID()]
	a.validationStarted[flow.UUID()] = true

	return flow, shouldValidate, nil
}
