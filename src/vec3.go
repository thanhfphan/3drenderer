package src

import (
	"math"
)

type Vec3 struct {
	X, Y, Z float64
}

func (v *Vec3) Rotate(vr *Vec3) *Vec3 {
	return v.RotateX(vr.X).RotateY(vr.Y).RotateZ(vr.Z)
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

func (v *Vec3) Clone() *Vec3 {
	return &Vec3{X: v.X, Y: v.Y, Z: v.Z}
}

func (v *Vec3) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v *Vec3) UnitVector() *Vec3 {
	result := v.Clone()
	length := v.Magnitude()
	if length != 0 {
		result.X /= length
		result.Y /= length
		result.Z /= length
	}

	return result
}

func (v *Vec3) Normalize() *Vec3 {
	length := v.Magnitude()
	if length != 0 {
		v.X /= length
		v.Y /= length
		v.Z /= length
	}
	return v
}

func Vec3Add(a, b *Vec3) *Vec3 {
	return &Vec3{
		X: a.X + b.X,
		Y: a.Y + b.Y,
		Z: a.Z + b.Z,
	}
}

func Vec3Sub(a, b *Vec3) *Vec3 {
	return &Vec3{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}

func Vec3Mul(a *Vec3, scalar float64) *Vec3 {
	return &Vec3{
		X: a.X * scalar,
		Y: a.Y * scalar,
	}
}

func Vec3Div(a *Vec3, scalar float64) *Vec3 {
	if scalar == 0 {
		return a
	}

	return Vec3Mul(a, 1/scalar)
}

func Vec3Dot(a, b *Vec3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func Vec3Cross(a, b *Vec3) *Vec3 {
	return &Vec3{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}
