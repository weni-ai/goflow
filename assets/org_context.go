package assets

import (
	"fmt"

	"github.com/nyaruka/gocommon/uuids"
)

type OrgContext interface {
	Context() string
	ChannelUUID() ChannelUUID
	ProjectUUID() uuids.UUID
}

type OrgContextReference struct {
	Context     string     `json:"context"`
	UUID        string     `json:"uuid"`
	ProjectUUID uuids.UUID `json:"project_uuid"`
}

func NewOrgContextReference(orgContext string, projectUUID uuids.UUID) *OrgContextReference {
	return &OrgContextReference{Context: orgContext, ProjectUUID: projectUUID}
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
