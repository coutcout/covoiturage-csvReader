// Package messaging defines messages sended in http response
package messaging

// Default response message
type SingleResponseMessage struct {
	Message string
	Errors  []string
}

// File import description
type FileImportResponseMessage struct {
	Filename       string
	Imported       bool
	NbLineImported int
	Errors         []string
}

// Import description
type FileImportData struct {
	TotalFilesImported int
	NbFilesWithErrors  int
	NbFilesSucceded    int
	NbLineImported     int
}

// Message used when files are imported
type MultipleResponseMessage struct {
	Files []FileImportResponseMessage
	Data  FileImportData
}

