package models

type IModel interface {
	GetRequireFiles() []string
}

func Create(architecture string) IModel {
	switch architecture {
	case "matcha":
		return &Matcha{}
	}
	return nil
}
