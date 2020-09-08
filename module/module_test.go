package module

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const testExeData = `path	github.com/mitchellh/golicense
mod	github.com/mitchellh/golicense	(devel)	
dep	github.com/fatih/color	v1.7.0	h1:DkWD4oS2D8LGGgTQ6IvwJJXSL5Vp2ffcQg58nFV38Ys=
dep	github.com/mattn/go-colorable	v0.0.9	h1:UVL0vNpWh04HeJXV0KLcaT7r06gOH2l4OW6ddYRUIY4=
dep	github.com/mattn/go-isatty	v0.0.4	h1:bnP0vzxcAdeI1zdubAl5PjU6zsERjGZb7raWodagDYs=
dep	github.com/rsc/goversion	v1.2.0	h1:zVF4y5ciA/rw779S62bEAq4Yif1cBc/UwRkXJ2xZyT4=
dep	github.com/rsc/goversion/v12	v12.0.0	h1:zVF4y5ciA/rw779S62bEAq4Yif1cBc/UwRkXJ2xZyT4=
`

const replacement = `
path	github.com/gohugoio/hugo
mod	github.com/gohugoio/hugo	(devel)
dep	github.com/markbates/inflect	v1.0.0
=>	github.com/markbates/inflect	v0.0.0-20171215194931-a12c3aec81a6	h1:LZhVjIISSbj8qLf2qDPP0D8z0uvOWAW5C85ly5mJW6c=
`

const testExeData2 = `path	github.com/JoakimSoderberg/go-license-finder
mod	github.com/JoakimSoderberg/go-license-finder	(devel)
dep	github.com/dgryski/go-minhash	v0.0.0-20190315135803-ad340ca03076	h1:EB7M2v8Svo3kvIDy+P1YDE22XskDQP+TEYGzeDwPAN4=
dep	github.com/ekzhu/minhash-lsh	v0.0.0-20171225071031-5c06ee8586a1	h1:/7G7q8SDJdrah5jDYqZI8pGFjSqiCzfSEO+NgqKCYX0=
dep	github.com/emirpasic/gods	v1.12.0	h1:QAUIPSaCu4G+POclxeqb3F+WPpdKqFGlw36+yOzGlrg=
dep	github.com/go-enry/go-license-detector/v4	v4.0.0
=>	github.com/JoakimSoderberg/go-license-detector/v4	v4.0.0-20200827131053-a8ed0b9cb40a	h1:YOPawvrqnDbtX+T+oM2b6UMNOCMqo+DoP6A6qzsfVHI=
dep	github.com/go-git/gcfg	v1.5.0	h1:Q5ViNfGF8zFgyJWPqYwA7qGFoMTEiBmdlkcfRmpIMa4=
`

func TestParseExeData(t *testing.T) {
	cases := []struct {
		Name     string
		Input    string
		Expected []Module
		Error    string
	}{
		{
			"typical (from golicense itself)",
			testExeData,
			[]Module{
				Module{
					Path:    "github.com/fatih/color",
					Version: "v1.7.0",
					Hash:    "h1:DkWD4oS2D8LGGgTQ6IvwJJXSL5Vp2ffcQg58nFV38Ys=",
				},
				Module{
					Path:    "github.com/mattn/go-colorable",
					Version: "v0.0.9",
					Hash:    "h1:UVL0vNpWh04HeJXV0KLcaT7r06gOH2l4OW6ddYRUIY4=",
				},
				Module{
					Path:    "github.com/mattn/go-isatty",
					Version: "v0.0.4",
					Hash:    "h1:bnP0vzxcAdeI1zdubAl5PjU6zsERjGZb7raWodagDYs=",
				},
				Module{
					Path:    "github.com/rsc/goversion",
					Version: "v1.2.0",
					Hash:    "h1:zVF4y5ciA/rw779S62bEAq4Yif1cBc/UwRkXJ2xZyT4=",
				},
				Module{
					Path:    "github.com/rsc/goversion/v12",
					Version: "v12.0.0",
					Hash:    "h1:zVF4y5ciA/rw779S62bEAq4Yif1cBc/UwRkXJ2xZyT4=",
				},
			},
			"",
		},

		{
			"replacement syntax",
			strings.TrimSpace(replacement),
			[]Module{
				Module{
					Path:    "github.com/markbates/inflect",
					Version: "v1.0.0",
					Hash:    "h1:LZhVjIISSbj8qLf2qDPP0D8z0uvOWAW5C85ly5mJW6c=",
					Replace: &Module{
						Path:    "",
						Version: "v0.0.0-20171215194931-a12c3aec81a6",
						Hash:    "h1:LZhVjIISSbj8qLf2qDPP0D8z0uvOWAW5C85ly5mJW6c=",
					},
				},
			},
			"",
		},
		{
			Name:  "from gobindep",
			Input: testExeData2,
			Expected: []Module{
				{
					Path:    "github.com/dgryski/go-minhash",
					Version: "v0.0.0-20190315135803-ad340ca03076",
					Hash:    "h1:EB7M2v8Svo3kvIDy+P1YDE22XskDQP+TEYGzeDwPAN4=",
				},
				{
					Path:    "github.com/ekzhu/minhash-lsh",
					Version: "v0.0.0-20171225071031-5c06ee8586a1",
					Hash:    "h1:/7G7q8SDJdrah5jDYqZI8pGFjSqiCzfSEO+NgqKCYX0=",
				},
				{
					Path:    "github.com/emirpasic/gods",
					Version: "v1.12.0",
					Hash:    "h1:QAUIPSaCu4G+POclxeqb3F+WPpdKqFGlw36+yOzGlrg=",
				},
				{
					Path:    "github.com/go-enry/go-license-detector/v4",
					Version: "v4.0.0",
					Hash:    "h1:YOPawvrqnDbtX+T+oM2b6UMNOCMqo+DoP6A6qzsfVHI=",
					Replace: &Module{
						Path:    "github.com/JoakimSoderberg/go-license-detector/v4",
						Version: "v4.0.0-20200827131053-a8ed0b9cb40a",
						Hash:    "h1:YOPawvrqnDbtX+T+oM2b6UMNOCMqo+DoP6A6qzsfVHI=",
					},
				},
				{
					Path:    "github.com/go-git/gcfg",
					Version: "v1.5.0",
					Hash:    "h1:Q5ViNfGF8zFgyJWPqYwA7qGFoMTEiBmdlkcfRmpIMa4=",
				},
			},
			Error: "",
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)
			actual, err := ParseExeData(tt.Input)
			if tt.Error != "" {
				require.Error(err)
				require.Contains(err.Error(), tt.Error)
				return
			}
			require.NoError(err)
			for i, mod := range tt.Expected {
				require.Equal(mod.Path, actual[i].Path)
			}
		})
	}
}
