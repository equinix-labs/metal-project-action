package main

import (
	b64 "encoding/base64"
	"fmt"
	"log"
	"net/url"
	"os"

	action "github.com/displague/metal-project-action/internal"
)

func main() {
	projectName := os.Getenv("INPUTS_PROJECTNAME")
	if projectName == "" {
		projectName = action.GenProjectName(os.Getenv("GITHUB_SHA"))
	}

	apiToken := os.Getenv("INPUTS_USERTOKEN")
	if apiToken == "" {
		apiToken = os.Getenv("METAL_AUTH_TOKEN")
		if apiToken == "" {
			log.Fatal("Either `with.userToken` or `env.METAL_AUTH_TOKEN` must be supplied")
		}

	}
	a, err := action.NewAction(apiToken, os.Getenv("INPUTS_ORGANIZATIONID"), projectName)
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

	for _, v := range map[string]string{
		"projectToken":               p.APIToken,
		"projectSSHPrivateKeyBase64": sshPrivateBase64,
	} {
		fmt.Printf("::add-mask::%s\n", url.QueryEscape(v))
	}

	for k, v := range map[string]string{
		"projectID":                  p.Project.ID,
		"projectName":                p.Project.Name,
		"projectToken":               p.APIToken,
		"projectSSHPrivateKeyBase64": sshPrivateBase64,
		"projectSSHPublicKey":        sshPublicKey,
		"organizationID":             p.Project.Organization.ID,
	} {
		fmt.Printf("::set-output name=%s::%s\n", k, url.QueryEscape(v))
	}

	for k, v := range map[string]string{
		"METAL_PROJECT_ID":             p.Project.ID,
		"METAL_PROJECT_NAME":           p.Project.Name,
		"METAL_PROJECT_TOKEN":          p.APIToken,
		"METAL_SSH_PRIVATE_KEY_BASE64": sshPrivateBase64,
		"METAL_SSH_PUBLIC_KEY":         sshPublicKey,
		"METAL_ORGANIZATION_ID":        p.Project.Organization.ID,
	} {
		fmt.Fprintf(envFile, "%s<<EOS\n%s\nEOS\n", k, v)
	}
}
