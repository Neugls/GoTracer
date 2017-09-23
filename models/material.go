package models

import (
	"github.com/DheerendraRathor/GoTracer/utils"
	"math/rand"
)

type Material interface {
	Scatter(*Ray, *HitRecord) (bool, *Vector, *Ray)
}

type BaseMaterial struct {
	Albedo *Vector
}

type Lambertian struct {
	*BaseMaterial
}

func NewLambertian(albedo *Vector) *Lambertian {
	return &Lambertian{
		&BaseMaterial{
			Albedo: albedo,
		},
	}
}

func (l *Lambertian) Scatter(ray *Ray, hitRecord *HitRecord) (bool, *Vector, *Ray) {
	pN := hitRecord.P.Copy().Add(hitRecord.N)
	targetPoint := pN.Add(RandomPointInUnitSphere())
	scattered := Ray{
		Origin:    hitRecord.P,
		Direction: targetPoint.Subtract(hitRecord.P),
	}
	return true, l.Albedo, &scattered
}

type Metal struct {
	*BaseMaterial
	fuzz float64
}

func NewMetal(albedo *Vector, fuzz float64) *Metal {
	return &Metal{
		BaseMaterial: &BaseMaterial{
			albedo,
		},
		fuzz: fuzz,
	}
}

func (m Metal) Scatter(ray *Ray, hitRecord *HitRecord) (bool, *Vector, *Ray) {
	reflected := ray.Direction.Copy().Reflect(hitRecord.N)
	reflected.Normalize()

	scattered := Ray{
		hitRecord.P,
		reflected.Add(RandomPointInUnitSphere().Scale(m.fuzz)),
	}
	shouldScatter := reflected.Dot(hitRecord.N) > 0
	return shouldScatter, m.Albedo, &scattered
}

type Dielectric struct {
	*BaseMaterial
	RefIndex float64
}

func NewDielectric(albedo *Vector, r float64) *Dielectric {
	return &Dielectric{
		BaseMaterial: &BaseMaterial{
			Albedo: albedo,
		},
		RefIndex: r,
	}
}

func (d *Dielectric) Scatter(ray *Ray, hitRecord *HitRecord) (bool, *Vector, *Ray) {
	reflected := ray.Direction.Copy().Reflect(hitRecord.N)
	var outwardNormal *Vector
	var ni, nt float64 = 1, 1
	var cosine, reflectionProb float64
	if ray.Direction.Dot(hitRecord.N) > 0 {
		outwardNormal = hitRecord.N.Copy().Negate()
		ni = d.RefIndex
		nt = 1
		cosine = d.RefIndex * ray.Direction.Dot(hitRecord.N) * ray.Direction.Norm()
	} else {
		outwardNormal = hitRecord.N
		ni = 1
		nt = d.RefIndex
		cosine = -ray.Direction.Dot(hitRecord.N) * ray.Direction.Norm()
	}

	var scattered *Ray
	willRefract, refractedVec := ray.Direction.Copy().Refract(outwardNormal, ni, nt)
	if willRefract {
		reflectionProb = utils.Schlick(cosine, 1, d.RefIndex)
		scattered = &Ray{hitRecord.P, refractedVec}
	} else {
		reflectionProb = 1.0
	}

	if rand.Float64() < reflectionProb {
		scattered = &Ray{hitRecord.P, reflected}
	}

	return true, d.Albedo, scattered
}
