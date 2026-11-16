// Copyright (c) 2022 Enflame. All Rights Reserved.

package collector

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	podresourcesapi "k8s.io/kubelet/pkg/apis/podresources/v1alpha1"
)

const enflameDeviceResourceName = "enflame.com/gcu"
const enflameVdeviceResourceName = "enflame.com/vgcu"

const (
	connectionTimeout   = 10 * time.Second
	podResourcesMaxSize = 1024 * 1024 * 16 // 16 Mb
)

// 存储map[uuid]index
var gpuUUIDMap = make(map[string]uint)

type devicePodInfo struct {
	name      string
	namespace string
	container string
}

// Helper function that creates a map of pod info for each device
func createDevicePodMap(devicePods podresourcesapi.ListPodResourcesResponse, clusterMap bool) map[string]devicePodInfo {
	deviceToPodMap := make(map[string]devicePodInfo)
	currentResourceName := enflameDeviceResourceName
	if clusterMap {
		currentResourceName = enflameVdeviceResourceName
	}

	for _, pod := range devicePods.GetPodResources() {
		for _, container := range pod.GetContainers() {
			for _, device := range container.GetDevices() {
				if device.GetResourceName() == currentResourceName {
					podInfo := devicePodInfo{
						name:      pod.GetName(),
						namespace: pod.GetNamespace(),
						container: container.GetName(),
					}
					for _, uuid := range device.GetDeviceIds() {
						deviceToPodMap[uuid] = podInfo
						/*
							index, err := strconv.Atoi(uuid)
							if err == nil {
								deviceToPodMap[index] = podInfo
							}
						*/
					}
				}
			}
		}
	}
	return deviceToPodMap
}

func getDevicePodInfo(socket string, clusterMap bool) (map[string]devicePodInfo, error) {
	_, err := os.Stat(socket)
	if err != nil {
		fmt.Println("get device pod info errored ", err)
		if os.IsNotExist(err) {
			deviceToPodMap := make(map[string]devicePodInfo)
			return deviceToPodMap, nil
		}
		return nil, err
	}
	devicePods, err := getListOfPods(socket)
	if err != nil {
		fmt.Println("failed to get devices Pod information:", err)
		return nil, fmt.Errorf("failed to get devices Pod information: %v", err)
	}
	return createDevicePodMap(*devicePods, clusterMap), nil

}

func connectToServer(socket string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, socket, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(podResourcesMaxSize)),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failure connecting to %s: %v", socket, err)
	}
	return conn, nil
}

func getListOfPods(socket string) (*podresourcesapi.ListPodResourcesResponse, error) {
	conn, err := connectToServer(socket)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := podresourcesapi.NewPodResourcesListerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	resp, err := client.List(ctx, &podresourcesapi.ListPodResourcesRequest{})
	if err != nil {
		return nil, fmt.Errorf("failure getting pod resources %v", err)
	}
	return resp, nil
}
