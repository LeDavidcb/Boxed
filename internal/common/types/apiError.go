package types

// Code errors
const (
	// Auth Related errors.
	AuthTokenExpired             = "AUTH_TOKEN_EXPIRED"
	AuthTokenInvalid             = "AUTH_TOKEN_INVALID"
	AuthTokenMissing             = "AUTH_TOKEN_MISSING"
	AuthInvalidCredentials       = "AUTH_INVALID_CREDENTIALS"
	RefreshTokenMissing          = "REFRESH_TOKEN_MISSING"
	RefreshTokenExpiredOrInvalid = "REFRESH_TOKEN_NOT_VALID"

	// Validation related errors.
	UserEmailAlreadyExists = "USER_EMAIL_ALREADY_EXISTS"

	// File related errords.
	FileUploadFailed     = "FILE_UPLOAD_FAILED"
	ResourceDeleteFailed = "RESOURCE_DELETE_FAILED"
	FileFetchFailed      = "FILE_FETCH_FAILED"
	ResourceNotFound     = "RESOURCE_NOT_FOUND"

	// Server generic errors.
	InternalServerError = "INTERNAL_SERVER_ERROR"
	DatabaseError       = "DATABASE_ERROR"
	InvalidRequest      = "INVALID_REQUEST"
	InvalidFormat       = "INVALID_FORMAT"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
