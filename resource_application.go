package main

import (
	"log"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kawasima/bouncr-client-go"
)

func resourceBouncrApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationCreate,
		Read:   resourceApplicationRead,
		Update: resourceApplicationUpdate,
		Delete: resourceApplicationDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"pass_to": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"virtual_path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"top_page": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"realm": &schema.Schema {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Required: true,
						},
						"url": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)
	input := &bouncr.ApplicationCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		PassTo:      d.Get("pass_to").(string),
		VirtualPath: d.Get("virtual_path").(string),
		TopPage:     d.Get("top_page").(string),
	}

	application, err := client.CreateApplication(input)
	if err != nil {
		return err
	}
	d.SetId(application.Name)

	if attr := d.Get("realm").([]interface{}); len(attr) > 0 {
		for _, v := range attr {
			val, ok := v.(map[string]interface{})
			if ok && val != nil {
				realmRequest := &bouncr.RealmCreateRequest{
					Name:        val["name"].(string),
					Description: val["description"].(string),
					URL:         val["url"].(string),
				}
				_, err := client.CreateRealm(application.Name, realmRequest)
				if err != nil {
					return err
				}
			}
		}

	}

	return resourceApplicationRead(d, meta)
}

func resourceApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)

	application, err := client.FindApplication(d.Get("name").(string))
	if err != nil {
		return err
	}

	d.Set("id", application.ID)
	d.Set("name", application.Name)
	d.Set("description", application.Description)
	d.Set("pass_to", application.PassTo)
	d.Set("virtual_path", application.VirtualPath)
	d.Set("top_page", application.TopPage)

	return nil
}

func resourceApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)
	input := &bouncr.ApplicationUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		PassTo:      d.Get("pass_to").(string),
		VirtualPath: d.Get("virtual_path").(string),
		TopPage:     d.Get("top_page").(string),
	}

	application, err := client.UpdateApplication(d.Get("name").(string), input)
	if err != nil {
		return err
	}
	d.SetId(application.Name)

	log.Printf("[DEBUG] application %q updated.", d.Id())
	return resourceApplicationRead(d, meta)
}

func resourceApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)

	err := client.DeleteApplication(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] bouncr application %q deleted.", d.Id())
	d.SetId("")
	return nil
}
