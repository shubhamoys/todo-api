package constants

type ErrorMessage struct {
	AppStatusCode int
	Message       struct {
		Log  string
		Dev  string
		User string
	}
}

var ErrorConstants = struct {
	EndpointNotExist ErrorMessage
	MissingField     ErrorMessage
	InvalidField     ErrorMessage
	DuplicateEntity  ErrorMessage
	DatabaseError    ErrorMessage
	EntityNotFound   ErrorMessage
	NoEntityData     ErrorMessage
	UnkownError      ErrorMessage
	InvalidInput     ErrorMessage
	Unauthorized     ErrorMessage
	ValidationError  ErrorMessage
}{
	EndpointNotExist: ErrorMessage{
		AppStatusCode: 1000000,
		Message: struct {
			Log  string
			Dev  string
			User string
		}{
			Log:  "{endpoint} does not exist",
			Dev:  "This endpoint {endpoint} does not exist",
			User: "Oops! You have reached a dead end",
		},
	},
	MissingField: ErrorMessage{
		AppStatusCode: 1000001,
		Message: struct {
			Log  string
			Dev  string
			User string
		}{
			Log:  "{field} required",
			Dev:  "`{field}` should be present in query",
			User: "Invalid request. No {field} found.",
		},
	},
	InvalidField: ErrorMessage{
		AppStatusCode: 1000002,
		Message: struct {
			Log  string
			Dev  string
			User string
		}{
			Log:  "{field} invalid",
			Dev:  "`{field}` is invalid",
			User: "Invalid {field}",
		},
	},
	DuplicateEntity: ErrorMessage{
		AppStatusCode: 1000004,
		Message: struct {
			Log  string
			Dev  string
			User string
		}{
			Log:  "{entity} duplicate",
			Dev:  "`{entity}` is duplicate",
			User: "Duplicate {entity}",
		},
	},
	DatabaseError: ErrorMessage{
		AppStatusCode: 1000005,
		Message: struct {
			Log  string
			Dev  string
			User string
		}{
			Log:  "DB Error occurred: {message}",
			Dev:  "DB error occurred: {message}",
			User: "Something went wrong",
		},
	},
	EntityNotFound: ErrorMessage{
		AppStatusCode: 1000006,
		Message: struct {
			Log  string
			Dev  string
			User string
		}{
			Log:  "{entity} does not exist",
			Dev:  "`{entity}` does not exist",
			User: "{entity} does not exist",
		},
	},
	NoEntityData: ErrorMessage{
		AppStatusCode: 1000007,
		Message: struct {
			Log  string
			Dev  string
			User string
		}{
			Log:  "no {entity} data",
			Dev:  "no {entity} data",
			User: "no {entity} data",
		},
	},
	UnkownError: ErrorMessage{
		AppStatusCode: 1000008,
		Message: struct {
			Log  string
			Dev  string
			User string
		}{
			Log:  "err occurred: {message}",
			Dev:  "err occurred: {message}",
			User: "err occurred: {message}",
		},
	},
	InvalidInput: ErrorMessage{
		AppStatusCode: 1000009,
		Message: struct {
			Log  string
			Dev  string
			User string
		}{
			Log:  "Invalid input data provided",
			Dev:  "Input validation failed",
			User: "Invalid input data provided: {message}",
		},
	},
	Unauthorized: ErrorMessage{
		AppStatusCode: 1000010,
		Message: struct {
			Log  string
			Dev  string
			User string
		}{
			Log:  "Unauthorized access",
			Dev:  "User not authorized to access {resource}",
			User: "You don't have permission to access this {resource}",
		},
	},
	ValidationError: ErrorMessage{
		AppStatusCode: 1000011,
		Message: struct {
			Log  string
			Dev  string
			User string
		}{
			Log:  "Validation error",
			Dev:  "Validation failed: {message}",
			User: "The provided data is invalid: {message}",
		},
	},
}
