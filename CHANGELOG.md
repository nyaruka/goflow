v0.10.0
----------
 * Empty/blank values should clear name/fields/language/timezone
 * Split SetContactAProperty into new actions for name, language and timezone
 * Improve error messages from struct validation
 * Result input should be nullable if there is no input
 * Timeout rule should use timeout as value

v0.9.10
----------
 * Result input should be nullable if there is no input
 * Fix HasWaitTimedOut test
 * Waits should retain their timeout value

v0.9.9
----------
 * Use fuzzy number parsing for tests, stricter for type conversion
 * Migrate endpoint of flowserver can take include_ui param
 * Add migration of notes
 * Include contact name on migration of contact references

v0.9.8
----------
 * Fix goreleasing from travis

v0.9.7
----------
 * Migrate contact.groups to be a CSV list
 * Support channels with UUIDs which aren't UUID4

v0.9.6
----------
 * Migrate templates in legacy HasState and HasWard tests
 * Fix migration of has_email tests
 * Add JSON util functions for marshalling without HTML escaping
 * Add support for searching for locations by path

v0.9.5
----------
 * Router test docstring examples should always include use of .match
 * Use left,top instead of x,y
