v0.10.16
----------
 * Don't log 400 responses to sentry

v0.10.15
----------
 * Migrate expressions in webhook header values
 * Allow dynamic searches to query language and created_on
 * Expose contact.created_on in expressions
 * Migrate rulesets where there is an explicit Other category

v0.10.14
----------
 * Fix @contact.id not being migrated in legacy flows

v0.10.13
----------
 * Migrate @contact.<scheme> expressions to @(format_urn(contact.urns.<scheme>)) so there's no error if such a URN doesn't exist
 * Fix not saving a result when router takes default exit

v0.10.12
----------
 * Merge pull request #308 from nyaruka/mock_webhooks
 * Add webhook mocks to engine config and check when making webhook calls if there is a matching mock
 * Fix migration of @extra.flow in legacy expressions
 * Fix mapping of @flow.contact in legacy expressions

v0.10.11
----------
 * Don't log errors twice

v0.10.10
----------
 * Add all_groups flag to remove_contact_groups action
 * Add sentry for reporting errors

v0.10.9
----------
 * Fix migrating legacy @flow.foo.text/time expressions

v0.10.8
----------
 * Handle legacy flows where things which are supposed to be translation dicts aren't

v0.10.7
----------
 * Expose flow.revision in context and webhook payloads

v0.10.6
----------
 * Fix panic when field value is nil

v0.10.5
----------
 * Implement redacting URNs
 * Add numeric ID to contacts

v0.10.4
----------
 * run.input.text should be optional

v0.10.3
----------
 * Allow webhook calls to be mocked and add tests
 * Don't blow up if legacy stickies have floating point positions
 * Update CLEAN() to match legacy behaviour
 * More fixes when evaluating migrated legacy template tests

v0.10.2
----------
 * Drop requirement for input.uuid to be a valid UUIDv4 since it comes from msg.uuid which isn't always valid UUID4

v0.10.1
----------
 * Improve parser error messages and add tests for error messages
 * Update to latest gocommon which fixes tel URN which are shortcodes failing validation

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
