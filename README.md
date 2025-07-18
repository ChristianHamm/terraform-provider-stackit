<div align="center">
<br>
<img src=".github/images/stackit-logo.svg" alt="STACKIT logo" width="50%"/>
<br>
<br>
</div>

# STACKIT Terraform Provider

[![Go Report Card](https://goreportcard.com/badge/github.com/stackitcloud/terraform-provider-stackit)](https://goreportcard.com/report/github.com/stackitcloud/terraform-provider-stackit) [![GitHub Release](https://img.shields.io/github/v/release/stackitcloud/terraform-provider-stackit)](https://registry.terraform.io/providers/stackitcloud/stackit/latest) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/stackitcloud/terraform-provider-stackit) [![GitHub License](https://img.shields.io/github/license/stackitcloud/terraform-provider-stackit)](https://www.apache.org/licenses/LICENSE-2.0)

This project is the official [Terraform Provider](https://registry.terraform.io/providers/stackitcloud/stackit/latest/docs) for [STACKIT](https://www.stackit.de/en/), which allows you to manage STACKIT resources through Terraform.

## Getting Started

To install the [STACKIT Terraform Provider](https://registry.terraform.io/providers/stackitcloud/stackit/latest/docs), copy and paste this code into your Terraform configuration. Then, run `terraform init`.

```hcl
terraform {
  required_providers {
    stackit = {
      source = "stackitcloud/stackit"
      version = "X.X.X"
    }
  }
}

provider "stackit" {
  # Configuration options
}
```

Check one of the examples in the [examples](examples/) folder.

## Authentication

To authenticate, you will need a [service account](https://docs.stackit.cloud/stackit/en/service-accounts-134415819.html). Create it in the [STACKIT Portal](https://portal.stackit.cloud/) and assign the necessary permissions to it, e.g. `project.owner`. There are multiple ways to authenticate:

- Key flow (recommended)
- Token flow

When setting up authentication, the provider will always try to use the key flow first and search for credentials in several locations, following a specific order:

1. Explicit configuration, e.g. by setting the field `service_account_key_path` in the provider block (see example below)
2. Environment variable, e.g. by setting `STACKIT_SERVICE_ACCOUNT_KEY_PATH` or `STACKIT_SERVICE_ACCOUNT_KEY`
3. Credentials file

   The provider will check the credentials file located in the path defined by the `STACKIT_CREDENTIALS_PATH` env var, if specified,
   or in `$HOME/.stackit/credentials.json` as a fallback.
   The credentials file should be a JSON and each credential should be set using the name of the respective environment variable, as stated below in each flow. Example:

   ```json
   {
     "STACKIT_SERVICE_ACCOUNT_TOKEN": "foo_token",
     "STACKIT_SERVICE_ACCOUNT_KEY_PATH": "path/to/sa_key.json"
   }
   ```

### Key flow

    The following instructions assume that you have created a service account and assigned the necessary permissions to it, e.g. `project.owner`.

To use the key flow, you need to have a service account key, which must have an RSA key-pair attached to it.

When creating the service account key, a new pair can be created automatically, which will be included in the service account key. This will make it much easier to configure the key flow authentication in the [STACKIT Terraform Provider](https://github.com/stackitcloud/terraform-provider-stackit), by just providing the service account key.

**Optionally**, you can provide your own private key when creating the service account key, which will then require you to also provide it explicitly to the [STACKIT Terraform Provider](https://github.com/stackitcloud/terraform-provider-stackit), additionally to the service account key. Check the STACKIT Knowledge Base for an [example of how to create your own key-pair](https://docs.stackit.cloud/stackit/en/usage-of-the-service-account-keys-in-stackit-175112464.html#UsageoftheserviceaccountkeysinSTACKIT-CreatinganRSAkey-pair).

To configure the key flow, follow this steps:

1.  Create a service account key:

- Use the [STACKIT Portal](https://portal.stackit.cloud/): go to the `Service Accounts` tab, choose a `Service Account` and go to `Service Account Keys` to create a key. For more details, see [Create a service account key](https://docs.stackit.cloud/stackit/en/create-a-service-account-key-175112456.html)

2.  Save the content of the service account key by copying it and saving it in a JSON file.

    The expected format of the service account key is a **JSON** with the following structure:

```json
{
  "id": "uuid",
  "publicKey": "public key",
  "createdAt": "2023-08-24T14:15:22Z",
  "validUntil": "2023-08-24T14:15:22Z",
  "keyType": "USER_MANAGED",
  "keyOrigin": "USER_PROVIDED",
  "keyAlgorithm": "RSA_2048",
  "active": true,
  "credentials": {
    "kid": "string",
    "iss": "my-sa@sa.stackit.cloud",
    "sub": "uuid",
    "aud": "string",
    (optional) "privateKey": "private key when generated by the SA service"
  }
}
```

3. Configure the service account key for authentication in the provider by following one of the alternatives below:

   - setting the fields in the provider block: `service_account_key` or `service_account_key_path`
   - setting the environment variable: `STACKIT_SERVICE_ACCOUNT_KEY_PATH` or `STACKIT_SERVICE_ACCOUNT_KEY`
     - ensure the set the service account key in `STACKIT_SERVICE_ACCOUNT_KEY` is correctly formatted. Use e.g.
       `$ export STACKIT_SERVICE_ACCOUNT_KEY=$(cat ./service-account-key.json)`
   - setting `STACKIT_SERVICE_ACCOUNT_KEY_PATH` in the credentials file (see above)

> **Optionally, only if you have provided your own RSA key-pair when creating the service account key**, you also need to configure your private key (takes precedence over the one included in the service account key, if present). **The private key must be PEM encoded** and can be provided using one of the options below:
>
> - setting the field in the provider block: `private_key` or `private_key_path`
> - setting the environment variable: `STACKIT_PRIVATE_KEY_PATH` or `STACKIT_PRIVATE_KEY`
> - setting `STACKIT_PRIVATE_KEY_PATH` in the credentials file (see above)

### Token flow

Using this flow is less secure since the token is long-lived. You can provide the token in several ways:

1. Setting the field `service_account_token` in the provider
2. Setting the environment variable `STACKIT_SERVICE_ACCOUNT_TOKEN`
3. Setting it in the credentials file (see above)

## Backend configuration

To keep track of your Terraform state, you can configure an [S3 backend](https://developer.hashicorp.com/terraform/language/settings/backends/s3) using [STACKIT Object Storage](https://docs.stackit.cloud/stackit/en/object-storage-s3-compatible-71009778.html).

To do so, you need an Object Storage [S3 bucket](https://docs.stackit.cloud/stackit/en/basic-concept-objectstorage-71009785.html#BasicConceptObjectStorage-Buckets) and [credentials](https://docs.stackit.cloud/stackit/en/basic-concept-objectstorage-71009785.html#BasicConceptObjectStorage-Credentials) to access it. If you need to create them, check [Getting Started - Object Storage](https://docs.stackit.cloud/stackit/en/getting-started-objectstorage-71009792.html).

Once you have everything setup, you can configure the backend by adding the following block to your Terraform configuration:

```hcl
terraform {
  backend "s3" {
    bucket = "BUCKET_NAME"
    key    = "path/to/key"
    endpoints = {
      s3 = "https://object.storage.eu01.onstackit.cloud"
    }
    region                      = "eu01"
    skip_credentials_validation = true
    skip_region_validation      = true
    skip_s3_checksum            = true
    skip_requesting_account_id  = true
    secret_key                  = "SECRET_KEY"
    access_key                  = "ACCESS_KEY"
  }
}
```

Note: AWS specific checks must be skipped as they do not work on STACKIT. For details on what those validations do, see [here](https://developer.hashicorp.com/terraform/language/settings/backends/s3#configuration).

## Opting into Beta Resources

To use beta resources in the STACKIT Terraform provider, follow these steps:

1. **Provider Configuration Option**

   Set the `enable_beta_resources` option in the provider configuration. This is a boolean attribute that can be either `true` or `false`.

   ```hcl
   provider "stackit" {
     default_region        = "eu01"
     enable_beta_resources = true
   }
   ```

2. **Environment Variable**

   Set the `STACKIT_TF_ENABLE_BETA_RESOURCES` environment variable to `"true"` or `"false"`. Other values will be ignored and will produce a warning.

   ```sh
   export STACKIT_TF_ENABLE_BETA_RESOURCES=true
   ```

> **Note**: The environment variable takes precedence over the provider configuration option. This means that if the `STACKIT_TF_ENABLE_BETA_RESOURCES` environment variable is set to a valid value (`"true"` or `"false"`), it will override the `enable_beta_resources` option specified in the provider configuration.

For more details, please refer to the [beta resources configuration guide](https://registry.terraform.io/providers/stackitcloud/stackit/latest/docs/guides/opting_into_beta_resources).

## Opting into Experiments

Experiments are features that are even less mature and stable than Beta Resources. While there is some assumed stability in beta resources, will have to expect breaking changes while using experimental resources. Experimental Resources do not come with any support or warranty.

To enable experiments set the experiments field in the provider definition:

```hcl
provider "stackit" {
  default_region        = "eu01"
  experiments           = ["iam", "routing-tables", "network"]
}
```

### Available Experiments

#### `iam`

Enables IAM management features in the Terraform provider. The underlying IAM API is expected to undergo a redesign in the future, which leads to it being considered experimental.

#### `routing-tables`

This feature enables experimental routing table capabilities in the Terraform Provider, available only to designated SNAs at this time.

#### `network`

The `stackit_network` provides the fields `region` and `routing_table_id` when the experiment flag `network` is set. 
The underlying API is not stable yet and could change in the future.  
If you don't need these fields, don't set the experiment flag `network`, to use the stable api.

## Acceptance Tests

Terraform acceptance tests are run using the command `make test-acceptance-tf`. For all services,

- The env var `TF_ACC_PROJECT_ID` must be set with the ID of the STACKIT test project to test it.
- Authentication is set as usual.
- Optionally, the env var `TF_ACC_XXXXXX_CUSTOM_ENDPOINT` (where `XXXXXX` is the uppercase name of the service) can be set to use endpoints other than the default value.
- There are some acceptance test where it is needed to provide additional parameters, some of them have default values in order to run normally without manual interaction. Those default values can be overwritten (see testutils.go for a full list.)

Additionally:

- For the Resource Manager service:
  - A service account with permissions to create and delete projects is required
  - The env var `TF_ACC_TEST_PROJECT_SERVICE_ACCOUNT_EMAIL` must be set as the email of the service account
  - The env var `TF_ACC_TEST_PROJECT_SERVICE_ACCOUNT_TOKEN` must be set as a valid token of the service account. Can also be set in the credentials file used by authentication (see [Authentication](#authentication) for more details)
  - The env var `TF_ACC_PROJECT_ID` is ignored

**WARNING:** Acceptance tests will create real resources, which may incur in costs.

## Migration

For guidance on how to migrate to using this provider, please see our [Migration Guide](./MIGRATION.md).

## Reporting Issues

If you encounter any issues or have suggestions for improvements, please open an issue in the [repository](https://github.com/stackitcloud/terraform-provider-stackit/issues).

## Contribute

Your contribution is welcome! For more details on how to contribute, refer to our [Contribution Guide](./CONTRIBUTION.md).

## Release creation

See the [release documentation](./RELEASE.md) for further information.

## License

Apache 2.0

## Useful Links

- [STACKIT Terraform Provider](https://registry.terraform.io/providers/stackitcloud/stackit/latest/docs)

- [STACKIT Portal](https://portal.stackit.cloud/)

- [STACKIT](https://www.stackit.de/en/)

- [STACKIT Knowledge Base](https://docs.stackit.cloud/stackit/en/knowledge-base-85301704.html)

- [STACKIT CLI](https://github.com/stackitcloud/stackit-cli/tree/main)
