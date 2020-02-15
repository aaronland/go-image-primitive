package primitive

import (
	"bufio"
	"bytes"
	"context"
	"github.com/aaronland/go-image-resize"
	pr "github.com/fogleman/primitive/primitive"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	_ "log"
	"math"
	"math/rand"
	"runtime"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type PrimitiveOptions struct {
	Alpha      int
	Mode       int
	Iterations int
	Size       int
	ScaleMax   int
	Animated   bool
}

func RandomMode() int {
	return randomInt(1, 7)
}

func NewDefaultPrimitiveOptions() *PrimitiveOptions {

	opts := &PrimitiveOptions{
		Alpha:      0,
		Mode:       1,
		Iterations: 15,
		Size:       0,
		ScaleMax:   256,
		Animated:   false,
	}

	return opts
}

func TransformImage(ctx context.Context, im image.Image, opts *PrimitiveOptions) (image.Image, error) {

	bounds := im.Bounds()
	dims := bounds.Max

	alpha := opts.Alpha
	mode := opts.Mode
	size := opts.Size

	if size == 0 {
		h := float64(dims.Y)
		w := float64(dims.X)
		max := math.Max(h, w)
		size = int(max)
	}

	scale_max := opts.ScaleMax

	if dims.Y > scale_max || dims.X > scale_max {

		new_im, err := resize.ResizeImageMax(ctx, im, scale_max)

		if err != nil {
			return nil, err
		}

		im = new_im
	}

	workers := runtime.NumCPU()

	bg := pr.MakeColor(pr.AverageImageColor(im))
	model := pr.NewModel(im, bg, size, workers)

	for i := 1; i <= opts.Iterations; i++ {

		// t1 := time.Now()
		// log.Printf("begin step %d at %v\n", i, t1)

		model.Step(pr.ShapeType(mode), alpha, workers)

		// log.Printf("finished step %d in %v\n", i, time.Since(t1))

	}

	if opts.Animated {

		g := gif.GIF{}

		frames := model.Frames(0.001)

		delay := 25
		lastDelay := delay * 10

		for i, src := range frames {

			// the original code in pr/utils.go
			// dst := image.NewPaletted(src.Bounds(), palette.Plan9)
			// draw.Draw(dst, dst.Rect, src, image.ZP, draw.Src)

			// https://groups.google.com/forum/#!topic/golang-nuts/28Kk1FfG5XE
			// https://github.com/golang/go/blob/master/src/image/gif/writer.go#L358-L366

			gif_opts := gif.Options{
				NumColors: 256,
				Drawer:    draw.FloydSteinberg,
				Quantizer: nil,
			}

			dst := image.NewPaletted(src.Bounds(), palette.Plan9[:gif_opts.NumColors])
			gif_opts.Drawer.Draw(dst, dst.Rect, src, image.ZP)

			g.Image = append(g.Image, dst)

			if i == len(frames)-1 {
				g.Delay = append(g.Delay, lastDelay)
			} else {
				g.Delay = append(g.Delay, delay)
			}
		}

		out := new(bytes.Buffer)
		err := gif.EncodeAll(out, &g)

		if err != nil {
			return nil, err
		}

		new_im, _, err := image.Decode(bufio.NewReader(out))
		return new_im, err
	}

	return model.Context.Image(), nil
}

func randomInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
