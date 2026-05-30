package models

import "time"

// UserRole represents the role of a user
type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

// UserSettings represents user preferences
type UserSettings struct {
	ReceiveEmails bool `json:"receiveEmails"`
}

// User represents a user in the system
type User struct {
	Key               string       `json:"_key,omitempty"`
	ID                string       `json:"_id,omitempty"`
	Rev               string       `json:"_rev,omitempty"`
	Username          string       `json:"username" binding:"required"`
	InvalidArtistName bool         `json:"invalidArtistName,omitempty"`
	Email             string       `json:"email" binding:"required,email"`
	PasswordMissing   bool         `json:"passwordMissing,omitempty"`
	ValidatedEmail    bool         `json:"validatedEmail"`
	Role              UserRole     `json:"role"`
	Settings          UserSettings `json:"settings"`
	CreatedAt         time.Time    `json:"createdAt"`
}

// ClientUser represents a user with expanded relations for API responses
type ClientUser struct {
	User
	SubscribedTo []ClientSubscriber `json:"subscribedTo"`
	Artist       *Artist            `json:"artist,omitempty"`
}

// FlatUser represents a minimal user for references
type FlatUser struct {
	Key      string `json:"_key,omitempty"`
	ID       string `json:"_id,omitempty"`
	Rev      string `json:"_rev,omitempty"`
	Username string `json:"username"`
}

// Social represents a connected social account
type Social struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
}

// AuthCredentials represents OAuth credentials
type AuthCredentials struct {
	AccessToken  string `json:"accessToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int    `json:"expiresIn"`
	RefreshToken string `json:"refreshToken"`
}

// SocialWithCredentials represents a social account with OAuth credentials
type SocialWithCredentials struct {
	Social
	Credentials AuthCredentials `json:"credentials"`
}

// UserWithCredentials represents a user with sensitive credential data (for internal use only)
type UserWithCredentials struct {
	User
	Password              string                           `json:"password"`
	EmailUnsubscribeToken string                           `json:"emailUnsubscribeToken"`
	Socials               map[string]SocialWithCredentials `json:"socials"`
}

// RegisterUserRequest represents the request body for user registration
type RegisterUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginUserRequest represents the request body for user login
type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// SetUsernameRequest represents the request body for setting username
type SetUsernameRequest struct {
	Username string `json:"username" binding:"required"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Username string        `json:"username,omitempty"`
	Email    string        `json:"email,omitempty"`
	Settings *UserSettings `json:"settings,omitempty"`
}

type CreateAvatarUploadRequest struct {
	FileName string `json:"fileName" binding:"required"`
	FileType string `json:"fileType" binding:"required"`
	FileSize int64  `json:"fileSize" binding:"required"`
}

type CreateAvatarUploadResponse struct {
	PresignedUrl string `json:"presignedUrl"`
	StoragePath  string `json:"storagePath"`
}
