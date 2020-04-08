# go svg server
Simple SVG Web Server Draws SparkLines

## Quick Start Guide

```shell script
docker build --tag kentilini/sparkline:latest .
docker run  --name svg_sparkline -p8080:8080 kentilini/sparkline:latest -fixParams -gzip -isAlwaysImg
```

And visit [http://localhost:8080/spkln?d=2,4,5,2,18.5,12,34](http://localhost:8080/spkln?d=2,4,5,2,18.5,12,34)

| Param | Type      | Example                                    | Desription                                                 |
|-------|-----------|--------------------------------------------|------------------------------------------------------------|
| w     | int       | w=100                                      | width of output img                                        |
| h     | int       | h=30                                       | height of output img                                       |
| m     | int       | m=4                                        | left, right, top offset of img                             |
| d     | []float32 | d=1,2,4,5.662,7                            | coma separated array of sparkline dots                     |
| lc    | string    | lc=red                                     | encoded stroke-line color, !be careful no force validation |
| lw    | int       | lw=2                                       | stroke-line width                                          |
| f     | bool      | f=1                                        | flag to fill area under sparkline                          |
| fc    | string    | fc=rgba%28219%2C%2059%2C%20158%2C%200.3%29 | encoded fill collor                                        |

---

## Examples

[Minimal http://localhost:8080/spkln?d=2,4,5,2,18.5,12,34](http://localhost:8080/spkln?d=2,4,5,2,18.5,12,34)

[Fill](http://localhost:8080/spkln?d=2,4,5,2,18.5,12,34&f=1)

[CustomWidth](http://localhost:8080/spkln?d=2,4,5,2,18.5,12,34&f=1&w=450)

[CustomColor](http://localhost:8080/spkln?d=2,4,5,2,18.5,12,34&f=1&w=450&lc=blue)


Try to mix params and verify result!