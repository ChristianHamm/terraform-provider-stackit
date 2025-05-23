package observability

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stackitcloud/stackit-sdk-go/core/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/observability"
)

func TestMapFields(t *testing.T) {
	tests := []struct {
		description string
		input       *observability.Job
		expected    Model
		isValid     bool
	}{
		{
			"default_ok",
			&observability.Job{
				JobName: utils.Ptr("name"),
			},
			Model{
				Id:             types.StringValue("pid,iid,name"),
				ProjectId:      types.StringValue("pid"),
				InstanceId:     types.StringValue("iid"),
				Name:           types.StringValue("name"),
				MetricsPath:    types.StringNull(),
				Scheme:         types.StringValue(""),
				ScrapeInterval: types.StringNull(),
				ScrapeTimeout:  types.StringNull(),
				SAML2:          types.ObjectNull(saml2Types),
				BasicAuth:      types.ObjectNull(basicAuthTypes),
				Targets:        types.ListNull(types.ObjectType{AttrTypes: targetTypes}),
			},
			true,
		},
		{
			description: "values_ok",
			input: &observability.Job{
				JobName:     utils.Ptr("name"),
				MetricsPath: utils.Ptr("/m"),
				BasicAuth: &observability.BasicAuth{
					Password: utils.Ptr("p"),
					Username: utils.Ptr("u"),
				},
				Params:         &map[string][]string{"saml2": {"disabled"}, "x": {"y", "z"}},
				Scheme:         observability.JOBSCHEME_HTTP.Ptr(),
				ScrapeInterval: utils.Ptr("1"),
				ScrapeTimeout:  utils.Ptr("2"),
				SampleLimit:    utils.Ptr(int64(17)),
				StaticConfigs: &[]observability.StaticConfigs{
					{
						Labels:  &map[string]string{"k1": "v1"},
						Targets: &[]string{"url1"},
					},
					{
						Labels:  &map[string]string{"k2": "v2", "k3": "v3"},
						Targets: &[]string{"url1", "url3"},
					},
					{
						Labels:  nil,
						Targets: &[]string{},
					},
				},
			},
			expected: Model{
				Id:             types.StringValue("pid,iid,name"),
				ProjectId:      types.StringValue("pid"),
				InstanceId:     types.StringValue("iid"),
				Name:           types.StringValue("name"),
				MetricsPath:    types.StringValue("/m"),
				Scheme:         types.StringValue(string(observability.JOBSCHEME_HTTP)),
				ScrapeInterval: types.StringValue("1"),
				ScrapeTimeout:  types.StringValue("2"),
				SampleLimit:    types.Int64Value(17),
				SAML2: types.ObjectValueMust(saml2Types, map[string]attr.Value{
					"enable_url_parameters": types.BoolValue(false),
				}),
				BasicAuth: types.ObjectValueMust(basicAuthTypes, map[string]attr.Value{
					"username": types.StringValue("u"),
					"password": types.StringValue("p"),
				}),
				Targets: types.ListValueMust(types.ObjectType{AttrTypes: targetTypes}, []attr.Value{
					types.ObjectValueMust(targetTypes, map[string]attr.Value{
						"urls": types.ListValueMust(types.StringType, []attr.Value{types.StringValue("url1")}),
						"labels": types.MapValueMust(types.StringType, map[string]attr.Value{
							"k1": types.StringValue("v1"),
						}),
					}),
					types.ObjectValueMust(targetTypes, map[string]attr.Value{
						"urls": types.ListValueMust(types.StringType, []attr.Value{types.StringValue("url1"), types.StringValue("url3")}),
						"labels": types.MapValueMust(types.StringType, map[string]attr.Value{
							"k2": types.StringValue("v2"),
							"k3": types.StringValue("v3"),
						}),
					}),
					types.ObjectValueMust(targetTypes, map[string]attr.Value{
						"urls":   types.ListValueMust(types.StringType, []attr.Value{}),
						"labels": types.MapNull(types.StringType),
					}),
				}),
			},
			isValid: true,
		},
		{
			"response_nil_fail",
			nil,
			Model{},
			false,
		},
		{
			"no_resource_id",
			&observability.Job{},
			Model{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			state := &Model{
				ProjectId:  tt.expected.ProjectId,
				InstanceId: tt.expected.InstanceId,
			}
			err := mapFields(context.Background(), tt.input, state)
			if !tt.isValid && err == nil {
				t.Fatalf("Should have failed")
			}
			if tt.isValid && err != nil {
				t.Fatalf("Should not have failed: %v", err)
			}
			if tt.isValid {
				diff := cmp.Diff(state, &tt.expected)
				if diff != "" {
					t.Fatalf("Data does not match: %s", diff)
				}
			}
		})
	}
}

