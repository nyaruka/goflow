v0.250.0 (2025-06-30)
-------------------------
 * Add UUIDs to events and remove UUID from msg and airtime transfer

v0.249.2 (2025-06-25)
-------------------------
 * For reals this time

v0.249.1 (2025-06-25)
-------------------------
 * Fix stuttering in waits/hints package too

v0.249.0 (2025-06-25)
-------------------------
 * Remove stuttering

v0.248.5 (2025-06-25)
-------------------------
 * Remove has_intent, has_top_intent and has_category from documentation

v0.248.4 (2025-06-24)
-------------------------
 * Make Contact.status required

v0.248.3 (2025-06-24)
-------------------------
 * Make it easier to embed a contact in JSON
 * Cleanup after recent IVR related changes

v0.248.2 (2025-06-24)
-------------------------
 * Convert remaining resume types to be event based

v0.248.1 (2025-06-21)
-------------------------
 * Cleanup environment stuff in resumes

v0.248.0 (2025-06-21)
-------------------------
 * Simplify starting new sessions
 * Remove environment and contact from triggers and resumes

v0.247.10 (2025-06-20)
-------------------------
 * Remove calls from triggers and instead require passing via engine methods

v0.247.9 (2025-06-20)
-------------------------
 * Add flows.NewCallUUID

v0.247.8 (2025-06-20)
-------------------------
 * Make calls objects with real channels and add UUID

v0.247.7 (2025-06-19)
-------------------------
 * Rework campaign event triggers to use events

v0.247.6 (2025-06-19)
-------------------------
 * Remove group field on campaign assets

v0.247.5 (2025-06-19)
-------------------------
 * Rename campaign event to campaign point
 * Model campaigns as actual assets

v0.247.4 (2025-06-18)
-------------------------
 * Convert optin triggers to use real events

v0.247.3 (2025-06-18)
-------------------------
 * Rework ticket triggers to use real events

v0.247.2 (2025-06-18)
-------------------------
 * Stop writing trigger.msg and resume.msg

v0.247.1 (2025-06-17)
-------------------------
 * Record flows entered on sprints

v0.247.0 (2025-06-17)
-------------------------
 * Rework msg triggers and resumes to take a msg_received event instead of a msg

v0.246.6 (2025-05-21)
-------------------------
 * Add Sprint.IsInitial

v0.246.5 (2025-05-21)
-------------------------
 * Add Session.CreatedOn

v0.246.4 (2025-05-19)
-------------------------
 * Add Trigger.SetCall

v0.246.3 (2025-05-19)
-------------------------
 * Fix reading of trigger without contact

v0.246.2 (2025-05-19)
-------------------------
 * Add SetContact to Trigger interface

v0.246.1 (2025-05-19)
-------------------------
 * Allow triggers to be persisted without the contact

v0.246.0 (2025-05-19)
-------------------------
 * Remove support for Facebook topics

v0.245.1 (2025-05-14)
-------------------------
 * Change user references to be by UUID with email being the variable form

v0.245.0 (2025-05-13)
-------------------------
 * Add UUID to user assets

v0.244.0 (2025-05-08)
-------------------------
 * Remove ID from MsgIn and add externalID as param to NewMsgIn

v0.243.3 (2025-05-06)
-------------------------
 * Tweak spec docs to include limits
 * Rework all tests using mocked time or rng to call new test.MockUniverse function
 * Use v7 UUIDs for airtime transfers and tickets as well
 * Update dependencies

v0.243.2 (2025-05-06)
-------------------------
 * Use v7 UUIDs for session UUIDs
 * Add functions to generate UUIDs for msgs, sessions, runs and contacts to ensure consistency

v0.243.1 (2025-05-05)
-------------------------
 * Use v7 UUIDs for messages and runs

v0.243.0 (2025-05-05)
-------------------------
 * Relax UUID validation to allow for v4 and v7 UUIDs
 * Convert event marshalling tests to use snapshots
 * Fix all the places we unnecessarily use json.RawMessage
 * Hide .extra on results from auto-complete

v0.242.1 (2025-04-30)
-------------------------
 * Don't create result when call_resthook action has no subscribers

v0.242.0 (2025-04-29)
-------------------------
 * Migrate open ticket and transfer airtime actions to move result generation to the router

v0.241.0 (2025-04-25)
-------------------------
 * Include locals in flow inspection

v0.240.5 (2025-04-24)
-------------------------
 * Update airtime transfer action to not require a result

v0.240.4 (2025-04-23)
-------------------------
 * Update dependencies
 * Require languages to be valid ISO codes not just 3 letter strings

v0.240.3 (2025-04-23)
-------------------------
 * Augment 14.2 migration to cleanup invalid languages

v0.240.2 (2025-04-22)
-------------------------
 * Rename flows.Inspection to flows.Info

v0.240.1 (2025-04-22)
-------------------------
 * Add node and language counts to inspection results

v0.240.0 (2025-04-22)
-------------------------
 * Add basic validation of flow localization
 * Add 14.2 flow migration to rename open_ticket.body to open_ticket.note
 * Simplify extracting results and dependencies during inspection

v0.239.0 (2025-04-21)
-------------------------
 * Make result optional on open ticket actions but always write a local that can be used for splitting

v0.238.2 (2025-04-15)
-------------------------
 * Amend 14.1 migration to trim quick replies to 10

v0.238.1 (2025-04-14)
-------------------------
 * Relax length limit on quick replies since they can be expressions but add limits elsewhere

v0.238.0 (2025-04-14)
-------------------------
 * Add 14.1 migration to make webhook nodes split on @webhook and fix quick replies that are too long

v0.237.0 (2025-04-14)
-------------------------
 * Limit groups and label actions to 100 items
 * Limit argument count on router cases
 * Add validation for flow names and number of router cases

v0.236.4 (2025-04-11)
-------------------------
 * For validation errors that relate to limits, include the count in the message

v0.236.3 (2025-04-10)
-------------------------
 * Amend 14.0 migration to also fix empty category names

v0.236.2 (2025-04-10)
-------------------------
 * Move dtone airtime service out of this repo and use mock in tests

v0.236.1 (2025-04-09)
-------------------------
 * Remove luis and bothub classification services

v0.236.0 (2025-04-09)
-------------------------
 * Refactor webhooks to 1) to clear run.webhook in case of error and 2) let run.webhook be more easily restored

v0.235.4 (2025-04-09)
-------------------------
 * Allow users to manage the whatsapp channels by using the preferred channel action

v0.235.3 (2025-04-08)
-------------------------
 * Simplify naming llm_prompt to prompt
 * Tweak llm_prompt so onus is on the caller to format args

v0.235.2 (2025-04-07)
-------------------------
 * Fix invalid expires values in 14.0.0 migration

v0.235.1 (2025-04-07)
-------------------------
 * Put flows.NewAssetsEnvironment back as we need it for non-flow contact modifications

v0.235.0 (2025-04-07)
-------------------------
 * Bump spec version to 14.0.0
 * Add call_llm action with corresponding llm_called event.
 * Add local variables to runs and expose in context
 * Add @node.categories as the list of category names on the current node

