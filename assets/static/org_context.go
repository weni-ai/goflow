package static

import "github.com/nyaruka/goflow/assets"

type OrgContext struct {
	Context_     string             `json:"context" validate:"required,context"`
	ChannelUUID_ assets.ChannelUUID `json:"channel_uuid"`
}

func NewOrgContext(context string, channelUUID assets.ChannelUUID) assets.OrgContext {
	return &OrgContext{
		Context_:     context,
		ChannelUUID_: channelUUID,
	}
}

func (c *OrgContext) Context() string { return c.Context_ }

func (c *OrgContext) ChannelUUID() assets.ChannelUUID { return c.ChannelUUID_ }
