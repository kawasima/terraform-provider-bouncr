package provider

import (
	"context"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kawasima/bouncr-client-go"
)

var _ provider.Provider = &bouncrProvider{}

type bouncrProvider struct{}

type bouncrProviderModel struct {
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	BaseURL      types.String `tfsdk:"base_url"`
}

func New() provider.Provider {
	return &bouncrProvider{}
}

func (p *bouncrProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "bouncr"
}

func (p *bouncrProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Optional:    true,
				Description: "Bouncr OAuth2 Client ID. Can also be set via BOUNCR_CLIENT_ID env var.",
			},
			"client_secret": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Bouncr OAuth2 Client Secret. Can also be set via BOUNCR_CLIENT_SECRET env var.",
			},
			"base_url": schema.StringAttribute{
				Optional:    true,
				Description: "Bouncr base URL. Can also be set via BOUNCR_URL env var.",
			},
		},
	}
}

func (p *bouncrProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	clientID := os.Getenv("BOUNCR_CLIENT_ID")
	clientSecret := os.Getenv("BOUNCR_CLIENT_SECRET")
	baseURL := os.Getenv("BOUNCR_URL")

	var data bouncrProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.ClientID.IsNull() {
		clientID = data.ClientID.ValueString()
	}
	if !data.ClientSecret.IsNull() {
		clientSecret = data.ClientSecret.ValueString()
	}
	if !data.BaseURL.IsNull() {
		baseURL = data.BaseURL.ValueString()
	}

	if clientID == "" {
		resp.Diagnostics.AddError("Missing Client ID", "The provider client_id must be set in the configuration or via the BOUNCR_CLIENT_ID environment variable.")
		return
	}
	if clientSecret == "" {
		resp.Diagnostics.AddError("Missing Client Secret", "The provider client_secret must be set in the configuration or via the BOUNCR_CLIENT_SECRET environment variable.")
		return
	}

	client := bouncr.NewClient(clientID, clientSecret)
	if baseURL != "" {
		u, err := url.Parse(baseURL)
		if err != nil {
			resp.Diagnostics.AddError("Invalid Base URL", err.Error())
			return
		}
		client.BaseURL = u
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *bouncrProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUserResource,
		NewApplicationResource,
		NewGroupResource,
		NewRoleResource,
		NewPermissionResource,
		NewAssignmentsResource,
		NewOidcProviderResource,
		NewOidcApplicationResource,
	}
}

func (p *bouncrProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}
