# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### About This Fork

This is a maintained fork of the original [brainly/terraform-provider-redshift](https://github.com/brainly/terraform-provider-redshift) provider, which has been deprecated.

**Key Features:**
- ✅ AWS Redshift Serverless support
- ✅ Backward compatibility with regular Redshift clusters
- ✅ Active maintenance and updates
- ✅ All original features: users, groups, schemas, grants, datashares, and more

### Major Enhancements (vs. Original)

**Redshift Serverless Support:**
- Added `type = "serverless"` parameter to provider configuration
- Use `pg_user_info` instead of `svl_user_info` for Serverless user operations
- Handle Serverless-specific column differences (missing `sessiontimeout`, `syslogaccess`)
- Fix permission denied errors when reading user information in Serverless

**Other Improvements:**
- Updated dependencies and Go SDK versions
- Improved documentation and examples
- Better error handling and validation

---

## Historical Changes (Inherited from Original Provider)

For historical changes from the original brainly/terraform-provider-redshift provider, see:
https://github.com/brainly/terraform-provider-redshift/blob/master/CHANGELOG.md

**Notable features from original provider:**
- GRANT permissions (tables, schemas, databases, functions, procedures, languages)
- Default privileges management
- User and group management with full attribute support
- Datashare and datashare privilege support
- Cross-account temporary credentials
- Support for quoted identifiers and case-sensitive names
