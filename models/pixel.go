package models

import "math"

func (p *Vector) R() float64 {
	return p.Data[0]
}

func (p *Vector) G() float64 {
	return p.Data[1]
}

func (p *Vector) B() float64 {
	return p.Data[2]
}

func (p *Vector) Gamma2() {
	x, y, z := math.Sqrt(p.R()), math.Sqrt(p.G()), math.Sqrt(p.B())
	p.Data = []float64{x, y, z}
}

func (p *Vector) UInt8Pixel() *Uint8Pixel {
	p.Scale(255.99)
	return &Uint8Pixel{
		uint8(p.R()),
		uint8(p.G()),
		uint8(p.B()),
	}
}

type Uint8Pixel struct {
	R uint8
	G uint8
	B uint8
}
