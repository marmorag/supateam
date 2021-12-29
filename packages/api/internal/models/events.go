package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventKind string

const (
	KindGrandPrix    EventKind = "Grand Prix"
	KindEquipe       EventKind = "Equipe"
	KindEntrainement EventKind = "Entrainement"
)

type Event struct {
	Id          primitive.ObjectID   `bson:"_id" json:"id,omitempty"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Date        primitive.DateTime   `json:"date"`
	Duration    int                  `json:"duration"`
	Kind        EventKind            `json:"kind"`
	Teams       []primitive.ObjectID `json:"teams"`
	Players     []primitive.ObjectID `json:"players"`
}

type CreateEventRequest struct {
	Title       string               `validate:"required"`
	Description string               `validate:""`
	Date        primitive.DateTime   `validate:"required"`
	Duration    int                  `validate:"required"`
	Kind        EventKind            `validate:"required"`
	Teams       []primitive.ObjectID `validate:""`
	Players     []primitive.ObjectID `validate:""`
}

type UpdateEventRequest struct {
	Title       string               `validate:"required"`
	Description string               `validate:""`
	Date        primitive.DateTime   `validate:"required"`
	Duration    int                  `validate:"required"`
	Kind        EventKind            `validate:"required"`
	Teams       []primitive.ObjectID `validate:""`
	Players     []primitive.ObjectID `validate:""`
}

func include(oids []primitive.ObjectID, search primitive.ObjectID) bool {
	for _, oid := range oids {
		if oid == search {
			return true
		}
	}
	return false
}

func (e Event) HasParticipation(participation Participation) bool {
	if participation.Team == primitive.NilObjectID {
		return include(e.Players, participation.Player)
	}
	return include(e.Teams, participation.Team)
}
