package utils

import (
	"open-crm/core/models"
	"slices"
)

var AuthPermissions = map[models.Role][]models.Permission{
	models.RoleSuperadmin: {models.PermCreateUser, models.PermDeleteUser, models.PermUpdateUser, models.PermViewUser},
	models.RoleAdmin:      {models.PermViewUser},
}

// Utils for permissions
// Local wrapper type for models.User
type UserWrapper struct {
	models.User
}

// Method for checking Role
func (u *UserWrapper) HasRole(role string) bool {
	roles := ConvertRoles(u.Roles) // Ensure u.Roles is converted to []string
	return slices.Contains(roles, string(role))
}

// Method for checking permissions
func RoleHasPermission(role models.Role, perm models.Permission) bool {
	perms, ok := AuthPermissions[role]
	if !ok {
		return false
	}

	return slices.Contains(perms, perm)
}

// Convert []Role into text[]
func ConvertRoles(roles string) []string {
	strRoles := make([]string, len(roles))
	for i, r := range roles {
		strRoles[i] = string(r)
	}
	return strRoles
}
