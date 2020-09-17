package aws

import (
	"os"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/customerrors"
	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// RemoteFileSpecification describes which files you want to copy from your instances
type RemoteFileSpecification struct {
	AsgNames               []string            //ASGs where our instances will be
	RemotePathToFileFilter map[string][]string //A map of the files to fetch, where the keys are directories on the remote host and the values are filters for what files to fetch from the directory. The filters support bash-style wildcards.
	UseSudo                bool
	SshUser                string
	KeyPair                *Ec2Keypair
	LocalDestinationDir    string //base path where to store downloaded artifacts locally. The final path of each resource will include the ip of the host and the name of the immediate parent folder.
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
	publicIp, err := GetPublicIpOfEc2InstanceE(t, instanceID, awsRegion)
	if err != nil {
		return "", err
	}

	host := ssh.Host{
		SshUserName: sshUserName,
		SshKeyPair:  keyPair.KeyPair,
		Hostname:    publicIp,
	}

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
	publicIp, err := GetPublicIpOfEc2InstanceE(t, instanceID, awsRegion)
	if err != nil {
		return nil, err
	}

	host := ssh.Host{
		SshUserName: sshUserName,
		SshKeyPair:  keyPair.KeyPair,
		Hostname:    publicIp,
	}

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

	instanceIdToContents := map[string]string{}

	for _, instanceID := range instanceIDs {
		contents, err := FetchContentsOfFileFromInstanceE(t, awsRegion, sshUserName, keyPair, instanceID, useSudo, filePath)
		if err != nil {
			return nil, err
		}
		instanceIdToContents[instanceID] = contents
	}

	return instanceIdToContents, err
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

	instanceIdToFilePathToContents := map[string]map[string]string{}

	for _, instanceID := range instanceIDs {
		contents, err := FetchContentsOfFilesFromInstanceE(t, awsRegion, sshUserName, keyPair, instanceID, useSudo, filePaths...)
		if err != nil {
			return nil, err
		}
		instanceIdToFilePathToContents[instanceID] = contents
	}

	return instanceIdToFilePathToContents, err
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
	publicIp, err := GetPublicIpOfEc2InstanceE(t, instanceID, awsRegion)

	if err != nil {
		return err
	}

	host := ssh.Host{
		Hostname:    publicIp,
		SshUserName: sshUserName,
		SshKeyPair:  keyPair.KeyPair,
	}

	finalLocalDestDir := filepath.Join(localDirectory, publicIp, filepath.Base(remoteDirectory))

	if !files.FileExists(finalLocalDestDir) {
		os.MkdirAll(finalLocalDestDir, 0755)
	}

	scpOptions := ssh.ScpDownloadOptions{
		RemoteHost:      host,
		RemoteDir:       remoteDirectory,
		LocalDir:        finalLocalDestDir,
		FileNameFilters: filenameFilters,
	}

	return ssh.ScpDirFromE(t, scpOptions, useSudo)
}

// FetchFilesFromAsgs looks up the EC2 Instances in all the ASGs given in the RemoteFileSpecification,
// looks up the public IPs of those EC2 Instances, connects to each Instance via SSH using the given
// username and Key Pair, downloads the files matching filenameFilters at the given
// remoteDirectory (using sudo if useSudo is true), and stores the files locally at
// localDirectory/<publicip>/<remoteFolderName>
func FetchFilesFromAsgs(t testing.TestingT, awsRegion string, spec RemoteFileSpecification) {
	err := FetchFilesFromAsgsE(t, awsRegion, spec)

	if err != nil {
		t.Fatal(err)
	}
}

// FetchFilesFromAsgsE looks up the EC2 Instances in all the ASGs given in the RemoteFileSpecification,
// looks up the public IPs of those EC2 Instances, connects to each Instance via SSH using the given
// username and Key Pair, downloads the files matching filenameFilters at the given
// remoteDirectory (using sudo if useSudo is true), and stores the files locally at
// localDirectory/<publicip>/<remoteFolderName>
func FetchFilesFromAsgsE(t testing.TestingT, awsRegion string, spec RemoteFileSpecification) error {
	errorsOccurred := []error{}

	for _, curAsg := range spec.AsgNames {
		for curRemoteDir, fileFilters := range spec.RemotePathToFileFilter {

			instanceIDs, err := GetInstanceIdsForAsgE(t, curAsg, awsRegion)
			if err != nil {
				errorsOccurred = append(errorsOccurred, err)
			} else {
				for _, instanceID := range instanceIDs {
					err = FetchFilesFromInstanceE(t, awsRegion, spec.SshUser, spec.KeyPair, instanceID, spec.UseSudo, curRemoteDir, spec.LocalDestinationDir, fileFilters)

					if err != nil {
						errorsOccurred = append(errorsOccurred, err)
					}
				}
			}
		}
	}
	return customerrors.NewMultiError(errorsOccurred...)
}
