v0.33.9
----------
 * Ingore broken recording dicts on legacy say actions
 * Add is_error as a regular function and has_error as a router case test
 * Fix JSONing time values

v0.33.8
----------
 * Expressions refactor

v0.33.7
----------
 * Ignore errors migrating legacy expressions
 * Regexes for migrating context references

v0.33.6
----------
 * Ignore empty webhook headers during flow migration and validate that header namews are valid during validation
 * ReadFlow should maintain UI.. but as raw JSON
 * Flowrunner improvements

v0.33.5
----------
 * Fix routing after a timeout

v0.33.4
----------
 * Ensure an error from a resume is logged to the run

v0.33.3
----------
 * Fix auditing context refs like foo.0

v0.33.2
----------
 * Add NewActivatedMsgWait

v0.33.1
----------
 * Simplify resuming sessions so we only look at the wait on the actual node
 * Use relative timeout value in activated wait and msg_wait event

v0.33.0
----------
 * Add category_uuid to timeouts on waits, migrate legacy timeout rules to that

v0.32.1
----------
 * Operand of group split should be @contact.groups
 * Migrate a ruleset of type contact_field that splits on @contact.groups to an expression split
 * Add hint to msg_wait event, move waits into router package
 * Update engine to look for wait on router instead of node

v0.32.0
----------
 * Convert almost all complex types to be represented in expressions as simple XDicts
 * Add functional programming basics
 * Add template assets to goflow

v0.31.3
----------
 * Add check to call_resthook that payload is valid JSON

v0.31.2
----------
 * CallResthookAction should error if it can't evaluiate the payload template
 * Resthook payload should still be valid when contact URN can't be formatted

v0.31.1
----------
 * Generate better error message when resthook payload is not valid JSON

v0.31.0
----------
 * Better error message when marshalling a run
 * Use dict() function to simplify default webhook payload
 * Convert @contact.groups to be only excellent primitives
 * Add extract and dict as excellent functions

v0.30.4
----------
 * Add @fields as top-level shortcut to contact fields as map
 * Add @urns as dict of highest-priority URN by scheme
 * Make location parsing more forgiving

v0.30.3
----------
 * Bug fix: switch router should use category from first matching rule
 * Stringify maps with {...} and arrays with [...]

v0.30.2
----------
 * Record exit UUIDs coming from waits in validated flow definition
 * Match characters  intended to be combined with another character to support Thai, Bengali and Burmese properly
 * Extract and save result categories during validation
 * Add validation that node has > 0 exits, routers have > 0 categories, and categories have an exit

v0.30.1
----------
 * Don't try to validate a subflow which is missing

v0.30.0
----------
 * Fix HasDate tests to compare dates in env timezone
 * Return missing assets from SessionAssets.Validate
 * Replace runtime loop detection with an engine limit on steps per sprint (default 100)

v0.29.11
----------
 * Don't trim whitespace on input to has_pattern test

v0.29.10
----------
 * Fix migration of the @step.attacthments array in legacy expressions
 * Merge result infos by key so caller doesn't have to know how to do that

v0.29.9
----------
 * Fix migration of datetime + time in legacy expressions

v0.29.8
----------
 * Fix migration of legacy flows that don't have entry set

v0.29.7
----------
 * Change random router to return raw random value as result value
 * Return results as list of name/key objects during flow validation
 * Order nodes by y during flow migration
 * Convert attachment URLs to absolute during flow migration

v0.29.6
----------
 * Migrate name only legacy label and group references
 * Add more tests for has_pattern

v0.29.5
----------
 * Fix tests

v0.29.4
----------
 * Flow name shouldn't be required (matches any other asset type)

v0.29.3
----------
 * Use presence of flow_type to determine if flow is in legacy format and handle legacy flows with missing metadata section

v0.29.2
----------
 * Allow flow validation without assets

v0.29.1
----------
 * Store a map of result keys to result names in a validated flow definition

v0.29.0
----------
 * call_resthook action should generate result even if there are no subscribers
 * If a resthook call returns a success and a 410, use the success as the result
 * Validating flow should add dependencies and result_names to definition
 * Rework flow validation so dependency checking happens centrally and not in each action
 * Add util methods for enumerating and rewriting templats in group and label references

