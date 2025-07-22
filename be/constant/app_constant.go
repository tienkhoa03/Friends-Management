package constant

type ResponseStatus int
type Headers int
type General int

const (
	Success ResponseStatus = iota + 1
	DataNotFound
	Invalidemailorpassword
	UnknownError
	InvalidRequest
	Unauthorized
	StatusForbidden
)

func (r ResponseStatus) GetResponseStatus() string {
	return [...]string{"SUCCESS", "DATA_NOT_FOUND", "Invalid email or password", "UNKNOWN_ERROR", "INVALID_REQUEST", "UNAUTHORIZED", "StatusForbidden"}[r-1]
}

func (r ResponseStatus) GetResponseMessage() string {
	return [...]string{"Success", "Data Not Found", "Invalid email or password", "Unknown Error", "Invalid Request", "Unauthorized", "StatusForbidden"}[r-1]
}
