package prompts

// IPrompts is the interface in which to interact with the user and accept
// input.
type IPrompts interface {
	UsernamePassword() (string, string, error)
	PHI() error
	YesNo(msg string) error
}

// SPrompts is a concrete implementation of IPrompts
type SPrompts struct{}

// New returns a new instance of IPrompts
func New() IPrompts {
	return &SPrompts{}
}

// UsernamePassword prompts a user to enter their username and password.
func (p *SPrompts) UsernamePassword() (string, string, error) {
	return "", "", nil
}

// PHI prompts a user to accept liability for downloading PHI to their local
// machine.
func (p *SPrompts) PHI() error {
	return nil
}

// PHI prompts a user to accept liability for downloading PHI to their local
// machine.
func (p *SPrompts) YesNo(msg string) error {
	return nil
}
