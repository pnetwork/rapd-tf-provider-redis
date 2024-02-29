/*
Copyright © 2023 PNTL <license@pentium.network>
*/

package provider

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pnetwork/rapd-tf-provider-redis/internal/provider/user"
)

// Ensure RedisProvider satisfies various provider interfaces.
var _ provider.Provider = &RedisProvider{}

// RedisProvider defines the provider implementation.
type RedisProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	version string
}

// RedisProviderModel describes the provider data model.
type RedisProviderModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Host     types.String `tfsdk:"host"`
	Port     types.Int64  `tfsdk:"port"`
}

func (p *RedisProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "redis"
	resp.Version = p.version
}

func (p *RedisProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				MarkdownDescription: "redis username",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "redis password",
				Required:            true,
				Sensitive:           true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "redis host",
				Required:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "redis port",
				Optional:            true,
			},
		},
	}
}

func (p *RedisProvider) Configure(
	ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse,
) {

	var data RedisProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	port := data.Port.ValueInt64()
	if data.Port.IsNull() {
		port = 6379
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", data.Host.ValueString(), port),
		Password: data.Password.ValueString(),
		Username: data.Username.ValueString(),
	})
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *RedisProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		user.NewUserResource,
	}
}

func (p *RedisProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &RedisProvider{
			version: version,
		}
	}
}