func TestToCreatePayload(t *testing.T) {
	tests := []struct {
		description    string
		input          *Model
		inputSAML2     *saml2Model
		inputBasicAuth *basicAuthModel
		inputTargets   []targetModel
		expected       *observability.CreateScrapeConfigPayload
		isValid        bool
	}{
		{
			"basic_ok",
			&Model{
				MetricsPath: types.StringValue("/metrics"),
			},
			&saml2Model{},
			&basicAuthModel{},
			[]targetModel{},
			&observability.CreateScrapeConfigPayload{
				MetricsPath: utils.Ptr("/metrics"),
				// Defaults
				Scheme:         observability.CREATESCRAPECONFIGPAYLOADSCHEME_HTTP.Ptr(),
				ScrapeInterval: utils.Ptr("5m"),
				ScrapeTimeout:  utils.Ptr("2m"),
				SampleLimit:    utils.Ptr(float64(5000)),
				StaticConfigs:  &[]observability.CreateScrapeConfigPayloadStaticConfigsInner{},
				Params:         &map[string]any{"saml2": []string{"enabled"}},
			},
			true,
		},
		{
			"ok - false enable_url_parameters",
			&Model{
				MetricsPath: types.StringValue("/metrics"),
				Name:        types.StringValue("Name"),
			},
			&saml2Model{
				EnableURLParameters: types.BoolValue(false),
			},
			&basicAuthModel{},
			[]targetModel{},
			&observability.CreateScrapeConfigPayload{
				MetricsPath: utils.Ptr("/metrics"),
				JobName:     utils.Ptr("Name"),
				Params:      &map[string]any{"saml2": []string{"disabled"}},
				// Defaults
				Scheme:         observability.CREATESCRAPECONFIGPAYLOADSCHEME_HTTP.Ptr(),
				ScrapeInterval: utils.Ptr("5m"),
				ScrapeTimeout:  utils.Ptr("2m"),
				SampleLimit:    utils.Ptr(float64(5000)),
				StaticConfigs:  &[]observability.CreateScrapeConfigPayloadStaticConfigsInner{},
			},
			true,
		},
		{
			"ok -  true enable_url_parameters",
			&Model{
				MetricsPath: types.StringValue("/metrics"),
				Name:        types.StringValue("Name"),
			},
			&saml2Model{
				EnableURLParameters: types.BoolValue(true),
			},
			&basicAuthModel{},
			[]targetModel{},
			&observability.CreateScrapeConfigPayload{
				MetricsPath: utils.Ptr("/metrics"),
				JobName:     utils.Ptr("Name"),
				Params:      &map[string]any{"saml2": []string{"enabled"}},
				// Defaults
				Scheme:         observability.CREATESCRAPECONFIGPAYLOADSCHEME_HTTP.Ptr(),
				ScrapeInterval: utils.Ptr("5m"),
				ScrapeTimeout:  utils.Ptr("2m"),
				SampleLimit:    utils.Ptr(float64(5000)),
				StaticConfigs:  &[]observability.CreateScrapeConfigPayloadStaticConfigsInner{},
			},
			true,
		},
		{
			"ok - with basic auth",
			&Model{
				MetricsPath: types.StringValue("/metrics"),
				Name:        types.StringValue("Name"),
			},
			&saml2Model{},
			&basicAuthModel{
				Username: types.StringValue("u"),
				Password: types.StringValue("p"),
			},
			[]targetModel{},
			&observability.CreateScrapeConfigPayload{
				MetricsPath: utils.Ptr("/metrics"),
				JobName:     utils.Ptr("Name"),
				BasicAuth: &observability.CreateScrapeConfigPayloadBasicAuth{
					Username: utils.Ptr("u"),
					Password: utils.Ptr("p"),
				},
				// Defaults
				Scheme:         observability.CREATESCRAPECONFIGPAYLOADSCHEME_HTTP.Ptr(),
				ScrapeInterval: utils.Ptr("5m"),
				ScrapeTimeout:  utils.Ptr("2m"),
				SampleLimit:    utils.Ptr(float64(5000)),
				StaticConfigs:  &[]observability.CreateScrapeConfigPayloadStaticConfigsInner{},
				Params:         &map[string]any{"saml2": []string{"enabled"}},
			},
			true,
		},
		{
			"ok - with targets",
			&Model{
				MetricsPath: types.StringValue("/metrics"),
				Name:        types.StringValue("Name"),
			},
			&saml2Model{},
			&basicAuthModel{},
			[]targetModel{
				{
					URLs:   types.ListValueMust(types.StringType, []attr.Value{types.StringValue("url1")}),
					Labels: types.MapValueMust(types.StringType, map[string]attr.Value{"k1": types.StringValue("v1")}),
				},
				{
					URLs:   types.ListValueMust(types.StringType, []attr.Value{types.StringValue("url1"), types.StringValue("url3")}),
					Labels: types.MapValueMust(types.StringType, map[string]attr.Value{"k2": types.StringValue("v2"), "k3": types.StringValue("v3")}),
				},
				{
					URLs:   types.ListValueMust(types.StringType, []attr.Value{}),
					Labels: types.MapNull(types.StringType),
				},
			},
			&observability.CreateScrapeConfigPayload{
				MetricsPath: utils.Ptr("/metrics"),
				JobName:     utils.Ptr("Name"),
				StaticConfigs: &[]observability.CreateScrapeConfigPayloadStaticConfigsInner{
					{
						Targets: &[]string{"url1"},
						Labels:  &map[string]interface{}{"k1": "v1"},
					},
					{
						Targets: &[]string{"url1", "url3"},
						Labels:  &map[string]interface{}{"k2": "v2", "k3": "v3"},
					},
					{
						Targets: &[]string{},
						Labels:  &map[string]interface{}{},
					},
				},
				// Defaults
				Scheme:         observability.CREATESCRAPECONFIGPAYLOADSCHEME_HTTP.Ptr(),
				ScrapeInterval: utils.Ptr("5m"),
				ScrapeTimeout:  utils.Ptr("2m"),
				SampleLimit:    utils.Ptr(float64(5000)),
				Params:         &map[string]any{"saml2": []string{"enabled"}},
			},
			true,
		},
		{
			"nil_model",
			nil,
			nil,
			nil,
			nil,
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output, err := toCreatePayload(context.Background(), tt.input, tt.inputSAML2, tt.inputBasicAuth, tt.inputTargets)
			if !tt.isValid && err == nil {
				t.Fatalf("Should have failed")
			}
			if tt.isValid && err != nil {
				t.Fatalf("Should not have failed: %v", err)
			}
			if tt.isValid {
				diff := cmp.Diff(output, tt.expected)
				if diff != "" {
					t.Fatalf("Data does not match: %s", diff)
				}
			}
		})
	}
}

