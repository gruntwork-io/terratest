package aws

import (
	"context"
	"os"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"
)

// RemoteFileSpecification describes which files you want to copy from your instances
type RemoteFileSpecification struct {
	RemotePathToFileFilter map[string][]string // A map of the files to fetch, where the keys are directories on the remote host and the values are filters for what files to fetch from the directory. The filters support bash-style wildcards.
	KeyPair                *Ec2Keypair
	SshUser                string   //nolint:staticcheck,revive // preserving existing field name
	LocalDestinationDir    string   // base path where to store downloaded artifacts locally. The final path of each resource will include the ip of the host and the name of the immediate parent folder.
	AsgNames               []string // ASGs where our instances will be
	UseSudo                bool
}

// FetchContentsOfFileFromInstanceContextE looks up the public IP address of the EC2 Instance with the given ID, connects to
// the Instance via SSH using the given username and Key Pair, fetches the contents of the file at the given path
// (using sudo if useSudo is true), and returns the contents of that file as a string.
// The ctx parameter supports cancellation and timeouts.
func FetchContentsOfFileFromInstanceContextE(t testing.TestingT, ctx context.Context, awsRegion string, sshUserName string, keyPair *Ec2Keypair, instanceID string, useSudo bool, filePath string) (string, error) {
	publicIP, err := GetPublicIPOfEc2InstanceContextE(t, ctx, instanceID, awsRegion)
	if err != nil {
		return "", err
	}

	host := ssh.Host{
		SshUserName: sshUserName,
		SshKeyPair:  keyPair.KeyPair,
		Hostname:    publicIP,
	}

	return ssh.FetchContentsOfFileContextE(t, ctx, &host, useSudo, filePath)
}

// FetchContentsOfFileFromInstanceContext looks up the public IP address of the EC2 Instance with the given ID, connects to
// the Instance via SSH using the given username and Key Pair, fetches the contents of the file at the given path
// (using sudo if useSudo is true), and returns the contents of that file as a string.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchContentsOfFileFromInstanceContext(t testing.TestingT, ctx context.Context, awsRegion string, sshUserName string, keyPair *Ec2Keypair, instanceID string, useSudo bool, filePath string) string {
	t.Helper()

	out, err := FetchContentsOfFileFromInstanceContextE(t, ctx, awsRegion, sshUserName, keyPair, instanceID, useSudo, filePath)
	require.NoError(t, err)

	return out
}

// FetchContentsOfFileFromInstance looks up the public IP address of the EC2 Instance with the given ID, connects to
// the Instance via SSH using the given username and Key Pair, fetches the contents of the file at the given path
// (using sudo if useSudo is true), and returns the contents of that file as a string.
//
// Deprecated: Use [FetchContentsOfFileFromInstanceContext] instead.
func FetchContentsOfFileFromInstance(t testing.TestingT, awsRegion string, sshUserName string, keyPair *Ec2Keypair, instanceID string, useSudo bool, filePath string) string {
	t.Helper()

	return FetchContentsOfFileFromInstanceContext(t, context.Background(), awsRegion, sshUserName, keyPair, instanceID, useSudo, filePath)
}

// FetchContentsOfFileFromInstanceE looks up the public IP address of the EC2 Instance with the given ID, connects to
// the Instance via SSH using the given username and Key Pair, fetches the contents of the file at the given path
// (using sudo if useSudo is true), and returns the contents of that file as a string.
//
// Deprecated: Use [FetchContentsOfFileFromInstanceContextE] instead.
func FetchContentsOfFileFromInstanceE(t testing.TestingT, awsRegion string, sshUserName string, keyPair *Ec2Keypair, instanceID string, useSudo bool, filePath string) (string, error) {
	return FetchContentsOfFileFromInstanceContextE(t, context.Background(), awsRegion, sshUserName, keyPair, instanceID, useSudo, filePath)
}

