package main
import (
    "os"
    "fmt"
    "flag"
    "strings"
    "github.com/aws/aws-sdk-go/aws/awserr"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sts"
)
type Credentials interface {
    getCredentials() map
}

type roleAssumptions struct {
    session, externalId, roleArn string
    duration int
    stdout bool
}

type orgAssumptions struct {
    session, externalId, accountNumber string
    duration int
    stdout bool
}

func (r roleAssumptions) getCredentials() {}
func (r orgAssumptions) getCredentials() {}

func parseInput() assumptionArguments {
    sessionStr := []string{"aws-cred-gen-go", os.Getenv("USER")}
    assumeCmd := flag.NewFlagSet("assume", flag.ExitOnError)
    stdoutPtr := assumeCmd.Bool("stdout", false, "Print credentials to stdout")
    sessionNamePtr := assumeCmd.String("session", strings.Join(sessionStr, "-"), "Session to use when assuming into AWS")
    externalIdPtr := assumeCmd.String("external-id", "", "External ID to use when assuming into target role")
    roleArnPtr := assumeCmd.String("role-arn", "", "Target role to assume into")
    durationPtr := assumeCmd.Int("duration", 3600, "Duration of credentials. (Default=3600)")

    assumeRoleCmd := flag.NewFlagSet("assume-org", flag.ExitOnError)
    accountNumberPtr := assumeRoleCmd.String("account-number", "", "Session to use when assuming into AWS")
    roleStdOutPtr := assumeRoleCmd.Bool("stdout", false, "Print credentials to stdout")
    roleSessionNamePtr := assumeRoleCmd.String("session", strings.Join(sessionStr, "-"), "Session to use when assuming into AWS")
    roleExternalIdPtr := assumeRoleCmd.String("external-id", "", "External ID to use when assuming into target role")
    roleDurationPtr := assumeRoleCmd.Int("duration", 3600, "Duration of credentials. (Default=3600)")

    if len(os.Args) < 2 {
        fmt.Println("Expected 'assume' or 'assume-org' subcommands")
        flag.Usage()
        os.Exit(1)
    }

    profilePtr := flag.String("profile", "default", "AWS Profile to use as entrypoint")
    targetProfilePtr := flag.String("target-profile", "", "AWS Profile to save new temporary credentials in")


    switch os.Args[1] {

    case "assume":
        assumeCmd.Parse(os.Args[2:])
        fmt.Println("assume")
        fmt.Println("stdout: ", *stdoutPtr)
        fmt.Println("sessionName: ", *sessionNamePtr)
        fmt.Println("externalId: ", *externalIdPtr)
        fmt.Println("roleArn: ", *roleArnPtr)
        fmt.Println("duration: ", *durationPtr)
        fmt.Println(assumeCmd.Args())
    case "assume-org":
        assumeRoleCmd.Parse(os.Args[2:])
        fmt.Println("assume-org")
        fmt.Println("Account Number: ", *accountNumberPtr)
        fmt.Println("Stdout: ", *roleStdOutPtr)
        fmt.Println("Session Name: ", *roleSessionNamePtr)
        fmt.Println("Duration: ", *roleDurationPtr)
        fmt.Println("externalId: ", *roleExternalIdPtr)
        fmt.Println(assumeRoleCmd.Args())
    default:
        fmt.Println("Invalid Subcommand!")
        flag.Usage()
        os.Exit(3)
    }

    flag.Parse()

    fmt.Println("Profile: ", *profilePtr)
    fmt.Println("Target Profile: ", *targetProfilePtr)

    if (*profilePtr == ""){
        fmt.Println("profile: Missing Value!")
        os.Exit(1)
    }
    if (*targetProfilePtr == "") {
        fmt.Println("target-profile: Missing Value!")
        os.Exit(1)
    }
}

func main() {
    parseInput()

    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))
    svc := sts.New(sess)
    input := &sts.GetCallerIdentityInput{}
    result, err := svc.GetCallerIdentity(input)
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            default:
                fmt.Println(aerr.Error())
            }
        } else {
            // Print the error, cast err to awserr.Error to get the Code and
            // Message from an error.
            fmt.Println(err.Error())
        }
        return
    }

    fmt.Println(result)
}
