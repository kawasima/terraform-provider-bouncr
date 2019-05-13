package main

import (
	"log"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kawasima/bouncr-client-go"
)

func resourceBouncrGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,
		Exists: resourceGroupExists,
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
			"members": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{ Type: schema.TypeString },
				Set:      schema.HashString,
			},
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)
	input := &bouncr.GroupCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	group, err := client.CreateGroup(input)
	if err != nil {
		return err
	}
	d.SetId(group.Name)

	if attr := d.Get("members").(*schema.Set); attr.Len() > 0 {
		users := expandStringSet(attr)
		client.AddUsersToGroup(group.Name, &users)
	}

	return resourceGroupRead(d, meta)
}

func resourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)

	group, err := client.FindGroup(d.Id())
	if err != nil {
		return err
	}

	d.Set("id", group.ID)
	d.Set("name", group.Name)
	d.Set("description", group.Description)

	return nil
}

func resourceGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bouncr.Client)

	group, err := client.FindGroup(d.Id())

	if err != nil {
		return false, err
	}
	return bool(group.Name != ""), nil
}

func resourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)
	input := &bouncr.GroupUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	group, err := client.UpdateGroup(d.Get("name").(string), input)
	if err != nil {
		return err
	}
	d.SetId(group.Name)

	log.Printf("[DEBUG] group %q updated.", d.Id())
	return resourceGroupRead(d, meta)

}

func resourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)

	err := client.DeleteGroup(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] bouncr group %q deleted.", d.Id())
	d.SetId("")
	return nil
}
