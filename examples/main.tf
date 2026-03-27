terraform {
  required_providers {
    bouncr = {
      source = "kawasima/bouncr"
    }
  }
}

variable "bouncr_client_id" {
  type        = string
  description = "Bouncr OAuth2 Client ID"
}

variable "bouncr_client_secret" {
  type        = string
  sensitive   = true
  description = "Bouncr OAuth2 Client Secret"
}

variable "bouncr_base_url" {
  type        = string
  default     = "http://localhost:3000"
  description = "Bouncr base URL"
}

provider "bouncr" {
  client_id     = var.bouncr_client_id
  client_secret = var.bouncr_client_secret
  base_url      = var.bouncr_base_url
}

# ----------------------------------------------------------------
# Permissions
# ----------------------------------------------------------------
resource "bouncr_permission" "user_read" {
  name        = "user:read"
  description = "Read user information"
}

resource "bouncr_permission" "user_write" {
  name        = "user:write"
  description = "Write user information"
}

resource "bouncr_permission" "admin" {
  name        = "admin"
  description = "Full admin access"
}

# ----------------------------------------------------------------
# Roles
# ----------------------------------------------------------------
resource "bouncr_role" "viewer" {
  name        = "viewer"
  description = "Read-only role"
  permissions = [
    bouncr_permission.user_read.name,
  ]
}

resource "bouncr_role" "editor" {
  name        = "editor"
  description = "Editor role"
  permissions = [
    bouncr_permission.user_read.name,
    bouncr_permission.user_write.name,
  ]
}

resource "bouncr_role" "administrator" {
  name        = "administrator"
  description = "Administrator role"
  permissions = [
    bouncr_permission.admin.name,
  ]
}

# ----------------------------------------------------------------
# Users
# ----------------------------------------------------------------
resource "bouncr_user" "alice" {
  account  = "alice"
  password = "P@ssw0rd!"
  user_profiles = {
    name  = "Alice"
    email = "alice@example.com"
  }
}

resource "bouncr_user" "bob" {
  account = "bob"
  user_profiles = {
    name  = "Bob"
    email = "bob@example.com"
  }
}

# ----------------------------------------------------------------
# Application with Realms
# ----------------------------------------------------------------
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

  realm {
    name        = "admin"
    description = "Admin realm"
    url         = "/myapp/admin/.*"
  }
}

# ----------------------------------------------------------------
# Groups
# ----------------------------------------------------------------
resource "bouncr_group" "developers" {
  name        = "developers"
  description = "Development team"
  members = [
    bouncr_user.alice.account,
    bouncr_user.bob.account,
  ]
}

resource "bouncr_group" "admins" {
  name        = "admins"
  description = "Administrators"
  members = [
    bouncr_user.alice.account,
  ]
}

# ----------------------------------------------------------------
# Assignments (Group + Role + Realm)
# ----------------------------------------------------------------
resource "bouncr_assignments" "main" {
  assignment {
    group = bouncr_group.developers.name
    role  = bouncr_role.editor.name
    realm = "default"
  }

  assignment {
    group = bouncr_group.admins.name
    role  = bouncr_role.administrator.name
    realm = "admin"
  }
}

# ----------------------------------------------------------------
# OIDC Provider
# ----------------------------------------------------------------
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

# ----------------------------------------------------------------
# OIDC Application
# ----------------------------------------------------------------
resource "bouncr_oidc_application" "spa" {
  name         = "myspa"
  description  = "Single Page Application"
  home_uri     = "https://spa.example.com"
  callback_uri = "https://spa.example.com/callback"
  grant_types  = ["authorization_code"]
}
