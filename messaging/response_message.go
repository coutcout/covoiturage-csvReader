package messaging

type SingleResponseMessage struct{
	Message string
	Errors []string
}

type FileImportResponseMessage struct {
	Filename string
	Imported bool
	Message string
	Errors []string
}

type FileImportData struct {
	TotalImported int
	NbErrors int
	NbSucceded int
}

type MultipleResponseMessage struct{
	Files []FileImportResponseMessage
	Data FileImportData
}