// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package cli

import (
	"context"
	"github.com/google/wire"
	"github.com/windmilleng/tilt/internal/build"
	"github.com/windmilleng/tilt/internal/demo"
	"github.com/windmilleng/tilt/internal/docker"
	"github.com/windmilleng/tilt/internal/dockercompose"
	"github.com/windmilleng/tilt/internal/dockerfile"
	"github.com/windmilleng/tilt/internal/engine"
	"github.com/windmilleng/tilt/internal/hud"
	"github.com/windmilleng/tilt/internal/hud/server"
	"github.com/windmilleng/tilt/internal/k8s"
	"github.com/windmilleng/tilt/internal/store"
	"time"
)

// Injectors from wire.go:

func wireDemo(ctx context.Context, branch demo.RepoBranch) (demo.Script, error) {
	v := provideClock()
	renderer := hud.NewRenderer(v)
	headsUpDisplay, err := hud.NewDefaultHeadsUpDisplay(renderer)
	if err != nil {
		return demo.Script{}, err
	}
	envOrError := k8s.ProvideEnvOrError(ctx)
	client, err := k8s.ProvideK8sClient(ctx, envOrError)
	if err != nil {
		return demo.Script{}, err
	}
	podWatcher := engine.NewPodWatcher(client)
	env := k8s.ProideEnv(envOrError)
	nodeIP, err := k8s.DetectNodeIP(ctx, env)
	if err != nil {
		return demo.Script{}, err
	}
	serviceWatcher := engine.NewServiceWatcher(client, nodeIP)
	reducer := _wireReducerValue
	storeLogActionsFlag := provideLogActions()
	storeStore := store.NewStore(reducer, storeLogActionsFlag)
	podLogManager := engine.NewPodLogManager(client)
	portForwardController := engine.NewPortForwardController(client)
	fsWatcherMaker := engine.ProvideFsWatcherMaker()
	timerMaker := engine.ProvideTimerMaker()
	watchManager := engine.NewWatchManager(fsWatcherMaker, timerMaker)
	syncletManager := engine.NewSyncletManager(client)
	syncletBuildAndDeployer := engine.NewSyncletBuildAndDeployer(syncletManager)
	cli, err := docker.DefaultClient(ctx, env)
	if err != nil {
		return demo.Script{}, err
	}
	containerUpdater := build.NewContainerUpdater(cli)
	analytics, err := provideAnalytics()
	if err != nil {
		return demo.Script{}, err
	}
	localContainerBuildAndDeployer := engine.NewLocalContainerBuildAndDeployer(containerUpdater, analytics)
	console := build.DefaultConsole()
	writer := build.DefaultOut()
	labels := _wireLabelsValue
	dockerImageBuilder := build.NewDockerImageBuilder(cli, console, writer, labels)
	imageBuilder := build.DefaultImageBuilder(dockerImageBuilder)
	cacheBuilder := build.NewCacheBuilder(cli)
	engineUpdateModeFlag := provideUpdateModeFlag()
	updateMode, err := engine.ProvideUpdateMode(engineUpdateModeFlag, env)
	if err != nil {
		return demo.Script{}, err
	}
	clock := build.ProvideClock()
	imageBuildAndDeployer := engine.NewImageBuildAndDeployer(imageBuilder, cacheBuilder, client, env, analytics, updateMode, clock)
	dockerComposeClient := dockercompose.NewDockerComposeClient()
	imageAndCacheBuilder := engine.NewImageAndCacheBuilder(imageBuilder, cacheBuilder, updateMode)
	dockerComposeBuildAndDeployer := engine.NewDockerComposeBuildAndDeployer(dockerComposeClient, cli, imageAndCacheBuilder, clock)
	buildOrder := engine.DefaultBuildOrder(syncletBuildAndDeployer, localContainerBuildAndDeployer, imageBuildAndDeployer, dockerComposeBuildAndDeployer, env, updateMode)
	compositeBuildAndDeployer := engine.NewCompositeBuildAndDeployer(buildOrder)
	buildController := engine.NewBuildController(compositeBuildAndDeployer)
	imageReaper := build.NewImageReaper(cli)
	imageController := engine.NewImageController(imageReaper)
	globalYAMLBuildController := engine.NewGlobalYAMLBuildController(client)
	configsController := engine.NewConfigsController()
	dockerComposeEventWatcher := engine.NewDockerComposeEventWatcher(dockerComposeClient)
	dockerComposeLogManager := engine.NewDockerComposeLogManager(dockerComposeClient)
	profilerManager := engine.NewProfilerManager()
	analyticsReporter := engine.ProvideAnalyticsReporter(analytics, storeStore)
	upper := engine.NewUpper(ctx, headsUpDisplay, podWatcher, serviceWatcher, storeStore, podLogManager, portForwardController, watchManager, buildController, imageController, globalYAMLBuildController, configsController, dockerComposeEventWatcher, dockerComposeLogManager, profilerManager, syncletManager, analyticsReporter)
	script := demo.NewScript(upper, headsUpDisplay, client, env, storeStore, branch)
	return script, nil
}

var (
	_wireReducerValue = engine.UpperReducer
	_wireLabelsValue  = dockerfile.Labels{}
)

