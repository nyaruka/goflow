package legacy

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition/legacy/expressions"

	"github.com/shopspring/decimal"
)

// template that matches the JSON payload sent by legacy webhooks
const legacyWebhookPayload = `@(json(object(
  "contact", object("uuid", contact.uuid, "name", contact.name, "urn", contact.urn),
  "flow", run.flow,
  "path", run.path,
  "results", foreach_value(results, extract_object, "category", "category_localized", "created_on", "input", "name", "node_uuid", "value"),
  "run", object("uuid", run.uuid, "created_on", run.created_on),
  "input", if(
    input,
    object(
      "attachments", foreach(input.attachments, attachment_parts),
      "channel", input.channel,
      "created_on", input.created_on,
      "text", input.text,
      "type", input.type,
      "urn", if(
        input.urn,
        object(
          "display", default(format_urn(input.urn), ""),
          "path", urn_parts(input.urn).path,
          "scheme", urn_parts(input.urn).scheme
        ),
        null
      ),
      "uuid", input.uuid
    ),
    null
  ),
  "channel", default(input.channel, null)
)))`

type migratedExit map[string]interface{}

func newExit(uuid uuids.UUID, destinationUUID uuids.UUID) migratedExit {
	d := map[string]interface{}{"uuid": uuid}
	if destinationUUID != "" {
		d["destination_uuid"] = destinationUUID
	}

	return migratedExit(d)
}

func (e migratedExit) UUID() uuids.UUID {
	return e["uuid"].(uuids.UUID)
}

type migratedNode map[string]interface{}

func newNode(uuid uuids.UUID, actions []migratedAction, router migratedRouter, exits []migratedExit) migratedNode {
	d := map[string]interface{}{
		"uuid":  uuid,
		"exits": exits,
	}
	if len(actions) > 0 {
		d["actions"] = actions
	}
	if router != nil {
		d["router"] = router
	}

	return migratedNode(d)
}

func (n migratedNode) UUID() uuids.UUID {
	return n["uuid"].(uuids.UUID)
}

type migratedLocalization map[envs.Language]map[uuids.UUID]map[string][]string

func (l migratedLocalization) addTranslation(lang envs.Language, itemUUID uuids.UUID, property string, translated []string) {
	_, found := l[lang]
	if !found {
		l[lang] = make(map[uuids.UUID]map[string][]string)
	}
	langTranslations := l[lang]

	_, found = langTranslations[itemUUID]
	if !found {
		langTranslations[itemUUID] = make(map[string][]string)
	}

	langTranslations[itemUUID][property] = translated
}

func (l migratedLocalization) addTranslationMap(baseLanguage envs.Language, mapped Translations, uuid uuids.UUID, property string) string {
	var inBaseLanguage string
	for language, item := range mapped {
		expression, _ := expressions.MigrateTemplate(item, nil)
		if language == baseLanguage {
			inBaseLanguage = expression
		} else if language != "base" {
			l.addTranslation(language, uuid, property, []string{expression})
		}
	}

	return inBaseLanguage
}

func (l migratedLocalization) addTranslationMultiMap(baseLanguage envs.Language, mapped map[envs.Language][]string, uuid uuids.UUID, property string) []string {
	var inBaseLanguage []string
	for language, items := range mapped {
		templates := make([]string, len(items))
		for i := range items {
			expression, _ := expressions.MigrateTemplate(items[i], nil)
			templates[i] = expression
		}
		if language != baseLanguage {
			l.addTranslation(language, uuid, property, templates)
		} else {
			inBaseLanguage = templates
		}
	}
	return inBaseLanguage
}

type migratedCategory map[string]interface{}

func newCategory(uuid uuids.UUID, name string, exitUUID uuids.UUID) migratedCategory {
	d := map[string]interface{}{"uuid": uuid}
	if name != "" {
		d["name"] = name
	}
	if exitUUID != "" {
		d["exit_uuid"] = exitUUID
	}

	return migratedCategory(d)
}

func (c migratedCategory) UUID() uuids.UUID {
	return c["uuid"].(uuids.UUID)
}

type migratedCase map[string]interface{}

func newCase(uuid uuids.UUID, type_ string, arguments []string, categoryUUID uuids.UUID) migratedCase {
	d := map[string]interface{}{
		"uuid":          uuid,
		"type":          type_,
		"category_uuid": categoryUUID,
	}
	if len(arguments) > 0 {
		d["arguments"] = arguments
	}

	return migratedCase(d)
}

