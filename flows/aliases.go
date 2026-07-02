package flows

import (
	"github.com/nyaruka/goflow/events"
)

// The types below moved to the events package so that it can be used without depending on this package.
// They are aliased here because they're as much a part of the flows API as they are the events wire format.

type Event = events.Event
type EventUUID = events.EventUUID
type EventLogger = events.EventLogger

type NodeUUID = events.NodeUUID
type ExitUUID = events.ExitUUID
type StepUUID = events.StepUUID
type InputUUID = events.InputUUID

type RunUUID = events.RunUUID
type RunStatus = events.RunStatus

type ContactUUID = events.ContactUUID
type ContactStatus = events.ContactStatus
type ContactReference = events.ContactReference

type SessionUUID = events.SessionUUID
type SessionHistory = events.SessionHistory

type TicketUUID = events.TicketUUID
type TicketStatus = events.TicketStatus
type TicketEnvelope = events.TicketEnvelope

type CallUUID = events.CallUUID
type CallEnvelope = events.CallEnvelope
type CallStatus = events.CallStatus

type Dial = events.Dial
type DialStatus = events.DialStatus

type HTTPLog = events.HTTPLog
type HTTPLogWithoutTime = events.HTTPLogWithoutTime
type HTTPLogCallback = events.HTTPLogCallback
type HTTPLogger = events.HTTPLogger
type HTTPLogStatusResolver = events.HTTPLogStatusResolver

type LLMResponse = events.LLMResponse
type AirtimeTransfer = events.AirtimeTransfer

type Result = events.Result
type Value = events.Value

type Hint = events.Hint

type BaseMsg = events.BaseMsg
type MsgIn = events.MsgIn
type MsgOut = events.MsgOut
type MsgContent = events.MsgContent
type MsgTemplating = events.MsgTemplating
type TemplatingComponent = events.TemplatingComponent
type TemplatingVariable = events.TemplatingVariable
type QuickReply = events.QuickReply
type UnsendableReason = events.UnsendableReason

type BroadcastUUID = events.BroadcastUUID
type BroadcastTranslations = events.BroadcastTranslations

const (
	RunStatusActive      = events.RunStatusActive
	RunStatusCompleted   = events.RunStatusCompleted
	RunStatusWaiting     = events.RunStatusWaiting
	RunStatusFailed      = events.RunStatusFailed
	RunStatusExpired     = events.RunStatusExpired
	RunStatusInterrupted = events.RunStatusInterrupted

	ContactStatusActive   = events.ContactStatusActive
	ContactStatusBlocked  = events.ContactStatusBlocked
	ContactStatusStopped  = events.ContactStatusStopped
	ContactStatusArchived = events.ContactStatusArchived

	TicketStatusOpen   = events.TicketStatusOpen
	TicketStatusClosed = events.TicketStatusClosed

	CallStatusSuccess         = events.CallStatusSuccess
	CallStatusConnectionError = events.CallStatusConnectionError
	CallStatusResponseError   = events.CallStatusResponseError
	CallStatusSubscriberGone  = events.CallStatusSubscriberGone

	DialStatusAnswered = events.DialStatusAnswered
	DialStatusNoAnswer = events.DialStatusNoAnswer
	DialStatusBusy     = events.DialStatusBusy
	DialStatusFailed   = events.DialStatusFailed

	UnsendableReasonNoRoute         = events.UnsendableReasonNoRoute
	UnsendableReasonContactBlocked  = events.UnsendableReasonContactBlocked
	UnsendableReasonContactStopped  = events.UnsendableReasonContactStopped
	UnsendableReasonContactArchived = events.UnsendableReasonContactArchived

	MaxAttachmentLength      = events.MaxAttachmentLength
	MaxQuickReplyTextLength  = events.MaxQuickReplyTextLength
	MaxQuickReplyExtraLength = events.MaxQuickReplyExtraLength

	RedactionMask = events.RedactionMask
)

var (
	NewEventUUID        = events.NewEventUUID
	NewRunUUID          = events.NewRunUUID
	NewContactUUID      = events.NewContactUUID
	NewSessionUUID      = events.NewSessionUUID
	NewTicketUUID       = events.NewTicketUUID
	NewCallUUID         = events.NewCallUUID
	NewBroadcastUUID    = events.NewBroadcastUUID
	NewContactReference = events.NewContactReference

	NewDial               = events.NewDial
	NewHTTPLog            = events.NewHTTPLog
	NewHTTPLogWithoutTime = events.NewHTTPLogWithoutTime
	HTTPStatusFromCode    = events.HTTPStatusFromCode

	NewResult = events.NewResult
	NewValue  = events.NewValue

	NewMsgIn         = events.NewMsgIn
	NewMsgOut        = events.NewMsgOut
	NewIVRMsgOut     = events.NewIVRMsgOut
	NewMsgTemplating = events.NewMsgTemplating

	EmptyHistory = events.EmptyHistory
)
