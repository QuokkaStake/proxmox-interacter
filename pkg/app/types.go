package app

import "main/pkg/types"

type ContainerActionRender struct {
	Action    string
	Container types.Container
}

type ContainerScaleRender struct {
	Container   types.Container
	Config      *types.ContainerConfig
	ScaleParams types.ScaleMatcher
}

type ContainerErrorRender struct {
	Error        error
	ClusterInfos types.ClusterInfos
}

type ContainerInfoRender struct {
	Container   types.Container
	Config      *types.ContainerConfig
	ConfigError error
}
