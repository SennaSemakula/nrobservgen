package hcl

type Variable struct {
	Name    string      `hcl:"name"`
	Type    string      `hcl:"type"`
	Default interface{} `hcl:"default"`
}

type Provider struct {
	Name   string `hcl:"name"`
	Source string `hcl:"source"`
}
type Backend struct {
	Name   string `hcl:"name"`
	Key    string `hcl:"key"`
	Bucket string `hcl:"bucket"`
}

type Module struct {
	Name      string
	Source    string
	Version   string
	Variables []*Variable
	LocalVars []*Variable
}
