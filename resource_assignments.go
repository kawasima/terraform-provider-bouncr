package main

import (
	"github.com/rs/xid"
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
	guid := xid.New()
	d.SetId(guid.String())
	return nil
}

func resourceAssignmentsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bouncr.Client)
	assignments := []map[string]interface{}{}
	if attr := d.Get("assignment").([]interface{}); len(attr) > 0 {
		for _,v := range attr {
			val, ok := v.(map[string]interface{})
			if ok && val != nil {
				assignmentRequest := &bouncr.AssignmentRequest{
					Group: bouncr.IdObject{ Name: val["group"].(string) },
					Role:  bouncr.IdObject{ Name: val["role"].(string) },
					Realm: bouncr.IdObject{ Name: val["realm"].(string) },
				}
				assignment, err := client.FindAssignment(assignmentRequest)
				if err != nil {
					return err
				}
				assignments = append(assignments, map[string]interface{}{
					"group": assignment.Group.Name,
					"role":  assignment.Role.Name,
					"realm": assignment.Realm.Name,
				})
			}
		}
		d.Set("assignment", assignments)
	}
	return nil
}

func resourceAssignmentsUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAssignmentsDelete(d *schema.ResourceData, meta interface{}) error {
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

	err := client.DeleteAssignments(&assignmentsRequest)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
