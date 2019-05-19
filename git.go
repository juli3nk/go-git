package git

import (
	"net/url"
	"time"

	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type Git struct {
	URL        string
	Config     Config
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

func (g *Git) Init() error {
	r, err := git.PlainInit(".", false)
	if err != nil {
		return err
	}

	g.Repository = r

	return nil
}

func (g *Git) RemoteAdd(name string) error {
	_, err := g.Repository.CreateRemote(&config.RemoteConfig{
		Name: name,
		URLs: []string{g.URL},
	})

	return err
}

func (g *Git) RemoteRemove(name string) error {
	err := g.Repository.DeleteRemote(name)

	return err
}

func (g *Git) Clone(path string) error {
	directory := path
	if len(path) == 0 {
		directory = "test"
	}

	cloneOpts := &git.CloneOptions{
		URL:  g.URL,
		Auth: g.Config.Auth,
	}

	if err := cloneOpts.Validate(); err != nil {
		return err
	}

	r, err := git.PlainClone(directory, false, cloneOpts)
	if err != nil {
		return err
	}

	g.Repository = r

	return nil
}

func (g *Git) Checkout(ref string) error {
	w, err := g.Repository.Worktree()
	if err != nil {
		return err
	}

	hash := plumbing.NewHash(ref)

	w.Checkout(&git.CheckoutOptions{
		Hash: hash,
	})

	return nil
}

func (g *Git) Status() (git.Status, error) {
	w, err := g.Repository.Worktree()
	if err != nil {
		return nil, err
	}

	return w.Status()
}

func (g *Git) Add(path string) error {
	w, err := g.Repository.Worktree()
	if err != nil {
		return err
	}

	if _, err = w.Add(path); err != nil {
		return err
	}

	return nil
}

func (g *Git) Remove(path string) error {
	w, err := g.Repository.Worktree()
	if err != nil {
		return err
	}

	if _, err = w.Remove(path); err != nil {
		return err
	}

	return nil
}

func (g *Git) Commit(msg string) error {
	w, err := g.Repository.Worktree()
	if err != nil {
		return err
	}

	author := g.Config.User
	author.When = time.Now()

	opt := git.CommitOptions{
		Author: &author,
	}

	if _, err = w.Commit(msg, &opt); err != nil {
		return err
	}

	return nil
}

func (g *Git) Push(remoteName string) error {
	opt := git.PushOptions{
		RemoteName: remoteName,
		Auth:       g.Config.Auth,
	}

	err := g.Repository.Push(&opt)

	return err
}
