package seeder

import "go.mongodb.org/mongo-driver/bson/primitive"

type Seeder interface {
	Seed() error
}

type Name string

const (
	Base     Name = "base"
	Dummy    Name = "dummy"
	HttpTest Name = "http-test"
)

var Mapping = map[Name]Seeder{
	Base:     BaseSeeder{},
	Dummy:    DummySeeder{},
	HttpTest: HttpTestSeeder{},
}

func MustObjectIdFromHex(hex string) primitive.ObjectID {
	oID, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		panic(err)
	}
	return oID
}
