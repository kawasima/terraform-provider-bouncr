package main

import (
	"log"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kawasima/bouncr-client-go"
)

func resourceBouncrOidcProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceOidcProviderCreate,
		Read:   resourceOidcProviderRead,
		Update: resourceOidcProviderUpdate,
		Delete: resourceOidcProviderDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"client_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"client_secret": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"scope": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"response_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"authorization_endpoint": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"token_endpoint": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"token_endpoint_auth_method": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceOidcProviderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)
	input := &bouncr.OidcProviderCreateRequest{
		Name:         d.Get("name").(string),
		ClientId:     d.Get("client_id").(string),
		ClientSecret: d.Get("client_secret").(string),
		Scope:        d.Get("scope").(string),
		ResponseType: d.Get("response_type").(string),
		AuthorizationEndpoint:   d.Get("authorization_endpoint").(string),
		TokenEndpoint:           d.Get("token_endpoint").(string),
		TokenEndpointAuthMethod: d.Get("token_endpoint_auth_method").(string),
	}

	oidcProvider, err := client.CreateOidcProvider(input)
	if err != nil {
		return err
	}
	d.SetId(oidcProvider.Name)

	return resourceOidcProviderRead(d, meta)
}

func resourceOidcProviderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)

	oidcProvider, err := client.FindOidcProvider(d.Id())
	if err != nil {
		return err
	}

	d.Set("id", oidcProvider.ID)
	d.Set("name", oidcProvider.Name)
	d.Set("client_id", oidcProvider.ClientId)
	d.Set("client_secret", oidcProvider.ClientSecret)
	d.Set("scope", oidcProvider.Scope)
	d.Set("response_type", oidcProvider.ResponseType)
	d.Set("authorization_endpoint", oidcProvider.AuthorizationEndpoint)
	d.Set("token_endpoint", oidcProvider.TokenEndpoint)
	d.Set("token_endpoint_auth_method", oidcProvider.TokenEndpointAuthMethod)

	return nil
}

func resourceOidcProviderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)
	input := &bouncr.OidcProviderUpdateRequest{
		Name:         d.Get("name").(string),
		ClientId:     d.Get("client_id").(string),
		ClientSecret: d.Get("client_secret").(string),
		Scope:        d.Get("scope").(string),
		ResponseType: d.Get("response_type").(string),
		AuthorizationEndpoint:   d.Get("authorization_endpoint").(string),
		TokenEndpoint:           d.Get("token_endpoint").(string),
		TokenEndpointAuthMethod: d.Get("token_endpoint_auth_method").(string),
	}

	oidcProvider, err := client.UpdateOidcProvider(d.Get("name").(string), input)
	if err != nil {
		return err
	}
	d.SetId(oidcProvider.Name)

	log.Printf("[DEBUG] oidcProvider %q updated.", d.Id())
	return resourceOidcProviderRead(d, meta)

}

func resourceOidcProviderDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)

	err := client.DeleteOidcProvider(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] bouncr oidcProvider %q deleted.", d.Id())
	d.SetId("")
	return nil
}
