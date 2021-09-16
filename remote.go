package git

import (
	"github.com/go-git/go-git"
	"github.com/go-git/go-git/config"
	"github.com/go-git/go-git/plumbing"
	"github.com/go-git/go-git/storage/memory"
)

func (g *Git) NewRemote(name string) error {
	config := config.RemoteConfig{
		Name: name,
		URLs: []string{g.URL},
	}

	storer := memory.NewStorage()

	remote := git.NewRemote(storer, &config)

	g.Remote = remote

	return nil
}

func (g *Git) LsRemote() ([]*plumbing.Reference, error) {
	opts := git.ListOptions{
		Auth: g.Config.Auth,
	}

	return g.Remote.List(&opts)
}
