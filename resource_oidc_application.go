package main

import (
	"log"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kawasima/bouncr-client-go"
)

func resourceBouncrOidcApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceOidcApplicationCreate,
		Read:   resourceOidcApplicationRead,
		Update: resourceOidcApplicationUpdate,
		Delete: resourceOidcApplicationDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"home_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"callback_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"permissions": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{ Type: schema.TypeString },
				Set:      schema.HashString,
			},
		},
	}
}

func resourceOidcApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)
	input := &bouncr.OidcApplicationCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		HomeURL:     d.Get("home_url").(string),
		CallbackURL: d.Get("callback_url").(string),
		Permissions: d.Get("permissions").([]string),
	}

	oidcApplication, err := client.CreateOidcApplication(input)
	if err != nil {
		return err
	}
	d.SetId(oidcApplication.Name)

	return resourceOidcApplicationRead(d, meta)
}

func resourceOidcApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)

	oidcApplication, err := client.FindOidcApplication(d.Id())
	if err != nil {
		return err
	}

	d.Set("id", oidcApplication.ID)
	d.Set("name", oidcApplication.Name)
	d.Set("description", oidcApplication.Description)
	d.Set("home_url", oidcApplication.HomeURL)
	d.Set("callback_url", oidcApplication.CallbackURL)

	return nil
}

func resourceOidcApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)
	input := &bouncr.OidcApplicationUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		HomeURL:     d.Get("home_url").(string),
		CallbackURL: d.Get("callback_url").(string),
		Permissions: d.Get("permissions").([]string),
	}

	oidcApplication, err := client.UpdateOidcApplication(d.Get("name").(string), input)
	if err != nil {
		return err
	}
	d.SetId(oidcApplication.Name)

	log.Printf("[DEBUG] oidcApplication %q updated.", d.Id())
	return resourceOidcApplicationRead(d, meta)

}

func resourceOidcApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)

	err := client.DeleteOidcApplication(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] bouncr oidcApplication %q deleted.", d.Id())
	d.SetId("")
	return nil
}
