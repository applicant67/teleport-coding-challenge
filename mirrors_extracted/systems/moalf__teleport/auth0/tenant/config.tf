# Auth0 Tenant Configuration

# General Settings
resource "auth0_tenant" "main" {
  friendly_name         = "Alfredo's Auth0"
  idle_session_lifetime = 2
  session_lifetime      = 2

  customize_mfa_in_postlogin_action = true
}

# Attack Protection Settings
resource "auth0_attack_protection" "main" {
  breached_password_detection {
    enabled                      = true
    method                       = "standard"
    shields                      = ["block", "user_notification", "admin_notification"]
    admin_notification_frequency = ["immediately"]
  }

  brute_force_protection {
    enabled      = true
    max_attempts = 5
    shields      = ["block"]
  }

  suspicious_ip_throttling {
    enabled = true
    shields = ["block", "admin_notification"]

    pre_login {
      max_attempts = 50
    }

    pre_user_registration {
      max_attempts = 50
    }
  }
}

# Configure MFA (WebAuthn Biometrics + recovery codes)
resource "auth0_guardian" "mfa" {
  policy = "all-applications"

  webauthn_platform {
    enabled = true
  }
}
