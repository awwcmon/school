package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Teacher struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Age       int                `bson:"age" json:"age"`
	BornAt    time.Time          `bson:"born_at" json:"bornAt"`
	Job       *Job               `bson:"job" json:"job"`
	Books     []Book             `bson:"books" json:"books"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}
type Book struct {
	fileId   string `bson:"file_id" json:"fileId"`
	fileName string `bson:"file_name" json:"fileName"`
	fileSize int64  `bson:"file_size" json:"fileSize"`
}

// TableName table name
func (m *Teacher) TableName() string {
	return "teacher"
}

type Job struct {
	Title  string `bson:"title" json:"title"`
	Salary int    `bson:"salary" json:"salary"`
}
