package helm

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"sigs.k8s.io/yaml"

	"github.com/gruntwork-io/go-commons/errors"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/testing"

	"os"

	"github.com/gonvenience/ytbx"
	"github.com/homeport/dyff/pkg/dyff"
)

// RenderTemplate runs `helm template` to render the template given the provided options and returns stdout/stderr from
// the template command. If you pass in templateFiles, this will only render those templates. This function will fail
// the test if there is an error rendering the template.
func RenderTemplate(t testing.TestingT, options *Options, chartDir string, releaseName string, templateFiles []string, extraHelmArgs ...string) string {
	out, err := RenderTemplateE(t, options, chartDir, releaseName, templateFiles, extraHelmArgs...)
	require.NoError(t, err)
	return out
}

// RenderTemplateE runs `helm template` to render the template given the provided options and returns stdout/stderr from
// the template command. If you pass in templateFiles, this will only render those templates.
func RenderTemplateE(t testing.TestingT, options *Options, chartDir string, releaseName string, templateFiles []string, extraHelmArgs ...string) (string, error) {
	// First, verify the charts dir exists
	absChartDir, err := filepath.Abs(chartDir)
	if err != nil {
		return "", errors.WithStackTrace(err)
	}
	if !files.FileExists(chartDir) {
		return "", errors.WithStackTrace(ChartNotFoundError{chartDir})
	}

	// check chart dependencies
	if options.BuildDependencies {
		if _, err := RunHelmCommandAndGetOutputE(t, options, "dependency", "build", chartDir); err != nil {
			return "", errors.WithStackTrace(err)
		}
	}

	// Now construct the args
	// We first construct the template args
	args := []string{}
	if options.KubectlOptions != nil && options.KubectlOptions.Namespace != "" {
		args = append(args, "--namespace", options.KubectlOptions.Namespace)
	}
	args, err = getValuesArgsE(t, options, args...)
	if err != nil {
		return "", err
	}
	for _, templateFile := range templateFiles {
		// validate this is a valid template file
		absTemplateFile := filepath.Join(absChartDir, templateFile)
		if !strings.HasPrefix(templateFile, "charts") && !files.FileExists(absTemplateFile) {
			return "", errors.WithStackTrace(TemplateFileNotFoundError{Path: templateFile, ChartDir: absChartDir})
		}

		// Note: we only get the abs template file path to check it actually exists, but the `helm template` command
		// expects the relative path from the chart.
		args = append(args, "--show-only", templateFile)
	}
	// deal extraHelmArgs
	args = append(args, extraHelmArgs...)

	// ... and add the name and chart at the end as the command expects
	args = append(args, releaseName, chartDir)

	// Finally, call out to helm template command
	return RunHelmCommandAndGetStdOutE(t, options, "template", args...)
}

// RenderTemplate runs `helm template` to render a *remote* chart  given the provided options and returns stdout/stderr from
// the template command. If you pass in templateFiles, this will only render those templates. This function will fail
// the test if there is an error rendering the template.
func RenderRemoteTemplate(t testing.TestingT, options *Options, chartURL string, releaseName string, templateFiles []string, extraHelmArgs ...string) string {
	out, err := RenderRemoteTemplateE(t, options, chartURL, releaseName, templateFiles, extraHelmArgs...)
	require.NoError(t, err)
	return out
}

// RenderTemplate runs `helm template` to render a *remote* helm chart  given the provided options and returns stdout/stderr from
// the template command. If you pass in templateFiles, this will only render those templates.
func RenderRemoteTemplateE(t testing.TestingT, options *Options, chartURL string, releaseName string, templateFiles []string, extraHelmArgs ...string) (string, error) {
	// Now construct the args
	// We first construct the template args
	args := []string{}
	if options.KubectlOptions != nil && options.KubectlOptions.Namespace != "" {
		args = append(args, "--namespace", options.KubectlOptions.Namespace)
	}
	args, err := getValuesArgsE(t, options, args...)
	if err != nil {
		return "", err
	}
	for _, templateFile := range templateFiles {
		// As the helm command fails if a non valid template is given as input
		// we do not check if the template file exists or not as we do for local charts
		// as it would add unecessary networking calls
		args = append(args, "--show-only", templateFile)
	}
	// deal extraHelmArgs
	args = append(args, extraHelmArgs...)

	// ... and add the helm chart name, the remote repo and chart URL at the end
	args = append(args, releaseName, "--repo", chartURL)

	// Finally, call out to helm template command
	return RunHelmCommandAndGetStdOutE(t, options, "template", args...)
}

// UnmarshalK8SYaml is the same as UnmarshalK8SYamlE, but will fail the test if there is an error.
func UnmarshalK8SYaml(t testing.TestingT, yamlData string, destinationObj interface{}) {
	require.NoError(t, UnmarshalK8SYamlE(t, yamlData, destinationObj))
}

