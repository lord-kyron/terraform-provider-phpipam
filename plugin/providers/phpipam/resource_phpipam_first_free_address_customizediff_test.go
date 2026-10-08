package phpipam

import (
	"context"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// buildRawConfig turns a plain values map into a cty.Value conforming to the
// resource's implied schema type, so it can be used as InstanceState.RawConfig
// - mirroring what the real Terraform protocol layer does before calling
// into CustomizeDiff. Fields not present in values are left null, matching an
// attribute that's genuinely absent from the HCL config.
func buildRawConfig(t *testing.T, res *schema.Resource, values map[string]interface{}) cty.Value {
	t.Helper()

	attrTypes := res.CoreConfigSchema().ImpliedType().AttributeTypes()
	vals := make(map[string]cty.Value, len(attrTypes))

	for name, at := range attrTypes {
		v, ok := values[name]
		if !ok {
			vals[name] = cty.NullVal(at)
			continue
		}

		switch {
		case at == cty.String:
			vals[name] = cty.StringVal(v.(string))
		case at == cty.Number:
			vals[name] = cty.NumberIntVal(int64(v.(int)))
		case at == cty.Bool:
			vals[name] = cty.BoolVal(v.(bool))
		case at.IsListType() || at.IsSetType():
			raw := v.([]interface{})
			elemType := at.ElementType()
			elems := make([]cty.Value, len(raw))
			for i, e := range raw {
				switch elemType {
				case cty.Number:
					elems[i] = cty.NumberIntVal(int64(e.(int)))
				case cty.String:
					elems[i] = cty.StringVal(e.(string))
				default:
					t.Fatalf("unsupported element type for %s: %s", name, elemType.FriendlyName())
				}
			}
			switch {
			case len(elems) == 0:
				vals[name] = cty.ListValEmpty(elemType)
			case at.IsSetType():
				vals[name] = cty.SetVal(elems)
			default:
				vals[name] = cty.ListVal(elems)
			}
		default:
			t.Fatalf("unsupported type for %s: %s", name, at.FriendlyName())
		}
	}

	return cty.ObjectVal(vals)
}

// TestResourcePHPIPAMFirstFreeAddressCustomizeDiff exercises
// resourcePHPIPAMFirstFreeAddressCustomizeDiff directly against the SDK's
// diffing engine, without needing a live PHPIPAM server. It guards against
// two past regressions: a genuine subnet_id change must still force a
// replacement, and a config-only migration to an equivalent subnet_ids list
// must not.
func TestResourcePHPIPAMFirstFreeAddressCustomizeDiff(t *testing.T) {
	res := resourcePHPIPAMFirstFreeAddress()

	priorStateAttrs := map[string]interface{}{
		"subnet_id":   139,
		"subnet_ids":  []interface{}{139},
		"ip_address":  "10.10.15.1",
		"address_id":  5000,
		"hostname":    "tf-test.cust1.local",
		"description": "Terraform test address",
	}

	cases := []struct {
		name            string
		config          map[string]interface{}
		wantRequiresNew bool
	}{
		{
			name: "changing subnet_id alone still forces a replacement",
			config: map[string]interface{}{
				"subnet_id":   140,
				"hostname":    "tf-test.cust1.local",
				"description": "Terraform test address",
			},
			wantRequiresNew: true,
		},
		{
			name: "moving to subnet_ids that still contains the current subnet is a no-op",
			config: map[string]interface{}{
				"subnet_ids":  []interface{}{139, 140},
				"hostname":    "tf-test.cust1.local",
				"description": "Terraform test address",
			},
			wantRequiresNew: false,
		},
		{
			name: "moving to subnet_ids without the current subnet still forces a replacement",
			config: map[string]interface{}{
				"subnet_ids":  []interface{}{140, 150},
				"hostname":    "tf-test.cust1.local",
				"description": "Terraform test address",
			},
			wantRequiresNew: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, res.Schema, priorStateAttrs)
			d.SetId("5000")
			state := d.State()
			state.RawConfig = buildRawConfig(t, res, tc.config)

			cfg := terraform.NewResourceConfigRaw(tc.config)

			diff, err := res.SimpleDiff(context.Background(), state, cfg, nil)
			if err != nil {
				t.Fatalf("diff error: %s", err)
			}

			requiresNew := diff != nil && diff.RequiresNew()
			if requiresNew != tc.wantRequiresNew {
				t.Fatalf("wantRequiresNew=%v, got %v; diff: %#v", tc.wantRequiresNew, requiresNew, diff)
			}
		})
	}
}
