module imagetools

go 1.22

replace imagecore => ./image_processing_core

require (
	golang.org/x/exp v0.0.0-20240404231335-c0f41cb1a7a0 // indirect
	golang.org/x/image v0.15.0 // indirect
)

require (
	gopkg.in/yaml.v2 v2.4.0
	imagecore v1.0.0
)
