package testutils

type MockShell struct {
	OutputToReturn    string
	ErrorToReturn     error
	ReceivedArguments []string
}

func (s *MockShell) RunWithArgs(args ...string) (string, error) {
	s.ReceivedArguments = args
	return s.OutputToReturn, s.ErrorToReturn
}
