package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Password  string             `bson:"password" json:"password"`
	RealName  string             `bson:"realName" json:"realName"`
	Roles     []string           `bson:"roles" json:"roles"`
	Username  string             `bson:"username" json:"username"`
	Codes     []string           `bson:"codes" json:"codes"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

// TableName table name
func (m *User) TableName() string {
	return "user"
}