func (c migratedCase) UUID() uuids.UUID {
	return c["uuid"].(uuids.UUID)
}

type migratedTimeout map[string]interface{}

func newTimeout(seconds int, categoryUUID uuids.UUID) migratedTimeout {
	return migratedTimeout(map[string]interface{}{"seconds": seconds, "category_uuid": categoryUUID})
}

type migratedHint map[string]interface{}

func newAudioHint() migratedHint {
	return migratedHint(map[string]interface{}{"type": "audio"})
}

func newImageHint() migratedHint {
	return migratedHint(map[string]interface{}{"type": "image"})
}

func newVideoHint() migratedHint {
	return migratedHint(map[string]interface{}{"type": "video"})
}

func newLocationHint() migratedHint {
	return migratedHint(map[string]interface{}{"type": "location"})
}

func newFixedDigitsHint(count int) migratedHint {
	return migratedHint(map[string]interface{}{"type": "digits", "count": count})
}

func newTerminatedDigitsHint(terminatedBy string) migratedHint {
	return migratedHint(map[string]interface{}{"type": "digits", "terminated_by": terminatedBy})
}

type migratedWait map[string]interface{}

func newMsgWait(timeout migratedTimeout, hint migratedHint) migratedWait {
	d := map[string]interface{}{"type": "msg"}
	if timeout != nil {
		d["timeout"] = timeout
	}
	if hint != nil {
		d["hint"] = hint
	}

	return migratedWait(d)
}

type migratedRouter map[string]interface{}

func newSwitchRouter(wait migratedWait, resultName string, categories []migratedCategory, operand string, cases []migratedCase, defaultCategory uuids.UUID) migratedRouter {
	d := map[string]interface{}{
		"type":                  "switch",
		"categories":            categories,
		"operand":               operand,
		"cases":                 cases,
		"default_category_uuid": defaultCategory,
	}
	if wait != nil {
		d["wait"] = wait
	}
	if resultName != "" {
		d["result_name"] = resultName
	}

	return migratedRouter(d)
}

func newRandomRouter(resultName string, categories []migratedCategory) migratedRouter {
	d := map[string]interface{}{
		"type":       "random",
		"categories": categories,
	}
	if resultName != "" {
		d["result_name"] = resultName
	}

	return migratedRouter(d)
}

type migratedAction map[string]interface{}

func (a migratedAction) UUID() uuids.UUID {
	return a["uuid"].(uuids.UUID)
}

func newAddContactGroupsAction(uuid uuids.UUID, groups []*assets.GroupReference) migratedAction {
	return migratedAction(map[string]interface{}{
		"uuid":   uuid,
		"type":   "add_contact_groups",
		"groups": groups,
	})
}

func newAddContactURNAction(uuid uuids.UUID, scheme string, path string) migratedAction {
	return migratedAction(map[string]interface{}{
		"uuid":   uuid,
		"type":   "add_contact_urn",
		"scheme": scheme,
		"path":   path,
	})
}

func newAddInputLabelsAction(uuid uuids.UUID, labels []*assets.LabelReference) migratedAction {
	return migratedAction(map[string]interface{}{
		"uuid":   uuid,
		"type":   "add_input_labels",
		"labels": labels,
	})
}

func newCallResthookAction(uuid uuids.UUID, resthook string, resultName string) migratedAction {
	d := map[string]interface{}{
		"uuid":     uuid,
		"type":     "call_resthook",
		"resthook": resthook,
	}
	if resultName != "" {
		d["result_name"] = resultName
	}

	return migratedAction(d)
}

func newCallWebhookAction(uuid uuids.UUID, method string, url string, headers map[string]string, body string, resultName string) migratedAction {
	d := map[string]interface{}{
		"uuid":   uuid,
		"type":   "call_webhook",
		"method": method,
		"url":    url,
	}
	if len(headers) > 0 {
		d["headers"] = headers
	}
	if body != "" {
		d["body"] = body
	}
	if resultName != "" {
		d["result_name"] = resultName
	}

	return migratedAction(d)
}

func newEnterFlowAction(uuid uuids.UUID, flow *assets.FlowReference, terminal bool) migratedAction {
	d := map[string]interface{}{
		"uuid": uuid,
		"type": "enter_flow",
		"flow": flow,
	}
	if terminal {
		d["terminal"] = terminal
	}

	return migratedAction(d)
}

func newPlayAudioAction(uuid uuids.UUID, audioURL string) migratedAction {
	return migratedAction(map[string]interface{}{
		"uuid":      uuid,
		"type":      "play_audio",
		"audio_url": audioURL,
	})
}

