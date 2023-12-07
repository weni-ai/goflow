package assets

import (
	"fmt"

	"github.com/nyaruka/gocommon/uuids"
)

type OrgContext interface {
	Context() string
	ChannelUUID() ChannelUUID
}

type OrgContextReference struct {
	Context string `json:"context"`
	UUID    string `json:"uuid"`
}

func NewOrgContextReference(orgContext string) *OrgContextReference {
	return &OrgContextReference{Context: orgContext}
}

func (r *OrgContextReference) Type() string {
	return "org_context"
}

func (r *OrgContextReference) GenericUUID() uuids.UUID {
	return uuids.UUID(r.UUID)
}

func (r *OrgContextReference) Identity() string {
	return string(r.UUID)
}

func (r *OrgContextReference) Variable() bool {
	return false
}

func (r *OrgContextReference) String() string {
	return fmt.Sprintf("%s[uuid=%s,context=%s]", r.Type(), r.Identity(), r.Context)
}

var _ UUIDReference = (*OrgContextReference)(nil)
