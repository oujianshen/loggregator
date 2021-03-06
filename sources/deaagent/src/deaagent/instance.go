package deaagent

import (
	"deaagent/loggregatorclient"
	"github.com/cloudfoundry/gosteno"
	"logMessage"
	"path/filepath"
	"strconv"
)

type instance struct {
	applicationId       string
	spaceId             string
	wardenJobId         uint64
	wardenContainerPath string
	index               uint64
}

func (instance *instance) identifier() string {
	return filepath.Join(instance.wardenContainerPath, "jobs", strconv.FormatUint(instance.wardenJobId, 10))
}

func (inst *instance) startListening(loggregatorClient loggregatorclient.LoggregatorClient, logger *gosteno.Logger) {
	newLoggingStream(inst, loggregatorClient, logger, logMessage.LogMessage_OUT).listen()
	newLoggingStream(inst, loggregatorClient, logger, logMessage.LogMessage_ERR).listen()
}
