package models

import (
	"math"
	"math/rand"
)

type Camera struct {
	LowerLeftCorner *Vector
	Origin          *Vector
	Horizontal      *Vector
	Vertical        *Vector
	LensRadius      float64
	U, V, W         *Vector
}

func (c *Camera) RayAt(u, v float64) *Ray {

	var randomPoint, offset, yOffset, origin, rayDirection, tempHorizontal, tempVertical *Vector

	randomPoint = RandomPointInUnitDisk()
	randomPoint.Scale(c.LensRadius)

	offset = c.U.Copy()
	yOffset = c.V.Copy()

	offset.Scale(randomPoint.X())
	yOffset.Scale(randomPoint.Y())

	offset.Add(yOffset)

	origin = c.Origin.Copy()
	origin.Add(offset)

	rayDirection = c.LowerLeftCorner.Copy()
	tempHorizontal = c.Horizontal.Copy()
	tempVertical = c.Vertical.Copy()

	tempHorizontal.Scale(u)
	tempVertical.Scale(v)

	rayDirection.Add(tempHorizontal).Add(tempVertical)
	rayDirection.Subtract(origin)

	return &Ray{
		origin,
		rayDirection,
	}
}

func NewCamera(lookFrom, lookAt *Vector, vup *Vector, vfov, aspect, aperture, focus float64) *Camera {
	theta := vfov * math.Pi / 180
	half_height := math.Tan(theta / 2)

	//half_height *= wVector.Length()
	half_width := aspect * half_height

	var u, v, w, llc *Vector

	w = lookFrom.Copy()
	w.Subtract(lookAt).Normalize()

	u = vup
	u.Cross(w).Normalize()

	v = w.Copy()
	v.Cross(u)

	llc = lookFrom.Copy().Subtract(
		u.Copy().Scale(half_width * focus),
	).Subtract(
		v.Copy().Scale(half_height * focus),
	).Subtract(
		w.Copy().Scale(focus),
	)

	return &Camera{
		LowerLeftCorner: llc,
		Horizontal:      u.Copy().Scale(2 * half_width * focus),
		Vertical:        v.Copy().Scale(2 * half_height * focus),
		Origin:          lookFrom,
		LensRadius:      aperture / 2,
		U:               u,
		V:               v,
		W:               w,
	}
}

func RandomPointInUnitDisk() *Vector {
	var p, x *Vector
	for {
		p = NewVector(rand.Float64(), rand.Float64(), 0)
		p.Scale(2)
		x = NewVector(1, 1, 0)
		p.Subtract(x)
		if p.Dot(p) < 1.0 {
			break
		}
	}

	return p
}
