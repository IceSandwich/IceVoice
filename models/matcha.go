package models

type Matcha struct {
	IModel
}

func (m *Matcha) GetRequireFiles() []string {
	return []string{"gogogog"}
}
