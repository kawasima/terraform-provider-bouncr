package main

import (
	"fmt"
	"testing"
	"strings"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/kawasima/bouncr-client-go"
)

func TestBouncrApplication_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBouncrApplicationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckBouncrApplicationConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"bouncr_application.testapp", "name", "testapp5"),
				),
			},
		},

	})
}

func testAccCheckBouncrApplicationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*bouncr.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bouncr_application" {
			continue
		}

		_, err := client.FindApplication(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Application stil exists")
		}
		if !strings.Contains(err.Error(), "404") {
			return err
		}
	}

	return nil
}

const testAccCheckBouncrApplicationConfig_basic = `
resource "bouncr_application" "testapp" {
  name         = "testapp5"
  description  = "testapp"
  pass_to      = "http://localhost:3002/testapp"
  virtual_path = "/testapp5"
  top_page     = "http://localhost:3000/testapp"

  realm {
    name        = "realm2"
    description = "This is realm1"
    url         = ".*"
  }
}
`
