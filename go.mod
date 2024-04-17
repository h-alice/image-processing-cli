module imagetools

go 1.22

replace imagecore => ./image_processing_core

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/exp v0.0.0-20240404231335-c0f41cb1a7a0 // indirect
	golang.org/x/image v0.15.0 // indirect
)

require (
	github.com/urfave/cli/v2 v2.27.1
	gopkg.in/yaml.v2 v2.4.0
	imagecore v1.0.0
)
