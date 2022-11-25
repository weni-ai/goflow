package flows

import "github.com/nyaruka/goflow/assets"

// ExternalService represents an third party service integration.
type ExternalService struct {
	assets.ExternalService
}

// NewExternalService returns a new external service object from the given external service asset.
func NewExternalService(asset assets.ExternalService) *ExternalService {
	return &ExternalService{ExternalService: asset}
}

// Asset returns the underlying asset
func (e *ExternalService) Asset() assets.ExternalService { return e.ExternalService }

// Reference returns a reference to this external service
func (e *ExternalService) Reference() *assets.ExternalServiceReference {
	return assets.NewExternalServiceReference(e.UUID(), e.Name())
}

// ExternalServiceAssets provides access to all external services assets
type ExternalServiceAssets struct {
	byUUID map[assets.ExternalServiceUUID]*ExternalService
}

// NewExternalServiceAssets creates a new set of external service assets
func NewExternalServiceAssets(externalServices []assets.ExternalService) *ExternalServiceAssets {
	s := &ExternalServiceAssets{
		byUUID: make(map[assets.ExternalServiceUUID]*ExternalService, len(externalServices)),
	}
	for _, asset := range externalServices {
		s.byUUID[asset.UUID()] = NewExternalService(asset)
	}
	return s
}

// Get returns the external service with the given UUID
func (s *ExternalServiceAssets) Get(uuid assets.ExternalServiceUUID) *ExternalService {
	return s.byUUID[uuid]
}
