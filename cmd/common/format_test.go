package common

import (
	"github.com/Aptomi/aptomi/pkg/api"
	"github.com/Aptomi/aptomi/pkg/config"
	"github.com/Aptomi/aptomi/pkg/engine/apply/action"
	"github.com/Aptomi/aptomi/pkg/engine/apply/action/component"
	"github.com/Aptomi/aptomi/pkg/engine/resolve"
	"github.com/Aptomi/aptomi/pkg/lang"
	"github.com/Aptomi/aptomi/pkg/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormat_Text(t *testing.T) {
	cfg := &config.Client{Output: Text}

	{
		// with policy changes
		result := makePolicyUpdateResult(true)
		data, err := Format(cfg.Output, true, result)
		assert.Nil(t, err, "Format should work without error")
		expected := "Policy Generation\tAction Plan                                                     \n41 -> 42         \tUpdate Instances                                                \n                 \t  [*] cluster#ns#contract#context#keysresolved#component        \n                 \tCreate Instances                                                \n                 \t  [+] cluster#ns#contract#context#keysresolved#component        \n                 \tDestroy Instances                                               \n                 \t  [-] cluster#ns#contract#context#keysresolved#component        \n                 \tRemove Consumers                                                \n                 \t  [<] cluster#ns#contract#context#keysresolved#component = depId\n                 \tAdd Consumers                                                   \n                 \t  [>] cluster#ns#contract#context#keysresolved#component = depId\n                 \t                                                                "
		if !assert.Equal(t, expected, string(data), "Format should return expected table") {
			t.Log("Expected:\n", expected)
			t.Log("Found:\n", string(data))
			t.Fail()
		}
	}
	{
		// without policy changes
		result := makePolicyUpdateResult(false)
		data, err := Format(cfg.Output, true, result)
		assert.Nil(t, err, "Format should work without error")
		expected := "Policy Generation\tAction Plan                                                     \n42               \tUpdate Instances                                                \n                 \t  [*] cluster#ns#contract#context#keysresolved#component        \n                 \tCreate Instances                                                \n                 \t  [+] cluster#ns#contract#context#keysresolved#component        \n                 \tDestroy Instances                                               \n                 \t  [-] cluster#ns#contract#context#keysresolved#component        \n                 \tRemove Consumers                                                \n                 \t  [<] cluster#ns#contract#context#keysresolved#component = depId\n                 \tAdd Consumers                                                   \n                 \t  [>] cluster#ns#contract#context#keysresolved#component = depId\n                 \t                                                                "
		if !assert.Equal(t, expected, string(data), "Format should return expected table") {
			t.Log("Expected:\n", expected)
			t.Log("Found:\n", string(data))
			t.Fail()
		}
	}

	{
		// empty set of actions
		result := &api.PolicyUpdateResult{
			PolicyGeneration: 42,
			PlanAsText:       action.NewPlanAsText(),
		}
		data, err := Format(cfg.Output, true, result)
		assert.Nil(t, err, "Format should work without error")
		expected := "Policy Generation\tAction Plan\n42               \t(none)     "
		if !assert.Equal(t, expected, string(data), "Format should return expected table") {
			t.Log("Expected:\n", expected)
			t.Log("Found:\n", string(data))
			t.Fail()
		}
	}
}

func makePolicyUpdateResult(policyChanged bool) *api.PolicyUpdateResult {
	key := resolve.NewComponentInstanceKey(
		&lang.Cluster{Metadata: lang.Metadata{Name: "cluster"}},
		&lang.Contract{Metadata: lang.Metadata{Name: "contract", Namespace: "ns"}},
		&lang.Context{Name: "context"},
		[]string{"keysresolved"},
		&lang.Service{Metadata: lang.Metadata{Name: "service"}},
		&lang.ServiceComponent{Name: "component"},
	)

	paramsPrev := util.NestedParameterMap{"name": "valuePrev"}
	params := util.NestedParameterMap{"name": "value"}

	result := &api.PolicyUpdateResult{
		PolicyGeneration: 42,
		PolicyChanged:    policyChanged,
		PlanAsText: &action.PlanAsText{
			Actions: []util.NestedParameterMap{
				component.NewCreateAction(key.GetKey(), params).DescribeChanges(),
				component.NewUpdateAction(key.GetKey(), paramsPrev, params).DescribeChanges(),
				component.NewDeleteAction(key.GetKey(), paramsPrev).DescribeChanges(),
				component.NewAttachDependencyAction(key.GetKey(), "depId").DescribeChanges(),
				component.NewDetachDependencyAction(key.GetKey(), "depId").DescribeChanges(),
				component.NewEndpointsAction(key.GetKey()).DescribeChanges(),
			},
		},
	}
	return result
}
