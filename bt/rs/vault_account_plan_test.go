// Modified by Stegra AB for the Stegra-maintained distribution.
// SPDX-License-Identifier: Apache-2.0
package rs

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestVaultAccountComputedValuesUseStateForUnknown(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		resource resource.Resource
		strings  []string
		bools    []string
		int64s   []string
	}{
		{
			name:     "ssh",
			resource: &vaultSSHAccountResource{},
			strings:  []string{"public_key", "private_key_public_cert", "last_checkout_timestamp"},
			bools:    []string{"personal"},
			int64s:   []string{"owner_user_id"},
		},
		{
			name:     "username_password",
			resource: &vaultUsernamePasswordAccountResource{},
			strings:  []string{"last_checkout_timestamp"},
			bools:    []string{"personal"},
			int64s:   []string{"owner_user_id"},
		},
		{
			name:     "token",
			resource: &vaultTokenAccountResource{},
			strings:  []string{"last_checkout_timestamp"},
			bools:    []string{"personal"},
			int64s:   []string{"owner_user_id"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var resp resource.SchemaResponse
			test.resource.Schema(context.Background(), resource.SchemaRequest{}, &resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("schema returned diagnostics: %v", resp.Diagnostics)
			}

			for _, name := range test.strings {
				assertStringStatePreserved(t, resp.Schema, name)
			}
			for _, name := range test.bools {
				assertBoolStatePreserved(t, resp.Schema, name)
			}
			for _, name := range test.int64s {
				assertInt64StatePreserved(t, resp.Schema, name)
			}
		})
	}
}

func TestAccountJumpItemAssociationUsesStateForUnknown(t *testing.T) {
	t.Parallel()

	association := accountJumpItemAssociationSchema()
	if len(association.PlanModifiers) != 1 {
		t.Fatalf("expected one jump_item_association plan modifier, got %d", len(association.PlanModifiers))
	}

	attributeTypes := association.GetType().(types.ObjectType).AttrTypes
	stateValue := types.ObjectNull(attributeTypes)
	request := planmodifier.ObjectRequest{
		ConfigValue: types.ObjectNull(attributeTypes),
		PlanValue:   types.ObjectUnknown(attributeTypes),
		State:       existingResourceState(),
		StateValue:  stateValue,
	}
	response := planmodifier.ObjectResponse{PlanValue: request.PlanValue}
	association.PlanModifiers[0].PlanModifyObject(context.Background(), request, &response)

	if !response.PlanValue.Equal(stateValue) {
		t.Fatalf("jump_item_association plan value was not preserved: %s", response.PlanValue)
	}
}

func assertStringStatePreserved(t *testing.T, resourceSchema schema.Schema, name string) {
	t.Helper()

	attribute := resourceSchema.Attributes[name].(schema.StringAttribute)
	stateValue := types.StringValue("preserved")
	request := planmodifier.StringRequest{
		ConfigValue: types.StringNull(),
		PlanValue:   types.StringUnknown(),
		State:       existingResourceState(),
		StateValue:  stateValue,
	}
	response := planmodifier.StringResponse{PlanValue: request.PlanValue}
	for _, modifier := range attribute.PlanModifiers {
		modifier.PlanModifyString(context.Background(), request, &response)
		request.PlanValue = response.PlanValue
	}
	if !response.PlanValue.Equal(stateValue) {
		t.Fatalf("%s plan value was not preserved: %s", name, response.PlanValue)
	}
}

func assertBoolStatePreserved(t *testing.T, resourceSchema schema.Schema, name string) {
	t.Helper()

	attribute := resourceSchema.Attributes[name].(schema.BoolAttribute)
	stateValue := types.BoolValue(true)
	request := planmodifier.BoolRequest{
		ConfigValue: types.BoolNull(),
		PlanValue:   types.BoolUnknown(),
		State:       existingResourceState(),
		StateValue:  stateValue,
	}
	response := planmodifier.BoolResponse{PlanValue: request.PlanValue}
	for _, modifier := range attribute.PlanModifiers {
		modifier.PlanModifyBool(context.Background(), request, &response)
		request.PlanValue = response.PlanValue
	}
	if !response.PlanValue.Equal(stateValue) {
		t.Fatalf("%s plan value was not preserved: %s", name, response.PlanValue)
	}
}

func assertInt64StatePreserved(t *testing.T, resourceSchema schema.Schema, name string) {
	t.Helper()

	attribute := resourceSchema.Attributes[name].(schema.Int64Attribute)
	stateValue := types.Int64Value(42)
	request := planmodifier.Int64Request{
		ConfigValue: types.Int64Null(),
		PlanValue:   types.Int64Unknown(),
		State:       existingResourceState(),
		StateValue:  stateValue,
	}
	response := planmodifier.Int64Response{PlanValue: request.PlanValue}
	for _, modifier := range attribute.PlanModifiers {
		modifier.PlanModifyInt64(context.Background(), request, &response)
		request.PlanValue = response.PlanValue
	}
	if !response.PlanValue.Equal(stateValue) {
		t.Fatalf("%s plan value was not preserved: %s", name, response.PlanValue)
	}
}

func existingResourceState() tfsdk.State {
	return tfsdk.State{Raw: tftypes.NewValue(tftypes.String, "existing")}
}
