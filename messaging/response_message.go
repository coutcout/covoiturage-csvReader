package messaging

type SingleResponseMessage struct {
	Message string
	Errors  []string
}

type FileImportResponseMessage struct {
	Filename string
	Imported bool
	NbLineImported int
	Errors   []string
}

type FileImportData struct {
	TotalFilesImported int
	NbFilesWithErrors  int
	NbFilesSucceded    int
	NbLineImported	   int
}

type MultipleResponseMessage struct {
	Files []FileImportResponseMessage
	Data  FileImportData
}
