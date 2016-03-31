package main

import (
    "fmt"
    getopt "github.com/kesselborn/go-getopt"
    "os"
    "github.com/kakwa/cinp/proto/v1"
)

func main() {
    optionDefinition := getopt.Options{
        "description",
        getopt.Definitions{
            {"debug|d",  "debug mode", getopt.Optional | getopt.Flag, false},
            {"config|c", "config file", getopt.IsConfigFile | getopt.ExampleIsDefault, "/etc/cinp/client.ini"},
            {"nic|i",    "network interface name", getopt.Required, "toto"},
        },
    }

    options, _, _, e := optionDefinition.ParseCommandLine()

    help, wantsHelp := options["help"]

    if e != nil || wantsHelp {
        exit_code := 0

        switch {
        case wantsHelp && help.String == "usage":
            fmt.Print(optionDefinition.Usage())
        case wantsHelp && help.String == "help":
            fmt.Print(optionDefinition.Help())
        default:
            fmt.Println("**** Error: ", e.Error(), "\n", optionDefinition.Help())
            exit_code = e.ErrorCode
        }
        os.Exit(exit_code)
    }
    req, xid, err := v1.NewRequest(v1.Clear)
    fmt.Println("packet: ", req, "\nxid: ", xid)
    ans, xid, err := v1.NewAnswer(req, v1.Payload("eth0"))
    fmt.Println("packet: ", ans, "\nxid: ", xid, "\npayload :", v1.Payload("eth0"))
    fmt.Printf("%s\n", ans.Payload())

    if err != nil {
        os.Exit(1)
    }
}
