package assets

import "fmt"

type CodeActionID string

type CodeAction interface {
	ID() CodeActionID
	Name() string
}

type CodeActionReference struct {
	ID   CodeActionID `json:"id" validate:"required"`
	Name string       `json:"name"`
}

func NewCodeActionReference(id CodeActionID, name string) *CodeActionReference {
	return &CodeActionReference{ID: id, Name: name}
}

func (r *CodeActionReference) Type() string {
	return "code_action"
}

func (r *CodeActionReference) Identity() string {
	return string(r.ID)
}

func (r *CodeActionReference) String() string {
	return fmt.Sprintf("%s[id=%s, name=%s]", r.Type(), r.Identity(), r.Name)
}
