## 2.2.0 (Unreleased)
## 2.1.1 (September 25, 2019)

NOTES:

* The provider has switched to the standalone TF SDK, there should be no noticeable impact on compatibility. ([#54](https://github.com/terraform-providers/terraform-provider-tls/issues/54))

## 2.1.0 (August 16, 2019)

ENHANCEMENTS:

* Certificate renewal is now handled as a "replace" action in the plan, rather than by behaving as if the expired certificate had been deleted. Although the effective behavior remains unchanged, renewal will now appear as a `-/+` action in the plan, rather than just as a `+`. ([#34](https://github.com/terraform-providers/terraform-provider-tls/issues/34))
* Certificates can now have URIs as subject alternative names. ([#50](https://github.com/terraform-providers/terraform-provider-tls/issues/50))
* Certificates can now optionally have the Subject Key ID field populated. ([#31](https://github.com/terraform-providers/terraform-provider-tls/issues/31))

BUG FIXES:

* More of the private key arguments are now marked as "sensitive" so that Terraform will know to hide their values when showing plans and state in response to various commands. ([#48](https://github.com/terraform-providers/terraform-provider-tls/issues/48))
* In `tls_public_key`, don't panic if the PEM isn't valid PEM syntax at all. ([#40](https://github.com/terraform-providers/terraform-provider-tls/issues/40))

## 2.0.1 (April 30, 2019)

* This release includes an upgraded Terraform SDK, for the sake of aligning versions of the SDK amongst released providers, as we lead up to Core v0.12. This should have no noticeable impact on the provider.

## 2.0.0 (April 17, 2019)

IMPROVEMENTS:

* The provider is now compatible with Terraform v0.12, while retaining compatibility with prior versions.

## 1.2.0 (August 15, 2018)

FEATURES: 

* `tls_private_key` (both datasource and resource) include MD5 public key fingerprints as computed attributes.


BUG FIXES:
* `tls_cert_request` and `tls_self_signed_cert`: changes to `subject` now
  correctly force the recreation of the resource, instead of returning an error
  ([#18](https://github.com/terraform-providers/terraform-provider-tls/issues/18))

## 1.1.0 (March 09, 2018)

FEATURES:

* **New Data Source:** `tls_public_key`
  ([#11](https://github.com/terraform-providers/terraform-provider-tls/issues/11))

## 1.0.1 (November 09, 2017)

BUG FIXES:

* `tls_cert_request` and `tls_self_signed_cert` no longer cause a crash when
  `subject` isn't specified.
  ([#7](https://github.com/terraform-providers/terraform-provider-tls/issues/7))
* `tls_cert_request` and `tls_self_signed_cert` no longer generate empty-string
  values for various subject fields when they are not set in configuration.
  ([#10](https://github.com/terraform-providers/terraform-provider-tls/issues/10))

## 1.0.0 (September 15, 2017)

* No changes from 0.1.0; just adjusting to [the new version numbering
  scheme](https://www.hashicorp.com/blog/hashicorp-terraform-provider-versioning/).

## 0.1.0 (June 21, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider
  Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
