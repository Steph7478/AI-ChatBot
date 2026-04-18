package model

func (m *Model) Learn(input, response string) {
	m.Conversations[input] = response
	m.SaveConversation(input, response)
}

func (m *Model) LearnAndSave(input, response string) error {
	m.Learn(input, response)
	return m.SaveConversation(input, response)
}