// FetchContentsOfFilesFromInstanceContextE looks up the public IP address of the EC2 Instance with the given ID, connects to
// the Instance via SSH using the given username and Key Pair, fetches the contents of the files at the given paths
// (using sudo if useSudo is true), and returns a map from file path to the contents of that file as a string.
// The ctx parameter supports cancellation and timeouts.
func FetchContentsOfFilesFromInstanceContextE(t testing.TestingT, ctx context.Context, awsRegion string, sshUserName string, keyPair *Ec2Keypair, instanceID string, useSudo bool, filePaths ...string) (map[string]string, error) {
	publicIP, err := GetPublicIPOfEc2InstanceContextE(t, ctx, instanceID, awsRegion)
	if err != nil {
		return nil, err
	}

	host := ssh.Host{
		SshUserName: sshUserName,
		SshKeyPair:  keyPair.KeyPair,
		Hostname:    publicIP,
	}

	return ssh.FetchContentsOfFilesContextE(t, ctx, &host, useSudo, filePaths...)
}

// FetchContentsOfFilesFromInstanceContext looks up the public IP address of the EC2 Instance with the given ID, connects to
// the Instance via SSH using the given username and Key Pair, fetches the contents of the files at the given paths
// (using sudo if useSudo is true), and returns a map from file path to the contents of that file as a string.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchContentsOfFilesFromInstanceContext(t testing.TestingT, ctx context.Context, awsRegion string, sshUserName string, keyPair *Ec2Keypair, instanceID string, useSudo bool, filePaths ...string) map[string]string {
	t.Helper()

	out, err := FetchContentsOfFilesFromInstanceContextE(t, ctx, awsRegion, sshUserName, keyPair, instanceID, useSudo, filePaths...)
	require.NoError(t, err)

	return out
}

// FetchContentsOfFilesFromInstance looks up the public IP address of the EC2 Instance with the given ID, connects to
// the Instance via SSH using the given username and Key Pair, fetches the contents of the files at the given paths
// (using sudo if useSudo is true), and returns a map from file path to the contents of that file as a string.
//
// Deprecated: Use [FetchContentsOfFilesFromInstanceContext] instead.
func FetchContentsOfFilesFromInstance(t testing.TestingT, awsRegion string, sshUserName string, keyPair *Ec2Keypair, instanceID string, useSudo bool, filePaths ...string) map[string]string {
	t.Helper()

	return FetchContentsOfFilesFromInstanceContext(t, context.Background(), awsRegion, sshUserName, keyPair, instanceID, useSudo, filePaths...)
}

// FetchContentsOfFilesFromInstanceE looks up the public IP address of the EC2 Instance with the given ID, connects to
// the Instance via SSH using the given username and Key Pair, fetches the contents of the files at the given paths
// (using sudo if useSudo is true), and returns a map from file path to the contents of that file as a string.
//
// Deprecated: Use [FetchContentsOfFilesFromInstanceContextE] instead.
func FetchContentsOfFilesFromInstanceE(t testing.TestingT, awsRegion string, sshUserName string, keyPair *Ec2Keypair, instanceID string, useSudo bool, filePaths ...string) (map[string]string, error) {
	return FetchContentsOfFilesFromInstanceContextE(t, context.Background(), awsRegion, sshUserName, keyPair, instanceID, useSudo, filePaths...)
}

// FetchContentsOfFileFromAsgContextE looks up the EC2 Instances in the given ASG, looks up the public IPs of those EC2
// Instances, connects to each Instance via SSH using the given username and Key Pair, fetches the contents of the file
// at the given path (using sudo if useSudo is true), and returns a map from Instance ID to the contents of that file
// as a string.
// The ctx parameter supports cancellation and timeouts.
func FetchContentsOfFileFromAsgContextE(t testing.TestingT, ctx context.Context, awsRegion string, sshUserName string, keyPair *Ec2Keypair, asgName string, useSudo bool, filePath string) (map[string]string, error) {
	instanceIDs, err := GetInstanceIdsForAsgContextE(t, ctx, asgName, awsRegion)
	if err != nil {
		return nil, err
	}

	instanceIDToContents := map[string]string{}

	for _, instanceID := range instanceIDs {
		contents, err := FetchContentsOfFileFromInstanceContextE(t, ctx, awsRegion, sshUserName, keyPair, instanceID, useSudo, filePath)
		if err != nil {
			return nil, err
		}

		instanceIDToContents[instanceID] = contents
	}

	return instanceIDToContents, nil
}

