// Modified by Stegra AB for the Stegra-maintained distribution.
// SPDX-License-Identifier: Apache-2.0
package ds

import (
	"context"
	"fmt"
	"strings"

	"terraform-provider-sra/api"
	"terraform-provider-sra/bt/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &securityProviderDataSource{}
	_ datasource.DataSourceWithConfigure = &securityProviderDataSource{}
)

func newSecurityProviderDataSource() datasource.DataSource {
	return &securityProviderDataSource{}
}

type securityProviderDataSource struct {
	apiClient *api.APIClient
}

func (d *securityProviderDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_provider"
}

func (d *securityProviderDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Looks up exactly one Security Provider by name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The unique identifier assigned to the Security Provider.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The exact, case-insensitive name of the Security Provider.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the Security Provider.",
			},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the Security Provider is enabled.",
			},
			"user_authentication": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the Security Provider authenticates users.",
			},
			"group_lookup": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the Security Provider looks up user groups.",
			},
		},
	}
}

func (d *securityProviderDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.apiClient = req.ProviderData.(*api.APIClient)
}

func (d *securityProviderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config models.SecurityProvider
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	providers, err := api.ListItems[api.SecurityProvider](d.apiClient, map[string]string{"per_page": "100"})
	if err != nil {
		resp.Diagnostics.AddError("Unable to list Security Providers", err.Error())
		return
	}

	matches := make([]api.SecurityProvider, 0, 1)
	for _, provider := range providers {
		if strings.EqualFold(provider.Name, config.Name.ValueString()) {
			matches = append(matches, provider)
		}
	}

	if len(matches) != 1 {
		resp.Diagnostics.AddError(
			"Security Provider lookup failed",
			fmt.Sprintf("Expected exactly one Security Provider named %q, found %d.", config.Name.ValueString(), len(matches)),
		)
		return
	}

	provider := matches[0]
	if provider.ID == nil {
		resp.Diagnostics.AddError("Invalid Security Provider response", "The appliance returned a Security Provider without an id.")
		return
	}

	state := models.SecurityProvider{
		ID:                 types.Int64Value(int64(*provider.ID)),
		Name:               config.Name,
		Type:               types.StringValue(provider.Type),
		Enabled:            types.BoolValue(provider.Enabled),
		UserAuthentication: types.BoolValue(provider.UserAuthentication),
		GroupLookup:        types.BoolValue(provider.GroupLookup),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
