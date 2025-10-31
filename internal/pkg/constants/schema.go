package constants

// Schema constants
const (
	YAMLTypeString  = "string"
	YAMLTypeBoolean = "boolean"
	YAMLTypeArray   = "array"
	YAMLTypeObject  = "object"

	SectionStack         = "stack"
	SectionProject       = "project"
	SectionValidation    = "validation"
	SectionAdvanced      = "advanced"
	SectionVersionConfig = "version_config"
	SectionServiceConfig = "service_configuration"

	PropertyEnabled = "enabled"
	PropertyName    = "name"

	TemplateProjectName             = "{{project_name}}"
	TemplateOttoVersion             = "{{otto_version}}"
	TemplateConfigDocsURL           = "{{config_docs_url}}"
	TemplateServiceConfigURL        = "{{service_config_url}}"
	TemplateDefaultSkipWarnings     = "{{default_skip_warnings}}"
	TemplateDefaultAllowMultipleDBs = "{{default_allow_multiple_dbs}}"
	TemplateDefaultAutoStart        = "{{default_auto_start}}"
	TemplateDefaultPullLatest       = "{{default_pull_latest_images}}"
	TemplateDefaultCleanupRecreate  = "{{default_cleanup_on_recreate}}"

	YAMLIndent   = "  "
	YAMLComment  = "# %s\n"
	YAMLProperty = "%s%s: "
	YAMLSection  = "%s:\n"
	YAMLListItem = "    - %s\n"
	YAMLEnabled  = "  enabled:\n"
)
