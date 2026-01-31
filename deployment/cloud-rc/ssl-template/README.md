# SSL Certificate Directory

This directory is a placeholder for SAML Service Provider certificates.

## Certificates Generated on Server

When you run `scripts/generate-saml-certs.sh` on the cloud.rc instance, it creates:

- `sp.key` - Private key (600 permissions)
- `sp.crt` - Public certificate (644 permissions)

These files are created in `/var/fider/ssl/` on the server.

## Security

**NEVER** commit private keys (`.key` files) to version control.

The `.gitignore` in the parent directory excludes:
- `ssl/*.key`
- `ssl/*.crt`
- `ssl/*.pem`

## Usage

After generating certificates on the server:

1. The certificates are automatically mounted to the Docker container via docker-compose.yml
2. Fider reads them from the `SAML_SP_CERT_PATH` and `SAML_SP_KEY_PATH` environment variables
3. You need to send `sp.crt` (public certificate only) to UAB IT for SAML IdP registration