v0.234.3 (2025-04-01)
-------------------------
 * Add maxTokens param to LLMService.Response and remove unused TicketService

v0.234.2 (2025-03-31)
-------------------------
 * Tweak tokensused to be int64

v0.234.1 (2025-03-31)
-------------------------
 * LLM services should return tokens used as well as the response

v0.234.0 (2025-03-31)
-------------------------
 * Allow embedding of _ui bods in item translations

v0.233.1 (2025-03-29)
-------------------------
 * Tweak test LLM service

v0.233.0 (2025-03-28)
-------------------------
 * Make engine context.Context aware and pass down into services

v0.232.4 (2025-03-27)
-------------------------
 * Tweak LLM service implementation for testing to just echo input back

v0.232.3 (2025-03-27)
-------------------------
 * Reorg test services

v0.232.2 (2025-03-27)
-------------------------
 * Tweak signature of LLMService.Response

v0.232.1 (2025-03-27)
-------------------------
 * Add LLM service

v0.232.0 (2025-03-26)
-------------------------
 * Add support for LLM assets
 * Update phonenumbers library
 * Udpate to go 1.24

v0.231.1 (2025-03-12)
-------------------------
 * Allow flows.MsgContent to be unmarshalled from string

v0.231.0 (2025-03-11)
-------------------------
 * Add limits for nodes per flow, actions per node, categories and exits
 * Update deps

v0.230.0 (2025-02-27)
-------------------------
 * Add extra field to quick replies and support encoding as text using new line
 * Rewrite places we use non-constant format strings
 * Change flow expires after default to 3 days

v0.229.1 (2025-02-13)
-------------------------
 * Marshal quick replies as structs

v0.229.0 (2025-02-13)
-------------------------
 * Convert quick replies into structs

v0.228.1 (2025-01-29)
-------------------------
 * Add UUID to sprints

v0.228.0 (2025-01-24)
-------------------------
 * Make expires_on non-nil on waits
 * Add fallback values for expires after for messaging and voice flow types
 * Use time.Duration internally for flow wait expire value

v0.227.0 (2025-01-08)
-------------------------
 * Enforce a limit of 50 URNs per contact

v0.226.1 (2025-01-06)
-------------------------
 * Disallow empty categories
 * Allow category names of up to 64 chars

v0.226.0 (2024-12-16)
-------------------------
 * Remove waiting_exits from flow inspection as it's no longer needed

v0.225.8 (2024-12-12)
-------------------------
 * Completely relax requirements for result names

v0.225.7 (2024-12-11)
-------------------------
 * Allow more characters in result names

v0.225.6 (2024-12-11)
-------------------------
 * Allow more characters in result names

v0.225.5 (2024-12-10)
-------------------------
 * Update deps
 * Revert keying results by name rather than snakified name

v0.225.4 (2024-12-10)
-------------------------
 * Rework results so they're keyed by name rather than snakified name

v0.225.3 (2024-12-09)
-------------------------
 * Add patch flow migration to fix result references in expressions that need to be truncated

v0.225.2 (2024-12-09)
-------------------------
 * Allow reading of results in sessions with invalid names

v0.225.1 (2024-12-05)
-------------------------
 * Add custom validators for result name and category and make them match current floweditor validation

v0.225.0 (2024-12-03)
-------------------------
 * Add new 13.6 flow migration to truncate result names and categories

v0.224.1 (2024-12-02)
-------------------------
 * Remove unused fields from run result events
 * Only generate run_result_changed event when value or category changes

v0.224.0 (2024-12-02)
-------------------------
 * Update deps
 * Include previous result value on run_result_changed events

v0.223.0 (2024-11-12)
-------------------------
 * Change transfer_airtime to fail when no exact matching amount
 * Update deps

v0.222.5 (2024-09-30)
-------------------------
 * Fix inspection of contact queries using group = "" or group != ""

v0.222.4 (2024-09-19)
-------------------------
 * Add query to trigger_session and create_broadcast action constructors

v0.222.3 (2024-09-17)
-------------------------
 * Validate result name/category fields in actions that create results

v0.222.2 (2024-09-16)
-------------------------
 * Update dependencies
 * Update to go 1.23

v0.222.1 (2024-08-21)
-------------------------
 * Add workaround for numbers that pass urns.ParseNumber but fail phone URN valdiation
 * Drop support for localizing template variables in send_msg actions

v0.222.0 (2024-08-19)
-------------------------
 * Enable to send airtime to WhatsApp URNs
 * Remove old body field on ticket modifier

v0.221.1 (2024-08-06)
-------------------------
 * On second thoughts, better to just bundle the note with the open ticket event

v0.221.0 (2024-08-06)
-------------------------
 * Replace ticket body with initial note

v0.220.0 (2024-08-02)
-------------------------
 * Update to latest gocommon

v0.219.2 (2024-07-16)
-------------------------
 * Add validation of attachment values in flow definitions

v0.219.1 (2024-07-16)
-------------------------
 * Don't use template variables as attachments unless they're valid

v0.219.0 (2024-07-15)
-------------------------
 * Update to latest gocommon which changes contact ql conversion to ES

v0.218.3 (2024-07-04)
-------------------------
 * Change flos.NewMsgOut to take flows.MsgContent

v0.218.2 (2024-07-04)
-------------------------
 * Move code to generate templating into flows.Template

v0.218.1 (2024-07-03)
-------------------------
 * Remove namespace on template translations as it's not needed by engine and can be added by mailroom when queueing
 * Move functionality to generate preview content for template message from send_msg action to TemplateTranslation

v0.218.0 (2024-07-03)
-------------------------
 * Rework broadcast translation selection to merge across text, attachments and quick replies

v0.217.0 (2024-06-26)
-------------------------
 * Drop db support on flows.BroadcastTranslations

v0.216.4 (2024-06-14)
-------------------------
 * Revert adding support for sending airtime to WhatsApp

v0.216.3 (2024-06-12)
-------------------------
 * Use empty value for airtime failure result

v0.216.2 (2024-06-12)
-------------------------
 * Save airtime transfer transaction id as result and on event as external_id
 * Remove support for old template component types

v0.216.1 (2024-06-12)
-------------------------
 * Add transfer UUID field to airtime transfer

v0.216.0 (2024-06-11)
-------------------------
 * Support sending for WhatsApp URNs
 * Update DT One requests

