// Modified by Stegra AB for the Stegra-maintained distribution.
// SPDX-License-Identifier: Apache-2.0
package ds

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"terraform-provider-sra/api"
	"terraform-provider-sra/bt/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecurityProviderDataSourceIsRegistered(t *testing.T) {
	for _, factory := range DatasourceList() {
		var metadata datasource.MetadataResponse
		factory().Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "sra"}, &metadata)
		if metadata.TypeName == "sra_security_provider" {
			return
		}
	}

	t.Fatal("sra_security_provider was not registered in DatasourceList")
}

func TestSecurityProviderDataSourceRead(t *testing.T) {
	requests := 0
	server := securityProviderTestServer(t, func(w http.ResponseWriter, req *http.Request) {
		requests++
		assert.Equal(t, "100", req.URL.Query().Get("per_page"))
		_, err := w.Write([]byte(`[
			{"id":12,"name":"Staff","enabled":true,"type":"saml","user_authentication":true,"group_lookup":true},
			{"id":15,"name":"Nessco","enabled":true,"type":"saml","user_authentication":true,"group_lookup":true}
		]`))
		require.NoError(t, err)
	})
	t.Cleanup(server.Close)

	managed := &securityProviderDataSource{apiClient: securityProviderTestClient(t, server)}
	sch := securityProviderTestSchema(t, managed)
	config := securityProviderTestConfig(t, sch, "nessco")
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: config.Raw}}

	managed.Read(context.Background(), datasource.ReadRequest{Config: config}, &resp)
	require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)

	var state models.SecurityProvider
	require.False(t, resp.State.Get(context.Background(), &state).HasError())
	assert.Equal(t, int64(15), state.ID.ValueInt64())
	assert.Equal(t, "nessco", state.Name.ValueString())
	assert.Equal(t, "saml", state.Type.ValueString())
	assert.True(t, state.Enabled.ValueBool())
	assert.True(t, state.UserAuthentication.ValueBool())
	assert.True(t, state.GroupLookup.ValueBool())
	assert.Equal(t, 1, requests)
}

func TestSecurityProviderDataSourceRequiresExactlyOneMatch(t *testing.T) {
	for name, response := range map[string]string{
		"missing":   `[{"id":12,"name":"Staff"}]`,
		"duplicate": `[{"id":15,"name":"Nessco"},{"id":16,"name":"NESSCO"}]`,
	} {
		t.Run(name, func(t *testing.T) {
			server := securityProviderTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
				_, err := w.Write([]byte(response))
				require.NoError(t, err)
			})
			t.Cleanup(server.Close)

			managed := &securityProviderDataSource{apiClient: securityProviderTestClient(t, server)}
			sch := securityProviderTestSchema(t, managed)
			config := securityProviderTestConfig(t, sch, "Nessco")
			resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: config.Raw}}

			managed.Read(context.Background(), datasource.ReadRequest{Config: config}, &resp)
			require.True(t, resp.Diagnostics.HasError())
			assert.Equal(t, "Security Provider lookup failed", resp.Diagnostics.Errors()[0].Summary())
		})
	}
}

func securityProviderTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if req.URL.Path == "/oauth2/token" {
			_, err := w.Write([]byte(`{"token_type":"Bearer","expires_in":3600,"access_token":"test-token"}`))
			require.NoError(t, err)
			return
		}

		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, "/api/config/v1/security-provider", req.URL.Path)
		handler(w, req)
	}))
}

func securityProviderTestClient(t *testing.T, server *httptest.Server) *api.APIClient {
	t.Helper()
	clientID, clientSecret := "id", "secret"
	client, err := api.NewClient(server.URL, &clientID, &clientSecret)
	require.NoError(t, err)
	client.Product = api.ProductPRA
	return client
}

func securityProviderTestSchema(t *testing.T, managed *securityProviderDataSource) schema.Schema {
	t.Helper()
	var resp datasource.SchemaResponse
	managed.Schema(context.Background(), datasource.SchemaRequest{}, &resp)
	require.False(t, resp.Diagnostics.HasError())
	return resp.Schema
}

func securityProviderTestConfig(t *testing.T, sch schema.Schema, name string) tfsdk.Config {
	t.Helper()
	ctx := context.Background()
	model := models.SecurityProvider{
		ID:                 types.Int64Unknown(),
		Name:               types.StringValue(name),
		Type:               types.StringUnknown(),
		Enabled:            types.BoolUnknown(),
		UserAuthentication: types.BoolUnknown(),
		GroupLookup:        types.BoolUnknown(),
	}
	var object types.Object
	diags := tfsdk.ValueFrom(ctx, model, sch.Type(), &object)
	require.False(t, diags.HasError(), "%v", diags)
	raw, err := object.ToTerraformValue(ctx)
	require.NoError(t, err)
	return tfsdk.Config{Schema: sch, Raw: raw}
}
