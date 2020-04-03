package astgen_test

import (
	"go/printer"
	"go/token"
	"log"
	"os"

	"github.com/itchyny/astgen-go"
)

type X struct {
	x int
	y Y
	z *Z
}

type Y struct {
	y int
}

type Z struct {
	s string
	t map[string]int
}

func ExampleBuild() {
	x := &X{1, Y{2}, &Z{"hello", map[string]int{"x": 42}}}
	t, err := astgen.Build(x)
	if err != nil {
		log.Fatal(err)
	}
	err = printer.Fprint(os.Stdout, token.NewFileSet(), t)
	if err != nil {
		log.Fatal(err)
	}

	// Output:
	// &X{x: 1, y: Y{y: 2}, z: &Z{s: "hello", t: map[string]int{"x": 42}}}
}
