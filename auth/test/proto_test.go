package test

import (
	"testing"

	"github.com/abisalde/gprc-microservice/auth/pkg/ent/proto/entpb"
)

func TestUserProto(t *testing.T) {
	user := entpb.User{
		FirstName: "rotemtam",
		Email:     "rotemtam@example.com",
	}
	if user.FirstName != "rotemtam" {
		t.Fatal("expected user name to be rotemtam")
	}
	if user.Email != "rotemtam@example.com" {
		t.Fatal("expected email address to be rotemtam@example.com")
	}
}
