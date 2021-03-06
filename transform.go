package primitive

import (
	"context"
	"github.com/aaronland/go-image-transform"
	"image"
	"net/url"
	"strconv"
)

type PrimitiveTransformation struct {
	transform.Transformation
	options *PrimitiveOptions
}

func init() {

	ctx := context.Background()
	err := transform.RegisterTransformation(ctx, "primitive", NewPrimitiveTransformation)

	if err != nil {
		panic(err)
	}
}

func NewPrimitiveTransformation(ctx context.Context, str_url string) (transform.Transformation, error) {

	parsed, err := url.Parse(str_url)

	if err != nil {
		return nil, err
	}

	opts := NewDefaultPrimitiveOptions()

	query := parsed.Query()

	str_mode := query.Get("mode")
	str_size := query.Get("size")
	str_iterations := query.Get("iterations")

	if str_mode != "" {

		if str_mode == "random" {

			str_exclude := query["exclude_mode"]
			exclude := make([]int, len(str_exclude))
			
			for idx, str_mode := range str_exclude {

				i, err := strconv.Atoi(str_mode)

				if err != nil {
					return nil, err
				}

				exclude[idx] = i
			}

			for {

				m := RandomMode()
				ok := true

				if len(exclude) > 0 {
					
					for _, i := range exclude {
						if i == m {
							ok = false
							break
						}
					}
				}
				
				if ok {
					opts.Mode = m
					break
				}
			}
			
		} else {

			mode, err := strconv.Atoi(str_mode)

			if err != nil {
				return nil, err
			}

			opts.Mode = mode
		}
	}

	if str_size != "" {

		size, err := strconv.Atoi(str_size)

		if err != nil {
			return nil, err
		}

		opts.Size = size
	}

	if str_iterations != "" {

		if str_iterations == "random" {
			opts.Iterations = randomInt(10, 100)
		} else {
			iterations, err := strconv.Atoi(str_iterations)

			if err != nil {
				return nil, err
			}

			opts.Iterations = iterations
		}
	}

	tr := &PrimitiveTransformation{
		options: opts,
	}

	return tr, nil
}

func (tr *PrimitiveTransformation) Transform(ctx context.Context, im image.Image) (image.Image, error) {
	return TransformImage(ctx, im, tr.options)
}
