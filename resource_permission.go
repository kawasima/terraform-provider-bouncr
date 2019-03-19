package main

import (
	"log"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kawasima/bouncr-client-go"
)

func resourceBouncrPermission() *schema.Resource {
	return &schema.Resource{
		Create: resourcePermissionCreate,
		Read:   resourcePermissionRead,
		Update: resourcePermissionUpdate,
		Delete: resourcePermissionDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourcePermissionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)
	input := &bouncr.PermissionCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	permission, err := client.CreatePermission(input)
	if err != nil {
		return err
	}
	d.SetId(permission.Name)

	return resourcePermissionRead(d, meta)
}

func resourcePermissionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)

	permission, err := client.FindPermission(d.Id())
	if err != nil {
		return err
	}

	d.Set("id", permission.ID)
	d.Set("name", permission.Name)
	d.Set("description", permission.Description)

	return nil
}

func resourcePermissionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)
	input := &bouncr.PermissionUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	permission, err := client.UpdatePermission(d.Get("name").(string), input)
	if err != nil {
		return err
	}
	d.SetId(permission.Name)

	log.Printf("[DEBUG] permission %q updated.", d.Id())
	return resourcePermissionRead(d, meta)

}

func resourcePermissionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)

	err := client.DeletePermission(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] bouncr permission %q deleted.", d.Id())
	d.SetId("")
	return nil
}
