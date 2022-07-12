package seeder

type Seeder interface {
	Seed() error
}

type Name string

const (
	Base  Name = "base"
	Dummy Name = "dummy"
)

var Mapping = map[Name]Seeder{
	Base:  BaseSeeder{},
	Dummy: DummySeeder{},
}
