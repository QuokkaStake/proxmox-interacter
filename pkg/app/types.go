package app

import "main/pkg/types"

type ContainerActionRender struct {
	Action    string
	Container types.Container
}

type ContainerScaleRender struct {
	Container   types.Container
	ScaleParams types.ScaleMatcher
}
