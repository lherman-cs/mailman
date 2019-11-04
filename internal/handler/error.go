package handler 

type HTTPError struct {
	status int
	reason string
}

func (e *HTTPError) Error() string {
	return e.reason
}
