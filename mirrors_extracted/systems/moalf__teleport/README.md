# Teleport Challenge

This project implements workflow automation to securely configure, add web applications, and users to an [Auth0](https://auth0.com) tenant. See our [System Design](/docs/design.md) documentation.

## Tenant Configuration

The implementation resides in `auth0/tenant`. Its associated workflow is: `.github/workflows/tenant-configuration.yml`.

### Secrets

The tenant configuration workflow requires these values to be added as GitHub secrets:

- `AUTH0_DOMAIN`: Auth0 tenant domain
- `AUTH0_CLIENT_ID`: Client ID of an M2M application
- `AUTH0_CLIENT_SECRET`: Client Secret of the same M2M application

### Terraform Variables

For local testing with `terraform.tfvars` use:

```sh
auth0_domain        = "<YOUR_AUTH0_TENANT_DOMAIN>"
auth0_client_id     = "<YOUR_AUTH0_M2M_CLIENT_ID>"
auth0_client_secret = "<YOUR_AUTH0_M2M_CLIENT_SECRET>"
```