// FetchContentsOfFileFromAsgContext looks up the EC2 Instances in the given ASG, looks up the public IPs of those EC2
// Instances, connects to each Instance via SSH using the given username and Key Pair, fetches the contents of the file
// at the given path (using sudo if useSudo is true), and returns a map from Instance ID to the contents of that file
// as a string.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchContentsOfFileFromAsgContext(t testing.TestingT, ctx context.Context, awsRegion string, sshUserName string, keyPair *Ec2Keypair, asgName string, useSudo bool, filePath string) map[string]string {
	t.Helper()

	out, err := FetchContentsOfFileFromAsgContextE(t, ctx, awsRegion, sshUserName, keyPair, asgName, useSudo, filePath)
	require.NoError(t, err)

	return out
}

// FetchContentsOfFileFromAsg looks up the EC2 Instances in the given ASG, looks up the public IPs of those EC2
// Instances, connects to each Instance via SSH using the given username and Key Pair, fetches the contents of the file
// at the given path (using sudo if useSudo is true), and returns a map from Instance ID to the contents of that file
// as a string.
//
// Deprecated: Use [FetchContentsOfFileFromAsgContext] instead.
func FetchContentsOfFileFromAsg(t testing.TestingT, awsRegion string, sshUserName string, keyPair *Ec2Keypair, asgName string, useSudo bool, filePath string) map[string]string {
	t.Helper()

	return FetchContentsOfFileFromAsgContext(t, context.Background(), awsRegion, sshUserName, keyPair, asgName, useSudo, filePath)
}

// FetchContentsOfFileFromAsgE looks up the EC2 Instances in the given ASG, looks up the public IPs of those EC2
// Instances, connects to each Instance via SSH using the given username and Key Pair, fetches the contents of the file
// at the given path (using sudo if useSudo is true), and returns a map from Instance ID to the contents of that file
// as a string.
//
// Deprecated: Use [FetchContentsOfFileFromAsgContextE] instead.
func FetchContentsOfFileFromAsgE(t testing.TestingT, awsRegion string, sshUserName string, keyPair *Ec2Keypair, asgName string, useSudo bool, filePath string) (map[string]string, error) {
	return FetchContentsOfFileFromAsgContextE(t, context.Background(), awsRegion, sshUserName, keyPair, asgName, useSudo, filePath)
}

// FetchContentsOfFilesFromAsgContextE looks up the EC2 Instances in the given ASG, looks up the public IPs of those EC2
// Instances, connects to each Instance via SSH using the given username and Key Pair, fetches the contents of the files
// at the given paths (using sudo if useSudo is true), and returns a map from Instance ID to a map of file path to the
// contents of that file as a string.
// The ctx parameter supports cancellation and timeouts.
func FetchContentsOfFilesFromAsgContextE(t testing.TestingT, ctx context.Context, awsRegion string, sshUserName string, keyPair *Ec2Keypair, asgName string, useSudo bool, filePaths ...string) (map[string]map[string]string, error) {
	instanceIDs, err := GetInstanceIdsForAsgContextE(t, ctx, asgName, awsRegion)
	if err != nil {
		return nil, err
	}

	instanceIDToFilePathToContents := map[string]map[string]string{}

	for _, instanceID := range instanceIDs {
		contents, err := FetchContentsOfFilesFromInstanceContextE(t, ctx, awsRegion, sshUserName, keyPair, instanceID, useSudo, filePaths...)
		if err != nil {
			return nil, err
		}

		instanceIDToFilePathToContents[instanceID] = contents
	}

	return instanceIDToFilePathToContents, nil
}

