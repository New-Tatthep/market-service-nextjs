package constant

const (
	UserSystemName      = "system"
	ConnectorVersion1_0 = "1.0"

	ConnectorTypeImapName   = "connector.imap"
	ConnectorTypeManualName = "manual"
	ConnectorTypeOdooName   = "connector.odoo"

	FrequencyWeekName = "week"
	FrequencyDayName  = "day"

	ExtensionZipName = "zip"

	DocumentStatusWait            = "wait"
	DocumentStatusWaitForWorkFlow = "wait_for_workflow"
	DocumentStatusSuccess         = "success"
	DocumentStatusOdooProcess     = "odoo_process"
	DocumentStatusError           = "error"

	ActivityDocumentStatusWaitForReview = "wait_for_review"
	ActivityDocumentStatusSuccess       = "success"
	ActivityDocumentStatusError         = "error"

	TypeScheduleName = "schedule"
)

var MapTargetConnector = map[string]string{
	ConnectorTypeImapName: "source",
	ConnectorTypeOdooName: "destination",
}
