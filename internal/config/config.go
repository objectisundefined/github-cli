package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/BurntSushi/toml"
)

// Repository is the config for the repository at Github
type Repository struct {
	Owner string
	Name  string
}

func (r Repository) String() string {
	return fmt.Sprintf("%s/%s", r.Owner, r.Name)
}

// UnmarshalText implementes TOML interface.
func (r *Repository) UnmarshalText(text []byte) error {
	seps := strings.Split(string(text), "/")
	if len(seps) != 2 {
		return fmt.Errorf("repository must be owner/name format")
	}

	r.Owner, r.Name = seps[0], seps[1]
	return nil
}

// Slack configuration
type Slack struct {
	Token   string `toml:"token"`
	Channel string `toml:"channel"`
	User    string `toml:"user"`
}

// Config is the configuration for the tool
type Config struct {
	Account string       `toml:"account"`
	Token   string       `toml:"token"`
	Slack   Slack        `toml:"slack"`
	Repos   []Repository `toml:"repos"`
}

// NewConfigFromFile creates the configuration from file
func NewConfigFromFile(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)

	// If the file is not found, return an empty config
	if err != nil {
		return &Config{
			Account: "",
			Token:   "",
		}, nil
	}

	c := new(Config)

	if err = toml.Unmarshal(data, c); err != nil {
		return nil, err
	}

	return c, nil
}

// FindRepo finds a repository.
func (c *Config) FindRepo(owner string, name string) *Repository {
	for _, repo := range c.Repos {
		if len(owner) == 0 && repo.Name == name {
			return &repo
		}

		if repo.Name == name && repo.Owner == owner {
			return &repo
		}
	}

	return nil
}