v0.28.14
----------
 * Allow conversion of numbers to times
 * Migrate datetime+time to a minutes addition expression
 * Add tools.RefactorTemplate and tools.FindContextRefsInTemplate

v0.28.13
----------
 * Fix title to work with text which is uppercase
 * Change implementation of remove_first_word so that punctuation is preserved
 * Anything + TIME should migrate to replace_time(..)

v0.28.12
----------
 * Fix migration of datevalue+time

v0.28.11
----------
 * length(nil) == 0
 * Arrays should stringify as CSV
 * Maps should stringify as new line separated key: value pairs

v0.28.10
----------
 * Wrap results of date arithmetic in format_date

v0.28.9
----------
 * Fix calling length on a complex object that needs to be reduced

v0.28.8
----------
 * Fix migration of DAYS()

v0.28.7
----------
 * add accessor for msg in MsgResume

v0.28.6
----------
 * DATEVALUE should migrate to date() so it returns a date rather than a datetime
 * Add date() conversion function
 * Change datetime_from_parts to date_from_parts

v0.28.5
----------
 * fix resolving @parent or @child when they are nil

v0.28.4
----------
 * Add utils.Date and types.XDate

v0.28.3
----------
 * Don't localize fields not localized in legacy engine
 * Rework time parsing to accept hour only and more ISO8601 formats

v0.28.2
----------
 * allow single tls renegotiation

v0.28.1
----------
 * Allow looping back into a flow that wasn't started in this sprint

v0.28.0
----------
 * Rename from_epoch to datetime_from_epoch
 * Add time functions and has_time router test
 * Remove support for .0 indexing in excellent
 * Remove gomobile dependency

v0.27.9
----------
 * Don't treat identifiers as special case, parse them like all other expressions
 * When migrating expressions like flow.2factor, wrap non-name keys in ["..."]

v0.27.8
----------
 * Add better flow spec version handling and ability to peek at definitions to determin if they are legacy

v0.27.7
----------
 * Don't fake an ignored response body but record it in the event as ignored
 * If a webhook call doesn't return a content-type header, try to detect type

v0.27.6
----------
 * Verify parsed numbers in has_phone are valid

v0.27.5
----------
 * update contact PreferredChannel and PreferredURN to resolve first sendable destination

v0.27.4
----------
 * add country to mobile.NewEnvironment
 * move NewSession and ReadSession into Engine
 * engine.EngineBuilder -> engine.Builder because gofmt doesn't like stuttering
 * remove webhook mocking and reading engine coonfig from JSON

v0.27.3
----------
 * Fix migration of weekday() to add 1
 * Add resthook_called event

v0.27.2
----------
 * simplify reading environments from JSON
 * rename .Environment() to .Build()

v0.27.1
----------
 * Add max_value_length to environment and apply in name and field change modifiers

v0.27.0
----------
 * Add type to sessions (the type of the flow it was triggered with)

v0.26.0
----------
 * Update modifier loading so missing assets are reported
 * Remove logrus logging
 * Rework all session objects to record missing assets
 * Ensure error events are logged to sprint as well as run

v0.25.5
----------
 * Include status code in webhook events

v0.25.4
----------
 * Fix renderEventDoc to properly render JSON in markdown
 * Add full constructor for sprint

v0.25.3
----------
 * Remove flow server components which are no longer used

v0.25.2
----------
 * Simpler trigger constructors

v0.25.1
----------
 * Fix unmarshalling legacy say actions

v0.25.0
----------
 * New IVR events

v0.24.1
----------
 * Fix time-filling bug

v0.24.0
----------
 * Convert to go module to be used as library

v0.23.0
----------
 * Migrate api actions so that URL expressions are wrapped in url_encode()
 * Don't url encode msg attachment expressions automatically
 * Fix not being to read contact_field_changed events where value is null
 * Add default country to environment and use for has_phone tests
 * Fix add_contact_urn so that URN is normalized and trimmed before being added
 * Fix legacy_extra so it can handle root-level arrays
 * Replace contact_urn_added event with contact_urns_changed and fix tracking of channel affinity on URNs
 * Fix not supporting dymnamic groups based on name
 * Change send_msg to always send a message even if it can't resolve channel/URN
 * Fix 410 resthook response becoming result with no category
 * Switch to codecov
 * Validate all flows which are referenced in the current flow
 * Use new status of subscriber_gone when resthook call returns 410
 * Require exit UUIds to be unique across the entire flow
 * Improve number parsing

