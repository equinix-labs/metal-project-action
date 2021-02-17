package main

import (
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
		log.Fatal(err)
	}

	p, err := a.CreateProject()
	if err != nil {
		log.Fatal(err)
	}

	// TODO(displague) any way to create a project token through the API?
	// If so, make sure to ::add-mask:: before adding it to the output or env

	envFile, err := os.OpenFile(os.Getenv("GITHUB_ENV"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer envFile.Close()

	for _, v := range map[string]string{

		"projectToken":         p.APIToken,
		"projectSSHPrivateKey": p.SSHPrivateKey,
	} {
		fmt.Printf("::add-mask::%s\n", url.QueryEscape(v))
	}

	for k, v := range map[string]string{
		"projectID":            p.Project.ID,
		"projectName":          p.Project.Name,
		"projectToken":         p.APIToken,
		"projectSSHPrivateKey": p.SSHPrivateKey,
	} {
		fmt.Printf("::set-output name=%s::%s\n", k, url.QueryEscape(v))
	}

	for k, v := range map[string]string{
		"METAL_PROJECT_ID":      p.Project.ID,
		"METAL_PROJECT_NAME":    p.Project.Name,
		"METAL_PROJECT_TOKEN":   p.APIToken,
		"METAL_SSH_PRIVATE_KEY": p.SSHPrivateKey,
	} {
		fmt.Fprintf(envFile, "%s<<EOS\n%s\nEOS\n", k, v)
	}
}
