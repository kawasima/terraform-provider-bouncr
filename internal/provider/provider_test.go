package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func testAccProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"bouncr": providerserver.NewProtocol6WithError(New()),
	}
}

func TestProviderNew(t *testing.T) {
	p := New()
	if p == nil {
		t.Fatal("New() returned nil")
	}
}

func TestProviderMetadata(t *testing.T) {
	p := New()
	req := provider.MetadataRequest{}
	resp := &provider.MetadataResponse{}
	p.Metadata(context.Background(), req, resp)
	if resp.TypeName != "bouncr" {
		t.Errorf("TypeName = %q, want %q", resp.TypeName, "bouncr")
	}
}

func TestProviderSchema(t *testing.T) {
	p := New()
	req := provider.SchemaRequest{}
	resp := &provider.SchemaResponse{}
	p.Schema(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %v", resp.Diagnostics)
	}

	attrs := resp.Schema.Attributes
	for _, name := range []string{"client_id", "client_secret", "base_url"} {
		if _, ok := attrs[name]; !ok {
			t.Errorf("schema missing %q attribute", name)
		}
	}
}

func TestProviderResources(t *testing.T) {
	bp := New().(*bouncrProvider)
	resources := bp.Resources(context.Background())

	expectedCount := 8
	if len(resources) != expectedCount {
		t.Errorf("expected %d resources, got %d", expectedCount, len(resources))
	}
}
