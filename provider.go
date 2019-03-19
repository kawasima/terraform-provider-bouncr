package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("BOUNCR_ACCOUNT", nil),
				Description: "Bouncr account",
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("BOUNCR_PASSWORD", nil),
				Description: "Bouncr password",
			},
			"base_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.EnvDefaultFunc("BOUNCR_URL", nil),
				Description: "Bouncr base URL",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"bouncr_application": resourceBouncrApplication(),
			"bouncr_user":        resourceBouncrUser(),
			"bouncr_group":       resourceBouncrGroup(),
			"bouncr_role":        resourceBouncrRole(),
			"bouncr_permission":  resourceBouncrPermission(),
			"bouncr_assignments": resourceBouncrAssignments(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Account:  d.Get("account").(string),
		Password: d.Get("password").(string),
		BaseURL:  d.Get("base_url").(string),
	}

	return config.NewClient()
}
