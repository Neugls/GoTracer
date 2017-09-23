package models

type Ray struct {
	Origin    *Vector
	Direction *Vector
}

func (r *Ray) PointAtParameter(t float64) *Vector {
	dirVec := r.Direction.Copy().Scale(t)
	return dirVec.Add(r.Origin)
}
