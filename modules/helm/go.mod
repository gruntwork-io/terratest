module github.com/gruntwork-io/terratest/modules/helm

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/files v0.55.0
	github.com/gruntwork-io/terratest/modules/http-helper v0.55.0
	github.com/gruntwork-io/terratest/modules/k8s v0.55.0
	github.com/gruntwork-io/terratest/modules/logger v0.55.0
	github.com/gruntwork-io/terratest/modules/random v0.55.0
	github.com/gruntwork-io/terratest/modules/shell v0.55.0
	github.com/gruntwork-io/terratest/modules/testing v0.55.0
)

require (
	github.com/BurntSushi/toml v1.5.0 // indirect
	github.com/gonvenience/bunt v1.4.2 // indirect
	github.com/gonvenience/idem v0.0.2 // indirect
	github.com/gonvenience/neat v1.3.16 // indirect
	github.com/gonvenience/term v1.0.4 // indirect
	github.com/gonvenience/text v1.0.9 // indirect
	github.com/gonvenience/ytbx v1.4.7 // indirect
	github.com/homeport/dyff v1.10.2 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-ciede2000 v0.0.0-20170301095244-782e8c62fec3 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-ps v1.0.0 // indirect
	github.com/mitchellh/hashstructure v1.1.0 // indirect
	github.com/sergi/go-diff v1.4.0 // indirect
	github.com/texttheater/golang-levenshtein v1.0.1 // indirect
	github.com/virtuald/go-ordered-json v0.0.0-20170621173500-b18e6e673d74 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/term v0.34.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.28.4 // indirect
	k8s.io/apimachinery v0.28.4 // indirect
	k8s.io/client-go v0.28.4 // indirect
)

replace (
	github.com/gruntwork-io/terratest/modules/files => ../files
	github.com/gruntwork-io/terratest/modules/http-helper => ../http-helper
	github.com/gruntwork-io/terratest/modules/k8s => ../k8s
	github.com/gruntwork-io/terratest/modules/logger => ../logger
	github.com/gruntwork-io/terratest/modules/random => ../random
	github.com/gruntwork-io/terratest/modules/shell => ../shell
	github.com/gruntwork-io/terratest/modules/testing => ../testing
)
