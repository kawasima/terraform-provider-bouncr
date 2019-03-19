package main

import (
	"net/url"
	"github.com/hashicorp/terraform/helper/logging"
	"github.com/kawasima/bouncr-client-go"
)

type Config struct {
	Account  string
	Password string
	BaseURL  string
}

func (c *Config) NewClient() (*bouncr.Client, error) {
	client := bouncr.NewClient(c.Account, c.Password)
	if logging.IsDebugOrHigher() {
		client.Verbose = true
	}
	if c.BaseURL != "" {
		u, err :=url.Parse(c.BaseURL)
		if err != nil {
			panic("BOUNCR_URL is not url")
		}
		client.BaseURL = u
	}
	return client, nil
}