func TestToUpdatePayload(t *testing.T) {
	tests := []struct {
		description    string
		input          *Model
		inputSAML2     *saml2Model
		basicAuthModel *basicAuthModel
		inputTargets   []targetModel
		expected       *observability.UpdateScrapeConfigPayload
		isValid        bool
	}{
		{
			"basic_ok",
			&Model{
				MetricsPath: types.StringValue("/metrics"),
			},
			&saml2Model{},
			&basicAuthModel{},
			[]targetModel{},
			&observability.UpdateScrapeConfigPayload{
				MetricsPath: utils.Ptr("/metrics"),
				// Defaults
				Scheme:         observability.UPDATESCRAPECONFIGPAYLOADSCHEME_HTTP.Ptr(),
				ScrapeInterval: utils.Ptr("5m"),
				ScrapeTimeout:  utils.Ptr("2m"),
				SampleLimit:    utils.Ptr(float64(5000)),
				StaticConfigs:  &[]observability.UpdateScrapeConfigPayloadStaticConfigsInner{},
			},
			true,
		},
		{
			"ok -  true enable_url_parameters",
			&Model{
				MetricsPath: types.StringValue("/metrics"),
				Scheme:      types.StringValue("http"),
			},
			&saml2Model{
				EnableURLParameters: types.BoolValue(true),
			},
			&basicAuthModel{},
			[]targetModel{},
			&observability.UpdateScrapeConfigPayload{
				MetricsPath: utils.Ptr("/metrics"),
				// Defaults
				Scheme:         observability.UPDATESCRAPECONFIGPAYLOADSCHEME_HTTP.Ptr(),
				ScrapeInterval: utils.Ptr("5m"),
				ScrapeTimeout:  utils.Ptr("2m"),
				SampleLimit:    utils.Ptr(float64(5000)),
				StaticConfigs:  &[]observability.UpdateScrapeConfigPayloadStaticConfigsInner{},
				Params:         &map[string]any{"saml2": []string{"enabled"}},
			},
			true,
		},
		{
			"ok -  false enable_url_parameters",
			&Model{
				MetricsPath: types.StringValue("/metrics"),
				Scheme:      types.StringValue("http"),
			},
			&saml2Model{
				EnableURLParameters: types.BoolValue(false),
			},
			&basicAuthModel{},
			[]targetModel{},
			&observability.UpdateScrapeConfigPayload{
				MetricsPath: utils.Ptr("/metrics"),
				// Defaults
				Scheme:         observability.UPDATESCRAPECONFIGPAYLOADSCHEME_HTTP.Ptr(),
				ScrapeInterval: utils.Ptr("5m"),
				ScrapeTimeout:  utils.Ptr("2m"),
				SampleLimit:    utils.Ptr(float64(5000)),
				StaticConfigs:  &[]observability.UpdateScrapeConfigPayloadStaticConfigsInner{},
				Params:         &map[string]any{"saml2": []string{"disabled"}},
			},
			true,
		},
		{
			"ok - with basic auth",
			&Model{
				MetricsPath: types.StringValue("/metrics"),
				Name:        types.StringValue("Name"),
			},
			&saml2Model{},
			&basicAuthModel{
				Username: types.StringValue("u"),
				Password: types.StringValue("p"),
			},
			[]targetModel{},
			&observability.UpdateScrapeConfigPayload{
				MetricsPath: utils.Ptr("/metrics"),
				BasicAuth: &observability.CreateScrapeConfigPayloadBasicAuth{
					Username: utils.Ptr("u"),
					Password: utils.Ptr("p"),
				},
				// Defaults
				Scheme:         observability.UPDATESCRAPECONFIGPAYLOADSCHEME_HTTP.Ptr(),
				ScrapeInterval: utils.Ptr("5m"),
				ScrapeTimeout:  utils.Ptr("2m"),
				SampleLimit:    utils.Ptr(float64(5000)),
				StaticConfigs:  &[]observability.UpdateScrapeConfigPayloadStaticConfigsInner{},
			},
			true,
		},
		{
			"ok - with targets",
			&Model{
				MetricsPath: types.StringValue("/metrics"),
				Name:        types.StringValue("Name"),
			},
			&saml2Model{},
			&basicAuthModel{},
			[]targetModel{
				{
					URLs:   types.ListValueMust(types.StringType, []attr.Value{types.StringValue("url1")}),
					Labels: types.MapValueMust(types.StringType, map[string]attr.Value{"k1": types.StringValue("v1")}),
				},
				{
					URLs:   types.ListValueMust(types.StringType, []attr.Value{types.StringValue("url1"), types.StringValue("url3")}),
					Labels: types.MapValueMust(types.StringType, map[string]attr.Value{"k2": types.StringValue("v2"), "k3": types.StringValue("v3")}),
				},
				{
					URLs:   types.ListValueMust(types.StringType, []attr.Value{}),
					Labels: types.MapNull(types.StringType),
				},
			},
			&observability.UpdateScrapeConfigPayload{
				MetricsPath: utils.Ptr("/metrics"),
				StaticConfigs: &[]observability.UpdateScrapeConfigPayloadStaticConfigsInner{
					{
						Targets: &[]string{"url1"},
						Labels:  &map[string]interface{}{"k1": "v1"},
					},
					{
						Targets: &[]string{"url1", "url3"},
						Labels:  &map[string]interface{}{"k2": "v2", "k3": "v3"},
					},
					{
						Targets: &[]string{},
						Labels:  &map[string]interface{}{},
					},
				},
				// Defaults
				Scheme:         observability.UPDATESCRAPECONFIGPAYLOADSCHEME_HTTP.Ptr(),
				ScrapeInterval: utils.Ptr("5m"),
				ScrapeTimeout:  utils.Ptr("2m"),
				SampleLimit:    utils.Ptr(float64(5000)),
			},
			true,
		},
		{
			"nil_model",
			nil,
			nil,
			nil,
			nil,
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output, err := toUpdatePayload(context.Background(), tt.input, tt.inputSAML2, tt.basicAuthModel, tt.inputTargets)
			if !tt.isValid && err == nil {
				t.Fatalf("Should have failed")
			}
			if tt.isValid && err != nil {
				t.Fatalf("Should not have failed: %v", err)
			}
			if tt.isValid {
				diff := cmp.Diff(output, tt.expected)
				if diff != "" {
					t.Fatalf("Data does not match: %s", diff)
				}
			}
		})
	}
}
