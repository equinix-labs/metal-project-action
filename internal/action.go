// Portions taken from https://gist.github.com/devinodaniel/8f9b8a4f31573f428f29ec0e884e6673
package action

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	mrand "math/rand"
	"time"

	metal "github.com/equinix-labs/metal-go/metal/v1"
	"golang.org/x/crypto/ssh"
)

var version = "dev"

const (
	bitSize = 4096
	uaFmt   = "gh-action-metal-project/%s %s"
)

type action struct {
	client         *metal.APIClient
	label          string
	organizationID string
}

type Project struct {
	Project *metal.Project

	SSHPrivateKey string
	SSHPublicKey  string
	APIToken      string
}

// NewAction returns an action with a Packngo client
func NewAction(apiToken, organizationID, label string) (*action, error) {
	config := metal.NewConfiguration()
	config.AddDefaultHeader("X-Auth-Token", apiToken)
	config.UserAgent = fmt.Sprintf(uaFmt, version, config.UserAgent)
	client := metal.NewAPIClient(config)

	return &action{
		organizationID: organizationID,
		label:          label,
		client:         client,
	}, nil
}

// CreateProject
//
// Create an Equinix Metal project with API keys and project SSH Keys
func (a *action) CreateProject() (*Project, error) {
	// TODO(displague) can we use a project description with more fields?
	//projectDescription := os.Getenv("GITHUB_SERVER_URL") + "/" + os.Getenv("GITHUB_REPOSITORY") + " " + os.Getenv("GITHUB_SHA")
	createOpts := metal.ProjectCreateFromRootInput{
		Name: a.label,
	}

	if a.organizationID != "" {
		createOpts.OrganizationId = &a.organizationID
	}

	log.Println("Creating Project")
	project, _, err := a.client.ProjectsApi.CreateProject(context.Background()).ProjectCreateFromRootInput(createOpts).Execute()
	if err != nil {
		return nil, err
	}

	p := &Project{Project: project}

	log.Println("Creating Keys")
	for _, f := range []func(*metal.APIClient) error{
		p.createSSHKey,
		p.createAPIKey,
	} {
		if err := f(a.client); err != nil {
			return nil, err
		}
	}

	return p, nil
}

// GenProjectName will generate a
func GenProjectName(sha string) string {
	prefix := "sha"
	if sha == "" {
		prefix = "rnd"
		sha = RandomString(16)
	}
	return "GHACTION-" + prefix + sha
}

// generatePrivateKey creates a RSA Private Key of specified byte size
func generatePrivateKey() (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	log.Println("Private Key generated")
	return privateKey, nil
}

// generatePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func generatePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	log.Println("Public key generated")
	return pubKeyBytes, nil
}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

func (p *Project) createSSHKey(c *metal.APIClient) error {
	key, err := generatePrivateKey()
	if err != nil {
		return err
	}

	pubKeyBytes, err := generatePublicKey(&key.PublicKey)
	if err != nil {
		return err
	}
	pubKey := string(pubKeyBytes)

	createOpts := metal.SSHKeyCreateInput{
		Label: p.Project.Name,
		Key:   &pubKey,
	}

	_, _, err = c.SSHKeysApi.CreateProjectSSHKey(context.Background(), p.Project.GetId()).SSHKeyCreateInput(createOpts).Execute()
	if err != nil {
		return err
	}

	privateKeyBytes := encodePrivateKeyToPEM(key)
	p.SSHPrivateKey = string(privateKeyBytes)
	p.SSHPublicKey = string(pubKey)
	return nil
}

func (p *Project) createAPIKey(c *metal.APIClient) error {
	createOpts := metal.AuthTokenInput{
		Description: p.Project.Name,
	}

	log.Println("Creating Project API Key")
	apiKey, _, err := c.AuthenticationApi.CreateProjectAPIKey(context.Background(), p.Project.GetId()).AuthTokenInput(createOpts).Execute()
	if err != nil {
		return err
	}

	p.APIToken = apiKey.GetToken()
	return nil
}

func RandomString(size int) string {
	r := mrand.New(mrand.NewSource(time.Now().UnixNano()))

	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, size)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}
