package gumble // import "github.com/101Bas/gumble/gumble"

// ContextActions is a map of ContextActions.
type ContextActions map[string]*ContextAction

func (c ContextActions) create(action string) *ContextAction {
	contextAction := &ContextAction{
		Name: action,
	}
	c[action] = contextAction
	return contextAction
}
