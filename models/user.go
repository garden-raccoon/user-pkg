package models

import (
	proto "github.com/garden-raccoon/user-pkg/protocols/user"

	"github.com/gofrs/uuid"
)

type User struct {
	UserUUID uuid.UUID
	Username string
	Email    string
	UserType int
}

// UserFromProto is
func UserFromProto(pb *proto.User) *User {
	return &User{
		UserUUID: uuid.FromBytesOrNil(pb.UserUuid),
		Email:    pb.Email,
		Username: pb.Username,
		UserType: int(pb.UserType),
	}
}

func (u User) Proto() *proto.User {
	employer := &proto.User{
		UserUuid: u.UserUUID.Bytes(),
		Username: u.Username,
		Email:    u.Email,
		UserType: int64(u.UserType),
	}
	return employer
}
