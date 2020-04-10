package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/go-playground/validator"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"

	"github.com/imdario/mergo"
)

var (
	isGzip      *bool
	isFixParams *bool
	isAlwaysImg *bool
	gzipLevel   *int
)

var decoder *schema.Decoder
var validate *validator.Validate

type sparklineChart struct {
	//Id           string    `schema:"id"`
	//ClassNames   string    `schema:"cl"`
	ImgWidth  int       `schema:"w" validate:"gt=4"`
	ImgHeight int       `schema:"h" validate:"gt=4"`
	OffSet    int       `schema:"m" validate:"gte=0,ltfield=ImgWidth"`
	LineColor string    `schema:"lc" validate:"required,hexcolor|rgb|rgba"`
	LineWidth int       `schema:"lw" validate:"gt=0"`
	IsFilled  bool      `schema:"f"`
	FillColor string    `schema:"fc" validate:"required_with=IsFilled,omitempty,hexcolor|rgb|rgba"`
	Dots      []float32 `schema:"d" validate:"min=1,max=1000"`
}

type svgImage struct {
	Line    string
	Closure string
	Opts    *sparklineChart
}

func (c *sparklineChart) Draw(w io.Writer) error {
	result := "M "
	closure := fmt.Sprint(" V", c.ImgHeight, " L", c.OffSet, " ", c.ImgHeight, " Z")

	min, max := MinMax(c.Dots)
	if min == max {
		midY := float32(c.ImgHeight+c.OffSet) / 2
		result = fmt.Sprintf("M %d %f L %d %f", c.OffSet, midY, c.ImgWidth-c.OffSet, midY)

		svgTpl.ExecuteTemplate(w, "Sparkline", &svgImage{
			Line:    result,
			Closure: closure,
			Opts:    c,
		})

		return nil
	}

	step := float32(c.ImgWidth-2*c.OffSet) / float32(len(c.Dots)-1)
	scaleY := float32(c.ImgHeight-2*c.OffSet) / (max - min)
	inialX := float32(c.OffSet)
	inialY := float32(c.ImgHeight) - (c.Dots[0]-min)*scaleY - float32(c.OffSet)
	result = fmt.Sprint(result, inialX, " ", inialY)

	for _, value := range c.Dots[1:] {
		inialX += step
		inialY = float32(c.ImgHeight) - (value-min)*scaleY - float32(c.OffSet)
		result = fmt.Sprint(result, " L", inialX, " ", inialY)
	}

	svgTpl.ExecuteTemplate(w, "Sparkline", &svgImage{
		Line:    result,
		Closure: closure,
		Opts:    c,
	})

	return nil
}

var lineDefaults = sparklineChart{
	ImgWidth:  100,
	ImgHeight: 30,
	LineColor: "rgb(255,0,0)",
	LineWidth: 2,
	OffSet:    4,
	IsFilled:  false,
	FillColor: "rgba(219, 59, 158, 0.3)",
}

func init() {
	isGzip = flag.Bool("gzip", false, "enable gzip mode")
	gzipLevel = flag.Int("l", gzip.DefaultCompression, "gzip level")
	//extraRenderOpts
	isFixParams = flag.Bool("fixParams", false, "try to fix missing params")
	isAlwaysImg = flag.Bool("isAlwaysImg", false, "always return correct image")

	//Init singletons
	validate = validator.New()
	decoder = schema.NewDecoder()
}

func fromQuery(u url.Values) (*sparklineChart, error) {
	var p sparklineChart

	if err := decoder.Decode(&p, u); err != nil {
		return nil, fmt.Errorf("failed to decode params: %v", err)
	}

	if *isFixParams {
		mergo.Merge(&p, lineDefaults)
	}

	if err := validate.Struct(&p); err != nil {
		return nil, fmt.Errorf("parameter validation failed: %v", err)
	}

	return &p, nil
}

func main() {
	flag.Parse()

	if !*isGzip {
		*gzipLevel = gzip.NoCompression
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/spkln", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "image/svg+xml")

		c, err := fromQuery(r.URL.Query())

		if err != nil {
			if *isAlwaysImg {
				w.Write([]byte(blank))
			} else {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			return
		}

		c.Draw(w)
	})

	handler := handlers.CompressHandlerLevel(router, *gzipLevel)

	log.Fatal(http.ListenAndServe(":8080", handler))
}
