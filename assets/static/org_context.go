package static

import (
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
)

type OrgContext struct {
	Context_         string             `json:"context" validate:"required,context"`
	ChannelUUID_     assets.ChannelUUID `json:"channel_uuid"`
	ProjectUUID_     uuids.UUID         `json:"project_uuid"`
	HasVtexAds_      bool               `json:"vtex_ads"`
	HideUnavailable_ bool               `json:"hide_unavailable"`
	ExtraPrompt_     string             `json:"extra_prompt"`
}

func NewOrgContext(context string, channelUUID assets.ChannelUUID, projectUUID uuids.UUID, hasVtexAds bool, hideUnavailable bool, extraPrompt string) assets.OrgContext {
	return &OrgContext{
		Context_:         context,
		ChannelUUID_:     channelUUID,
		ProjectUUID_:     projectUUID,
		HasVtexAds_:      hasVtexAds,
		HideUnavailable_: hideUnavailable,
		ExtraPrompt_:     extraPrompt,
	}
}

func (c *OrgContext) Context() string { return c.Context_ }

func (c *OrgContext) ChannelUUID() assets.ChannelUUID { return c.ChannelUUID_ }

func (c *OrgContext) ProjectUUID() uuids.UUID { return c.ProjectUUID_ }

func (c *OrgContext) HasVtexAds() bool { return c.HasVtexAds_ }

func (c *OrgContext) HideUnavailable() bool { return c.HideUnavailable_ }

func (c *OrgContext) ExtraPrompt() string { return c.ExtraPrompt_ }
