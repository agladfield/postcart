package postmark

import (
	"errors"
	"fmt"
)

// Postmark error definitions
var (
	errPostmarkBadOrMissingAPIToken            = errors.New("bad or missing API token")
	errPostmarkMultipleErrorsOccurred          = errors.New("multiple errors occurred")
	errPostmarkResourceNotFound                = errors.New("resource not found")
	errPostmarkInvalidPaginationKey            = errors.New("invalid pagination key")
	errPostmarkMaintenance                     = errors.New("Postmark API is offline for maintenance")
	errPostmarkInvalidEmailRequest             = errors.New("invalid email request")
	errPostmarkSenderSignatureNotFound         = errors.New("sender signature not found")
	errPostmarkSenderSignatureNotConfirmed     = errors.New("sender signature not confirmed")
	errPostmarkInvalidJSON                     = errors.New("invalid JSON")
	errPostmarkInvalidRequestFields            = errors.New("invalid request fields")
	errPostmarkNotAllowedToSend                = errors.New("not allowed to send")
	errPostmarkInactiveRecipient               = errors.New("inactive recipient")
	errPostmarkJSONRequired                    = errors.New("JSON required")
	errPostmarkTooManyBatchMessages            = errors.New("too many batch messages")
	errPostmarkForbiddenAttachmentType         = errors.New("forbidden attachment type")
	errPostmarkAccountIsPending                = errors.New("account is pending")
	errPostmarkAccountMayNotSend               = errors.New("account may not send")
	errPostmarkSenderSignatureQueryException   = errors.New("sender signature query exception")
	errPostmarkSenderSignatureNotFoundByID     = errors.New("sender signature not found by id")
	errPostmarkNoUpdatedSenderSignatureData    = errors.New("no updated sender signature data received")
	errPostmarkCannotUsePublicDomain           = errors.New("cannot use a public domain")
	errPostmarkSenderSignatureAlreadyExists    = errors.New("sender signature already exists")
	errPostmarkDKIMAlreadyScheduledForRenewal  = errors.New("DKIM already scheduled for renewal")
	errPostmarkSenderSignatureAlreadyConfirmed = errors.New("sender signature already confirmed")
	errPostmarkDoNotOwnSenderSignature         = errors.New("you do not own this sender signature")
	errPostmarkDomainNotFound                  = errors.New("domain was not found")
	errPostmarkInvalidFieldsSupplied           = errors.New("invalid fields supplied")
	errPostmarkDomainAlreadyExists             = errors.New("domain already exists")
	errPostmarkDoNotOwnDomain                  = errors.New("you do not own this domain")
	errPostmarkNameRequiredForDomain           = errors.New("name is a required field to create a domain")
	errPostmarkNameTooLongForDomain            = errors.New("name field must be less than or equal to 255 characters")
	errPostmarkInvalidNameFormatForDomain      = errors.New("name format is invalid")
	errPostmarkMissingFieldForSenderSignature  = errors.New("missing a required field to create a sender signature")
	errPostmarkSenderSignatureFieldTooLong     = errors.New("a field in the sender signature request is too long")
	errPostmarkInvalidFieldValue               = errors.New("value for field is invalid")
	errPostmarkServerQueryException            = errors.New("server query exception")
	errPostmarkDuplicateInboundDomain          = errors.New("duplicate inbound domain")
	errPostmarkServerNameAlreadyExists         = errors.New("server name already exists")
	errPostmarkNoDeleteAccess                  = errors.New("you don’t have delete access")
	errPostmarkUnableToDeleteServer            = errors.New("unable to delete server")
	errPostmarkInvalidWebhookURL               = errors.New("invalid webhook URL")
	errPostmarkInvalidServerColor              = errors.New("invalid server color")
	errPostmarkServerNameMissingOrInvalid      = errors.New("server name missing or invalid")
	errPostmarkNoUpdatedServerData             = errors.New("no updated server data received")
	errPostmarkInvalidMXRecordForInboundDomain = errors.New("invalid MX record for inbound domain")
	errPostmarkInvalidInboundSpamThreshold     = errors.New("InboundSpamThreshold value is invalid")
	errPostmarkMessagesQueryException          = errors.New("messages query exception")
	errPostmarkMessageDoesNotExist             = errors.New("message doesn’t exist")
	errPostmarkCannotBypassBlockedInbound      = errors.New("could not bypass this blocked inbound message")
	errPostmarkCannotRetryFailedInbound        = errors.New("could not retry this failed inbound message")
	errPostmarkTriggerQueryException           = errors.New("trigger query exception")
	errPostmarkNoTriggerDataReceived           = errors.New("no trigger data received")
	errPostmarkInboundRuleAlreadyExists        = errors.New("this inbound rule already exists")
	errPostmarkUnableToRemoveInboundRule       = errors.New("unable to remove this inbound rule")
	errPostmarkInboundRuleNotFound             = errors.New("this inbound rule was not found")
	errPostmarkInvalidEmailOrDomain            = errors.New("not a valid email address or domain")
	errPostmarkStatsQueryException             = errors.New("stats query exception")
	errPostmarkBouncesQueryException           = errors.New("bounces query exception")
	errPostmarkBounceNotFound                  = errors.New("bounce was not found")
	errPostmarkBounceIDParameterRequired       = errors.New("BounceID parameter required")
	errPostmarkCannotActivateBounce            = errors.New("cannot activate bounce")
	errPostmarkTemplateQueryException          = errors.New("template query exception")
	errPostmarkTemplateNotFound                = errors.New("template not found")
	errPostmarkTemplateLimitExceeded           = errors.New("template limit would be exceeded")
	errPostmarkNoTemplateDataReceived          = errors.New("no template data received")
	errPostmarkRequiredTemplateFieldMissing    = errors.New("a required template field is missing")
	errPostmarkTemplateFieldTooLarge           = errors.New("template field is too large")
	errPostmarkInvalidTemplateField            = errors.New("a templated field is invalid")
	errPostmarkInvalidFieldInRequest           = errors.New("a field was included in the request that is not allowed")
	errPostmarkTemplateTypesMismatch           = errors.New("the template types don't match on the source and destination servers")
	errPostmarkLayoutTemplateCannotBeDeleted   = errors.New("the layout template cannot be deleted because it has dependent templates")
	errPostmarkInvalidLayoutContentPlaceholder = errors.New("the layout content placeholder must be present exactly once")
	errPostmarkInvalidMessageStreamType        = errors.New("invalid MessageStreamType")
	errPostmarkInvalidMessageStreamID          = errors.New("a valid ID must be provided")
	errPostmarkInvalidMessageStreamName        = errors.New("a valid Name must be provided")
	errPostmarkMessageStreamNameTooLong        = errors.New("the Name is too long, limited to 100 characters")
	errPostmarkMaxMessageStreamsReached        = errors.New("maximum number of message streams reached")
	errPostmarkMessageStreamNotFound           = errors.New("the message stream for the provided ID was not found")
	errPostmarkInvalidMessageStreamIDFormat    = errors.New("the ID must be a non-empty string starting with a letter")
	errPostmarkOnlyOneInboundStreamAllowed     = errors.New("a server can only have one inbound stream")
	errPostmarkCannotArchiveDefaultStreams     = errors.New("cannot archive the default transactional and inbound streams")
	errPostmarkDuplicateMessageStreamID        = errors.New("the ID provided already exists for this server")
	errPostmarkMessageStreamDescriptionTooLong = errors.New("the Description is too long, limited to 1000 characters")
	errPostmarkCannotUnarchiveStream           = errors.New("cannot unarchive this message stream anymore")
	errPostmarkReservedMessageStreamID         = errors.New("the ID must not start with the 'pm-' prefix")
	errPostmarkInvalidDescriptionFormat        = errors.New("the Description must not contain HTML tags")
	errPostmarkMessageStreamDoesNotExist       = errors.New("the MessageStream provided does not exist on this server")
	errPostmarkSendingNotSupported             = errors.New("sending is not supported on the supplied MessageStream")
	errPostmarkReservedID                      = errors.New("the ID 'all' is reserved")
	errPostmarkInvalidDataRemovalRequest       = errors.New("invalid data removal request")
	errPostmarkInvalidDataRemovalID            = errors.New("invalid data removal request ID")
	errPostmarkNoDataRemovalAccess             = errors.New("you don’t have data removal request access")
)

