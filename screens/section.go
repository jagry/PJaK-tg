package screens

func NewSection(base Base, name string) Section {
	return Section{Base: base, name: name}
}

type Section struct {
	Base
	name string
}
