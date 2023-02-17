package main

import (
	"fmt"

	"github.com/gruntwork-io/go-commons/entrypoint"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/urfave/cli"
)

const CustomUsageText = `Usage: pick-instance-type [OPTIONS] <REGION> <INSTANCE_TYPE> <INSTANCE_TYPE...> 

This tool takes in an AWS region and a list of EC2 instance types and returns the first instance type in the list that is available in all Availability Zones (AZs) in the given region, or exits with an error if no instance type is available in all AZs. This is useful because certain instance types, such as t2.micro, are not available in some of the newer AZs, while t3.micro is not available in some of the older AZs. If you have code that needs to run on a "small" instance across all AZs in many different regions, you can use this CLI tool to automatically figure out which instance type you should use.

Arguments:
   
  REGION           The AWS region in which to look up instance availability. E.g.: us-east-1. 
  INSTANCE_TYPE    One more more EC2 instance types. E.g.: t2.micro.


Options:

  --help            Show this help text and exit.

Example:

  pick-instance-type ap-northeast-2 t2.micro t3.micro 
`

func run(cliContext *cli.Context) error {
	region := cliContext.Args().First()
	if region == "" {
		return fmt.Errorf("You must specify an AWS region as the first argument")
	}

	instanceTypes := cliContext.Args().Tail()
	if len(instanceTypes) == 0 {
		return fmt.Errorf("You must specify at least one instance type")
	}

	// Create mock testing.T implementation so we can re-use Terratest methods
	t := MockTestingT{MockName: "pick-instance-type"}

	recommendedInstanceType, err := aws.GetRecommendedInstanceTypeE(t, region, instanceTypes)
	if err != nil {
		return err
	}

	// Print the recommended instance type to stdout
	fmt.Print(recommendedInstanceType)

	return nil
}

func main() {
	app := entrypoint.NewApp()
	cli.AppHelpTemplate = CustomUsageText
	entrypoint.HelpTextLineWidth = 120

	app.Name = "pick-instance-type"
	app.Author = "Gruntwork <www.gruntwork.io>"
	app.Description = `This tool takes in a list of EC2 instance types (e.g., "t2.micro", "t3.micro") and returns the first instance type in the list that is available in all Availability Zones (AZs) in the given AWS region, or exits with an error if no instance type is available in all AZs.`
	app.Action = run

	entrypoint.RunApp(app)
}

// MockTestingT is a mock implementation of testing.TestingT. All the functions are essentially no-ops. This allows us
// to use Terratest methods outside of a testing context (e.g., in a CLI tool).
type MockTestingT struct {
	MockName string
}

func (t MockTestingT) Fail()                                     {}
func (t MockTestingT) FailNow()                                  {}
func (t MockTestingT) Fatal(args ...interface{})                 {}
func (t MockTestingT) Fatalf(format string, args ...interface{}) {}
func (t MockTestingT) Error(args ...interface{})                 {}
func (t MockTestingT) Errorf(format string, args ...interface{}) {}
func (t MockTestingT) Name() string {
	return t.MockName
}
