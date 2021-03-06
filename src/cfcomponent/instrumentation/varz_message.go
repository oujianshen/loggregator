package instrumentation

import (
	"runtime"
)

type VarzMessage struct {
	Name string `json:"name"`
	NumCpus int `json:"numCPUS"`
	Contexts []Context `json:"contexts"`
}

func NewVarzMessage(name string, instrumentables []Instrumentable) *VarzMessage {
	contexts := make([]Context, len(instrumentables))
	for i, instrumentable := range instrumentables {
		contexts[i] = instrumentable.Emit()
	}
	return &VarzMessage{name, runtime.NumCPU(), contexts}
}
