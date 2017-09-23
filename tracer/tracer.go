package goTracer

import (
	"github.com/DheerendraRathor/GoTracer/models"
	"github.com/DheerendraRathor/GoTracer/utils"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"runtime"
	"sync"
)

var MaxRenderDepth int = 10
var NaN float64 = math.NaN()

func GoTrace(env *models.World, progress chan<- bool) {
	if env.Settings.RenderDepth > 0 {
		MaxRenderDepth = env.Settings.RenderDepth
	}

	camera := env.GetCamera()
	world := env.GetHitableList()

	width, height := env.Image.Width, env.Image.Height

	var renderWg sync.WaitGroup

	pngImage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	pngFile := utils.CreateNestedFile(env.Image.OutputFile)
	defer pngFile.Close()

	renderRoutines := env.Settings.RenderRoutines
	if renderRoutines <= 0 {
		renderRoutines = runtime.NumCPU()
	}
	renderer := make(chan bool, renderRoutines)
	defer close(renderer)

	imin, imax, jmin, jmax := env.Image.GetPatch()

	for i := imax - 1; i >= imin; i-- {
		for j := jmin; j < jmax; j++ {
			renderer <- true
			renderWg.Add(1)
			go func(i, j, samples int, camera *models.Camera, world *models.HitableList, pngImage *image.RGBA) {
				defer func() {
					<-renderer
					renderWg.Done()
				}()
				processPixel(i, j, width, height, samples, camera, world, pngImage)
				if env.Settings.ShowProgress {
					progress <- false
				}
			}(i, j, env.Image.Samples, camera, world, pngImage)
		}
	}
	renderWg.Wait()

	if env.Settings.ShowProgress {
		progress <- true
	}

	png.Encode(pngFile, pngImage)
}

func processPixel(i, j, imageWidth, imageHeight, sample int, camera *models.Camera, world *models.HitableList, pngImage *image.RGBA) {
	colorVector := models.NewVector(0, 0, 0)
	for s := 0; s < sample; s++ {
		randFloatu, randFloatv := rand.Float64(), rand.Float64()
		u, v := (float64(j)+randFloatu)/float64(imageWidth), (float64(i)+randFloatv)/float64(imageHeight)
		ray := camera.RayAt(u, v)
		temp := getColor(ray, world, 0)
		colorVector.Add(temp)
	}

	colorVector.Scale(1 / float64(sample)).Gamma2()

	uint8Pixel := colorVector.UInt8Pixel()
	rgba := color.RGBA{uint8Pixel.R, uint8Pixel.G, uint8Pixel.B, 255}
	pngImage.Set(j, imageHeight-i-1, rgba)
}

func getColor(r *models.Ray, world *models.HitableList, renderDepth int) *models.Vector {

	willHit, hitRecord := world.Hit(r, 0.0, math.MaxFloat64)
	if willHit {
		shouldScatter, attenuation, ray := hitRecord.Material.Scatter(r, hitRecord)
		if renderDepth < MaxRenderDepth && shouldScatter {
			colorVector := attenuation.Copy().ElemMult(getColor(ray, world, renderDepth+1))
			return colorVector
		}
	}

	unitDir := r.Direction.Copy().Normalize()
	t := 0.5 * (unitDir.Y() + 1.0)
	var startValue, endValue, startBlend, endBlend *models.Vector
	startValue = models.NewVector(1.0, 1.0, 1.0)
	endValue = models.NewVector(0.5, 0.7, 1.0)

	startBlend = startValue.Scale(1 - t)
	endBlend = endValue.Scale(t)
	return startBlend.Add(endBlend)
}
