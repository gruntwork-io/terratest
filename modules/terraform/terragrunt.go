package terraform

import (
	"github.com/gruntwork-io/terratest/modules/terragrunt"
)

func toTerragruntOptions(options Options) *terragrunt.Options {
	opt := terragrunt.Options{
		TerragruntBinary:         options.TerraformBinary,
		TerragruntDir:            options.TerraformDir,
		EnvVars:                  options.EnvVars,
		Logger:                   options.Logger,
		MaxRetries:               options.MaxRetries,
		TimeBetweenRetries:       options.TimeBetweenRetries,
		RetryableTerraformErrors: options.RetryableTerraformErrors,
		WarningsAsErrors:         options.WarningsAsErrors,
		BackendConfig:            options.BackendConfig,
		PluginDir:                options.PluginDir,
		Stdin:                    options.Stdin,
		Vars:                     options.Vars,
		VarFiles:                 options.VarFiles,
		SetVarsAfterVarFiles:     options.SetVarsAfterVarFiles,
		PlanFilePath:             options.PlanFilePath,
		Targets:                  options.Targets,
		Lock:                     options.Lock,
		LockTimeout:              options.LockTimeout,
		NoColor:                  options.NoColor,
		ExtraArgs: terragrunt.ExtraArgs{
			Apply:   options.ExtraArgs.Apply,
			Destroy: options.ExtraArgs.Destroy,
			Plan:    options.ExtraArgs.Plan,
		},
	}
	for _, v := range options.MixedVars {
		opt.MixedVars = append(opt.MixedVars, v)
	}
	return &opt
}
