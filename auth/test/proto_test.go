package test

import (
	"testing"

	"github.com/abisalde/grpc-microservice/auth/pkg/ent/proto/auth_pbuf"
)

func TestUserProto(t *testing.T) {
	user := auth_pbuf.User{
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
