# terraform-provider-bouncr

Terraform Provider for bouncr.

```
resource "bouncr_application" "foo" {
  name = "foo"
  url  = "foo"

  realm {
    name    = "foo_realm"
    pass_to = "http://foo.internal:3000"
  }
}

resource "bouncr_role" "foo_admin" {
  name = "foo-admin"
  permissions = [
    "${bouncr_permission.foo_create}",
    "${bouncr_permission.foo_read}",
  ]
}

resource "bouncr_role" "foo_user" {
  name = "foo-user"
  permissions = [
    "${bouncr_permission.foo_read}",
  ]
}

resource "bouncr_permission" "foo_create" {
  name = "foo:create"
}

resource "bouncr_permission" "foo_read" {
  name = "foo:read"
}

resource "bouncr_group" "foo_admins" {
  name        = "foo_admins"
  description = "The administrators for Foo"

  members = [
    ${data.bouncr_user.admin}
  ]
}

data "bouncr_user" "admin" {
  account = "admin"
}

resource "bouncr_assignments" "assign" {
  assignment {
    group = "${bouncr_group.foo_admins.id}"
    role  = "${bouncr_role.foo_admin.id}"
    realm = "${bouncr_application.realm[0].id}"
  }
}
```