v0.22.0
----------
 * Move input from run to session
 * Make contact.name .language .timezone omitted in JSON when empty
 * Make contact.created_on required
 * Simplify un/marshaling of typed objects
 * Allow waits to skip themselves, and have the msg wait skip itself if it's the first thing after a msg trigger
 * Improve flowrunner cmd to display more event types
 * Add new terminal option to start_flow actions
 * Ensure arguments to date router tests are migrated
 * Create new static asset source type for simpler testing
 * Add channel event trigger type

v0.21.4
----------
 * Rework transferto action to generate result like a webhook call

v0.21.3
----------
 * Add results from run summary on flow_action triggers to @legacy_extra
 * Fix migration of @flow by itself

v0.21.2
----------
 * Fix migration of calls to HOUR(...)
 * Support formatting/parsing of decimal values with configurable digit separators

v0.21.1
----------
 * Re-evaluate and correct contact groups at the start/resume of a session
 * Add run.modified_on

v0.21.0
----------
 * Change sigature of ReevaluateDynamicGroups to not require a session
 * Fix language selection for flow localization 

v0.20.1
----------
 * Fix run expirations
 * Rename webhook_called.time_taken -> elapsed_ms

v0.20.0
----------
 * Replace caller events with resumes
 * Record time taken in webhook_called events

v0.19.2
----------
 * Only generate a contact_field_changed event if a value has actually changed

v0.19.1
----------
 * Fix blowing up when contact doesn't have a value for a contact field

v0.19.0
----------
 * Add missing regex_match func and migration of string literals in legacy flows
 * Allow contacts to be loaded in a mode that ignores missing groups or fields
 * Disallow creating contacts with empty field values
 * Generate a JSON listing of functions
 * Cleanup field value code so that we no longer need empty values
 * Fix group reevaluation when contact has no value for a text field
 * Cleanup function docstrings and fix from_epoch

v0.18.0
----------
 * Use default value for router on migrated webhook ruleset in case resthook didn't have any subscribers
 * Add resthook slug to webhook_called events if it exists
 * Generate a groups changed event when dynamic groups are re-evaluated
 * Contact field changed events should have the entire value objects
 * Non-caller events should only ever be added to runs and not applied
 * Merged contact_groups_add/removed into contact_groups_changed
 * Only generate events when state has actually changed
 * Move event functiionality into the actions that generate them

v0.17.1
----------
 * Use result name when populating @legacy_extra instead of just .webhook

v0.17.0
----------
 * Remove Connection Error as a separate exit for webhook routers
 * call_webhook and call_resthook now save a result and run.Webhook is removed
 * @input and @results as shortcuts to @run.input and @run.results
 * Router tests can now return extra to be added to the result

v0.17.0
----------
 * Remove "Connection Error" as a separate exit for webhook routers
 * call_webhook and call_resthook now save a result and run.Webhook is removed
 * @input and @results as shortcuts to @run.input and @run.results
 * Router tests can now return extra to be added to the result

v0.16.0
----------
 * Adds @legacy_extra to context to mimic @extra in legacy flows

v0.15.2
----------
 * Fix determining whether asset serversource supports a particular type
 * Fix parsing of numbers when string contains uppercase letters

v0.15.1
----------
 * Fill in current time when parsing dates during tests
 * Add ui config options for composed split migrations

v0.15.0
----------
 * Great assets refactor

v0.14.7
----------
 * Replicate channel matching logic from RP

v0.14.6
----------
 * Don't strip URN params when parsing

v0.14.5
----------
 * Migrate legacy webhook calls to POST if they don't have a method set

v0.14.4
----------
 * Simplify string tokenization and add more tests

v0.14.3
----------
 * Fix message attachment expressions not being URL encoded
 * Fix parsing of locations

v0.14.2
----------
 * Small refactor of how we migrate rules to cases/exits that fixes an ordering bug and add option to not collapse exits

v0.14.1
----------
 * Fix sentry integration
 * Fix FIELD() blowing up when using space as separator

v0.14.0
----------
 * Documentation improvements and make language required in flow definition
 * Add flow type and validate that actions only occur in supported flow types

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
