v0.13.3
----------
 * Fix scanner treating parentheses inside string literals as expression boundaries
 * Fix not correctly scanning excellent identifiers followed by periods
 * Add migrate tab to flowserver index page for easy flow migration testing

v0.13.2
----------
 * Refactor assets functionality

v0.13.1
----------
 * to_epoch() becomes epoch() and returns fractional seconds
 * Improved docuentation
 * Add Heroku deployment support

v0.13.0
----------
 * Wrap multiple asset responses in a results object

v0.12.2
----------
 * Update UI node type names
 * Add all trigger types to docs so that we test parsing of the examples
 * Add campaign as a trigger type

v0.12.1
----------
 * Improve error messages when trigger can't be read due to asset load failure

v0.12.0
----------
 * Migrate flows individually

v0.11.0
----------
 * Rework assets so whether or not they are managed as sets is configured at the type level
 * LocationHiearchy assets should be managed like a set like other asset types
 * Migrate empty flows
 * TransferTo airtime transfer action
 * Simplify server handler methods and improve validation errors involving fields on composed types
 * Dynamic type for actions, events, routers etc
 * Make func names more consistent and implement same typing system for waits

v0.10.21
----------
 * Tweak resthook_called events so the event itself has status=success for HTTP 410 responses

v0.10.20
----------
 * Migrate ruleset_types to editor types in _ui
 * Use a sequential time source for flow tests instead of replacing time values with placeholders
 * Add call_resthook action type and migrate from legacy resthook rulesets
 * Improve error message when switch router test function returns error
 * Add support for making expressions evaluate to themselves on error
 * Migrate "" escape sequences in string literals in legacy expressions to \"
 * Add support for \n \" sequences in Excellent string literals

v0.10.19
----------
 * Add format_date function which only takes date (non-time) formatting chars
 * Fix encoding of contact name in webhook payloads
 * Fix logging of panics to sentry

v0.10.18
----------
 * Encode spaces as %20 in URL expressions
 * Match legacy behavior for @contact

v0.10.17
----------
 * Add proper migrations for when word_* functions have final by_spaces param, and also collapse decremented values if they are literals (e.g. 2- 1)
 * Update word(), word_slice() and word_count() to take a final optional param called delimiters
 * Don't throw validation error if add_input_labels or add_contact_groups has zero groups/labels

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
