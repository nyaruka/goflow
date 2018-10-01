# Triggers

Triggers start a new session with the flow engine. They describe why the session is being started and provide parameters which can
be accessed in expressions.

<div class="triggers">
{{ .triggerDocs }}
</div>

# Resumes

Resumes resume an existing session with the flow engine and describe why the session is being resumed.

<div class="resumes">
{{ .resumeDocs }}
</div>

# Events

Events are the output of a flow run and represent instructions to the engine container on what actions should be taken due to the flow execution.
All templates in events have been evaluated and can be used to create concrete messages, contact updates, emails etc by the container.

<div class="events">
{{ .eventDocs }}
</div>
