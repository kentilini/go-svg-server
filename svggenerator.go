package main

const (
	XML_HEADING   string = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><!DOCTYPE svg PUBLIC \"-//W3C//DTD SVG 20010904//EN\" \"http://www.w3.org/TR/2001/REC-SVG-20010904/DTD/svg10.dtd\">"
	SVG_OPEN_TAG         = "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" xml:space=\"preserve\" width=\"%d\" height=\"%d\">"
	SVG_CLOSE_TAG        = "</svg>"
	SVG_PATH_LINE        = "<path d=\"%s\" stroke=\"%s\" stroke-width=\"%d\" fill=\"none\"/>"
	SVG_FILL_PATH        = "<path d=\"%s\" stroke=\"none\" fill=\"%s\"/>"
)
