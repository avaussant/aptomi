package k8s

import (
	"fmt"
	"github.com/Aptomi/aptomi/pkg/event"
	"github.com/Aptomi/aptomi/pkg/util"
	"github.com/Aptomi/aptomi/pkg/util/retry"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	api "k8s.io/client-go/pkg/api/v1"
	"k8s.io/kubernetes/pkg/kubectl/resource"
	"strings"
	"time"
)

// EndpointsForManifests returns endpoints for specified manifest
func (p *Plugin) EndpointsForManifests(deployName, targetManifest string, eventLog *event.Log) (map[string]string, error) {
	kubeClient, err := p.NewClient()
	if err != nil {
		return nil, err
	}

	helmKube := p.NewHelmKube(deployName, eventLog)

	infos, err := helmKube.BuildUnstructured(p.Namespace, strings.NewReader(targetManifest))
	if err != nil {
		return nil, err
	}

	endpoints := make(map[string]string)

	for _, info := range infos {
		if info.Mapping.GroupVersionKind.Kind == "Service" { // nolint: goconst

			endpointsErr := p.addEndpointsFromService(kubeClient, info, endpoints)
			if endpointsErr != nil {
				return nil, endpointsErr
			}
		}
	}

	return endpoints, nil
}

// addEndpointsFromService searches for the available endpoints in specified service and writes them into provided map
func (p *Plugin) addEndpointsFromService(kubeClient kubernetes.Interface, info *resource.Info, endpoints map[string]string) error {
	service, getErr := kubeClient.CoreV1().Services(info.Namespace).Get(info.Name, meta.GetOptions{})
	if getErr != nil {
		return getErr
	}

	// todo(slukjanov): support not only node ports
	if service.Spec.Type == api.ServiceTypeNodePort {
		for _, port := range service.Spec.Ports {
			sURL := fmt.Sprintf("%s:%d", p.ExternalAddress, port.NodePort)
			addEndpointsForServicePort(port, sURL, endpoints)
		}
	} else if service.Spec.Type == api.ServiceTypeLoadBalancer {
		ingress := service.Status.LoadBalancer.Ingress

		// wait for LB external IP to be provisioned
		ok := retry.Do(90, 10*time.Second, func() bool {
			service, getErr = kubeClient.CoreV1().Services(info.Namespace).Get(info.Name, meta.GetOptions{})
			if getErr != nil {
				panic(fmt.Sprintf("Error while getting Service %s in namespace %s", info.Name, info.Namespace))
			}

			ingress = service.Status.LoadBalancer.Ingress
			if ingress == nil {
				return false
			}

			externalAddress := ""
			for _, entry := range ingress {
				if entry.Hostname != "" {
					externalAddress = entry.Hostname
				} else if entry.IP != "" {
					externalAddress = entry.IP
				}
				if externalAddress == "" {
					panic(fmt.Sprintf("Got empty LoadBalancerIngress for Service %s in namespace %s", info.Name, info.Namespace))
				} else {
					// handle only first ingress entry for LB
					break
				}
			}

			for _, port := range service.Spec.Ports {
				sURL := fmt.Sprintf("%s:%d", externalAddress, port.Port)
				addEndpointsForServicePort(port, sURL, endpoints)
			}

			return true
		})

		if ingress == nil || !ok {
			return fmt.Errorf("unable to get endpoints for Service type LoadBalancer (%s in %s)", info.Name, info.Name)
		}
	}

	return nil
}

func addEndpointsForServicePort(port api.ServicePort, sURL string, endpoints map[string]string) {
	// todo(slukjanov): could we somehow detect real schema? I think no :(
	if util.StringContainsAny(port.Name, "https") {
		sURL = "https://" + sURL
	} else if util.StringContainsAny(port.Name, "ui", "rest", "http", "grafana", "service") {
		sURL = "http://" + sURL
	}
	name := port.Name
	if len(name) == 0 {
		name = port.TargetPort.String()
	}
	endpoints[name] = sURL
}
