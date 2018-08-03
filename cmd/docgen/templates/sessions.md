## Event Definitions

Events are the output of a flow run and represent instructions to the engine container on what actions should be taken due to the flow execution.
All templates in events have been evaluated and can be used to create concrete messages, contact updates, emails etc by the container.

<div class="events">
{{ .eventDocs }}
</div>

## Trigger Types

Triggers are the entities which can trigger a new session with the flow engine.

<div class="triggers">
{{ .triggerDocs }}
</div>
