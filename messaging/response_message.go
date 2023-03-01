package messaging

type ResponseMessage struct{
	StatusCode int
	Message string
	Errors []string
}