func wireThreads(ctx context.Context) (Threads, error) {
	v := provideClock()
	renderer := hud.NewRenderer(v)
	headsUpDisplay, err := hud.NewDefaultHeadsUpDisplay(renderer)
	if err != nil {
		return Threads{}, err
	}
	envOrError := k8s.ProvideEnvOrError(ctx)
	client, err := k8s.ProvideK8sClient(ctx, envOrError)
	if err != nil {
		return Threads{}, err
	}
	podWatcher := engine.NewPodWatcher(client)
	env := k8s.ProideEnv(envOrError)
	nodeIP, err := k8s.DetectNodeIP(ctx, env)
	if err != nil {
		return Threads{}, err
	}
	serviceWatcher := engine.NewServiceWatcher(client, nodeIP)
	reducer := _wireReducerValue
	storeLogActionsFlag := provideLogActions()
	storeStore := store.NewStore(reducer, storeLogActionsFlag)
	podLogManager := engine.NewPodLogManager(client)
	portForwardController := engine.NewPortForwardController(client)
	fsWatcherMaker := engine.ProvideFsWatcherMaker()
	timerMaker := engine.ProvideTimerMaker()
	watchManager := engine.NewWatchManager(fsWatcherMaker, timerMaker)
	syncletManager := engine.NewSyncletManager(client)
	syncletBuildAndDeployer := engine.NewSyncletBuildAndDeployer(syncletManager)
	cli, err := docker.DefaultClient(ctx, env)
	if err != nil {
		return Threads{}, err
	}
	containerUpdater := build.NewContainerUpdater(cli)
	analytics, err := provideAnalytics()
	if err != nil {
		return Threads{}, err
	}
	localContainerBuildAndDeployer := engine.NewLocalContainerBuildAndDeployer(containerUpdater, analytics)
	console := build.DefaultConsole()
	writer := build.DefaultOut()
	labels := _wireLabelsValue
	dockerImageBuilder := build.NewDockerImageBuilder(cli, console, writer, labels)
	imageBuilder := build.DefaultImageBuilder(dockerImageBuilder)
	cacheBuilder := build.NewCacheBuilder(cli)
	engineUpdateModeFlag := provideUpdateModeFlag()
	updateMode, err := engine.ProvideUpdateMode(engineUpdateModeFlag, env)
	if err != nil {
		return Threads{}, err
	}
	clock := build.ProvideClock()
	imageBuildAndDeployer := engine.NewImageBuildAndDeployer(imageBuilder, cacheBuilder, client, env, analytics, updateMode, clock)
	dockerComposeClient := dockercompose.NewDockerComposeClient()
	imageAndCacheBuilder := engine.NewImageAndCacheBuilder(imageBuilder, cacheBuilder, updateMode)
	dockerComposeBuildAndDeployer := engine.NewDockerComposeBuildAndDeployer(dockerComposeClient, cli, imageAndCacheBuilder, clock)
	buildOrder := engine.DefaultBuildOrder(syncletBuildAndDeployer, localContainerBuildAndDeployer, imageBuildAndDeployer, dockerComposeBuildAndDeployer, env, updateMode)
	compositeBuildAndDeployer := engine.NewCompositeBuildAndDeployer(buildOrder)
	buildController := engine.NewBuildController(compositeBuildAndDeployer)
	imageReaper := build.NewImageReaper(cli)
	imageController := engine.NewImageController(imageReaper)
	globalYAMLBuildController := engine.NewGlobalYAMLBuildController(client)
	configsController := engine.NewConfigsController()
	dockerComposeEventWatcher := engine.NewDockerComposeEventWatcher(dockerComposeClient)
	dockerComposeLogManager := engine.NewDockerComposeLogManager(dockerComposeClient)
	profilerManager := engine.NewProfilerManager()
	analyticsReporter := engine.ProvideAnalyticsReporter(analytics, storeStore)
	upper := engine.NewUpper(ctx, headsUpDisplay, podWatcher, serviceWatcher, storeStore, podLogManager, portForwardController, watchManager, buildController, imageController, globalYAMLBuildController, configsController, dockerComposeEventWatcher, dockerComposeLogManager, profilerManager, syncletManager, analyticsReporter)
	headsUpServer := server.ProvideHeadsUpServer(storeStore)
	threads := provideThreads(headsUpDisplay, upper, headsUpServer)
	return threads, nil
}

func wireK8sClient(ctx context.Context) (k8s.Client, error) {
	envOrError := k8s.ProvideEnvOrError(ctx)
	client, err := k8s.ProvideK8sClient(ctx, envOrError)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// wire.go:

var K8sWireSet = wire.NewSet(k8s.ProvideEnvOrError, k8s.ProideEnv, k8s.DetectNodeIP, k8s.ProvideK8sClient)

var BaseWireSet = wire.NewSet(
	K8sWireSet, docker.DefaultClient, wire.Bind(new(docker.Client), new(docker.Cli)), dockercompose.NewDockerComposeClient, build.NewImageReaper, engine.DeployerWireSet, engine.NewPodLogManager, engine.NewPortForwardController, engine.NewBuildController, engine.NewPodWatcher, engine.NewServiceWatcher, engine.NewImageController, engine.NewConfigsController, engine.NewDockerComposeEventWatcher, engine.NewDockerComposeLogManager, engine.NewProfilerManager, provideClock, hud.NewRenderer, hud.NewDefaultHeadsUpDisplay, provideLogActions, store.NewStore, wire.Bind(new(store.RStore), new(store.Store)), engine.NewUpper, provideAnalytics, engine.ProvideAnalyticsReporter, provideUpdateModeFlag, engine.NewWatchManager, engine.ProvideFsWatcherMaker, engine.ProvideTimerMaker, server.ProvideHeadsUpServer, provideThreads,
)

type Threads struct {
	hud    hud.HeadsUpDisplay
	upper  engine.Upper
	server server.HeadsUpServer
}

func provideThreads(h hud.HeadsUpDisplay, upper engine.Upper, server2 server.HeadsUpServer) Threads {
	return Threads{h, upper, server2}
}

func provideClock() func() time.Time {
	return time.Now
}
