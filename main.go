package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/packethost/packngo"
)

func main() {
	// PACKET_AUTH_TOKEN should be already set
	// the client will use it by default (?)
	client, err := packngo.NewClient()
	if err != nil {
		panic(err)
	}

	projectName := os.Getenv("INPUTS_PROJECTNAME")
	if projectName == "" {
		// TODO(displague) use a random string
		sha := os.Getenv("GITHUB_SHA")
		if sha == "" {
			sha = RandomString(16)
		}
		projectName = "GHACTION-" + sha

		// TODO(displague) no way to set a description? is "customdata" a description?
		// projectDescription := os.Getenv("GITHUB_SERVER_URL") + "/" + os.Getenv("GITHUB_REPOSITORY") + " " + os.Getenv("GITHUB_SHA")
	}

	createOpts := &packngo.ProjectCreateRequest{
		Name: projectName,
	}

	project, _, err := client.Projects.Create(createOpts)
	if err != nil {
		panic(err)
	}

	// TODO(displague) any way to create a project token through the API?

	for k, v := range map[string]string{
		"projectID":   project.ID,
		"projectName": project.Name,
	} {
		fmt.Printf("::set-output name=%s::%s\n", k, v)
	}
}

func RandomString(size int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, size)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}
