package main

import (
	"flag"
	"fmt"
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
	var gpg bool
	// Version flags
	var dir, format string

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
	deployCommand.BoolVar(&gpg, "gpg", false, "use GPG 2 instead of openssl for decryption")
	// Version
	versionCommand.StringVar(&dir, "dir", "", `path to base directory`)
	versionCommand.StringVar(&format, "format", "", `format of git version:

%M - major version
%m - minor version
%P - patch version
%C - commit count since last tag
%S[|NUMBER] - HEAD sha1, where NUMBER is optional desired length of hash (default: 7)
%B - hybrid patch number (B = semver_minor_version*100 + commit_count)

Default: v%M.%m.%P+%C-%S -> v3.5.0+66-bbb06b1
`)

	flag.Usage = func() {
		fmt.Println("Usage for Janus:")
		fmt.Println("  $ janus deploy -to builds.etcdevteam.com/go-ethereum/version -file geth.zip -key .gcloud.json")
		fmt.Println("  $ janus version -format 'v%M.%m.%P+%C-%S'")
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
			fmt.Println("--to requires an argument")
			flag.Usage()
			os.Exit(1)
		}
		if files == "" {
			fmt.Println("--files requires an argument")
			flag.Usage()
			os.Exit(1)
		}
		if key == "" {
			fmt.Println("--key requires an argument")
			flag.Usage()
			os.Exit(1)
		}

		// Handle deploy.
		// -- Will check for existing file(s) to upload, will return error if not exists.
		if e := gcp.SendToGCP(to, files, key, gpg); e != nil {
			fmt.Println("Failed to deploy:")
			fmt.Println(e)
			os.Exit(1)
		}
	} else
	// Version
	if versionCommand.Parsed() {
		v := gitvv.GetVersion(format, dir)
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
