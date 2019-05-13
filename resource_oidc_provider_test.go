package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/kawasima/bouncr-client-go"
)

func TestBouncrOidcProvider_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBouncrOidcProviderDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckBouncrOidcProviderConfig_basic,
				Check: resource.ComposeTestCheckFunc(
				),
			},
		},

	})
}

func testAccCheckBouncrOidcProviderDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*bouncr.Client)

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "bouncr_oidc_provider":
			_, err := client.FindRole(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("OidcProvider stil exists")
			}
		}
	}

	return nil
}

const testAccCheckBouncrOidcProviderConfig_basic = `
resource "bouncr_oidc_provider" "google" {
  name            = "Google"
  client_id       = "google-no-client-id"
  client_secret   = "google-no-client-secret"
  scope           = "openid"
  response_type   = "code"

  authorization_endpoint     = "https://accounts.google.com/o/oauth2/v2/auth"
  token_endpoint             = "https://www.googleapis.com/oauth2/v4/token"
  token_endpoint_auth_method = "POST"
}
`
