package model

// LError is localized error that contains Localization Message Code and
// extractable template
// Check *Locale.Error for details
type LError struct {
	Code     string
	Template map[string]any

	Err error
}

func NewLError(code string, tpl map[string]any, err error) *LError {
	return &LError{
		Code:     code,
		Template: tpl,
		Err:      err,
	}
}

// Error returns stringified error,
// this is backup method that shouldn't happen for Localized Error in most cases
func (e *LError) Error() string {
	if e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

// LErrorExpandable sub-type that marks that value need to be expanded to Localized Message
type LErrorExpandable string