// FetchContentsOfFilesFromAsgContext looks up the EC2 Instances in the given ASG, looks up the public IPs of those EC2
// Instances, connects to each Instance via SSH using the given username and Key Pair, fetches the contents of the files
// at the given paths (using sudo if useSudo is true), and returns a map from Instance ID to a map of file path to the
// contents of that file as a string.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchContentsOfFilesFromAsgContext(t testing.TestingT, ctx context.Context, awsRegion string, sshUserName string, keyPair *Ec2Keypair, asgName string, useSudo bool, filePaths ...string) map[string]map[string]string {
	t.Helper()

	out, err := FetchContentsOfFilesFromAsgContextE(t, ctx, awsRegion, sshUserName, keyPair, asgName, useSudo, filePaths...)
	require.NoError(t, err)

	return out
}

// FetchContentsOfFilesFromAsg looks up the EC2 Instances in the given ASG, looks up the public IPs of those EC2
// Instances, connects to each Instance via SSH using the given username and Key Pair, fetches the contents of the files
// at the given paths (using sudo if useSudo is true), and returns a map from Instance ID to a map of file path to the
// contents of that file as a string.
//
// Deprecated: Use [FetchContentsOfFilesFromAsgContext] instead.
func FetchContentsOfFilesFromAsg(t testing.TestingT, awsRegion string, sshUserName string, keyPair *Ec2Keypair, asgName string, useSudo bool, filePaths ...string) map[string]map[string]string {
	t.Helper()

	return FetchContentsOfFilesFromAsgContext(t, context.Background(), awsRegion, sshUserName, keyPair, asgName, useSudo, filePaths...)
}

// FetchContentsOfFilesFromAsgE looks up the EC2 Instances in the given ASG, looks up the public IPs of those EC2
// Instances, connects to each Instance via SSH using the given username and Key Pair, fetches the contents of the files
// at the given paths (using sudo if useSudo is true), and returns a map from Instance ID to a map of file path to the
// contents of that file as a string.
//
// Deprecated: Use [FetchContentsOfFilesFromAsgContextE] instead.
func FetchContentsOfFilesFromAsgE(t testing.TestingT, awsRegion string, sshUserName string, keyPair *Ec2Keypair, asgName string, useSudo bool, filePaths ...string) (map[string]map[string]string, error) {
	return FetchContentsOfFilesFromAsgContextE(t, context.Background(), awsRegion, sshUserName, keyPair, asgName, useSudo, filePaths...)
}

// FetchFilesFromInstanceContextE looks up the EC2 Instances in the given ASG, looks up the public IPs of those EC2
// Instances, connects to each Instance via SSH using the given username and Key Pair, downloads the files
// matching filenameFilters at the given remoteDirectory (using sudo if useSudo is true), and stores the files locally
// at localDirectory/<publicip>/<remoteFolderName>.
// The ctx parameter supports cancellation and timeouts.
func FetchFilesFromInstanceContextE(t testing.TestingT, ctx context.Context, awsRegion string, sshUserName string, keyPair *Ec2Keypair, instanceID string, useSudo bool, remoteDirectory string, localDirectory string, filenameFilters []string) error {
	publicIP, err := GetPublicIPOfEc2InstanceContextE(t, ctx, instanceID, awsRegion)
	if err != nil {
		return err
	}

	host := ssh.Host{
		Hostname:    publicIP,
		SshUserName: sshUserName,
		SshKeyPair:  keyPair.KeyPair,
	}

	finalLocalDestDir := filepath.Join(localDirectory, publicIP, filepath.Base(remoteDirectory))

	if !files.FileExists(finalLocalDestDir) {
		if err := os.MkdirAll(finalLocalDestDir, 0755); err != nil { //nolint:mnd // standard directory permissions
			return err
		}
	}

	//nolint:staticcheck,contextcheck // ScpDirFromE has no Context variant yet
	scpOptions := ssh.SCPDownloadOptions{
		RemoteHost:      host,
		RemoteDir:       remoteDirectory,
		LocalDir:        finalLocalDestDir,
		FileNameFilters: filenameFilters,
	}

	//nolint:staticcheck,contextcheck // ScpDirFromE has no Context variant yet
	return ssh.ScpDirFromE(t, scpOptions, useSudo)
}

