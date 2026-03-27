# terraform-provider-bouncr

Terraform Provider for [Bouncr](https://github.com/kawasima/bouncr) - an identity and access management system.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23 (to build the provider)

## Building the Provider

```sh
git clone https://github.com/kawasima/terraform-provider-bouncr.git
cd terraform-provider-bouncr
make build
```

## Installing the Provider

```sh
make install
```

This installs the provider to `~/.terraform.d/plugins/registry.terraform.io/kawasima/bouncr/<VERSION>/<OS>_<ARCH>/`.

To use the locally installed provider, create `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "kawasima/bouncr" = "<HOME>/.terraform.d/plugins/registry.terraform.io/kawasima/bouncr/0.2.0/<OS>_<ARCH>"
  }
  direct {}
}
```

## Provider Configuration

```hcl
terraform {
  required_providers {
    bouncr = {
      source = "kawasima/bouncr"
    }
  }
}

provider "bouncr" {
  client_id     = var.bouncr_client_id
  client_secret = var.bouncr_client_secret
  base_url      = "http://localhost:8080"
}
```

### Arguments

| Name | Description | Default | Environment Variable |
|------|-------------|---------|---------------------|
| `client_id` | OAuth2 Client ID | - | `BOUNCR_CLIENT_ID` |
| `client_secret` | OAuth2 Client Secret | - | `BOUNCR_CLIENT_SECRET` |
| `base_url` | Bouncr API base URL | `http://localhost:3000` | `BOUNCR_URL` |

## Resources

### bouncr_permission

```hcl
resource "bouncr_permission" "user_read" {
  name        = "user:read"
  description = "Read user information"
}
```

#### Arguments

| Name | Required | Description |
|------|----------|-------------|
| `name` | Yes | Permission name |
| `description` | Yes | Permission description |

### bouncr_role

```hcl
resource "bouncr_role" "editor" {
  name        = "editor"
  description = "Editor role"
  permissions = [
    bouncr_permission.user_read.name,
    bouncr_permission.user_write.name,
  ]
}
```

#### Arguments

| Name | Required | Description |
|------|----------|-------------|
| `name` | Yes | Role name |
| `description` | Yes | Role description |
| `permissions` | No | Set of permission names to attach |

### bouncr_user

```hcl
resource "bouncr_user" "alice" {
  account  = "alice"
  password = "P@ssw0rd!"
  user_profiles = {
    name  = "Alice"
    email = "alice@example.com"
  }
}
```

#### Arguments

| Name | Required | Description |
|------|----------|-------------|
| `account` | Yes | User account name |
| `password` | No | Password (sensitive) |
| `user_profiles` | No | Map of profile attributes |

### bouncr_application

```hcl
resource "bouncr_application" "myapp" {
  name         = "myapp"
  description  = "My Application"
  pass_to      = "http://localhost:8080"
  virtual_path = "/myapp"
  top_page     = "/myapp/"

  realm {
    name        = "default"
    description = "Default realm"
    url         = "/myapp/.*"
  }
}
```

#### Arguments

| Name | Required | Description |
|------|----------|-------------|
| `name` | Yes | Application name |
| `description` | Yes | Application description |
| `pass_to` | Yes | Backend URL to proxy to |
| `virtual_path` | Yes | Virtual path for the application |
| `top_page` | Yes | Top page path |
| `realm` | No | Realm block (repeatable) |

#### Realm Block

| Name | Required | Description |
|------|----------|-------------|
| `name` | Yes | Realm name |
| `description` | Yes | Realm description |
| `url` | Yes | URL pattern for the realm |

### bouncr_group

```hcl
resource "bouncr_group" "developers" {
  name        = "developers"
  description = "Development team"
  members = [
    bouncr_user.alice.account,
    bouncr_user.bob.account,
  ]
}
```

#### Arguments

| Name | Required | Description |
|------|----------|-------------|
| `name` | Yes | Group name |
| `description` | Yes | Group description |
| `members` | No | Set of user account names |

### bouncr_assignments

```hcl
resource "bouncr_assignments" "main" {
  assignment {
    group = bouncr_group.developers.name
    role  = bouncr_role.editor.name
    realm = "default"
  }
}
```

#### Arguments

| Name | Required | Description |
|------|----------|-------------|
| `assignment` | Yes | Assignment block (repeatable) |

#### Assignment Block

| Name | Required | Description |
|------|----------|-------------|
| `group` | Yes | Group name |
| `role` | Yes | Role name |
| `realm` | Yes | Realm name |

### bouncr_oidc_provider

```hcl
resource "bouncr_oidc_provider" "google" {
  name                       = "google"
  client_id                  = "google-client-id"
  client_secret              = "google-client-secret"
  scope                      = "openid email profile"
  response_type              = "code"
  authorization_endpoint     = "https://accounts.google.com/o/oauth2/v2/auth"
  token_endpoint             = "https://oauth2.googleapis.com/token"
  token_endpoint_auth_method = "client_secret_basic"
  redirect_uri               = "http://localhost:3000/oauth2/callback"
}
```

#### Arguments

| Name | Required | Description |
|------|----------|-------------|
| `name` | Yes | Provider name |
| `client_id` | Yes | OIDC Client ID |
| `client_secret` | Yes | OIDC Client Secret (sensitive) |
| `scope` | Yes | OAuth2 scopes |
| `response_type` | Yes | OAuth2 response type |
| `authorization_endpoint` | Yes | Authorization endpoint URL |
| `token_endpoint` | No | Token endpoint URL |
| `token_endpoint_auth_method` | Yes | Token endpoint auth method |
| `redirect_uri` | Yes | Redirect URI |
| `pkce_enabled` | No | Enable PKCE |
| `jwks_uri` | No | JWKS endpoint URI |
| `issuer` | No | Issuer identifier |

### bouncr_oidc_application

```hcl
resource "bouncr_oidc_application" "spa" {
  name         = "myspa"
  description  = "Single Page Application"
  home_uri     = "https://spa.example.com"
  callback_uri = "https://spa.example.com/callback"
  grant_types  = ["authorization_code"]
}
```

#### Arguments

| Name | Required | Description |
|------|----------|-------------|
| `name` | Yes | Application name |
| `description` | Yes | Application description |
| `home_uri` | No | Home URI |
| `callback_uri` | No | Callback URI |
| `grant_types` | Yes | Set of OAuth2 grant types |
| `permissions` | No | Set of permission names |

## Import

All resources support `terraform import`:

```sh
terraform import bouncr_permission.example "permission-name"
terraform import bouncr_role.example "role-name"
terraform import bouncr_user.example "account-name"
terraform import bouncr_application.example "application-name"
terraform import bouncr_group.example "group-name"
terraform import bouncr_oidc_provider.example "provider-name"
terraform import bouncr_oidc_application.example "application-name"
```

## Development

### Running Tests

```sh
make test
```

### Acceptance Tests

Requires a running Bouncr server:

```sh
make testacc
```

### Formatting

```sh
make fmt
```

## License

Apache License 2.0
