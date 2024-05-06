package model

import "time"

type Base struct {
	ID        string    `json:"id,omitempty" bson:"_id" validate:"uuid"`
	CreatedAt time.Time `json:"createAt,omitempty" bson:"createdAt"`
	UpdatedAt time.Time `json:"updateAt,omitempty" bson:"updatedAt"`
}