// FetchFilesFromInstanceContext looks up the EC2 Instances in the given ASG, looks up the public IPs of those EC2
// Instances, connects to each Instance via SSH using the given username and Key Pair, downloads the files
// matching filenameFilters at the given remoteDirectory (using sudo if useSudo is true), and stores the files locally
// at localDirectory/<publicip>/<remoteFolderName>.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchFilesFromInstanceContext(t testing.TestingT, ctx context.Context, awsRegion string, sshUserName string, keyPair *Ec2Keypair, instanceID string, useSudo bool, remoteDirectory string, localDirectory string, filenameFilters []string) {
	t.Helper()

	err := FetchFilesFromInstanceContextE(t, ctx, awsRegion, sshUserName, keyPair, instanceID, useSudo, remoteDirectory, localDirectory, filenameFilters)
	require.NoError(t, err)
}

// FetchFilesFromInstance looks up the EC2 Instances in the given ASG, looks up the public IPs of those EC2
// Instances, connects to each Instance via SSH using the given username and Key Pair, downloads the files
// matching filenameFilters at the given remoteDirectory (using sudo if useSudo is true), and stores the files locally
// at localDirectory/<publicip>/<remoteFolderName>
//
// Deprecated: Use [FetchFilesFromInstanceContext] instead.
func FetchFilesFromInstance(t testing.TestingT, awsRegion string, sshUserName string, keyPair *Ec2Keypair, instanceID string, useSudo bool, remoteDirectory string, localDirectory string, filenameFilters []string) {
	t.Helper()

	FetchFilesFromInstanceContext(t, context.Background(), awsRegion, sshUserName, keyPair, instanceID, useSudo, remoteDirectory, localDirectory, filenameFilters)
}

// FetchFilesFromInstanceE looks up the EC2 Instances in the given ASG, looks up the public IPs of those EC2
// Instances, connects to each Instance via SSH using the given username and Key Pair, downloads the files
// matching filenameFilters at the given remoteDirectory (using sudo if useSudo is true), and stores the files locally
// at localDirectory/<publicip>/<remoteFolderName>
//
// Deprecated: Use [FetchFilesFromInstanceContextE] instead.
func FetchFilesFromInstanceE(t testing.TestingT, awsRegion string, sshUserName string, keyPair *Ec2Keypair, instanceID string, useSudo bool, remoteDirectory string, localDirectory string, filenameFilters []string) error {
	return FetchFilesFromInstanceContextE(t, context.Background(), awsRegion, sshUserName, keyPair, instanceID, useSudo, remoteDirectory, localDirectory, filenameFilters)
}

// FetchFilesFromAsgsPContextE looks up the EC2 Instances in all the ASGs given in the RemoteFileSpecification,
// looks up the public IPs of those EC2 Instances, connects to each Instance via SSH using the given
// username and Key Pair, downloads the files matching filenameFilters at the given
// remoteDirectory (using sudo if useSudo is true), and stores the files locally at
// localDirectory/<publicip>/<remoteFolderName>. This variant accepts a pointer to RemoteFileSpecification
// to avoid copying the large struct.
// The ctx parameter supports cancellation and timeouts.
func FetchFilesFromAsgsPContextE(t testing.TestingT, ctx context.Context, awsRegion string, spec *RemoteFileSpecification) error {
	errorsOccurred := new(multierror.Error)

	for _, curAsg := range spec.AsgNames {
		for curRemoteDir, fileFilters := range spec.RemotePathToFileFilter {
			instanceIDs, err := GetInstanceIdsForAsgContextE(t, ctx, curAsg, awsRegion)
			if err != nil {
				errorsOccurred = multierror.Append(errorsOccurred, err)
			} else {
				for _, instanceID := range instanceIDs {
					err = FetchFilesFromInstanceContextE(t, ctx, awsRegion, spec.SshUser, spec.KeyPair, instanceID, spec.UseSudo, curRemoteDir, spec.LocalDestinationDir, fileFilters)
					if err != nil {
						errorsOccurred = multierror.Append(errorsOccurred, err)
					}
				}
			}
		}
	}

	return errorsOccurred.ErrorOrNil()
}

