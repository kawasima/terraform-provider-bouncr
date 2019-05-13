package main

import (
	"log"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kawasima/bouncr-client-go"
)

func resourceBouncrRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceRoleCreate,
		Read:   resourceRoleRead,
		Update: resourceRoleUpdate,
		Delete: resourceRoleDelete,
		Exists: resourceRoleExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
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

func resourceRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)
	input := &bouncr.RoleCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	role, err := client.CreateRole(input)
	if err != nil {
		return err
	}
	d.SetId(role.Name)

	if attr := d.Get("permissions").(*schema.Set); attr.Len() > 0 {
		permissions := expandStringSet(attr)
		client.AddPermissionsToRole(role.Name, &permissions)
	}
	return resourceRoleRead(d, meta)
}

func resourceRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)

	role, err := client.FindRole(d.Id())
	if err != nil {
		return err
	}

	d.Set("id", role.ID)
	d.Set("name", role.Name)
	d.Set("description", role.Description)

	return nil
}

func resourceRoleExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bouncr.Client)

	role, err := client.FindRole(d.Id())

	if err != nil {
		return false, err
	}
	return bool(role.Name != ""), nil
}

func resourceRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)
	input := &bouncr.RoleUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	role, err := client.UpdateRole(d.Get("name").(string), input)
	if err != nil {
		return err
	}
	d.SetId(role.Name)

	log.Printf("[DEBUG] role %q updated.", d.Id())
	return resourceRoleRead(d, meta)

}

func resourceRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)

	err := client.DeleteRole(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] bouncr role %q deleted.", d.Id())
	d.SetId("")
	return nil
}
