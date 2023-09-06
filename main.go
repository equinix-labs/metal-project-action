package main

import (
	b64 "encoding/base64"
	"fmt"
	"log"
	"net/url"
	"os"

	action "github.com/equinix-labs/metal-project-action/internal"
)

func main() {
	projectName := os.Getenv("INPUT_PROJECTNAME")
	if projectName == "" {
		projectName = action.GenProjectName(os.Getenv("GITHUB_SHA"))
	}

	apiToken := os.Getenv("INPUT_USERTOKEN")
	if apiToken == "" {
		log.Fatal("You must provide an auth token in `with.userToken` must be supplied")
	}

	enableBGP := os.Getenv("INPUT_ENABLEBGP") == "true"

	a, err := action.NewAction(apiToken, os.Getenv("INPUT_ORGANIZATIONID"), projectName, enableBGP)
	if err != nil {
		log.Fatal("Could not create client action", err)
	}

	if err != nil {
		log.Fatal("Could not create temp file", err)
	}

	if err != nil {
		log.Fatal("Could not close temp file", err)
	}

	p, err := a.CreateProject()
	if err != nil {
		log.Fatal("Could not create project", err)
	}

	sshPrivateBase64 := b64.StdEncoding.EncodeToString([]byte(p.SSHPrivateKey))
	sshPublicKey := b64.StdEncoding.EncodeToString([]byte(p.SSHPublicKey))

	envFile, err := os.OpenFile(os.Getenv("GITHUB_ENV"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Could not open env file", err)
	}
	defer envFile.Close()

	outputFile, err := os.OpenFile(os.Getenv("GITHUB_OUTPUT"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		//nolint:gocritic
		log.Fatal("Could not open output file", err)
	}
	defer outputFile.Close()

	for _, v := range map[string]string{
		"projectToken":               p.APIToken,
		"projectSSHPrivateKeyBase64": sshPrivateBase64,
	} {
		fmt.Printf("::add-mask::%s\n", url.QueryEscape(v))
	}

	for k, v := range map[string]string{
		"projectID":                  p.Project.GetId(),
		"projectName":                p.Project.GetName(),
		"projectToken":               p.APIToken,
		"projectSSHPrivateKeyBase64": sshPrivateBase64,
		"projectSSHPublicKey":        sshPublicKey,
		"organizationID":             p.Project.Organization.GetId(),
	} {
		fmt.Fprintf(outputFile, "%s=%s\n", k, url.QueryEscape(v))
	}
}
