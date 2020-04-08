package main

import (
	"compress/gzip"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"

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

type SparkLineParams struct {
	//Id           string    `schema:"id"`
	//ClassNames   string    `schema:"cl"`
	ImgWidth  int       `schema:"w" validate:"gt=4"`
	ImgHeight int       `schema:"h" validate:"gt=4"`
	OffSet    int       `schema:"m" validate:"gte=0,ltfield=ImgWidth"`
	LineColor string    `schema:"lc" validate:"hexcolor|rgb|rgba"`
	LineWidth int       `schema:"lw" validate:"gt=0"`
	IsFilled  bool      `schema:"f"`
	FillColor string    `schema:"fc" validate:"hexcolor|rgb|rgba"`
	Dots      []float32 `schema:"d" validate:"min=2,max=1000"`
}

var lineDefaults = SparkLineParams{
	ImgWidth:  100,
	ImgHeight: 30,
	LineColor: "red",
	LineWidth: 2,
	OffSet:    4,
	IsFilled:  false,
}

func parseSparkLineRequestParams(rawParams map[string][]string) (SparkLineParams, error) {
	var lineParams SparkLineParams
	err := decoder.Decode(&lineParams, rawParams)
	return lineParams, err
}

func createLineStringFromDots(dots []float32, width, height, offset int) (string, string, error) {
	if dots == nil || len(dots) < 2 {
		return "", "", errors.New("len(dots): dots array must contain at least 2 elements")
	}
	if width <= 2*offset || height <= offset {
		return "", "", errors.New(fmt.Sprintf("imageDimensions: %dx%d width must be > 2 * offset and height must be gt offset", width, height))
	}
	result := "M "
	closure := fmt.Sprint(" V", height, " L", offset, " ", height, " Z")

	min, max := MinMax(dots)
	if min == max {
		//ToDo: deal with one dot
		midY := float32(height+offset) / 2
		result = fmt.Sprintf("M %d %f L %d %f", offset, midY, width-offset, midY)
		return result, closure, nil
	}
	step := float32(width-2*offset) / float32(len(dots)-1)
	scaleY := float32(height-2*offset) / (max - min)
	inialX := float32(offset)
	inialY := float32(height) - (dots[0]-min)*scaleY - float32(offset)

	result = fmt.Sprint(result, inialX, " ", inialY)

	for _, value := range dots[1:] {
		inialX += step
		inialY = float32(height) - (value-min)*scaleY - float32(offset)
		result = fmt.Sprint(result, " L", inialX, " ", inialY)
	}

	return result, closure, nil
}

func drawSparkLine(w http.ResponseWriter, r *http.Request) {
	lineParams := parseAndValidateRequestParams(r)

	//Generate SVG
	w.Header().Add("content-type", "image/svg+xml")
	fmt.Fprint(w, XML_HEADING)
	fmt.Fprintf(w, SVG_OPEN_TAG, lineParams.ImgWidth, lineParams.ImgHeight)
	sprklinePath, closurePath, _ := createLineStringFromDots(lineParams.Dots, lineParams.ImgWidth, lineParams.ImgHeight, lineParams.OffSet)

	fmt.Fprintf(w, SVG_PATH_LINE, sprklinePath, lineParams.LineColor, lineParams.LineWidth)

	if lineParams.IsFilled || lineParams.FillColor != "" {
		fillColor := "rgba(219, 59, 158, 0.3)"
		if lineParams.FillColor != "" {
			fillColor = lineParams.FillColor
		}
		fmt.Fprintf(w, SVG_FILL_PATH, fmt.Sprint(sprklinePath, closurePath), fillColor)
	}

	fmt.Fprint(w, SVG_CLOSE_TAG)
}

func parseAndValidateRequestParams(r *http.Request) SparkLineParams {
	var lineParams, err = parseSparkLineRequestParams(r.URL.Query())

	if err != nil {
		log.Println("Request params parse errors:")
		log.Println(err)
	}

	if *isFixParams {
		mergo.Merge(&lineParams, lineDefaults)
	}

	err = validate.Struct(&lineParams)
	log.Println(lineParams)

	if err != nil {
		log.Println("Request params validation errors:")

		for _, fieldErr := range err.(validator.ValidationErrors) {
			if *isAlwaysImg {
				if fieldErr.Namespace() == "SparkLineParams.Dots" {
					lineParams.Dots = []float32{0, 0}
				}
			}
			//ToDo: add extra checks
			log.Println(fieldErr)
		}

	}
	return lineParams
}

func main() {
	initEnvVars()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/spkln", drawSparkLine).Methods("GET")
	var handler http.Handler = router
	if !*isGzip {
		*gzipLevel = gzip.NoCompression
	}
	handler = handlers.CompressHandlerLevel(router, *gzipLevel)

	log.Fatal(http.ListenAndServe(":8080", handler))
}

func initEnvVars() {
	isGzip = flag.Bool("gzip", false, "enable gzip mode")
	gzipLevel = flag.Int("l", gzip.DefaultCompression, "gzip level")
	//extraRenderOpts
	isFixParams = flag.Bool("fixParams", false, "try to fix missing params")
	isAlwaysImg = flag.Bool("isAlwaysImg", true, "always return correct image")
	flag.Parse()

	//Init singletons
	validate = validator.New()
	decoder = schema.NewDecoder()
}
