package model

import "time"

type Base struct {
	ID        string    `bson:"_id" validate:"uuid"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}
