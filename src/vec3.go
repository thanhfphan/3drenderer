package src

import (
	"math"
)

type Vec3 struct {
	x, y, z float64
}

func (v *Vec3) RotateX(angle float64) *Vec3 {
	return &Vec3{
		x: v.x,
		y: v.y*math.Cos(angle) - v.z*math.Sin(angle),
		z: v.y*math.Sin(angle) + v.z*math.Cos(angle),
	}
}

func (v *Vec3) RotateY(angle float64) *Vec3 {
	return &Vec3{
		x: v.x*math.Cos(angle) - v.z*math.Sin(angle),
		y: v.y,
		z: v.x*math.Sin(angle) + v.z*math.Cos(angle),
	}
}

func (v *Vec3) RotateZ(angle float64) *Vec3 {
	return &Vec3{
		x: v.x*math.Cos(angle) - v.y*math.Sin(angle),
		y: v.x*math.Sin(angle) + v.y*math.Cos(angle),
		z: v.z,
	}
}
