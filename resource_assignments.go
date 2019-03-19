package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kawasima/bouncr-client-go"
)

func resourceBouncrAssignments() *schema.Resource {
	return &schema.Resource{
		Create: resourceAssignmentsCreate,
		Read:   resourceAssignmentsRead,
		Update: resourceAssignmentsUpdate,
		Delete: resourceAssignmentsDelete,

		Schema: map[string]*schema.Schema{
			"assignment": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group": {
							Type: schema.TypeString,
							Required: true,
						},
						"role": {
							Type: schema.TypeString,
							Required: true,
						},
						"realm": {
							Type: schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}

}

func resourceAssignmentsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)

	assignmentsRequest := []bouncr.AssignmentRequest{}
	if attr := d.Get("assignment").([]interface{}); len(attr) > 0 {
		for _, v := range attr {
			val, ok := v.(map[string]interface{})
			if ok && val != nil {
				assignmentRequest := &bouncr.AssignmentRequest{
					Group: bouncr.IdObject{ Name: val["group"].(string) },
					Role:  bouncr.IdObject{ Name: val["role"].(string) },
					Realm: bouncr.IdObject{ Name: val["realm"].(string) },
				}
				assignmentsRequest = append(assignmentsRequest, *assignmentRequest)
			}
		}
	}

	_, err := client.CreateAssignments(&assignmentsRequest)
	if err != nil {
		return err
	}
	return resourceAssignmentsRead(d, meta)
}

func resourceAssignmentsRead(d *schema.ResourceData, meta interface{}) error {
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

func resourceAssignmentsUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAssignmentsDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
