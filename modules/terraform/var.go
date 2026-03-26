package terraform

type Var interface {
	Args() []string
	internal()
}

func VarInline(name string, value any) Var {
	return varInline{name: name, value: value}
}

type varInline struct {
	value any
	name  string
}

func (vi varInline) Args() []string {
	m := map[string]any{vi.name: vi.value}
	return formatTerraformArgs(m, "-var", true, false)
}
func (vi varInline) internal() {}

func VarFile(path string) Var {
	return varFile(path)
}

type varFile string

func (vf varFile) Args() []string {
	return []string{"-var-file", string(vf)}
}
func (vf varFile) internal() {}
