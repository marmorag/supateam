package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ParticipationStatus string

const (
	ParticipationUnknown  ParticipationStatus = "unknown"
	ParticipationAccepted ParticipationStatus = "accepted"
	ParticipationRejected ParticipationStatus = "rejected"
)

type Participation struct {
	Id     primitive.ObjectID  `bson:"_id" json:"id,omitempty"`
	Event  primitive.ObjectID  `json:"event"`
	Player primitive.ObjectID  `json:"player"`
	Team   primitive.ObjectID  `json:"team"`
	Status ParticipationStatus `json:"status"`
}

type ParticipationLong struct {
	Id     primitive.ObjectID  `bson:"_id" json:"id,omitempty"`
	Event  primitive.ObjectID  `json:"event"`
	Player []User              `json:"player"`
	Team   []Team              `json:"team"`
	Status ParticipationStatus `json:"status"`
}

type CreateParticipationRequest struct {
	Event  primitive.ObjectID  `validate:"required"`
	Player primitive.ObjectID  `validate:"required"`
	Team   primitive.ObjectID  `validate:""`
	Status ParticipationStatus `validate:"required"`
}

type UpdateParticipationRequest struct {
	Event  primitive.ObjectID  `validate:"required"`
	Player primitive.ObjectID  `validate:"required"`
	Team   primitive.ObjectID  `validate:""`
	Status ParticipationStatus `validate:"required"`
}

func IncludeObject(participations []Participation, objectID primitive.ObjectID) bool {
	for _, participation := range participations {
		if participation.Team == objectID || participation.Player == objectID {
			return true
		}
	}
	return false
}
