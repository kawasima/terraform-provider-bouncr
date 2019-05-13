package main

import (
	"log"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kawasima/bouncr-client-go"
)

func resourceBouncrUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Exists: resourceUserExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"account": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"password": &schema.Schema {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_profiles": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)
	input := &bouncr.UserCreateRequest{
		"account":      d.Get("account").(string),
	}
	for key, value := range d.Get("user_profiles").(map[string]interface{}) {
		(*input)[key] = value
	}

	user, err := client.CreateUser(input)
	if err != nil {
		return err
	}

	if p := d.Get("password").(string); p != "" {
		_, err := client.CreatePasswordCredential(&bouncr.PasswordCredentialCreateRequest{
			Account:  user.Account,
			Password: p,
		})
		if err != nil {
			return err
		}
	}
	d.SetId(user.Account)
	return resourceUserRead(d, meta)
}

func resourceUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)

	user, err := client.FindUser(d.Get("account").(string))
	if err != nil {
		return err
	}

	d.Set("id", user.ID)
	d.Set("account", user.Account)
	d.Set("user_profiles", user.UserProfiles)

	return nil
}

func resourceUserExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bouncr.Client)

	user, err := client.FindUser(d.Id())

	if err != nil {
		return false, err
	}
	return bool(user.Account != ""), nil
}

func resourceUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)
	input := &bouncr.UserUpdateRequest{
		"account":     d.Get("account").(string),
	}
	for key, value := range d.Get("user_profiles").(map[string]interface{}) {
		(*input)[key] = value
	}

	user, err := client.UpdateUser(d.Get("account").(string), input)
	if err != nil {
		return err
	}
	d.SetId(user.Account)

	log.Printf("[DEBUG] user %q updated.", d.Id())
	return resourceUserRead(d, meta)

}

func resourceUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)

	err := client.DeleteUser(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] bouncr user %q deleted.", d.Id())
	d.SetId("")
	return nil
}