v0.215.1 (2024-06-10)
-------------------------
 * Support template component types like header/* and body/*

v0.215.0 (2024-06-06)
-------------------------
 * Support template header components of type image, video or document and use as attachments on preview message

v0.214.6 (2024-06-05)
-------------------------
 * Update goreleaser config to v2

v0.214.5 (2024-05-24)
-------------------------
 * Make date parsing stricter

v0.214.4 (2024-05-22)
-------------------------
 * Ensure QueryError always has a code and use `syntax` as code for all parser errors

v0.214.3 (2024-05-21)
-------------------------
 * Tweak type of extra on query errors

v0.214.2 (2024-05-21)
-------------------------
 * Use std library errors

v0.214.1 (2024-05-20)
-------------------------
 * Move elastic utils to gocommon

v0.214.0 (2024-05-17)
-------------------------
 * Rewrite contactql query conversion to ES to use simple maps instead of unmaintained library

v0.213.3 (2024-05-17)
-------------------------
 * Merge pull request #1255 from nyaruka/migration_fix
 * Fix migrating directly from 13.2 to 13.5

v0.213.2 (2024-05-16)
-------------------------
 * Tweak releaser config

v0.213.1 (2024-05-16)
-------------------------
 * Switch to use goreleaser and update actions to latest versions

v0.213.0 (2024-05-15)
-------------------------
 * Normalize contactql queries with type prefixes

v0.212.0 (2024-05-15)
-------------------------
 * Add support for contactql conditions with properties prefixed by their type
 * Update to latest ANTLR4

v0.211.2 (2024-05-09)
-------------------------
 * Update gocommom

v0.211.1 (2024-05-09)
-------------------------
 * Update to latest gocommon

v0.211.0 (2024-05-08)
-------------------------
 * Update to latest gocommon
 * Remove old urn modifier which was replaced by the urns modifier

v0.210.0 (2024-04-25)
-------------------------
 * Add new 13.5 flow spec which changes send_msg to use single list of variables

v0.209.2 (2024-04-17)
-------------------------
 * Remove deprecated params in templating

v0.209.1 (2024-04-16)
-------------------------
 * Re-include legacy .params field on components temporarily

v0.209.0 (2024-04-16)
-------------------------
 * Rework the format of MsgTemplating

v0.208.0 (2024-04-15)
-------------------------
 * Update to latest decimal library that fixes .Pow
 * Rework template assets to use variables instead of component specific param lists

v0.207.3 (2024-03-20)
-------------------------
 * Remove no longer used templating.params

v0.207.2 (2024-03-20)
-------------------------
 * Include component name on templating

v0.207.1 (2024-03-20)
-------------------------
 * Drop support for variables in display parts of template components

v0.207.0 (2024-03-20)
-------------------------
 * Rework template components... again

v0.206.4 (2024-03-19)
-------------------------
 * Add type and display fields to temmplate components

v0.206.3 (2024-03-15)
-------------------------
 * Truncate quick replies for template preview messages

v0.206.2 (2024-03-14)
-------------------------
 * Update to go 1.22
 * Mark @child.run.status as deprecated in context
 * Reorg of expression refactoring code into its own package

v0.206.1 (2024-03-05)
-------------------------
 * Remove no longer used uuid property from msg.templating

v0.206.0 (2024-03-05)
-------------------------
 * Instead of supporting both variables and components in send_msg actions, add 13.4 migration

v0.205.0 (2024-03-05)
-------------------------
 * Support excluding deprecated object properties from marshalling
 * Deprecate undocumented multi-match result fields
 * Update autocompletion docs for results extra

v0.204.1 (2024-03-04)
-------------------------
 * Update docs for @webhook

v0.204.0 (2024-03-04)
-------------------------
 * Make @webhook an object with headers and json properties

v0.203.0 (2024-02-27)
-------------------------
 * Rework semd_msg templating to have a simpler structure in flow definitions and fix enumerating templates and localizable strings

v0.202.1 (2024-02-23)
-------------------------
 * Limit use of reflect nil checking

v0.202.0 (2024-02-20)
-------------------------
 * Add warnings from template eval functions and recorded as events
 * Add excellent.Evaluator and move top level evaluation functions to methods on that
 * Rework xvalues to be passed by pointer and have a base struct
 * Rework router reading to be like action reading (and allow for enumeration of types)

v0.201.0 (2024-02-12)
-------------------------
 * Replace template translation content and params by the components fields

v0.200.1 (2024-01-30)
-------------------------
 * Support more WA templates components

v0.200.0 (2024-01-22)
-------------------------
 * Rework template assets to use params instead of variables

v0.199.0 (2024-01-17)
-------------------------
 * Replace flows.MsgTemplating variables by params

v0.198.0 (2024-01-02)
-------------------------
 * Update deps
 * Remove ticketers

v0.197.7 (2023-12-12)
-------------------------
 * Merge pull request #1198 from nyaruka/trim_attachments
 * Trim whitespace from evaluated message attachments

v0.197.6 (2023-11-27)
-------------------------
 * Update to latest gocommon
 * Rework template translation finding

v0.197.5 (2023-11-20)
-------------------------
 * Update deps

v0.197.4 (2023-11-13)
-------------------------
 * Tweak Arabic Variants collation
 * Update depedencies

v0.197.3 (2023-10-30)
-------------------------
 * Update to latest gocommon

v0.197.2 (2023-10-27)
-------------------------
 * Less panics, more validation errors for invalid flow defs

v0.197.1 (2023-10-12)
-------------------------
 * Update deps

v0.197.0 (2023-10-06)
-------------------------
 * Drop support for contact less sessions
 * Add devcontainer config

v0.196.2 (2023-10-02)
-------------------------
 * Don't treat ژ and ږ as equivalent in arabic variants collation

v0.196.1 (2023-09-27)
-------------------------
 * Remove deprecated collation alias

v0.196.0 (2023-09-26)
-------------------------
 * Update to go 1.21
 * Add more mappings to arabic collation

v0.195.4 (2023-09-23)
-------------------------
 * Fix inspectng optin refs

v0.195.3 (2023-09-21)
-------------------------
 * Rename send_optin to request_optin and optin_created to optin_requested
 * Add mores tests for arabic-farsi collation

v0.195.2 (2023-09-18)
-------------------------
 * Rename optin_sent to optin_created for consistency with msg_created and add channel/URN properties

v0.195.1 (2023-09-18)
-------------------------
 * Fix optin triggers marshaling with type ticket

v0.195.0 (2023-09-14)
-------------------------
 * Add send_optin action and optin_sent event
 * Add features to channel assets and define optins as a feature
 * Add trigger type for optins
 * Add OptIn as new asset type

v0.194.0 (2023-09-11)
-------------------------
 * Update to latest gocommon
 * Remove support for delegate channels

v0.193.1 (2023-09-04)
-------------------------
 * Update to latest gocommon

v0.193.0 (2023-09-04)
-------------------------
 * Simplify naming in po package to avoid stuttering
 * Switch to i18n package from gocommon

v0.192.1 (2023-08-31)
-------------------------
 * Update to latest gocommon
 * Replace interface{} with any

v0.192.0 (2023-08-18)
-------------------------
 * Replace input cleaning with input collation

v0.191.0 (2023-08-17)
-------------------------
 * Add input cleaners

v0.190.1 (2023-08-15)
-------------------------
 * Merge pull request #1176 from nyaruka/split_env_types
 * Split flows.Environment into AssetsEnvironment and SessionEnvironment

v0.190.0 (2023-08-15)
-------------------------
 * Replace max_value_length on env with new engine options for max field and result chars
 * Simplify how we create a merged environment with properties from the contact

v0.189.0 (2023-08-11)
-------------------------
 * Add filter expression function

v0.188.0 (2023-08-09)
-------------------------
 * Update go to 1.20
 * Make numerical router tests work with Eastern Arabic digits

v0.187.0 (2023-07-20)
-------------------------
 * Update deps

v0.186.1 (2023-07-10)
-------------------------
 * Fix has_ward argument ordering

v0.186.0 (2023-05-24)
-------------------------
 * Add contains function

v0.185.2 (2023-05-03)
-------------------------
 * Support hyphen in locations name as one word

v0.185.1 (2023-05-02)
-------------------------
 * Fix matching location name with diacritics

v0.185.0 (2023-04-26)
-------------------------
 * Contacts only have one (open) ticket
 * Bump golang.org/x/net from 0.5.0 to 0.7.0

v0.184.4 (2023-02-27)
-------------------------
 * Sort should return error (instead of panic) if array has mixed types

v0.184.3 (2023-02-22)
-------------------------
 * Re-add matching by group name in start_session actions

v0.184.2 (2023-02-13)
-------------------------
 * Remove matching legacy vars against group names for start_session and send_broadcast actions

v0.184.1 (2023-02-08)
-------------------------
 * Tweak BroadcastTranslations.ForContact to also return language

v0.184.0 (2023-02-08)
-------------------------
 * Truncate quick replies that are too long
 * Add support for resolving broadcast translations for different contacts

v0.183.0 (2023-02-07)
-------------------------
 * Ignore importing split arguments translations
 * Convert StringSet functions to be generic

v0.182.0 (2023-02-02)
-------------------------
 * Reorganize ANTLR generated code as per recommendations in ANTLR README
 * Upgrade to latest ANTLR and rebuild grammars
 * Add support for contact queries in broadcasts

v0.181.0 (2023-02-01)
-------------------------
 * Update to validator v10
 * Update to nyaruka/null v2

v0.180.0 (2023-01-26)
-------------------------
 * Make envs.Language, envs.Country and envs.Locale play nice with NULLs when persisting
 * Fix docstring on ivr_created event
 * Fix wording in spec migration comments

v0.179.1 (2023-01-18)
-------------------------
 * Add migration versions to auto-generated docs

v0.179.0 (2023-01-18)
-------------------------
 * Add 13.2.0 flow migration to replace "base" as flow language with "und"

v0.178.2 (2023-01-11)
-------------------------
 * Fix and improve mobile bindings to work without wait objects on sessions

v0.178.1 (2022-12-16)
-------------------------
 * Tweak envs.NewLocale(..) to allow passing NilCountry

v0.178.0 (2022-12-16)
-------------------------
 * Replace TemplateTranslation.language and country fields with single locale field and also don't duplicate locale on MsgTemplating since it's now on the message itself
 * Derive contact country from tel urns if no preferred channel
 * Change msg_out.language to locale

v0.177.0 (2022-12-15)
-------------------------
 * Reduce msg/ivr_created localization blobs to a single language

v0.176.0 (2022-12-15)
-------------------------
 * Add localization outcome object to msg/ivr_created events to replace msg_out.text_language field
 * Simplify getting text translations in flow actions

v0.175.0 (2022-11-30)
-------------------------
 * Add a new issue type for use of legacy vars in start_session actions

v0.174.4 (2022-11-28)
-------------------------
 * Update deps

v0.174.3 (2022-11-18)
-------------------------
 * Update to latest gocommon

v0.174.2
----------
 * Add support for parsing attachments with unavailable: prefixes

v0.174.1
----------
 * Update to latest gocommon

v0.174.0
----------
 * Add optional revision number to flow refs, use on run.flow

v0.173.2
----------
 * Also fix case where translation just has different number of args

v0.173.1
----------
 * Ignore empty string translations in flow localization sections

v0.173.0
----------
 * Remove merging translation text arrays as it never makes sense

v0.172.6
----------
 * Fix bubbling failing runs in case of missing flow

v0.172.5
----------
 * Fix parsing of geo type attachments

v0.172.4
----------
 * Make utils.Attachment smarter about parsing raw URLs

v0.172.3
----------
 * Update dependencies

v0.172.2
----------
 * Dial wait expiry should be based on the maximum time the forwarded call might take

v0.172.1
----------
 * Set defaults for dial wait time limits

v0.172.0
----------
 * Add time limits to dial waits

v0.171.0
----------
 * Rename connection to call in IVR context

v0.170.0
----------
 * Add unsendable reason to MsgOut and set when contact status is not active

v0.169.0
----------
 * Update to latest gocommon
 * Add keys function to get object properties
 * Use go 1.19

v0.168.2
----------
 * Update to latest gocommon

v0.168.1
----------
 * Update to latest gocommon

v0.168.0
----------
 * Use HTTP log functionality from gocommon
 * Update parsed docstrings to use codeblock formating (tab) for examples

v0.167.0
----------
 * Update to latest gocommon

v0.166.0
----------
 * Add contact status to variable context

v0.165.5
----------
 * Update to latest gocommon which adds teams URNs

v0.165.4
----------
 * Another tweak to SessionBuilder so that it doesn't override assets

v0.165.3
----------
 * More tweaks to SessionBuilder

v0.165.2
----------
 * Allow passing of a flows.SessionAssets to test.SessionBuilder

v0.165.1
----------
 * Updated translations from Transifex

v0.165.0
----------
 * Add modifier to open a ticket

v0.164.3
----------
 * Modifer.Apply should return if it made a modification and caller responsible for calling ReevaluateGroups if so
 * Add Services param to Modifier.Apply

v0.164.2
----------
 * Pass session assets to service factories

v0.164.1
----------
 * Remove flow sessions from service factories

v0.164.0
----------
 * Rework services to not require sessions

v0.163.0
----------
 * Always truncate URLs in HTTP logs to 2048 chars

v0.162.1
----------
 * Support simplifying of queries than can't be parsed but can be constructed

v0.162.0
----------
 * Provide API for building contact queries programatically

v0.161.2
----------
 * Update ANTLR
 * Update to better maintained fork of go-mail

v0.161.1
----------
 * Ensure that failing a session doesn't leave runs in active/waiting state

v0.161.0
----------
 * Update to latest gocommon and phonenumbers

v0.160.0
----------
 * Add option to exclude contacts in a flow on start session action

v0.159.2
----------
 * Trim URLs in call_webhook actions

v0.159.1
----------
 * Fix not equals conditions in contact queries on fields that aren't set

v0.159.0
----------
 * Improve simplifying of contactql queries
 * Update direct dependencies except ANTLR4
 * Go 1.18

v0.158.1
----------
 * send_msg action should fallback to template trans in env default language if no trans found for contact language

v0.158.0
----------
 * Add status as a contact query attribute, disallowed for smart groups

v0.157.0
----------
 * Update to latest gocommon
 * Allow querying on whether group is set or not for consistency with other fields
 * Support contact queries on flow history

v0.156.1
----------
 * SessionAssets implementation of contactql.Resolver methods should return pure assets

v0.156.0
----------
 * Give flows.Flow a reference to their asset if they have one

v0.155.0
----------
 * Switch from flow to flow_id and groups to group_ids for ES queries

v0.154.0
----------
 * Give errors returned from Session.Resume codes

v0.153.0
----------
 * Add concat excellent function
 * Updated translations from Transifex

v0.152.0
----------
 * start_session actions should generate error event if flow asset missing

v0.151.0
----------
 * Add flow as contactql query attribute

v0.150.2
----------
 * Export events.BaseEvent so that callers can create their own events

v0.150.1
----------
 * Fix bug when we remove a contact from all static groups

v0.150.0
----------
 * If caller tries to resume with wrong resume type, don't fail session but error instead

v0.149.1
----------
 * Update to gocommon 1.17.1

v0.149.0
----------
 * Remove no longer used Run.expires_on

v0.148.0
----------
 * Add expiresOn to dial waits so all wait types have it

v0.147.0
----------
 * Add @trigger.campaign for campaign triggers
 * Only treat start_session legacy vars as tel URNs if they are parseable phone numbers

v0.146.1
----------
 * Update to latest gocommon

v0.146.0
----------
 * Rename FlowRun to Run
 * Update to latest gocommon

v0.145.0
----------
 * Add expires_on to msg_wait events
 * Remove activated wait objects on sessions, callers should use the wait events
 * Tweak validation error message for min and max tags when field isn't a slice

v0.144.3
----------
 * Fix tests broken by new scheme addition

v0.144.2
----------
 * Update to latest gocommon to get Instagram scheme type

v0.144.1
----------
 * Update to latest gocommon

v0.144.0
----------
 * Add extraction field to webhook_called events

v0.143.4
----------
 * Use WebhookCall.ResponseJSON for @webhook and @result.*.extra
 * Updated translations from Transifex

v0.143.3
----------
 * Add more options for customizing contact used by test.SessionBuilder

v0.143.2
----------
 * Add test.SessionBuilder to make it easier to build sessions for testing

v0.143.1
----------
 * Include node on segments, revert change to add segments without destinations

v0.143.0
----------
 * For random router results, input should be raw random number, value is the bucket
 * Include segments with no destination
 * Add operand and time to flows.Segment

v0.142.1
----------
 * Put back engine.NewSprint which mailroom uses for surveyor submissions

v0.142.0
----------
 * Add Segments() to Sprint which returns all complete segments in that sprint
 * Simplify error message that users see if they have label action with no input

v0.141.0
----------
 * Rework Context into Scope, expose functions via a root scope, and support shadowing
 * Cleanup function exposure in contexts and add more tests
 * Drop unused child.run.* and parent.run.* parts of the context except .status as used subflow splits
 * Add support for anonymous functions in excellent
 * Fix function equality/inequality for consistency
 * Use syntax tree for refactoring operations
 * Excellent evaluation should parse to syntax tree as first step
 * Let Excellent functions know their own name to make better error messages
 * Only msg resumes should set input, all other resumes clear it

v0.140.1
----------
 * Update locale files (adds empty cs and mn translations)

v0.140.0
----------
 * Add reverse excellent function

v0.139.1
----------
 * Limit webhook URLs to 2048 chars

v0.139.0
----------
 * Add sort() excellent function

v0.138.0
----------
 * Add engine property for maximum resumes per session

v0.137.0
----------
 * Simplify contactql queries after parsing

v0.136.5
----------
 * Update to latest gocommon
 * Update locale files

v0.136.4
----------
 * Improve validator error message with startswith tag

v0.136.3
----------
 * Tweak validation to work when struct doesn't use json tags

v0.136.2
----------
 * Fix trigger docs

v0.136.1
----------
 * Add Session.FindStep

v0.136.0
----------
 * Rework WebhookCall and HTTPLog to overlap as much as possible

v0.135.0
----------
 * Get rid of ticket subjects

v0.134.3
----------
 * Add number of retries to webhook_called events

v0.134.2
----------
 * Fix word_slice when passing custom delimiters

v0.134.1
----------
 * Re-evaluate dynamic groups after opening tickets

v0.134.0
----------
 * Add support for variable user references to open ticket actions

v0.133.1
----------
 * Update to latest gocommon/phonenumbers

v0.133.0
----------
 * If open ticket action doesn't specify a topic, default to General
 * Expose topic instead of subject in context for ticket objects

v0.132.1
----------
 * Update to latest gocommon and add webchat URN schemes

v0.132.0
----------
 * Require either a subject or a topic to open a ticket but not both
 * Add assignee as optional field to open ticket actions
 * Add topics to tickets

v0.131.1
----------
 * Move slot param for LUIS classifiers to last
 * Add util cmd for testing classifier services

v0.131.0
----------
 * Update to LUIS API v3.0

v0.130.2
----------
 * Add support for tickets queries in elastic

v0.130.1
----------
 * Also simplify converting queries to elastic

v0.130.0
----------
 * Simplify parsing contact queries
 * Add support for query property tickets

v0.129.0
----------
 * Rework contactql to separate query parsing, validation and evaluation

v0.128.0
----------
 * Add unique excellent function to get unique elements in an array

v0.127.0
----------
 * Update to latest gocommon

v0.126.2
----------
 * Updated translations from Transifex
 * Replace usages of soon to be deprecated ioutil
 * Allow Msg type triggers to have connections

v0.126.1
----------
 * Update locale files

v0.126.0
----------
 * Use latest gocommon, replace all nulls and escaped nulls when parsing bodies as JSON

v0.125.2
----------
 * Add assignee (optional) to ticket_opened events
 * Cleanup from linter suggestions

v0.125.1
----------
 * Strip out invalid UTF-8 from webhook responses before trying to convert to JSON

v0.125.0
----------
 * Update to latest gocommon

v0.124.4
----------
 * Make users more like contacts in expressions by giving them an always non-empty default and a first_name property

v0.124.3
----------
 * Fix inspecting user dependencies in flows

v0.124.2
----------
 * Change default property of user objects in expressions to be the name to match contacts

v0.124.1
----------
 * Fix remove_first_word when input contains non-ASCII

v0.124.0
----------
 * Add User assets and use for Ticket.Assignee and Trigger.user

v0.123.0
----------
 * Add SUM() excellent function
 * Remove default_language from envs and usa first item of allowed_languages as the default

v0.122.0
----------
 * Allow build failing on codecov uploads again
 * Get rid of ticket references

v0.121.0
----------
 * Don't generate separate completion/functions doc files

v0.120.1
----------
 * Tweak test.AssertEqualJSON to take msgAndArgs param like asserts library

v0.120.0
----------
 * Add ticket as property to @trigger in context
 * Add new ticket type trigger with a closed event

v0.119.0
----------
 * Remove legacy_extra issue type

v0.118.2
----------
 * Add @ticket to the root context as the last opened ticket

v0.118.1
----------
 * Add way to create new ticket reference instances and add contact tickets to editor autocompletion

v0.118.0
----------
 * Add contact tickets to expression context

v0.117.0
----------
 * Add WA template translations namespace

v0.116.1
----------
 * Use standard hypenated BCP47 locale codes consistently

v0.116.0
----------
 * Test on go 1.16.x
 * Update to latest gocommon datefmt and pass locale to all date formatting calls

v0.115.2
----------
 * Fix dtone API endpoint URL and use external IDs

v0.115.1
----------
 * Fix sometimes retrying successful SMTP sends

v0.115.0
----------
 * Add support for retrying SMTP sends

v0.114.0
----------
 * Fully implement airtime service using new DT One API

v0.113.3
----------
 * Fix resuming a parent run when flow is missing

v0.113.2
----------
 * Add last missing translations for es and pt-BR

v0.113.1
----------
 * Update Spanish locale and gocommon dependency

v0.113.0
----------
 * Don't blow up building context if node is null
 * Add contact language to resthook payload

v0.112.2
----------
 * Add accessor for URN on ActivatedDialWait

v0.112.1
----------
 * Log error event and skip when attachment is longer than 2048 limit

v0.112.0
----------
 * Include resume and node in migration expression parsing
 * Add dial types of waits and resumes

v0.111.0
----------
 * Move to ElasticSearch v7 clients (backwards incompatible change)

v0.110.2
----------
 * Remove forward_ivr action and ivr_forwarded event

v0.110.1
----------
 * Spanish translations from transifex

v0.110.0
----------
 * Combine the completion.json and functions.json editor support files into a single editor.json file
 * Remove generated docs from repo

v0.109.4
----------
 * Fix release workflow

v0.109.3
----------
 * Fix release workflow

v0.109.2
----------
 * Don't use fuzzy entries in po files

v0.109.1
----------
 * Add forward_ivr action and ivr_forwarded event

v0.109.0
----------
 * Add @node.(uuid|visit_count) to context

v0.108.0
----------
 * Rename messaging_passive to messaging_background

v0.107.2
----------
 * Disallow labeling actions in passive flows

v0.107.1
----------
 * Add float64 workaround for exponent expressions with non-integer exponents

v0.107.0
----------
 * Add new flow type for passive messaging flows
 * Update to gocommon v1.7.1 to get fix for phone number parsing

v0.106.3
----------
 * Engine evaluator for contact sql should support != x for number and datetime values

v0.106.2
----------
 * Update to latest gocommon

v0.106.1
----------
 * URN and channel modifiers should error with invalid URNs and channels

v0.106.0
----------
 * Getting channel for URN should always consider the role on the channels
 * Update to latest gocommon which adds rocketchat scheme
 * SetPreferredChannel only when the channel has the send role

v0.105.5
----------
 * Support sorting contacts by last seen on attribute

v0.105.4
----------
 * Add support for Bengali numerals in number tests

v0.105.3
----------
 * Add support for Eastern Arabic numerals in number tests

v0.105.2
----------
 * Clear a run's expiration when it exits
 * Unwind accumulated run expirations as child runs complete
 * Include country in msg templating on msg_created events

v0.105.1
----------
 * Update to latest gocommon v1.5.3

v0.105.0
----------
 * Bump some deps, test on go 1.15 and fix bug found by 1.15 compiler

v0.104.1
----------
 * Update to gocommon v1.5.1

v0.104.0
----------
 * Use dummy value to avoid sending empty emails
 * Rework smtpx package for sending emails in places besides flows
 * Don't parse numbers in scientific notation

v0.103.1
----------
 * Update to latest gocommon v1.5.0
 * Run environment's DefaultLanguage and DefaultLocale methods should use contact language

v0.103.0
----------
 * Update to latest gocommon
 * Update terminology around groups with queries

v0.102.1
----------
 * Add archived contact status

v0.102.0
----------
 * Update to latest gocommon

v0.101.2
----------
 * Add empty localizations for all the languages used in RapidPro

v0.101.1
----------
 * Fix test

v0.101.0
----------
 * Use language codes (e.g. en-us) rather than locale names (en_US) for docs directories

v0.100.1
----------
 * Add completed pt_BR translation

v0.100.0
----------
 * Add last_seen_on to contacts and expose in expressions and queries

v0.99.0
----------
 * Rework elastic query generation so that all errors are caught at parsing stage
 * Allow URN inequality in elastic searches

v0.98.0
----------
 * Rework error handling in contactql so more errors are caught during parsing and have associated codes

v0.97.0
----------
 * Re-add classifier_called events for backward compatibility
 * Groups modifier should generate error if asked to operate on blocked or stopped contact
 * Move modifiers package out of actions package
 * ContactQL parser errors should contain more info

v0.96.0
----------
 * Reorganize validation code so utils doesn't have to know about tags defined higher up
 * Clone the test session during doc generation so actions always start with the same session
 * Add action to change contact status
 * Add historical information to triggers about the session that triggered them and use to prevent looping

v0.95.1
----------
 * Improve documentation of call_webhook action

v0.95.0
----------
 * Use latest wit.ai API version
 * Allow searching with values containing single quotes
 * Add user and origin fields to manual triggers
 * Add builder for triggers
 * Pass language to bothub API calls

v0.94.2
----------
 * Use jsonx.Marshal consistently

v0.94.1
----------
 * Add IsQueryError util

v0.94.0
----------
 * Move all location stuff from utils to envs
 * Simplify resolving locations from environments
 * Refactor field modifiers to take raw values and location parsing to not require a session

v0.93.1
----------
 * Fix clearing all URNs

v0.93.0
----------
 * Add urns modifier to replace all the URNs on a contact

v0.92.0
----------
 * Move elastic functionality from mailroom

v0.91.1
----------
 * Fix clearing of fields

v0.91.0
----------
 * Move generic PO stuff into utils/i18n

v0.90.0
----------
 * Allow querying contacts by UUID
 * Move i18n package under flows to avoid confusion with locales package
 * Add completion to localized documentation

v0.89.0
----------
 * Tweak change language functionality to allow missing translations
 * Add country to template translations and use when resolving templates

v0.88.0
----------
 * Add support for localized documentation

v0.87.0
----------
 * Disallow opening tickets, starting sessions and sending broadcasts when doing batch start
 * Add ability to change the language of a flow
 * Update our format_datetime docs to properly show range of 01-24
 * Fix evaluation of legacy vars in other-contacts actions

v0.86.2
----------
 * Fix spelling of Readact

v0.86.1
----------
 * Do redaction of access keys from HTTP logs

v0.86.0
----------
 * Add open_ticket actions and ticket_opened events

v0.85.0
----------
 * Add new service_called event to be used for classifiers and ticketers etc

v0.84.0
----------
 * Replace contact blocked and stopped fields with status field
 * Rename blocked and stopped modifiers to contact status modifier

v0.83.1
----------
 * Fix anywhere we truncate strings to do it by rune

v0.83.0
----------
 * Add blocked and stopped modifiers and events
 * Add blocked and stopped fields to contact

v0.82.0
----------
 * Fix default to understand objects with defaults

v0.81.0
----------
 * Rework httpx to replace NewTrace with NewRequest+DoTrace
 * Separate out the header part of response traces from the body which won't always be valid UTF-8

v0.80.0
----------
 * ivr_created events should include language of translated text

v0.79.1
----------
 * Include 3-char language code as extra header in PO files

v0.79.0
----------
 * Add custom Source-Flows header to exported PO files
 * Make router categories inspectable
 * Importing of translations into flows

v0.78.1
----------
 * Add decode_html Excellent function
 * Start of i18n work
 * Prevent XText.Slice from panicking

v0.78.0
----------
 * Add support for extracting the "base" translation of a flow
 * Allow queries on URNs to check if they are set or not
 * Add Language.ToISO639_2()
 * Make flowrunner easier to use by defaulting to first flow in the assets
 * Default to current version in flowmigrate cmd
 * Rework group asset loading so that parsing is not deferred
 * Override environment country if contact has preferred channel with country

v0.77.4
----------
 * Fix loading flow assets that are new spec but also have metadata section

v0.77.3
----------
 * Update README

v0.77.2
----------
 * Update README

v0.77.1
----------
 * Fix not passing access config correctly to webhook services

v0.77.0
----------
 * Allow http services to be configured with a list of disallowed hosts

v0.76.3
----------
 * fix @legacy_extra issue on routers
 * update gomobile instructions

v0.76.2
----------
 * Add trim, trim_left and trim_right excellent functions

v0.76.1
----------
 * Sort issues by node order
 * Add issues to report on invalid regex and usage of @legacy_extra

v0.76.0
----------
 * Validate language codes in contact queries
 * Disallow group queries against group names that don't exist

v0.75.1
----------
 * Remove contacts from broken groups

v0.75.0
----------
 * Handle missing groups on contact creation
 * Fix != with multiple values and add support for group attribute in contact queries
 * Improve docs for operators

v0.74.0
----------
 * Use jsonx functions for all JSON marshal/unmarshal
 * Add support for removing a URN to the urns modifier
 * Groups modifier should log errors for dynamic groups
 * Rename snapshot flag to -update

v0.73.0
----------
 * Include translation language with missing dependency issues

v0.72.2
----------
 * Allow flow inspection without assets

v0.72.1
----------
 * Quote telephone numbers in contact queries
 * Tweak parsing of phone numbers in contact queries

v0.72.0
----------
 * Rename "problems" to "issues"

v0.71.3
----------
 * Implement missing_dependency as a type of problem
 * Add framework for checking for problems during flow inspection

v0.71.2
----------
 * Rework dependency and template extraction to include actions and routers

v0.71.1
----------
 * Channels with no country are implicitly international

v0.71.0
----------
 * Add field to channel assets which determines whether they should try to send internationally

v0.70.0
----------
 * Make cloning a flow definition more deterministic
 * Update actions to log error events when dependencies are missing
 * Interpret contact queries which are formatted phone numbers as tel = queries

v0.69.0
----------
 * Move JSON utils into their own package
 * Track node UUIDs of dependencies

v0.68.0
----------
 * Convert dependency inspection output to list of things with type attribute
 * Replace Flow.CheckDepedencies and CheckDependenciesRecursive with passing assets to Inspect

v0.67.1
----------
 * Update to gocommon v1.2.0

v0.67.0
----------
 * Rename Flow.Validate to Flow.CheckDependencies for clarity
 * Create error event when webhook response too big
 * Rework webhook calls to use same calling code as other HTTP services

v0.66.3
----------
 * Allow globals with empty values

v0.66.2
----------
 * Add mobile binding for IsVersionSupported
 * Re-add version check to ReadFlow

v0.66.1
----------
 * Match evaluation of contact queries in ES

v0.66.0
----------
 * Fix problems with contact searching and add support for URN as attribute

v0.65.0
----------
 * Ignore content-type headers and try to parse all webhook responses as JSON
 * Update ContactQL to interpret implicit conditions which are URNs as scheme=path

v0.64.11
----------
 * Limit the size of evaluated templates and truncate anything bigger

v0.64.10
----------
 * Stringify contactql like the queries they came from
 * Trim webhook_called request traces to 10K same as response traces
 * Only set extra on webhook result if less than 10000 bytes

v0.64.9
----------
 * Allow getting current context even for ended sessions

v0.64.8
----------
 * Fix another panic during context walking

v0.64.7
----------
 * Fix panic in context walking

v0.64.6
----------
 * Add support for marshaling XObjects with their defaults, and tool for walking the context to find objects

v0.64.5
----------
 * Fix creation of no-nil interface to nil structs in context

v0.64.4
----------
 * Make it easier to get current expression context of a waiting session

v0.64.3
----------
 * Allow webhook calls with GET method to have bodies

v0.64.2
----------
 * Include parent result references in flow inspection

v0.64.1
----------
 * Add support for jitter in webhook retries

v0.64.0
----------
 * Make http retrying available to all services which use HTTP
 * Fix parsing out relative date value during migration of date tests

v0.63.1
----------
 * Perform URL validation in call_webhook and skip action appropriately

v0.63.0
----------
 * Loosen email regex used by has_email test
 * Allow cloning of JSON flow definitions not tied to any spec version

v0.62.0
----------
 * Render email_sent events in flowrunner
 * Allow flowmigrate to take a target version argument
 * Implement 13.1 migration as adding UUID to semd_msg.templating

v0.61.0
----------
 * Implement email as a service

v0.60.1
----------
 * Fix re-evaluating dynamic groups when query references non-existent field

v0.60.0
----------
 * Add @globals to completion
 * Add topic to send_msg actions

v0.59.0
----------
 * Validate run summary JSON passed to flow action triggers
 * Expose keyword match on trigger in context

v0.58.0
----------
 * Give services their own HTTP clients
 * Allow webhook service to take a map of deafult header values

v0.57.0
----------
 * Add globals to evaluation context as @globals
 * Add global as new asset type

v0.56.3
----------
 * DTOne client should record http log for timeouts

v0.56.2
----------
 * Fix migrating save actions with URN fields

v0.56.1
----------
 * Tweak criteria for deciding whether to try reading a flow as legacy

v0.56.0
----------
 * Fix docstring for UPPER()
 * Rework ReadFlow to accept legacy flows too
 * Move legacy package inside flows/definition

v0.55.0
----------
 * Update start_session action to use escaping when evaluating the contact query
 * Add support for escaping expressions in templates

v0.54.3
----------
 * Relax requirement for field assets to have UUID set since engine doesn't use this

v0.54.2
----------
 * Fix naming in mobile bindings

v0.54.1
----------
 * Fix docstring for UPPER()

v0.54.0
----------
 * NewEnvironmentBuilder() -> envs.NewBuilder() to match engine.NewBuilder()
 * Include classifiers in flow dependency inspection

v0.53.1
----------
 * Add classification service for Bothub

v0.53.0
----------
 * Record arrays of http logs on classifier_called and airtime_transferred events

v0.52.1
----------
 * Modify grammar to allow result names that start with underscores

v0.52.0
----------
 * All service factory methods should return an error if service can't be returned
 * Rework airtime transfer nodes to function more like NLU nodes
 * Add classification service implementation for LUIS

v0.51.0
----------
 * Add NLU support: a classify action, a classification service and various router tests

v0.50.4
----------
 * Revert change to operands for media waits

v0.50.3
----------
 * Fix migrating operands on rulesets waiting for media
 * Change autocompletion type of related_run.results to any since we can't autocomplete it

v0.50.2
----------
 * Fix migration of localization when flow has unused base translations

v0.50.1
----------
 * Have a single HTTPClient on the engine instead of every service having its own

v0.50.0
----------
 * Fix formatting runsummary with missing flow
 * Add contact_query field to start_session actions
 * Rework services so they take a session and resolve to a provider that does the work

v0.49.0
----------
 * Rework webhook calling code as a service and fix not saving result when connection errors

v0.48.2
----------
 * Include sender and recipient in airtime events

v0.48.1
----------
 * Add .Source() to SessionAssets interface

v0.48.0
----------
 * Unexport things that no longer need to be exported now that we've ditched extensions, clean up names of typed things
 * Remove transferto extension functionality and instead have standard transfer_airtime action which defers to an airtime service

v0.47.3
----------
 * completions.json should include section for session-less contexts

v0.47.2
----------
 * Add FlowReference to FlowRun interface and add some more tests

v0.47.1
----------
 * Renamed errored statuses to failed, replace fatal error events with failure events

v0.47.0
----------
 * Allow loading of runs with missing flows
 * A terminal enter_flow action should leave existing runs as completed instead of interrupted
 * Make documented item titles into actual links so it's easier to get the link of a particular item in the docs

v0.46.0
----------
 * Add UUID to assets.Field

v0.45.2
----------
 * Fix parsing context references like foo.0

v0.45.1
----------
 * ContactSQL query parsing should error if URN schenme used when URN redaction is enabled, and validate fields

v0.45.0
----------
 * urn_parts should error for non-URNs and so Wrap migrated urn_parts expressions with default to catch errors
 * Migrate non-tel URN types using urn_parts(..).path
 * Redacted URNs should still have scheme, and format_urn should work for redacted URNs

v0.44.4
----------
 * Set redaction policy in visitor constructor for contactql

v0.44.3
----------
 * Fix parsing of implicit conditions in contactql

v0.44.2
----------
 * Add UUID() to Session interface

v0.44.1
----------
 * Make trigger.params null for trigger types that don't use it, non-null for those that do

v0.44.0
----------
 * Add UUID field to sessions
 * Rework trigger.params to be an XObject and always non-null in expressions
 * Implement a week_number function which matches Excel's WEEKNUM

v0.43.2
----------
 * rename voice trigger to be more consistent

v0.43.1
----------
 * add ivr flow trigger constructor

v0.43.0
----------
 * Allow array lookups like foo.0
 * More re-organization of utils code into smaller packages

v0.42.0
----------
 * Move Environment type and environment based date parsing to new envs package
 * Move Date and TimeOfDay types to new dates package
 * Do template and dependency enumeration by reflection

v0.41.18
----------
 * Drop current template rewriting functionality which isn't used and can't be used with migrations
 * Generate context map from docstrings

v0.41.17
----------
 * Add SetURN function to Msg
 * Fix localization UUID in test action holder flow
 * Update send_email action to allow localization of subject and body
 * Reorganize docgen code to make it easier to add new doc outputs

v0.41.16
----------
 * Allow setting channel on a non-tel URN if it doesn't have a channel
 * Deprecate parent.run and child.run in the context and move those fields up one level

v0.41.15
----------
 * Tweak to goflow interfaces to allow introspection into contactql
 * Include external ID of msg on input expression context

v0.41.14
----------
 * Fix start_session when create_contact is true

v0.41.13
----------
 * Fix index out of bounds panic when transation for item exists but has less strings than original

v0.41.12
----------
 * Fix format excellent function when passed a nil
 * Fix parsing of 12AM and 12PM times

v0.41.11
----------
 * parse_json should error for invalid JSON
 * Switch to faster json.Valid for checking JSON validity

v0.41.10
----------
 * Add foreach_value function to allow us to keep legacy webhook payloads the same
 * Make @results and @run.results the same

v0.41.9
----------
 * NOOP if there are no rceipients for start_session and send_broadcast
 * Update send_broadcast and start_session to accept a URN in legacy_vars

v0.41.8
----------
 * Use std lib function to check HTTP headers

v0.41.7
----------
 * Remapping UUIDs during cloning must include UUIDs which are values in arrays

v0.41.6
----------
 * has_group can take optional second parameter which is group name.. used only in dependency inspection

v0.41.5
----------
 * Fix cloning of UI sections

v0.41.4
----------
 * Check that numbers are actually valid in our has_phone test

v0.41.3
----------
 * Handle missing ruleset types in legacy flows

v0.41.2
----------
 * handle no media type in our migration

v0.41.1
----------
 * Include node UUIDs in result infos returned from flow inspection

v0.41.0
----------
 * Legacy flow migration should just ignore invalid actionset/rule destinations
 * Drop includeUI as an option for flow migration and just always include it

v0.40.3
----------
 * Re-organize test utils so they're all in the test package

v0.40.2
----------
 * Add FixedUUID4Generator for testing

v0.40.1
----------
 * Make recursion optional again during flow validation

v0.40.0
----------
 * Add temporary Flow.MarshalWithInfo to aid with moving mailroom to new endpoint
 * Split up flow inspection and dependency validation and don't embed inspection results in the flow definition

v0.39.4
----------
 * Add support for cloning flows using generic JSON representations

v0.39.3
----------
 * Handle malformed single message campaign event flows

v0.39.2
----------
 * Fix IsLegacyDefinition

v0.39.1
----------
 * (Re)allow inspecting without session assets

v0.39.0
----------
 * Allow igoring of missing flow assets like any other asset type
 * Simplify the expression used to emulate legacy webhook payloads
 * Do structural validation in ReadFlow so flow returned from that is always valid

v0.38.3
----------
 * Switch to new library for UUID generation
 * Expose current flow spec version in mobile bindings

v0.38.2
----------
 * Fix send_broadcast and start_session actions so telephone numbers are normalized

v0.38.1
----------
 * properly return template dependencies for flows

v0.38.0
----------
 * Rework creating and starting sessions so sessions no longer exist in limbo state between the two
 * Add @webhook as shortcut to .extra of last webhook result
 * Make URLJoin honor absolute urls

v0.37.3
----------
 * Omit empty variables in message templates
 * Router failing to pick category should be fatal error event - not hard error
 * Match number like .5 and 1OO (o's) as 1)
 * Don't include quick replies and attachments whih both error and evaluate to empty

v0.37.2
----------
 * fix text_slice for unicode strings

v0.37.1
----------
 * @fields defaults to table like @results
 * Add default values to @run, @parent, @parent.run, @child, @child.run

v0.37.0
----------
 * Add TemplateIncluder to ease inclusion of templates as strings, slices or maps
 * Remove left() and right(), replace with text_slice()
 * Re-add defaults for several context objects
 * Validate that switch case tests are registered XTESTs
 * Allow passing count of replacements to replace()
 * Allow empty results, fix empty category caltulations

v0.36.2
----------
 * accept text/javascript as a content type in webhooks

v0.36.1
----------
 * better error message for not being able to resume

v0.36.0
----------
 * Bump current flow spec version to 13

v0.35.0
----------
 * Migrate date_ tests to set delta value in UI config
 * is_text_eq -> has_only_text
 * @child and @parent should mirror root of context, not @run

v0.34.1
----------
 * Add format_results and migrate @flow to @(format_results(results))
 * Add type-aware format() function 

v0.34.0
----------
 * dict/keys -> object/properties
 * Split length() into count() and text_length()
 * Support != operator in contact queries
 * Extract operators into their own functions which can then be documented and have live examples

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
