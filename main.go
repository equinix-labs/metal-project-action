package main

import (
	"fmt"
	"io/ioutil"
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
		log.Fatal(err)
	}

	// Use GH Runner temp directory. TempFile uses TempDir if empty.
	tmp := os.Getenv("RUNNER_TEMP")

	f, err := ioutil.TempFile(tmp, "id_rsa_")
	if err != nil {
		log.Fatal(err)
	}
	sshPrivateFile := f.Name()
	sshPublicFile := sshPrivateFile + ".pub"

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}

	p, err := a.CreateProject(sshPrivateFile)
	if err != nil {
		log.Fatal(err)
	}

	envFile, err := os.OpenFile(os.Getenv("GITHUB_ENV"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer envFile.Close()

	for _, v := range map[string]string{
		"projectToken": p.APIToken,
	} {
		fmt.Printf("::add-mask::%s\n", url.QueryEscape(v))
	}

	for k, v := range map[string]string{
		"projectID":                p.Project.ID,
		"projectName":              p.Project.Name,
		"projectToken":             p.APIToken,
		"projectSSHPrivateKeyFile": sshPrivateFile,
		"projectSSHPublicKeyFile":  sshPublicFile,
		"organizationID":           p.Project.Organization.ID,
	} {
		fmt.Printf("::set-output name=%s::%s\n", k, url.QueryEscape(v))
	}

	for k, v := range map[string]string{
		"METAL_PROJECT_ID":           p.Project.ID,
		"METAL_PROJECT_NAME":         p.Project.Name,
		"METAL_PROJECT_TOKEN":        p.APIToken,
		"METAL_SSH_PRIVATE_KEY_FILE": sshPrivateFile,
		"METAL_SSH_PUBLIC_KEY_FILE":  sshPublicFile,
		"METAL_ORGANIZATION_ID":      p.Project.Organization.ID,
	} {
		fmt.Fprintf(envFile, "%s<<EOS\n%s\nEOS\n", k, v)
	}
}
