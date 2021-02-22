// Portions taken from https://gist.github.com/devinodaniel/8f9b8a4f31573f428f29ec0e884e6673
package action

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	mrand "math/rand"
	"time"

	"github.com/packethost/packngo"
	"golang.org/x/crypto/ssh"
)

const (
	bitSize = 4096
)

type action struct {
	client *packngo.Client

	label          string
	organizationID string
}

type Project struct {
	Project *packngo.Project

	SSHPrivateKey string
	APIToken      string
}

// NewAction returns an action with a Packngo client
func NewAction(apiToken, organizationID, label string) (*action, error) {
	client := packngo.NewClientWithAuth("metal-project-action", apiToken, nil)

	return &action{
		organizationID: organizationID,
		label:          label,
		client:         client,
	}, nil
}

// CreateProject
//
// Create an Equinix Metal project with API keys and project SSH Keys generated at {filename}
// and {filename}.pub
func (a *action) CreateProject(filename string) (*Project, error) {
	// TODO(displague) can we use a project description with more fields?
	//projectDescription := os.Getenv("GITHUB_SERVER_URL") + "/" + os.Getenv("GITHUB_REPOSITORY") + " " + os.Getenv("GITHUB_SHA")
	createOpts := &packngo.ProjectCreateRequest{
		Name:           a.label,
		OrganizationID: a.organizationID,
	}

	log.Println("Creating Project")
	project, _, err := a.client.Projects.Create(createOpts)
	if err != nil {
		return nil, err
	}

	p := &Project{Project: project}

	log.Println("Creating Keys")
	for _, f := range []func(*packngo.Client) error{
		func(c *packngo.Client) error {
			return p.createSSHKey(c, filename)
		},
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

// writeKeyFile writes keys to a file
func writeKeyFile(key []byte, filename string) error {
	err := ioutil.WriteFile(filename, key, 0600)
	if err != nil {
		log.Println("Could not write", err)
	}
	log.Println("Wrote", filename)
	return err
}

func (p *Project) createSSHKey(c *packngo.Client, filename string) error {
	key, err := generatePrivateKey()
	if err != nil {
		return err
	}

	if writeKeyFile(encodePrivateKeyToPEM(key), filename); err != nil {
		return err
	}

	pubKey, err := generatePublicKey(&key.PublicKey)
	if err != nil {
		return err
	}

	if err = writeKeyFile(pubKey, filename+".pub"); err != nil {
		return err
	}

	createOpts := &packngo.SSHKeyCreateRequest{
		Label:     p.Project.Name,
		ProjectID: p.Project.ID,
		Key:       string(pubKey),
	}

	_, _, err = c.SSHKeys.Create(createOpts)
	if err != nil {
		return err
	}

	privateKeyBytes := encodePrivateKeyToPEM(key)
	p.SSHPrivateKey = string(privateKeyBytes)
	return nil
}

func (p *Project) createAPIKey(c *packngo.Client) error {
	createOpts := &packngo.APIKeyCreateRequest{
		Description: p.Project.Name,
		ProjectID:   p.Project.ID,
	}

	log.Println("Creating Project API Key")
	apiKey, _, err := c.APIKeys.Create(createOpts)
	if err != nil {
		return err
	}

	p.APIToken = apiKey.Token
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
