package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/kawasima/bouncr-client-go"
)

func TestBouncrGroup_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBouncrGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckBouncrGroupConfig_basic,
				Check: resource.ComposeTestCheckFunc(
				),
			},
		},

	})
}

func testAccCheckBouncrGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*bouncr.Client)

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "bouncr_group":
			_, err := client.FindGroup(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Group stil exists")
			}
		case "bouncr_permission":
			_, err := client.FindUser(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("User stil exists")
			}
		}

	}

	return nil
}

const testAccCheckBouncrGroupConfig_basic = `
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
`
