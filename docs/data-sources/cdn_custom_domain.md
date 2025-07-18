---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "stackit_cdn_custom_domain Data Source - stackit"
subcategory: ""
description: |-
  CDN distribution data source schema.
  ~> This datasource is in beta and may be subject to breaking changes in the future. Use with caution. See our guide https://registry.terraform.io/providers/stackitcloud/stackit/latest/docs/guides/opting_into_beta_resources for how to opt-in to use beta resources.
---

# stackit_cdn_custom_domain (Data Source)

CDN distribution data source schema.

~> This datasource is in beta and may be subject to breaking changes in the future. Use with caution. See our [guide](https://registry.terraform.io/providers/stackitcloud/stackit/latest/docs/guides/opting_into_beta_resources) for how to opt-in to use beta resources.

## Example Usage

```terraform
data "stackit_cdn_custom_domain" "example" {
  project_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  distribution_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name            = "https://xxx.xxx"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `distribution_id` (String) CDN distribution ID
- `name` (String)
- `project_id` (String) STACKIT project ID associated with the distribution

### Read-Only

- `errors` (List of String) List of distribution errors
- `id` (String) Terraform's internal resource identifier. It is structured as "`project_id`,`distribution_id`".
- `status` (String) Status of the distribution
