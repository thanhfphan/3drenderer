package src

import (
	"math"
)

type Vec3 struct {
	X, Y, Z float64
}

func (v *Vec3) RotateX(angle float64) *Vec3 {
	return &Vec3{
		X: v.X,
		Y: v.Y*math.Cos(angle) - v.Z*math.Sin(angle),
		Z: v.Y*math.Sin(angle) + v.Z*math.Cos(angle),
	}
}

func (v *Vec3) RotateY(angle float64) *Vec3 {
	return &Vec3{
		X: v.X*math.Cos(angle) - v.Z*math.Sin(angle),
		Y: v.Y,
		Z: v.X*math.Sin(angle) + v.Z*math.Cos(angle),
	}
}

func (v *Vec3) RotateZ(angle float64) *Vec3 {
	return &Vec3{
		X: v.X*math.Cos(angle) - v.Y*math.Sin(angle),
		Y: v.X*math.Sin(angle) + v.Y*math.Cos(angle),
		Z: v.Z,
	}
}
