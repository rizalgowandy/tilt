package controllers

import (
	"github.com/google/wire"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/tilt-dev/tilt/internal/controllers/core/cmd"
	"github.com/tilt-dev/tilt/internal/controllers/core/filewatch"
)

var controllerSet = wire.NewSet(
	filewatch.NewController,
	cmd.WireSet,

	ProvideControllers,
)

func ProvideControllers(fileWatch *filewatch.Controller, cmd *cmd.Controller) []Controller {
	return []Controller{
		fileWatch,
		cmd,
	}
}

var WireSet = wire.NewSet(
	NewTiltServerControllerManager,

	NewScheme,
	NewControllerBuilder,
	NewClientBuilder,

	ProvideDeferredClient,
	wire.Bind(new(ctrlclient.Client), new(*DeferredClient)),

	controllerSet,
)
