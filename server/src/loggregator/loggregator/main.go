package main

import (
	"cfcomponent"
	"cfcomponent/instrumentation"
	"errors"
	"flag"
	"fmt"
	"github.com/cloudfoundry/go_cfmessagebus"
	"github.com/cloudfoundry/gosteno"
	"io/ioutil"
	"loggregator/agentlistener"
	"loggregator/sink"
	"os"
	"os/signal"
	"registrar"
	"runtime"
	"time"
)

type Config struct {
	Index                  uint
	ApiHost                string
	UaaVerificationKeyFile string
	SystemDomain           string
	NatsHost               string
	NatsPort               int
	NatsUser               string
	NatsPass               string
	VarzUser               string
	VarzPass               string
	VarzPort               uint32
	SourcePort             uint32
	WebPort                uint32
	LogFilePath            string
	decoder                sink.TokenDecoder
	mbusClient             cfmessagebus.MessageBus
}

func (c *Config) validate(logger *gosteno.Logger) (err error) {
	if c.VarzPass == "" || c.VarzUser == "" || c.VarzPort == 0 {
		return errors.New("Need VARZ username/password/port.")
	}
	if c.SystemDomain == "" {
		return errors.New("Need system domain to register with NATS")
	}
	uaaVerificationKey, err := ioutil.ReadFile(c.UaaVerificationKeyFile)
	if err != nil {
		return errors.New(fmt.Sprintf("Can not read UAA verification key from file %s: %s", c.UaaVerificationKeyFile, err))
	}
	c.decoder, err = sink.NewUaaTokenDecoder(uaaVerificationKey)
	if err != nil {
		return errors.New(fmt.Sprintf("Can not parse UAA verification key: %s", err))
	}

	c.mbusClient, err = cfmessagebus.NewMessageBus("NATS")
	if err != nil {
		return errors.New(fmt.Sprintf("Can not create message bus to NATS: %s", err))
	}
	c.mbusClient.Configure(c.NatsHost, c.NatsPort, c.NatsUser, c.NatsPass)
	c.mbusClient.SetLogger(logger)
	err = c.mbusClient.Connect()
	if err != nil {
		return errors.New(fmt.Sprintf("Could not connect to NATS: ", err.Error()))
	}
	return nil
}

var version = flag.Bool("version", false, "Version info")
var logFilePath = flag.String("logFile", "", "The agent log file, defaults to STDOUT")
var logLevel = flag.Bool("v", false, "Verbose logging")
var configFile = flag.String("config", "config/loggregator.json", "Location of the loggregator config json file")
var uaaVerificationKeyFile = flag.String("tokenFile", "config/uaa_token.pub", "Location of the loggregator's uaa public token file")

const versionNumber = `0.0.TRAVIS_BUILD_NUMBER`
const gitSha = `TRAVIS_COMMIT`

type LoggregatorServerHealthMonitor struct {
}

func (hm LoggregatorServerHealthMonitor) Ok() bool {
	return true
}

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("\n\nversion: %s\ngitSha: %s\n\n", versionNumber, gitSha)
		return
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	logger := cfcomponent.NewLogger(*logLevel, *logFilePath, "loggregator")

	config := &Config{SourcePort: 3456, WebPort: 8080, UaaVerificationKeyFile: *uaaVerificationKeyFile}
	err := cfcomponent.ReadConfigInto(config, *configFile)
	if err != nil {
		panic(err)
	}
	err = config.validate(logger)
	if err != nil {
		panic(err)
	}

	listener := agentlistener.NewAgentListener(fmt.Sprintf("0.0.0.0:%d", config.SourcePort), logger)
	incomingData := listener.Start()

	authorizer := sink.NewLogAccessAuthorizer(config.decoder, config.ApiHost)
	sinkServer := sink.NewSinkServer(incomingData, logger, fmt.Sprintf("0.0.0.0:%d", config.WebPort), authorizer, 30*time.Second)

	cfc, err := cfcomponent.NewComponent(
		config.SystemDomain,
		config.WebPort,
		"LoggregatorServer",
		config.Index,
		&LoggregatorServerHealthMonitor{},
		config.VarzPort,
		[]string{config.VarzUser, config.VarzPass},
		[]instrumentation.Instrumentable{listener, sinkServer},
	)

	if err != nil {
		panic(err)
	}

	r := registrar.NewRegistrar(config.mbusClient, logger)
	err = r.RegisterWithRouter(&cfc)
	if err != nil {
		panic(err)
	}

	err = r.RegisterWithCollector(cfc)
	if err != nil {
		panic(err)
	}

	cfc.StartMonitoringEndpoints()

	go sinkServer.Start()
	go cfcomponent.DumpGoroutineInfoOnCommand()

	killChan := make(chan os.Signal)
	signal.Notify(killChan, os.Kill)

	<-killChan
	r.UnregisterFromRouter(cfc)
}
