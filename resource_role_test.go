package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/kawasima/bouncr-client-go"
)

func TestBouncrRole_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBouncrRoleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckBouncrRoleConfig_basic,
				Check: resource.ComposeTestCheckFunc(
				),
			},
		},

	})
}

func testAccCheckBouncrRoleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*bouncr.Client)

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "bouncr_role":
			_, err := client.FindRole(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Role stil exists")
			}
		case "bouncr_permission":
			_, err := client.FindPermission(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Role stil exists")
			}
		}

	}

	return nil
}

const testAccCheckBouncrRoleConfig_basic = `
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
`
