package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ethereumproject/janus/gcp"
	"github.com/ethereumproject/janus/gitvv"
)

func main() {

	// Subcommands
	deployCommand := flag.NewFlagSet("deploy", flag.ExitOnError)
	versionCommand := flag.NewFlagSet("version", flag.ExitOnError)

	// Deploy flags
	var key, files, to string
	// Version flags
	var format string

	// Set up flags.
	//
	// Deploy
	deployCommand.StringVar(&to, "to", "", `directory path to deploy files to

the first directory in the given path is GCP <bucket>
files will be uploaded INTO this path

eg.
-to=builds.etcdevteam.com/go-ethereum/releases/v3.5.x \
-files=./dist/*.zip [eg. ./dist/geth-linux.zip, ./dist/geth-osx.zip]

--> builds.etcdevteam.com/go-ethereum/releases/v3.5.x/geth-linux.zip
--> builds.etcdevteam.com/go-ethereum/releases/v3.5.x/geth-osx.zip
`)
	deployCommand.StringVar(&files, "files", "", "file(s) to upload, allows globbing")
	deployCommand.StringVar(&key, "key", "", "service account json key file, may be encrypted OR decrypted")
	// Version
	versionCommand.StringVar(&format, "format", "", `format of git version:

%M - major version
%m - minor version
%P - patch version
%C - commit count since last tag
%S - HEAD sha1

Default: v%M.%m.%P+%C-%S -> v3.5.0+66-bbb06b1
`)

	flag.Usage = func() {
		log.Println("Usage for Janus:")
		log.Println("  $ janus deploy -to builds.etcdevteam.com/go-ethereum/version -file geth.zip -key .gcloud.json")
		log.Println("  $ janus version -format 'v%M.%m.%P+%C-%S'")
		flag.PrintDefaults()
	}

	// Ensure subcommand is used.
	if len(os.Args) < 2 {
		fmt.Println("'deploy' or 'version' subcommand is required")
		os.Exit(1)
	}

	// Parse subcommands.
	switch os.Args[1] {
	case "deploy":
		deployCommand.Parse(os.Args[2:])
	case "version":
		versionCommand.Parse(os.Args[2:])
	default:
		flag.Usage()
		os.Exit(1)
	}

	// Handle which command is used.
	//
	// Deploy
	if deployCommand.Parsed() {
		// Ensure required flags are set.
		if to == "" {
			log.Println("--to requires an argument")
			flag.Usage()
			os.Exit(1)
		}
		if files == "" {
			log.Println("--files requires an argument")
			flag.Usage()
			os.Exit(1)
		}
		if key == "" {
			log.Println("--key requires an argument")
			flag.Usage()
			os.Exit(1)
		}

		// Handle deploy.
		// -- Will check for existing file(s) to upload, will return error if not exists.
		if e := gcp.SendToGCP(to, files, key); e != nil {
			log.Println("Failed to deploy:")
			log.Fatalln(e)
		}
	} else
	// Version
	if versionCommand.Parsed() {
		v := gitvv.GetVersion(format)
		fmt.Print(v)
		os.Exit(0)
	} else
	// No command
	{
		// Must use a subcommand.
		flag.Usage()
		os.Exit(1)
	}
}
