package main

import (
    "fmt"
    getopt "github.com/kesselborn/go-getopt"
    "os"
)

func main() {
    optionDefinition := getopt.Options{
        "description",
        getopt.Definitions{
            {"debug|d",     "debug mode", getopt.Optional | getopt.Flag, false},
            {"config|c",    "config file", getopt.IsConfigFile | getopt.ExampleIsDefault, "/etc/cinp/server.ini"},
            {"nic|i",       "network interface name", getopt.Required, "toto"},
            {"forground|f", "run in forground", getopt.Optional | getopt.Flag, false},
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
}
