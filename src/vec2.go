package src

import "math"

type Vec2 struct {
	X, Y float64
}

func (v *Vec2) Clone() *Vec2 {
	return &Vec2{X: v.X, Y: v.Y}
}

func (v *Vec2) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vec2) UnitVector() *Vec2 {
	result := v.Clone()
	length := v.Magnitude()
	if length != 0 {
		result.X /= length
		result.Y /= length
	}

	return result
}

func (v *Vec2) Normalize() *Vec2 {
	length := v.Magnitude()
	if length != 0 {
		v.X /= length
		v.Y /= length
	}

	return v
}

func Vec2Add(a, b *Vec2) *Vec2 {
	return &Vec2{
		X: a.X + b.X,
		Y: a.Y + b.Y,
	}
}

func Vec2Sub(a, b *Vec2) *Vec2 {
	return &Vec2{
		X: a.X - b.Y,
		Y: a.Y - b.Y,
	}
}

func Vec2Mul(v *Vec2, scalar float64) *Vec2 {
	return &Vec2{
		X: v.X * scalar,
		Y: v.Y * scalar,
	}
}

func Vec2Div(v *Vec2, scalar float64) *Vec2 {
	if scalar == 0 {
		return v
	}

	return Vec2Mul(v, 1/scalar)
}

func Vec2Dot(a, b *Vec2) float64 {
	return a.X*b.X + a.Y*b.Y
}
