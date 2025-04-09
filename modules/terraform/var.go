package terraform

type Var interface {
	Args() []string
	internal()
}

func VarInline(name string, value interface{}) Var {
	return varInline{name: name, value: value}
}

type varInline struct {
	name  string
	value interface{}
}

func (vi varInline) Args() []string {
	m := map[string]interface{}{vi.name: vi.value}
	return formatTerraformArgs(m, "-var", true)
}
func (vi varInline) internal() {}

func VarFile(path string) Var {
	return varFile(path)
}

type varFile string

func (vf varFile) Args() []string {
	return []string{"-var-file", string(vf)}
}
func (vi varFile) internal() {}