func newRemoveContactGroupsAction(uuid uuids.UUID, groups []*assets.GroupReference, allGroups bool) migratedAction {
	d := map[string]interface{}{
		"uuid": uuid,
		"type": "remove_contact_groups",
	}
	if len(groups) > 0 {
		d["groups"] = groups
	}
	if allGroups {
		d["all_groups"] = allGroups
	}

	return migratedAction(d)
}

func newSayMsgAction(uuid uuids.UUID, text string, audioURL string) migratedAction {
	d := map[string]interface{}{
		"uuid": uuid,
		"type": "say_msg",
		"text": text,
	}
	if audioURL != "" {
		d["audio_url"] = audioURL
	}

	return migratedAction(d)
}

func newSendBroadcastAction(uuid uuids.UUID, text string, attachments []string, quickReplies []string, urns []urns.URN, contacts []*flows.ContactReference, groups []*assets.GroupReference, legacyVars []string) migratedAction {
	d := map[string]interface{}{
		"uuid": uuid,
		"type": "send_broadcast",
		"text": text,
	}
	if len(attachments) > 0 {
		d["attachments"] = attachments
	}
	if len(quickReplies) > 0 {
		d["quick_replies"] = quickReplies
	}
	if len(urns) > 0 {
		d["urns"] = urns
	}
	if len(contacts) > 0 {
		d["contacts"] = contacts
	}
	if len(groups) > 0 {
		d["groups"] = groups
	}
	if len(legacyVars) > 0 {
		d["legacy_vars"] = legacyVars
	}

	return migratedAction(d)
}

func newSendEmailAction(uuid uuids.UUID, addresses []string, subject string, body string) migratedAction {
	return migratedAction(map[string]interface{}{
		"uuid":      uuid,
		"type":      "send_email",
		"addresses": addresses,
		"subject":   subject,
		"body":      body,
	})
}

func newSendMsgAction(uuid uuids.UUID, text string, attachments []string, quickReplies []string, allURNs bool) migratedAction {
	d := map[string]interface{}{
		"uuid": uuid,
		"type": "send_msg",
		"text": text,
	}
	if len(attachments) > 0 {
		d["attachments"] = attachments
	}
	if len(quickReplies) > 0 {
		d["quick_replies"] = quickReplies
	}
	if allURNs {
		d["all_urns"] = allURNs
	}

	return migratedAction(d)
}

func newSetContactChannelAction(uuid uuids.UUID, channel *assets.ChannelReference) migratedAction {
	return migratedAction(map[string]interface{}{
		"uuid":    uuid,
		"type":    "set_contact_channel",
		"channel": channel,
	})
}

func newSetContactFieldAction(uuid uuids.UUID, field *assets.FieldReference, value string) migratedAction {
	return migratedAction(map[string]interface{}{
		"uuid":  uuid,
		"type":  "set_contact_field",
		"field": field,
		"value": value,
	})
}

func newSetContactLanguageAction(uuid uuids.UUID, language string) migratedAction {
	return migratedAction(map[string]interface{}{
		"uuid":     uuid,
		"type":     "set_contact_language",
		"language": language,
	})
}

func newSetContactNameAction(uuid uuids.UUID, name string) migratedAction {
	return migratedAction(map[string]interface{}{
		"uuid": uuid,
		"type": "set_contact_name",
		"name": name,
	})
}

func newStartSessionAction(uuid uuids.UUID, flow *assets.FlowReference, urns []urns.URN, contacts []*flows.ContactReference, groups []*assets.GroupReference, legacyVars []string, createContact bool) migratedAction {
	d := map[string]interface{}{
		"uuid": uuid,
		"type": "start_session",
		"flow": flow,
	}
	if len(urns) > 0 {
		d["urns"] = urns
	}
	if len(contacts) > 0 {
		d["contacts"] = contacts
	}
	if len(groups) > 0 {
		d["groups"] = groups
	}
	if len(legacyVars) > 0 {
		d["legacy_vars"] = legacyVars
	}
	if createContact {
		d["create_contact"] = createContact
	}

	return migratedAction(d)
}

func newTransferAirtimeAction(uuid uuids.UUID, amounts map[string]decimal.Decimal, resultName string) migratedAction {
	d := map[string]interface{}{
		"uuid":    uuid,
		"type":    "transfer_airtime",
		"amounts": amounts,
	}
	if resultName != "" {
		d["result_name"] = resultName
	}

	return migratedAction(d)
}
