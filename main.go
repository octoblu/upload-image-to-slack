package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/coreos/go-semver/semver"
	"github.com/fatih/color"
	"github.com/octoblu/upload-image-to-slack/slack"
	De "github.com/tj/go-debug"
)

var debug = De.Debug("upload-image-to-slack:main")

func main() {
	app := cli.NewApp()
	app.Name = "upload-image-to-slack"
	app.Version = version()
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "channel, c",
			EnvVar: "UITS_SLACK_CHANNEL",
			Usage:  "Slack Channel to post into",
		},
		cli.StringFlag{
			Name:   "filename, f",
			EnvVar: "UITS_FILENAME",
			Usage:  "File name to upload",
		},
		cli.StringFlag{
			Name:   "token, t",
			EnvVar: "UITS_SLACK_TOKEN",
			Usage:  "Slack Token",
		},
	}
	app.Run(os.Args)
}

func run(context *cli.Context) error {
	channel, filename, token := getOpts(context)
	content := bufio.NewReader(os.Stdin)

	client := slack.New(channel, token)
	err := client.Upload(filename, content)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return err
}

func getOpts(context *cli.Context) (string, string, string) {
	channel := context.String("channel")
	filename := context.String("filename")
	token := context.String("token")

	if channel == "" || filename == "" || token == "" {
		cli.ShowAppHelp(context)

		if channel == "" {
			color.Red("  Missing required flag --channel or UITS_SLACK_CHANNEL")
		}
		if filename == "" {
			color.Red("  Missing required flag --filename or UITS_FILENAME")
		}
		if token == "" {
			color.Red("  Missing required flag --token or UITS_SLACK_TOKEN")
		}
		os.Exit(1)
	}

	return channel, filename, token
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}
