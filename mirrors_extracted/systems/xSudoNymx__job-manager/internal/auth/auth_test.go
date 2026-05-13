package auth

import (
	"testing"
)

func TestRole_HasPermission(t *testing.T) {
	tests := []struct {
		role Role
		perm Permission
		want bool
	}{
		// Admin can do everything
		{RoleAdmin, PermStart, true},
		{RoleAdmin, PermStop, true},
		{RoleAdmin, PermStatus, true},
		{RoleAdmin, PermStreamOutput, true},

		// Read can only read
		{RoleRead, PermStart, false},
		{RoleRead, PermStop, false},
		{RoleRead, PermStatus, true},
		{RoleRead, PermStreamOutput, true},
	}

	for _, tt := range tests {
		name := string(tt.role) + "/" + string(tt.perm)
		t.Run(name, func(t *testing.T) {
			if got := tt.role.HasPermission(tt.perm); got != tt.want {
				t.Errorf("HasPermission() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseRole(t *testing.T) {
	tests := []struct {
		input   string
		want    Role
		wantErr bool
	}{
		{"Admin", RoleAdmin, false},
		{"admin", RoleAdmin, false},
		{"Read", RoleRead, false},
		{"read", RoleRead, false},
		{"invalid", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseRole(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthorize(t *testing.T) {
	admin := &Identity{Role: RoleAdmin}
	read := &Identity{Role: RoleRead}

	// Admin can start
	if err := Authorize(admin, PermStart); err != nil {
		t.Errorf("admin should be able to start: %v", err)
	}

	// Read cannot start
	if err := Authorize(read, PermStart); err == nil {
		t.Error("read should not be able to start")
	}

	// Read can get status
	if err := Authorize(read, PermStatus); err != nil {
		t.Errorf("read should be able to get status: %v", err)
	}
}
