package models

import (
	"math"
	"math/rand"
)

type Sphere struct {
	Center   *Vector
	Radius   float64
	Material Material
}

func NewSphere(x, y, z, r float64, material Material) *Sphere {
	return &Sphere{
		NewVector(x, y, z),
		r,
		material,
	}
}

func (s *Sphere) Hit(r *Ray, tmin, tmax float64) (bool, *HitRecord) {
	oc := r.Origin.Copy().Subtract(s.Center)
	var a, b, c, d float64
	a = r.Direction.Dot(r.Direction)
	b = 2.0 * oc.Dot(r.Direction)
	c = oc.Dot(oc) - s.Radius*s.Radius
	d = b*b - 4*a*c

	record := HitRecord{}
	if d > 0 {
		sqrtD := math.Sqrt(d)
		a2 := 2 * a
		temp := (-b - sqrtD) / a2
		if temp > tmin && temp < tmax {
			record.T = temp
			record.P = r.PointAtParameter(temp)
			record.N = record.P.Copy().Subtract(s.Center).Normalize()
			record.Material = s.Material
			return true, &record
		}
		temp = (-b + sqrtD) / a2
		if temp > tmin && temp < tmax {
			record.T = temp
			record.P = r.PointAtParameter(temp)
			record.N = record.P.Copy().Subtract(s.Center).Normalize()
			record.Material = s.Material
			return true, &record
		}
	}
	return false, nil
}

func RandomPointInUnitSphere() *Vector {
	var p, offset *Vector
	offset = NewVector(1, 1, 1)
	for {
		p = NewVector(rand.Float64(), rand.Float64(), rand.Float64()).Scale(2).Subtract(offset)

		if p.Dot(p) < 1.0 {
			break
		}
	}

	return p
}
