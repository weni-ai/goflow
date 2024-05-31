package static

import (
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
)

type OrgContext struct {
	Context_     string             `json:"context" validate:"required,context"`
	ChannelUUID_ assets.ChannelUUID `json:"channel_uuid"`
	ProjectUUID_ uuids.UUID         `json:"project_uuid"`
}

func NewOrgContext(context string, channelUUID assets.ChannelUUID, projectUUID uuids.UUID) assets.OrgContext {
	return &OrgContext{
		Context_:     context,
		ChannelUUID_: channelUUID,
		ProjectUUID_: projectUUID,
	}
}

func (c *OrgContext) Context() string { return c.Context_ }

func (c *OrgContext) ChannelUUID() assets.ChannelUUID { return c.ChannelUUID_ }

func (c *OrgContext) ProjectUUID() uuids.UUID { return c.ProjectUUID_ }
