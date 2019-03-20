package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/kawasima/bouncr-client-go"
)

func TestBouncrUser_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBouncrUserDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckBouncrUserConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"bouncr_user.user1", "account", "user1"),
				),
			},
		},

	})
}

func TestBouncrUser_WithPassword(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBouncrUserDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckBouncrUserConfig_withPassword,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"bouncr_user.user2", "account", "user2"),
				),
			},
		},

	})
}

func testAccCheckBouncrUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*bouncr.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bouncr_user" {
			continue
		}

		_, err := client.FindUser(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("User stil exists")
		}
	}

	return nil
}

const testAccCheckBouncrUserConfig_basic = `
resource "bouncr_user" "user1" {
  account = "user1"
  user_profiles = {
    email = "user1@example.com"
    name  = "Test User1"
  }
}
`

const testAccCheckBouncrUserConfig_withPassword = `
resource "bouncr_user" "user2" {
  account = "user2"
  password = "pass1234"
  user_profiles = {
    email = "user2@example.com"
    name  = "Test User2"
  }
}
`
