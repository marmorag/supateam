package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Team struct {
	Id      primitive.ObjectID   `bson:"_id" json:"id,omitempty"`
	Name    string               `json:"name,omitempty"`
	Members []primitive.ObjectID `json:"members"`
}

type CreateTeamRequest struct {
	Name    string               `validate:"required"`
	Members []primitive.ObjectID `validate:"required"`
}

type UpdateTeamRequest struct {
	Name    string               `validate:"required"`
	Members []primitive.ObjectID `validate:"required"`
}