// postmarkErrorForCode returns the corresponding error for a given Postmark error code
func postmarkErrorForCode(errCode int) error {
	switch errCode {
	case 10:
		return errPostmarkBadOrMissingAPIToken
	case 11:
		return errPostmarkMultipleErrorsOccurred
	case 12:
		return errPostmarkResourceNotFound
	case 13:
		return errPostmarkInvalidPaginationKey
	case 100:
		return errPostmarkMaintenance
	case 300:
		return errPostmarkInvalidEmailRequest
	case 400:
		return errPostmarkSenderSignatureNotFound
	case 401:
		return errPostmarkSenderSignatureNotConfirmed
	case 402:
		return errPostmarkInvalidJSON
	case 403:
		return errPostmarkInvalidRequestFields
	case 405:
		return errPostmarkNotAllowedToSend
	case 406:
		return errPostmarkInactiveRecipient
	case 409:
		return errPostmarkJSONRequired
	case 410:
		return errPostmarkTooManyBatchMessages
	case 411:
		return errPostmarkForbiddenAttachmentType
	case 412:
		return errPostmarkAccountIsPending
	case 413:
		return errPostmarkAccountMayNotSend
	case 500:
		return errPostmarkSenderSignatureQueryException
	case 501:
		return errPostmarkSenderSignatureNotFoundByID
	case 502:
		return errPostmarkNoUpdatedSenderSignatureData
	case 503:
		return errPostmarkCannotUsePublicDomain
	case 504:
		return errPostmarkSenderSignatureAlreadyExists
	case 505:
		return errPostmarkDKIMAlreadyScheduledForRenewal
	case 506:
		return errPostmarkSenderSignatureAlreadyConfirmed
	case 507:
		return errPostmarkDoNotOwnSenderSignature
	case 510:
		return errPostmarkDomainNotFound
	case 511:
		return errPostmarkInvalidFieldsSupplied
	case 512:
		return errPostmarkDomainAlreadyExists
	case 513:
		return errPostmarkDoNotOwnDomain
	case 514:
		return errPostmarkNameRequiredForDomain
	case 515:
		return errPostmarkNameTooLongForDomain
	case 516:
		return errPostmarkInvalidNameFormatForDomain
	case 520:
		return errPostmarkMissingFieldForSenderSignature
	case 521:
		return errPostmarkSenderSignatureFieldTooLong
	case 522:
		return errPostmarkInvalidFieldValue
	case 600:
		return errPostmarkServerQueryException
	case 602:
		return errPostmarkDuplicateInboundDomain
	case 603:
		return errPostmarkServerNameAlreadyExists
	case 604:
		return errPostmarkNoDeleteAccess
	case 605:
		return errPostmarkUnableToDeleteServer
	case 606:
		return errPostmarkInvalidWebhookURL
	case 607:
		return errPostmarkInvalidServerColor
	case 608:
		return errPostmarkServerNameMissingOrInvalid
	case 609:
		return errPostmarkNoUpdatedServerData
	case 610:
		return errPostmarkInvalidMXRecordForInboundDomain
	case 611:
		return errPostmarkInvalidInboundSpamThreshold
	case 700:
		return errPostmarkMessagesQueryException
	case 701:
		return errPostmarkMessageDoesNotExist
	case 702:
		return errPostmarkCannotBypassBlockedInbound
	case 703:
		return errPostmarkCannotRetryFailedInbound
	case 800:
		return errPostmarkTriggerQueryException
	case 809:
		return errPostmarkNoTriggerDataReceived
	case 810:
		return errPostmarkInboundRuleAlreadyExists
	case 811:
		return errPostmarkUnableToRemoveInboundRule
	case 812:
		return errPostmarkInboundRuleNotFound
	case 813:
		return errPostmarkInvalidEmailOrDomain
	case 900:
		return errPostmarkStatsQueryException
	case 1000:
		return errPostmarkBouncesQueryException
	case 1001:
		return errPostmarkBounceNotFound
	case 1002:
		return errPostmarkBounceIDParameterRequired
	case 1003:
		return errPostmarkCannotActivateBounce
	case 1100:
		return errPostmarkTemplateQueryException
	case 1101:
		return errPostmarkTemplateNotFound
	case 1105:
		return errPostmarkTemplateLimitExceeded
	case 1109:
		return errPostmarkNoTemplateDataReceived
	case 1120:
		return errPostmarkRequiredTemplateFieldMissing
	case 1121:
		return errPostmarkTemplateFieldTooLarge
	case 1122:
		return errPostmarkInvalidTemplateField
	case 1123:
		return errPostmarkInvalidFieldInRequest
	case 1125:
		return errPostmarkTemplateTypesMismatch
	case 1130:
		return errPostmarkLayoutTemplateCannotBeDeleted
	case 1131:
		return errPostmarkInvalidLayoutContentPlaceholder
	case 1221:
		return errPostmarkInvalidMessageStreamType
	case 1222:
		return errPostmarkInvalidMessageStreamID
	case 1223:
		return errPostmarkInvalidMessageStreamName
	case 1224:
		return errPostmarkMessageStreamNameTooLong
	case 1225:
		return errPostmarkMaxMessageStreamsReached
	case 1226:
		return errPostmarkMessageStreamNotFound
	case 1227:
		return errPostmarkInvalidMessageStreamIDFormat
	case 1228:
		return errPostmarkOnlyOneInboundStreamAllowed
	case 1229:
		return errPostmarkCannotArchiveDefaultStreams
	case 1230:
		return errPostmarkDuplicateMessageStreamID
	case 1231:
		return errPostmarkMessageStreamDescriptionTooLong
	case 1232:
		return errPostmarkCannotUnarchiveStream
	case 1233:
		return errPostmarkReservedMessageStreamID
	case 1234:
		return errPostmarkInvalidDescriptionFormat
	case 1235:
		return errPostmarkMessageStreamDoesNotExist
	case 1236:
		return errPostmarkSendingNotSupported
	case 1237:
		return errPostmarkReservedID
	case 1300:
		return errPostmarkInvalidDataRemovalRequest
	case 1301:
		return errPostmarkInvalidDataRemovalID
	case 1302:
		return errPostmarkNoDataRemovalAccess
	default:
		return errors.New("unknown Postmark error code")
	}
}

func errorWithMessage(code int, message string) error {
	err := postmarkErrorForCode(code)
	msg := message
	if msg == "" {
		msg = "no further message provided"
	}
	return fmt.Errorf("%w: %s", err, msg)
}
