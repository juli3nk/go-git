package git

import (
	"net/url"

	"golang.org/x/crypto/ssh"
	"github.com/go-git/go-git"
	"github.com/go-git/go-git/plumbing/object"
	"github.com/go-git/go-git/plumbing/transport"
	githttp "github.com/go-git/go-git/plumbing/transport/http"
	gitssh "github.com/go-git/go-git/plumbing/transport/ssh"
)

type Git struct {
	URL        string
	Config     Config
	Remote     *git.Remote
	Repository *git.Repository
}

type Config struct {
	User object.Signature
	Auth transport.AuthMethod
}

func New(url string) (*Git, error) {
	conf := Git{
		URL: url,
	}

	return &conf, nil
}

func (g *Git) SetConfigUser(name, email string) error {
	g.Config.User = object.Signature{
		Name:  name,
		Email: email,
	}

	return nil
}

func (g *Git) SetAuth(username, secretType, secret string) error {
	u, err := url.Parse(g.URL)
	if err != nil {
		return err
	}

	if u.Scheme == "http" || u.Scheme == "https" {
		a := &githttp.BasicAuth{
			Username: username,
			Password: secret,
		}

		g.Config.Auth = a
	}

	if u.Scheme == "ssh" {
		var a gitssh.AuthMethod

		if secretType == "password" {
			a = &gitssh.Password{
				User:     username,
				Password: secret,
			}
		}

		if secretType == "pubkey" {
			a, err = gitssh.NewPublicKeys(username, []byte(secret), "")
			if err != nil {
				return err
			}
		}

		a.(*gitssh.PublicKeys).HostKeyCallback = ssh.InsecureIgnoreHostKey()
		g.Config.Auth = a
	}

	return nil
}
