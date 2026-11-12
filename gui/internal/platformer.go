package internal

import (
	"image/color"
	"math"

	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/collision"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/key"
	"github.com/oakmound/oak/v4/scene"
)

const (
	Ground collision.Label = 1
)

func createCharacter(ctx *scene.Context) *entities.Entity {
	return entities.New(
		ctx,
		entities.WithRect(floatgeom.NewRect2WH(100, 100, 16, 32)),
		entities.WithColor(color.RGBA{255, 0, 0, 255}),
		entities.WithSpeed(floatgeom.Point2{3, 7}),
	)
}

func createPlatforms(ctx *scene.Context) {
	platforms := []floatgeom.Rect2{
		floatgeom.NewRect2WH(0, 400, 300, 20),
		floatgeom.NewRect2WH(100, 250, 30, 20),
		floatgeom.NewRect2WH(340, 300, 100, 20),
	}

	for _, p := range platforms {
		entities.New(ctx,
			entities.WithRect(p),
			entities.WithColor(color.RGBA{0, 0, 255, 255}),
			entities.WithLabel(Ground),
		)
	}
}

func applyControls(character *entities.Entity) {
	if oak.IsDown(key.A) {
		character.Delta[0] = -character.Speed.X()
	} else if oak.IsDown(key.D) {
		character.Delta[0] = character.Speed.X()
	} else {
		character.Delta[0] = (0)
	}
}

func onEnterEventCallback(
	character *entities.Entity,
	fallSpeed float64,
) event.Bindable[*entities.Entity, event.EnterPayload] {
	return func(c *entities.Entity, e event.EnterPayload) event.Response {
		applyControls(character)

		oldX, oldY := character.X(), character.Y()
		character.ShiftDelta()
		onGround := false

		hit := collision.HitLabel(character.Space, Ground)

		if hit != nil && !(oldY != character.Y() && oldY+character.H() > hit.Y()) {
			character.SetY(hit.Y() - character.H())
			character.Delta[1] = 0

			if oak.IsDown(key.Spacebar) {
				character.Delta[1] -= character.Speed.Y()
			}
			onGround = true
		} else {
			if character.Y() > 500 {
				character.Delta[1] = 0
				character.SetY(100)
				character.SetX(100)
			}

			character.Delta[1] += fallSpeed
		}

		if hit != nil {
			xover, yover := character.Space.Overlap(hit)

			if math.Abs(yover) > 1 {
				xbump := 1.0
				if xover > 0 {
					xbump = -1
				}
				character.SetX(oldX + xbump)
				if character.Delta.Y() < 0 {
					character.Delta[1] = 0
				}
			}

			if !onGround && math.Abs(xover) > 1 {
				character.SetY(oldY + 1)
				character.Delta[1] = fallSpeed
			}
		}

		return 0
	}
}

func goSceneStart(ctx *scene.Context) {
	character := createCharacter(ctx)

	const fallSpeed = .2

	eventCallback := onEnterEventCallback(character, fallSpeed)
	event.Bind(ctx, event.Enter, character, eventCallback)

	createPlatforms(ctx)
}

func addGoScene() {
	oak.AddScene("go", scene.Scene{Start: goSceneStart})
}

func RunPlatformer() {
	addGoScene()
	oak.Init("go")
}