// UnmarshalK8SYamlE can be used to take template outputs and unmarshal them into the corresponding client-go struct. For
// example, suppose you render the template into a Deployment object. You can unmarshal the yaml as follows:
//
// var deployment appsv1.Deployment
// UnmarshalK8SYamlE(t, renderedOutput, &deployment)
//
// At the end of this, the deployment variable will be populated.
func UnmarshalK8SYamlE(t testing.TestingT, yamlData string, destinationObj interface{}) error {
	// NOTE: the client-go library can only decode json, so we will first convert the yaml to json before unmarshaling
	jsonData, err := yaml.YAMLToJSON([]byte(yamlData))
	if err != nil {
		return errors.WithStackTrace(err)
	}
	err = json.Unmarshal(jsonData, destinationObj)
	if err != nil {
		return errors.WithStackTrace(err)
	}
	return nil
}

// UpdateSnapshot creates or updates the k8s manifest snapshot of a chart (e.g bitnami/nginx).
// It is one of the two functions needed to implement snapshot based testing for helm.
// see https://github.com/gruntwork-io/terratest/issues/1377
// A snapshot is used to compare the current manifests of a chart with the previous manifests.
// A global diff is run against the two snapshosts and the number of differences is returned.
func UpdateSnapshot(t testing.TestingT, options *Options, yamlData string, releaseName string) {
	require.NoError(t, UpdateSnapshotE(t, options, yamlData, releaseName))
}

// UpdateSnapshotE creates or updates the k8s manifest snapshot of a chart (e.g bitnami/nginx).
// It is one of the two functions needed to implement snapshot based testing for helm.
// see https://github.com/gruntwork-io/terratest/issues/1377
// A snapshot is used to compare the current manifests of a chart with the previous manifests.
// A global diff is run against the two snapshosts and the number of differences is returned.
// It will failed the test if there is an error while writing the manifests' snapshot in the file system
func UpdateSnapshotE(t testing.TestingT, options *Options, yamlData string, releaseName string) error {

	var snapshotDir = "__snapshot__"
	if options.SnapshotPath != "" {
		snapshotDir = options.SnapshotPath
	}
	// Create a directory if not exists
	if !files.FileExists(snapshotDir) {
		if err := os.Mkdir(snapshotDir, 0755); err != nil {
			return errors.WithStackTrace(err)
		}
	}

	filename := filepath.Join(snapshotDir, releaseName+".yaml")
	// Open a file in write mode
	file, err := os.Create(filename)
	if err != nil {
		return errors.WithStackTrace(err)
	}
	defer file.Close()

	// Write the k8s manifest into the file
	if _, err = file.WriteString(yamlData); err != nil {
		return errors.WithStackTrace(err)
	}

	if options.Logger != nil {
		options.Logger.Logf(t, "helm chart manifest written into file: %s", filename)
	}
	return nil
}

// DiffAgainstSnapshot compare the current manifests of a chart (e.g bitnami/nginx)
// with the previous manifests stored in the snapshot.
// see https://github.com/gruntwork-io/terratest/issues/1377
// It returns the number of difference between the two manifests or -1 in case of error
// It will fail the test if there is an error while reading or writing the two manifests in the file system
func DiffAgainstSnapshot(t testing.TestingT, options *Options, yamlData string, releaseName string) int {
	numberOfDiffs, err := DiffAgainstSnapshotE(t, options, yamlData, releaseName)
	require.NoError(t, err)
	return numberOfDiffs
}

// DiffAgainstSnapshotE compare the current manifests of a chart (e.g bitnami/nginx)
// with the previous manifests stored in the snapshot.
// see https://github.com/gruntwork-io/terratest/issues/1377
// It returns the number of difference between the manifests or -1 in case of error
func DiffAgainstSnapshotE(t testing.TestingT, options *Options, yamlData string, releaseName string) (int, error) {

	var snapshotDir = "__snapshot__"
	if options.SnapshotPath != "" {
		snapshotDir = options.SnapshotPath
	}

	// load the yaml snapshot file
	snapshot := filepath.Join(snapshotDir, releaseName+".yaml")
	from, err := ytbx.LoadFile(snapshot)
	if err != nil {
		return -1, errors.WithStackTrace(err)
	}

	// write the current manifest into a file as `dyff` does not support string input
	currentManifests := releaseName + ".yaml"
	file, err := os.Create(currentManifests)
	if err != nil {
		return -1, errors.WithStackTrace(err)
	}

	if _, err = file.WriteString(yamlData); err != nil {
		return -1, errors.WithStackTrace(err)
	}
	defer file.Close()
	defer os.Remove(currentManifests)

	to, err := ytbx.LoadFile(currentManifests)
	if err != nil {
		return -1, errors.WithStackTrace(err)
	}

	// compare the two manifests using `dyff`
	compOpt := dyff.KubernetesEntityDetection(false)

	// create a report
	report, err := dyff.CompareInputFiles(from, to, compOpt)
	if err != nil {
		return -1, errors.WithStackTrace(err)
	}

	// write any difference to stdout
	reportWriter := &dyff.HumanReport{
		Report:            report,
		DoNotInspectCerts: false,
		NoTableStyle:      false,
		OmitHeader:        false,
		UseGoPatchPaths:   false,
	}

	err = reportWriter.WriteReport(os.Stdout)
	if err != nil {
		return -1, errors.WithStackTrace(err)
	}
	// return the number of diffs to use in assertion while testing: 0 = no differences
	return len(reportWriter.Diffs), nil
}
