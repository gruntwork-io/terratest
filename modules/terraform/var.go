package terraform

// Var represents a Terraform variable assignment that can be rendered as
// command-line arguments. Use [VarInline] to pass an inline value with -var,
// or [VarFile] to reference a -var-file on disk.
type Var interface {
	// Args returns the command-line arguments that pass this variable to Terraform.
	Args() []string
	internal()
}

// VarInline returns a [Var] that passes the given name/value pair to Terraform
// via a -var flag.
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

// VarFile returns a [Var] that passes the file at the given path to Terraform
// via a -var-file flag.
func VarFile(path string) Var {
	return varFile(path)
}

type varFile string

func (vf varFile) Args() []string {
	return []string{"-var-file", string(vf)}
}
func (vf varFile) internal() {}
