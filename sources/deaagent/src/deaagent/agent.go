package deaagent

import (
	"deaagent/loggregatorclient"
	"github.com/cloudfoundry/gosteno"
	"io/ioutil"
	"os"
	"runtime"
	"time"
)

type agent struct {
	InstancesJsonFilePath string
	*gosteno.Logger
	lastModTime time.Time
}

const bufferSize = 4096

func NewAgent(instancesJsonFilePath string, logger *gosteno.Logger) *agent {
	return &agent{instancesJsonFilePath, logger, *(new(time.Time))}
}

func (agent *agent) Start(loggregatorClient loggregatorclient.LoggregatorClient) {
	newInstances := agent.watchInstancesJsonFileForChanges()
	for {
		instance := <-newInstances
		agent.Infof("Starting to listen to %v\n", instance.identifier())
		instance.startListening(loggregatorClient)
	}
}

func (agent *agent) readInstancesJsonFileIfChanged() ([]byte, bool) {
	var curModTime time.Time

	f, err := os.Open(agent.InstancesJsonFilePath)
	defer f.Close()
	if err != nil {
		agent.Warnf("Opening failed. %s\n", err)
		return nil, false
	}

	if fi, err := f.Stat(); err == nil {
		curModTime = fi.ModTime()
	} else {
		agent.Warnf("File stat() failed. %s\n", err)
		return nil, false
	}

	if agent.lastModTime.Before(curModTime) {
		if bytes, err := ioutil.ReadAll(f); err != nil {
			agent.Warnf("Reading failed. %s\n", err)
			return nil, false
		} else {
			agent.lastModTime = curModTime
			return bytes, true
		}
	}

	return nil, false
}

func (agent *agent) watchInstancesJsonFileForChanges() chan *instance {
	instancesChan := make(chan *instance)

	pollInstancesJson := func() {
		knownInstances := make(map[string]bool)

		for {
			runtime.Gosched()
			time.Sleep(1 * time.Millisecond)
			json, changed := agent.readInstancesJsonFileIfChanged()
			if !changed {
				continue
			}

			currentInstances, err := readInstances(json, agent.Logger)
			if err != nil {
				agent.Warnf("Failed parsing json %s: %v Trying again...\n", err, string(json))
				continue
			}

			for instanceIdentifier, _ := range knownInstances {
				_, present := currentInstances[instanceIdentifier]
				if present {
					continue
				}

				delete(knownInstances, instanceIdentifier)
				agent.Infof("Removing stale instance %v", instanceIdentifier)
			}

			for _, instance := range currentInstances {
				_, present := knownInstances[instance.identifier()]
				if present {
					continue
				}

				knownInstances[instance.identifier()] = true
				agent.Infof("Adding new instance %v", instance.identifier())
				instancesChan <- &instance
			}
		}
	}

	go pollInstancesJson()
	return instancesChan
}
