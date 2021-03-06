package fake

import (
	"github.com/Aptomi/aptomi/pkg/engine/resolve"
	"github.com/Aptomi/aptomi/pkg/event"
	"github.com/Aptomi/aptomi/pkg/external"
	"github.com/Aptomi/aptomi/pkg/lang"
	"github.com/Aptomi/aptomi/pkg/plugin"
	"github.com/Aptomi/aptomi/pkg/util"
	"time"
)

type noOpPlugin struct {
	sleepTime time.Duration
}

var _ plugin.ClusterPlugin = &noOpPlugin{}
var _ plugin.CodePlugin = &noOpPlugin{}

// NewNoOpClusterPlugin returns fake cluster plugin which does nothing, except sleeping a given time amount on every action
func NewNoOpClusterPlugin(sleepTime time.Duration) plugin.ClusterPlugin {
	return &noOpPlugin{
		sleepTime: sleepTime,
	}
}

// NewNoOpCodePlugin returns fake code plugin which does nothing, except sleeping a given time amount on every action
func NewNoOpCodePlugin(sleepTime time.Duration) plugin.CodePlugin {
	return &noOpPlugin{
		sleepTime: sleepTime,
	}
}

func (plugin *noOpPlugin) Validate() error {
	return nil
}

func (plugin *noOpPlugin) Cleanup() error {
	return nil
}

func (plugin *noOpPlugin) Create(deployName string, params util.NestedParameterMap, eventLog *event.Log) error {
	time.Sleep(plugin.sleepTime)
	return nil
}

func (plugin *noOpPlugin) Update(deployName string, params util.NestedParameterMap, eventLog *event.Log) error {
	time.Sleep(plugin.sleepTime)
	return nil
}

func (plugin *noOpPlugin) Destroy(deployName string, params util.NestedParameterMap, eventLog *event.Log) error {
	time.Sleep(plugin.sleepTime)
	return nil
}

func (plugin *noOpPlugin) Endpoints(deployName string, params util.NestedParameterMap, eventLog *event.Log) (map[string]string, error) {
	time.Sleep(plugin.sleepTime)
	return make(map[string]string), nil
}

func (plugin *noOpPlugin) Resources(deployName string, params util.NestedParameterMap, eventLog *event.Log) (plugin.Resources, error) {
	return nil, nil
}

func (plugin *noOpPlugin) Status(deployName string, params util.NestedParameterMap, eventLog *event.Log) (bool, error) {
	return true, nil
}

func (plugin *noOpPlugin) Process(desiredPolicy *lang.Policy, desiredState *resolve.PolicyResolution, externalData *external.Data, eventLog *event.Log) error {
	return nil
}
