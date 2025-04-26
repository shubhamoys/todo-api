package constants

// UserRoles defines all available user roles in the system
var UserRoles = struct {
	User       string
	Admin      string
	Superadmin string
}{
	User:       "User",
	Admin:      "Admin",
	Superadmin: "Superadmin",
}
