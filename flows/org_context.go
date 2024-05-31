package flows

import (
	"github.com/nyaruka/goflow/assets"
)

type OrgContext struct {
	assets.OrgContext
}

func NewOrgContext(asset assets.OrgContext) *OrgContext {
	return &OrgContext{
		OrgContext: asset,
	}
}

func NewOrgContextAssets(orgContexts []assets.OrgContext) *OrgContextAssets {
	s := &OrgContextAssets{
		byChannelUUID: make(map[assets.ChannelUUID]*OrgContext, len(orgContexts)),
	}
	for _, asset := range orgContexts {
		s.byChannelUUID[asset.ChannelUUID()] = NewOrgContext(asset)
	}
	return s
}

func (c *OrgContextAssets) GetByChannelUUID() *OrgContext {
	for _, c := range c.byChannelUUID {
		if c.OrgContext.Context() != "" {
			return c
		}
	}
	return nil
}

func (c *OrgContextAssets) GetProjectUUIDByChannelUUID() *OrgContext {
	for _, c := range c.byChannelUUID {
		if c.OrgContext.ProjectUUID() != "" {
			return c
		}
	}
	return nil
}

func (c *OrgContext) Asset() assets.OrgContext { return c.OrgContext }

// Reference returns a reference to this context
func (c *OrgContext) Reference() *assets.OrgContextReference {
	if c == nil {
		return nil
	}
	return assets.NewOrgContextReference(c.Context(), c.ProjectUUID())
}

type OrgContextAssets struct {
	byChannelUUID map[assets.ChannelUUID]*OrgContext
}
