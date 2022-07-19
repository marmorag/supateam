package seeder

import "fmt"

type DummySeeder struct{}

func (d DummySeeder) Seed() error {
	fmt.Println("dummy seeder")

	return nil
}
