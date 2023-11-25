package model

// LError is localized error that contains Localization Message Code and
// extractable template
// Check *Locale.Error for details
type LError struct {
	code     string
	template map[string]any

	err error
}

func NewLError(code string, tpl map[string]any, err error) *LError {
	return &LError{
		code:     code,
		template: tpl,
		err:      err,
	}
}

func (e *LError) GetCode() string {
	return e.code
}
func (e *LError) GetTemplate() map[string]any {
	return e.template
}
func (e *LError) GetError() error {
	return e.err
}

// Error returns stringified error,
// this is backup method that shouldn't happen for Localized Error in most cases
func (e *LError) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

// LErrorExpandable sub-type that marks that value need to be expanded to Localized Message
type LErrorExpandable string
