package models

import (
	"math"
	"math/rand"
)

type Sphere struct {
	Center   Vector
	Radius   float64
	Material Material
}

func NewSphere(x, y, z, r float64, material Material) *Sphere {
	return &Sphere{
		Vector{x, y, z},
		r,
		material,
	}
}

func (s *Sphere) Hit(r *Ray, tmin, tmax float64) (bool, *HitRecord) {
	oc := SubtractVectors(r.Origin, s.Center)
	var a, b, c, d float64
	a = VectorDotProduct(r.Direction, r.Direction)
	b = 2.0 * VectorDotProduct(oc, r.Direction)
	c = VectorDotProduct(oc, oc) - s.Radius*s.Radius
	d = b*b - 4*a*c

	record := HitRecord{}
	if d > 0 {
		sqrtD := math.Sqrt(d)
		a2 := 2 * a
		temp := (-b - sqrtD) / a2
		if temp > tmin && temp < tmax {
			record.T = temp
			record.P = r.PointAtParameter(temp)
			record.N = UnitVector(SubtractVectors(record.P, s.Center))
			record.Material = s.Material
			return true, &record
		}
		temp = (-b + sqrtD) / a2
		if temp > tmin && temp < tmax {
			record.T = temp
			record.P = r.PointAtParameter(temp)
			record.N = UnitVector(SubtractVectors(record.P, s.Center))
			record.Material = s.Material
			return true, &record
		}
	}
	return false, nil
}

func RandomPointInUnitSphere() Vector {
	var p, offset Vector
	offset = []float64{1, 1, 1}
	for {
		p = []float64{rand.Float64(), rand.Float64(), rand.Float64()}
		p.MultiplyScalar(2).Subtract(offset)
		if VectorDotProduct(p, p) < 1.0 {
			break
		}
	}

	return p
}