// FetchFilesFromAsgsPContext looks up the EC2 Instances in all the ASGs given in the RemoteFileSpecification,
// looks up the public IPs of those EC2 Instances, connects to each Instance via SSH using the given
// username and Key Pair, downloads the files matching filenameFilters at the given
// remoteDirectory (using sudo if useSudo is true), and stores the files locally at
// localDirectory/<publicip>/<remoteFolderName>. This variant accepts a pointer to RemoteFileSpecification
// to avoid copying the large struct.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchFilesFromAsgsPContext(t testing.TestingT, ctx context.Context, awsRegion string, spec *RemoteFileSpecification) {
	t.Helper()

	err := FetchFilesFromAsgsPContextE(t, ctx, awsRegion, spec)
	require.NoError(t, err)
}

// FetchFilesFromAsgsP looks up the EC2 Instances in all the ASGs given in the RemoteFileSpecification,
// looks up the public IPs of those EC2 Instances, connects to each Instance via SSH using the given
// username and Key Pair, downloads the files matching filenameFilters at the given
// remoteDirectory (using sudo if useSudo is true), and stores the files locally at
// localDirectory/<publicip>/<remoteFolderName>. This variant accepts a pointer to RemoteFileSpecification
// to avoid copying the large struct.
//
// Deprecated: Use [FetchFilesFromAsgsPContext] instead.
func FetchFilesFromAsgsP(t testing.TestingT, awsRegion string, spec *RemoteFileSpecification) {
	t.Helper()

	FetchFilesFromAsgsPContext(t, context.Background(), awsRegion, spec)
}

// FetchFilesFromAsgsPE looks up the EC2 Instances in all the ASGs given in the RemoteFileSpecification,
// looks up the public IPs of those EC2 Instances, connects to each Instance via SSH using the given
// username and Key Pair, downloads the files matching filenameFilters at the given
// remoteDirectory (using sudo if useSudo is true), and stores the files locally at
// localDirectory/<publicip>/<remoteFolderName>. This variant accepts a pointer to RemoteFileSpecification
// to avoid copying the large struct.
//
// Deprecated: Use [FetchFilesFromAsgsPContextE] instead.
func FetchFilesFromAsgsPE(t testing.TestingT, awsRegion string, spec *RemoteFileSpecification) error {
	return FetchFilesFromAsgsPContextE(t, context.Background(), awsRegion, spec)
}

// FetchFilesFromAsgs looks up the EC2 Instances in all the ASGs given in the
// RemoteFileSpecification, downloads the matching files from each instance,
// and stores them locally as described in [FetchFilesFromAsgsPContext].
//
// Deprecated: Use [FetchFilesFromAsgsPContext] instead.
//
//nolint:staticcheck,revive,gocritic // preserving deprecated function name
func FetchFilesFromAsgs(t testing.TestingT, awsRegion string, spec RemoteFileSpecification) {
	t.Helper()

	FetchFilesFromAsgsPContext(t, context.Background(), awsRegion, &spec)
}

// FetchFilesFromAsgsE is the error-returning equivalent of [FetchFilesFromAsgs].
//
// Deprecated: Use [FetchFilesFromAsgsPContextE] instead.
//
//nolint:staticcheck,revive,gocritic // preserving deprecated function name
func FetchFilesFromAsgsE(t testing.TestingT, awsRegion string, spec RemoteFileSpecification) error {
	return FetchFilesFromAsgsPContextE(t, context.Background(), awsRegion, &spec)
}
