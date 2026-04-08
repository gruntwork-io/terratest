package aws

import (
	"os"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/hashicorp/go-multierror"
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

// FetchContentsOfFileFromInstance looks up the public IP address of the EC2 Instance with the given ID, connects to
// the Instance via SSH using the given username and Key Pair, fetches the contents of the file at the given path
// (using sudo if useSudo is true), and returns the contents of that file as a string.
func FetchContentsOfFileFromInstance(t testing.TestingT, awsRegion string, sshUserName string, keyPair *Ec2Keypair, instanceID string, useSudo bool, filePath string) string {
	out, err := FetchContentsOfFileFromInstanceE(t, awsRegion, sshUserName, keyPair, instanceID, useSudo, filePath)
	if err != nil {
		t.Fatal(err)
	}

	return out
}

// FetchContentsOfFileFromInstanceE looks up the public IP address of the EC2 Instance with the given ID, connects to
// the Instance via SSH using the given username and Key Pair, fetches the contents of the file at the given path
// (using sudo if useSudo is true), and returns the contents of that file as a string.
func FetchContentsOfFileFromInstanceE(t testing.TestingT, awsRegion string, sshUserName string, keyPair *Ec2Keypair, instanceID string, useSudo bool, filePath string) (string, error) {
	publicIP, err := GetPublicIPOfEc2InstanceE(t, instanceID, awsRegion)
	if err != nil {
		return "", err
	}

	host := ssh.Host{
		SshUserName: sshUserName,
		SshKeyPair:  keyPair.KeyPair,
		Hostname:    publicIP,
	}

	//nolint:staticcheck // using deprecated ssh function
	return ssh.FetchContentsOfFileE(t, host, useSudo, filePath)
}

// FetchContentsOfFilesFromInstance looks up the public IP address of the EC2 Instance with the given ID, connects to
// the Instance via SSH using the given username and Key Pair, fetches the contents of the files at the given paths
// (using sudo if useSudo is true), and returns a map from file path to the contents of that file as a string.
func FetchContentsOfFilesFromInstance(t testing.TestingT, awsRegion string, sshUserName string, keyPair *Ec2Keypair, instanceID string, useSudo bool, filePaths ...string) map[string]string {
	out, err := FetchContentsOfFilesFromInstanceE(t, awsRegion, sshUserName, keyPair, instanceID, useSudo, filePaths...)
	if err != nil {
		t.Fatal(err)
	}

	return out
}

// FetchContentsOfFilesFromInstanceE looks up the public IP address of the EC2 Instance with the given ID, connects to
// the Instance via SSH using the given username and Key Pair, fetches the contents of the files at the given paths
// (using sudo if useSudo is true), and returns a map from file path to the contents of that file as a string.
func FetchContentsOfFilesFromInstanceE(t testing.TestingT, awsRegion string, sshUserName string, keyPair *Ec2Keypair, instanceID string, useSudo bool, filePaths ...string) (map[string]string, error) {
	publicIP, err := GetPublicIPOfEc2InstanceE(t, instanceID, awsRegion)
	if err != nil {
		return nil, err
	}

	host := ssh.Host{
		SshUserName: sshUserName,
		SshKeyPair:  keyPair.KeyPair,
		Hostname:    publicIP,
	}

	//nolint:staticcheck // using deprecated ssh function
	return ssh.FetchContentsOfFilesE(t, host, useSudo, filePaths...)
}

// FetchContentsOfFileFromAsg looks up the EC2 Instances in the given ASG, looks up the public IPs of those EC2
// Instances, connects to each Instance via SSH using the given username and Key Pair, fetches the contents of the file
// at the given path (using sudo if useSudo is true), and returns a map from Instance ID to the contents of that file
// as a string.
func FetchContentsOfFileFromAsg(t testing.TestingT, awsRegion string, sshUserName string, keyPair *Ec2Keypair, asgName string, useSudo bool, filePath string) map[string]string {
	out, err := FetchContentsOfFileFromAsgE(t, awsRegion, sshUserName, keyPair, asgName, useSudo, filePath)
	if err != nil {
		t.Fatal(err)
	}

	return out
}

// FetchContentsOfFileFromAsgE looks up the EC2 Instances in the given ASG, looks up the public IPs of those EC2
// Instances, connects to each Instance via SSH using the given username and Key Pair, fetches the contents of the file
// at the given path (using sudo if useSudo is true), and returns a map from Instance ID to the contents of that file
// as a string.
func FetchContentsOfFileFromAsgE(t testing.TestingT, awsRegion string, sshUserName string, keyPair *Ec2Keypair, asgName string, useSudo bool, filePath string) (map[string]string, error) {
	instanceIDs, err := GetInstanceIdsForAsgE(t, asgName, awsRegion)
	if err != nil {
		return nil, err
	}

	instanceIDToContents := map[string]string{}

	for _, instanceID := range instanceIDs {
		contents, err := FetchContentsOfFileFromInstanceE(t, awsRegion, sshUserName, keyPair, instanceID, useSudo, filePath)
		if err != nil {
			return nil, err
		}

		instanceIDToContents[instanceID] = contents
	}

	return instanceIDToContents, err
}

// FetchContentsOfFilesFromAsg looks up the EC2 Instances in the given ASG, looks up the public IPs of those EC2
// Instances, connects to each Instance via SSH using the given username and Key Pair, fetches the contents of the files
// at the given paths (using sudo if useSudo is true), and returns a map from Instance ID to a map of file path to the
// contents of that file as a string.
func FetchContentsOfFilesFromAsg(t testing.TestingT, awsRegion string, sshUserName string, keyPair *Ec2Keypair, asgName string, useSudo bool, filePaths ...string) map[string]map[string]string {
	out, err := FetchContentsOfFilesFromAsgE(t, awsRegion, sshUserName, keyPair, asgName, useSudo, filePaths...)
	if err != nil {
		t.Fatal(err)
	}

	return out
}

// FetchContentsOfFilesFromAsgE looks up the EC2 Instances in the given ASG, looks up the public IPs of those EC2
// Instances, connects to each Instance via SSH using the given username and Key Pair, fetches the contents of the files
// at the given paths (using sudo if useSudo is true), and returns a map from Instance ID to a map of file path to the
// contents of that file as a string.
func FetchContentsOfFilesFromAsgE(t testing.TestingT, awsRegion string, sshUserName string, keyPair *Ec2Keypair, asgName string, useSudo bool, filePaths ...string) (map[string]map[string]string, error) {
	instanceIDs, err := GetInstanceIdsForAsgE(t, asgName, awsRegion)
	if err != nil {
		return nil, err
	}

	instanceIDToFilePathToContents := map[string]map[string]string{}

	for _, instanceID := range instanceIDs {
		contents, err := FetchContentsOfFilesFromInstanceE(t, awsRegion, sshUserName, keyPair, instanceID, useSudo, filePaths...)
		if err != nil {
			return nil, err
		}

		instanceIDToFilePathToContents[instanceID] = contents
	}

	return instanceIDToFilePathToContents, err
}

// FetchFilesFromInstance looks up the EC2 Instances in the given ASG, looks up the public IPs of those EC2
// Instances, connects to each Instance via SSH using the given username and Key Pair, downloads the files
// matching filenameFilters at the given remoteDirectory (using sudo if useSudo is true), and stores the files locally
// at localDirectory/<publicip>/<remoteFolderName>
func FetchFilesFromInstance(t testing.TestingT, awsRegion string, sshUserName string, keyPair *Ec2Keypair, instanceID string, useSudo bool, remoteDirectory string, localDirectory string, filenameFilters []string) {
	err := FetchFilesFromInstanceE(t, awsRegion, sshUserName, keyPair, instanceID, useSudo, remoteDirectory, localDirectory, filenameFilters)
	if err != nil {
		t.Fatal(err)
	}
}

// FetchFilesFromInstanceE looks up the EC2 Instances in the given ASG, looks up the public IPs of those EC2
// Instances, connects to each Instance via SSH using the given username and Key Pair, downloads the files
// matching filenameFilters at the given remoteDirectory (using sudo if useSudo is true), and stores the files locally
// at localDirectory/<publicip>/<remoteFolderName>
func FetchFilesFromInstanceE(t testing.TestingT, awsRegion string, sshUserName string, keyPair *Ec2Keypair, instanceID string, useSudo bool, remoteDirectory string, localDirectory string, filenameFilters []string) error {
	publicIP, err := GetPublicIPOfEc2InstanceE(t, instanceID, awsRegion)
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

	scpOptions := ssh.SCPDownloadOptions{
		RemoteHost:      host,
		RemoteDir:       remoteDirectory,
		LocalDir:        finalLocalDestDir,
		FileNameFilters: filenameFilters,
	}

	//nolint:staticcheck // using deprecated ssh function
	return ssh.ScpDirFromE(t, scpOptions, useSudo)
}

// FetchFilesFromAsgsP looks up the EC2 Instances in all the ASGs given in the RemoteFileSpecification,
// looks up the public IPs of those EC2 Instances, connects to each Instance via SSH using the given
// username and Key Pair, downloads the files matching filenameFilters at the given
// remoteDirectory (using sudo if useSudo is true), and stores the files locally at
// localDirectory/<publicip>/<remoteFolderName>. This variant accepts a pointer to RemoteFileSpecification
// to avoid copying the large struct.
func FetchFilesFromAsgsP(t testing.TestingT, awsRegion string, spec *RemoteFileSpecification) {
	err := FetchFilesFromAsgsPE(t, awsRegion, spec)
	if err != nil {
		t.Fatal(err)
	}
}

// FetchFilesFromAsgsPE looks up the EC2 Instances in all the ASGs given in the RemoteFileSpecification,
// looks up the public IPs of those EC2 Instances, connects to each Instance via SSH using the given
// username and Key Pair, downloads the files matching filenameFilters at the given
// remoteDirectory (using sudo if useSudo is true), and stores the files locally at
// localDirectory/<publicip>/<remoteFolderName>. This variant accepts a pointer to RemoteFileSpecification
// to avoid copying the large struct.
func FetchFilesFromAsgsPE(t testing.TestingT, awsRegion string, spec *RemoteFileSpecification) error {
	errorsOccurred := new(multierror.Error)

	for _, curAsg := range spec.AsgNames {
		for curRemoteDir, fileFilters := range spec.RemotePathToFileFilter {
			instanceIDs, err := GetInstanceIdsForAsgE(t, curAsg, awsRegion)
			if err != nil {
				errorsOccurred = multierror.Append(errorsOccurred, err)
			} else {
				for _, instanceID := range instanceIDs {
					err = FetchFilesFromInstanceE(t, awsRegion, spec.SshUser, spec.KeyPair, instanceID, spec.UseSudo, curRemoteDir, spec.LocalDestinationDir, fileFilters)
					if err != nil {
						errorsOccurred = multierror.Append(errorsOccurred, err)
					}
				}
			}
		}
	}

	return errorsOccurred.ErrorOrNil()
}

// Deprecated: Use FetchFilesFromAsgsP instead, which accepts a pointer to RemoteFileSpecification.
//
//nolint:staticcheck,revive,gocritic // preserving deprecated function name
func FetchFilesFromAsgs(t testing.TestingT, awsRegion string, spec RemoteFileSpecification) {
	FetchFilesFromAsgsP(t, awsRegion, &spec)
}

// Deprecated: Use FetchFilesFromAsgsPE instead, which accepts a pointer to RemoteFileSpecification.
//
//nolint:staticcheck,revive,gocritic // preserving deprecated function name
func FetchFilesFromAsgsE(t testing.TestingT, awsRegion string, spec RemoteFileSpecification) error {
	return FetchFilesFromAsgsPE(t, awsRegion, &spec)
}
