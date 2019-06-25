package main

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock/cmd/minimock/practice"
	"github.com/stretchr/testify/require"
)

//go:generate go run minimock.go -g -i Target -o practice/ -s _mock.go

type Target interface {
	Shoot(ctx context.Context, projectile string) error
}

func Test_Target_Practice(t *testing.T) {
	m := practice.NewTargetMock(t)

	const projectile = "projectile"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	cancel()

	m.ShootMock.
		When(
			m.MinimockArg.MatchedContextContext(func(ctx context.Context) bool {
				_, deadlineIsSet := ctx.Deadline()
				return deadlineIsSet
			}),
			projectile,
		).
		Then(nil)

	require.NoError(t, m.Shoot(ctx, projectile))

	m.MinimockFinish()
}

//go:generate go run minimock.go -g -i AmbiguousSpecs -o practice/ -s _mock.go

type AmbiguousSpecs interface {
	TryIt(i1, i2, i3 int) error
}

func Test_AmbiguousSpecs_MixingMatchersAndValuesIsAlright(t *testing.T) {
	m := practice.NewAmbiguousSpecsMock(t)
	m.TryItMock.
		When(13, m.MinimockArg.MatchedInt(func(int) bool { return true }), 11).
		Then(nil)
	require.Nil(t, m.TryIt(13, 42, 11))
	m.MinimockFinish()
}

func Test_AmbiguousSpecs_MinimockBailsOnBadSpecs(t *testing.T) {
	m := practice.NewAmbiguousSpecsMock(t)
	// it's impossible to know if the matcher was specified for the first parameter or for the last one, so this should fail
	m.TryItMock.When(m.MinimockArg.MatchedInt(func(int) bool { return true }), 42, 0)
}
