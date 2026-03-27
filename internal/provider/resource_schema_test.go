package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestPermissionResourceSchema(t *testing.T) {
	r := NewPermissionResource()
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %v", resp.Diagnostics)
	}
	assertRequiredAttr(t, resp, "name")
	assertRequiredAttr(t, resp, "description")
}

func TestUserResourceSchema(t *testing.T) {
	r := NewUserResource()
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %v", resp.Diagnostics)
	}
	assertRequiredAttr(t, resp, "account")
	assertOptionalAttr(t, resp, "password")
	assertOptionalAttr(t, resp, "user_profiles")
}

func TestApplicationResourceSchema(t *testing.T) {
	r := NewApplicationResource()
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %v", resp.Diagnostics)
	}
	for _, name := range []string{"name", "description", "pass_to", "virtual_path", "top_page"} {
		assertRequiredAttr(t, resp, name)
	}
	if _, ok := resp.Schema.Blocks["realm"]; !ok {
		t.Error("schema missing realm block")
	}
}

func TestGroupResourceSchema(t *testing.T) {
	r := NewGroupResource()
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %v", resp.Diagnostics)
	}
	assertRequiredAttr(t, resp, "name")
	assertRequiredAttr(t, resp, "description")
	assertOptionalAttr(t, resp, "members")
}

func TestRoleResourceSchema(t *testing.T) {
	r := NewRoleResource()
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %v", resp.Diagnostics)
	}
	assertRequiredAttr(t, resp, "name")
	assertRequiredAttr(t, resp, "description")
	assertOptionalAttr(t, resp, "permissions")
}

func TestAssignmentsResourceSchema(t *testing.T) {
	r := NewAssignmentsResource()
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %v", resp.Diagnostics)
	}
	if _, ok := resp.Schema.Attributes["id"]; !ok {
		t.Error("schema missing id attribute")
	}
	if _, ok := resp.Schema.Blocks["assignment"]; !ok {
		t.Error("schema missing assignment block")
	}
}

func TestOidcProviderResourceSchema(t *testing.T) {
	r := NewOidcProviderResource()
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %v", resp.Diagnostics)
	}
	for _, name := range []string{"name", "client_id", "client_secret", "scope", "response_type", "authorization_endpoint", "token_endpoint_auth_method", "redirect_uri"} {
		assertRequiredAttr(t, resp, name)
	}
	for _, name := range []string{"token_endpoint", "pkce_enabled", "jwks_uri", "issuer"} {
		assertOptionalAttr(t, resp, name)
	}
}

func TestOidcApplicationResourceSchema(t *testing.T) {
	r := NewOidcApplicationResource()
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %v", resp.Diagnostics)
	}
	assertRequiredAttr(t, resp, "name")
	assertRequiredAttr(t, resp, "description")
	assertRequiredAttr(t, resp, "grant_types")
	assertOptionalAttr(t, resp, "home_uri")
	assertOptionalAttr(t, resp, "callback_uri")
	assertOptionalAttr(t, resp, "permissions")
}

func assertRequiredAttr(t *testing.T, resp *resource.SchemaResponse, name string) {
	t.Helper()
	attr, ok := resp.Schema.Attributes[name]
	if !ok {
		t.Errorf("schema missing %q attribute", name)
		return
	}
	if !attr.IsRequired() {
		t.Errorf("attribute %q should be required", name)
	}
}

func assertOptionalAttr(t *testing.T, resp *resource.SchemaResponse, name string) {
	t.Helper()
	attr, ok := resp.Schema.Attributes[name]
	if !ok {
		t.Errorf("schema missing %q attribute", name)
		return
	}
	if !attr.IsOptional() && !attr.IsComputed() {
		t.Errorf("attribute %q should be optional or computed", name)
	}
}
