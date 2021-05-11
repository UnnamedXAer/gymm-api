package mocks

type MockMailer struct{}

func (m *MockMailer) Send(recipients []string, subject, data []byte) {
}

func (m *MockMailer) Close() {

}
