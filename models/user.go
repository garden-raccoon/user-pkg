package models

import (
	proto "github.com/garden-raccoon/user-pkg/protocols/user"

	"github.com/gofrs/uuid"
)

type User struct {
	UserUUID  uuid.UUID
	Username  string
	Email     string
	UserType  int
	FirstName string
	LastName  string
}

type UpdateUserRequest struct {
	UserUUID  uuid.UUID
	Username  *string
	Email     *string
	FirstName *string
	LastName  *string
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

// Proto is
func Proto(u UpdateUserRequest) *proto.UpdateUserRequest {
	fields := &proto.UpdateUserRequest{UserUuid: u.UserUUID.Bytes()}

	if u.Email != nil {
		fields.Email = *u.Email
	}

	if u.Username != nil {
		fields.Username = *u.Username
	}

	if u.FirstName != nil {
		fields.FirstName = *u.FirstName
	}

	if u.LastName != nil {
		fields.LastName = *u.LastName
	}

	return fields
}

// UpdateUserRequestFromProto is
func UpdateUserRequestFromProto(pb *proto.UpdateUserRequest) *UpdateUserRequest {
	req := &UpdateUserRequest{
		UserUUID: uuid.FromBytesOrNil(pb.UserUuid),
	}
	if pb.Email != "" {
		req.Email = &pb.Email
	}

	if pb.Username != "" {
		req.Username = &pb.Username
	}

	if pb.FirstName != "" {
		req.FirstName = &pb.FirstName
	}

	if pb.LastName != "" {
		req.LastName = &pb.LastName
	}
	return req
}
