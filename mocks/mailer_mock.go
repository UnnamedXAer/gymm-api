package mocks

type MockMailer struct{}

func (m *MockMailer) Send(recipients []string, data []byte) {
}

func (m *MockMailer) Close() {

}
