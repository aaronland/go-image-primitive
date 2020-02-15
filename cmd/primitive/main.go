package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-image-cli"
	"github.com/aaronland/go-image-primitive"
	"image"
	"log"
	"path/filepath"
	"strings"
)

func main() {

	flag.Parse()

	ctx := context.Background()

	opts := primitive.NewDefaultPrimitiveOptions()

	cb := func(ctx context.Context, im image.Image, path string) (image.Image, string, error) {

		new_im, err := primitive.TransformImage(ctx, im, opts)

		if err != nil {
			return nil, "", err
		}

		root := filepath.Dir(path)

		fname := filepath.Base(path)
		ext := filepath.Ext(path)

		label := "primitive"

		short_name := strings.Replace(fname, ext, "", 1)
		new_name := fmt.Sprintf("%s-%s%s", short_name, label, ext)

		new_path := filepath.Join(root, new_name)

		return new_im, new_path, nil
	}

	paths := flag.Args()

	err := cli.Process(ctx, cb, paths...)

	if err != nil {
		log.Fatal(err)
	}

}
