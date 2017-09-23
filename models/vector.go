package models

import (
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/gonum"
	"math"
)

var blasf64 blas.Float64 = gonum.Implementation{}

type Vector struct {
	Data []float64
}

func (v *Vector) X() float64 {
	return v.Data[0]
}

func (v *Vector) Y() float64 {
	return v.Data[1]
}

func (v *Vector) Z() float64 {
	return v.Data[2]
}

func (v *Vector) Norm() float64 {
	return blasf64.Dnrm2(3, v.Data, 1)
}

func NewVector(x, y, z float64) *Vector {
	return &Vector{
		Data: []float64{x, y, z},
	}
}

func NewVectorFromSlice(data []float64) *Vector {
	if len(data) != 3 {
		panic("slice size must be 3 for Vector")
	}

	return &Vector{
		Data: data,
	}
}

func (v *Vector) Copy() *Vector {
	v1 := NewVector(0, 0, 0)
	blasf64.Dcopy(3, v.Data, 1, v1.Data, 1)
	return v1
}

// v = v.v1
func (v *Vector) Dot(v1 *Vector) float64 {
	return blasf64.Ddot(3, v.Data, 1, v1.Data, 1)
}

// v = v*v1 and return v
func (v *Vector) Cross(v1 *Vector) *Vector {
	tempx := v.Y()*v1.Z() - v.Z()*v1.Y()
	tempy := -v.X()*v1.Z() - v.Z()*v1.X()
	tempz := v.X()*v1.Y() - v.Y()*v1.X()

	v.Data = []float64{tempx, tempy, tempz}
	return v
}

// v *= t and return v
func (v *Vector) Scale(t float64) *Vector {
	blasf64.Dscal(3, t, v.Data, 1)
	return v
}

// v += v1 and return v
func (v *Vector) Add(v1 *Vector) *Vector {
	blasf64.Daxpy(3, 1, v1.Data, 1, v.Data, 1)
	return v
}

// v = v - v1 and return v
func (v *Vector) Subtract(v1 *Vector) *Vector {
	blasf64.Daxpy(3, -1, v1.Data, 1, v.Data, 1)
	return v
}

// v /= ||v|| and returns v
func (v *Vector) Normalize() *Vector {
	norm := v.Norm()
	v.Scale(1 / norm)
	return v
}

// Reflect v around n and return v
func (v *Vector) Reflect(n *Vector) *Vector {
	temp := n.Copy()
	temp.Scale(2 * v.Dot(n))
	v.Subtract(temp)
	return v
}

func (v *Vector) Negate() *Vector {
	return v.Scale(-1)
}

func (v *Vector) Refract(n *Vector, ni, nt float64) (bool, *Vector) {
	v.Normalize()
	cosθ := v.Dot(n)
	snellRatio := ni / nt
	discriminator := 1 - snellRatio*(1-cosθ*cosθ)
	if discriminator > 0 {
		v1 := v.Subtract(n.Copy().Scale(cosθ)).Scale(snellRatio)
		v2 := n.Copy().Scale(math.Sqrt(discriminator))
		return true, v1.Subtract(v2)
	}
	return false, v
}

func (v *Vector) ElemMult(v1 *Vector) *Vector {
	x, y, z := v.Data[0]*v1.Data[0], v.Data[1]*v1.Data[1], v.Data[2]*v1.Data[2]
	v.Data = []float64{x, y, z}
	return v
}
