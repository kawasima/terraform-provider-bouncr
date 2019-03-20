package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/kawasima/bouncr-client-go"
)

func TestBouncrAssignments_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBouncrAssignmentsDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckBouncrAssignmentsConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"bouncr_assignments.assign", "assignment.0.group", "group1"),
					resource.TestCheckResourceAttr(
						"bouncr_assignments.assign", "assignment.0.role", "role4"),
					resource.TestCheckResourceAttr(
						"bouncr_assignments.assign", "assignment.0.realm", "realm1"),
				),
			},
		},

	})
}

func testAccCheckBouncrAssignmentsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*bouncr.Client)

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "bouncr_group":
			_, err := client.FindGroup(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Group stil exists")
			}
		case "bouncr_user":
			_, err := client.FindUser(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("User stil exists")
			}
		case "bouncr_role":
			_, err := client.FindRole(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("User stil exists")
			}
		case "bouncr_permission":
			_, err := client.FindPermission(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("User stil exists")
			}
		case "bouncr_application":
			_, err := client.FindApplication(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("User stil exists")
			}
		}
	}

	return nil
}

const testAccCheckBouncrAssignmentsConfig_basic = `
resource "bouncr_application" "testapp" {
  name         = "testapp"
  description  = "testapp"
  pass_to      = "http://localhost:3002/testapp"
  virtual_path = "/testapp"
  top_page     = "http://localhost:3000/testapp"

  realm {
    name        = "realm1"
    description = "This is a realm1"
    url         = ".*"
  }
}

resource "bouncr_user" "user1" {
  account = "user1"
  user_profiles = {
    email = "user1@example.com"
    name  = "Test User1"
  }
}
resource "bouncr_group" "group1" {
  name        = "group1"
  description = "group1"

  members = [
    "${bouncr_user.user1.id}"
  ]
}

resource "bouncr_permission" "perm1" {
  name        = "perm4"
  description = "perm4"
}

resource "bouncr_role" "role1" {
  name        = "role4"
  description = "role4"

  permissions = [
    "${bouncr_permission.perm1.id}"
  ]
}

resource "bouncr_assignments" "assign" {
  assignment {
    role  = "${bouncr_role.role1.id}"
    group = "${bouncr_group.group1.id}"
    realm = "${bouncr_application.testapp.realm[0].name}"
  }
}
`
