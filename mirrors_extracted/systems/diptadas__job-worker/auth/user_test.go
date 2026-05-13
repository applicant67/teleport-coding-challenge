package auth

import "testing"

func TestVerifyPermissionForName(t *testing.T) {
	if err := VerifyPermissionForName("alice", PermissionReadOnly); err != nil {
		t.Errorf("expected permission verified, found %v", err)
	}
	if err := VerifyPermissionForName("alice", PermissionReadWrite); err != nil {
		t.Errorf("expected permission verified, found %v", err)
	}
	if err := VerifyPermissionForName("bob", PermissionReadOnly); err != nil {
		t.Errorf("expected permission verified, found %v", err)
	}
	if err := VerifyPermissionForName("bob", PermissionReadWrite); err == nil {
		t.Error("expected permission denied error")
	}
	if err := VerifyPermissionForName("unknown", PermissionReadWrite); err == nil {
		t.Error("expected user not found error")
	}
}
