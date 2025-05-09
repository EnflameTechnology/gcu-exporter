package collector

import (
	"testing"

	podresourcesapi "k8s.io/kubelet/pkg/apis/podresources/v1alpha1"
)

func TestCreateDevicePodMap(t *testing.T) {
	devicePods := podresourcesapi.ListPodResourcesResponse{
		PodResources: []*podresourcesapi.PodResources{
			{
				Name:      "pod1",
				Namespace: "namespace1",
				Containers: []*podresourcesapi.ContainerResources{
					{
						Name: "container1",
						Devices: []*podresourcesapi.ContainerDevices{
							{
								ResourceName: enflameDeviceResourceName,
								DeviceIds:    []string{"uuid1"},
							},
						},
					},
				},
			},
		},
	}
	deviceToPodMap := createDevicePodMap(devicePods, false)
	if deviceToPodMap == nil {
		t.Errorf("createDevicePodMap failed, got: nil, want: non-nil map")
	}
	if _, ok := deviceToPodMap["uuid1"]; !ok {
		t.Errorf("createDevicePodMap failed, got: no entry for uuid1, want: entry for uuid1")
	}
}
