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

func (mc *OrgContext) Context() string { return mc.Context_ }

func (mc *OrgContext) ChannelUUID() assets.ChannelUUID { return mc.ChannelUUID_ }
