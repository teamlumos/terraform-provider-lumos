package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ provider.Provider = &lumosAppStoreProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &lumosAppStoreProvider{
			version: version,
		}
	}
}

// lumosAppStoreProvider is the provider implementation.
type lumosAppStoreProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type lumosAppStoreProviderModel struct {
	APIToken types.String `tfsdk:"api_token"`
	BaseUrl  types.String `tfsdk:"base_url"`
}

func (p *lumosAppStoreProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "lumos"
	resp.Version = p.version
}

func (p *lumosAppStoreProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				Optional:    false,
				Required:    true,
				Sensitive:   true,
				Description: "The Lumos API Token.",
			},
			"base_url": schema.StringAttribute{
				Optional:    true,
				Required:    false,
				Sensitive:   false,
				Description: "The Lumos API URL.",
			},
		},
	}
}

func (p *lumosAppStoreProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	// Retrieve provider data from configuration
	var config lumosAppStoreProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.APIToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Lumos API Token",
			"The provider cannot create a Lumos API client as there is an unknown configuration value for the Lumos API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LUMOS_API_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	apiToken := os.Getenv("LUMOS_API_TOKEN")
	if !config.APIToken.IsNull() {
		apiToken = config.APIToken.ValueString()
	}

	baseUrl := "https://api.lumos.com"

	if !config.BaseUrl.IsNull() {
		baseUrl = config.BaseUrl.ValueString()
	}
	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Lumos API Token",
			"The provider cannot create a Lumos API client as there is an unknown configuration value for the Lumos API token. "+
				"Set the host value in the configuration or use the LUMOS_API_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	lumosClient, err := NewLumosAPIClient(baseUrl, apiToken)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Lumos API Client",
			"An unexpected error occurred when creating the Lumos API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Lumos Client Error: "+err.Error(),
		)
	}

	resp.DataSourceData = lumosClient
	resp.ResourceData = lumosClient
}

func (p *lumosAppStoreProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewRequestablePermissionDataSource,
		NewUserDataSource,
		NewAppstoreAppDataSource,
	}
}

func (p *lumosAppStoreProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRequestablePermissionResource,
	}
}
