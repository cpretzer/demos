package main

import (
	"fmt"
	metal "github.com/pulumi/pulumi-equinix-metal/sdk/v2/go/equinix"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"log"
)

const (
	NODE_PREFIX = "test"
	VM_SIZE     = "s3.xlarge.x86"
	DATACENTER  = "sv"
	NODE_OS     = "ubuntu_20_10"
	NODE_PLAN   = "hourly"
)

type Nodes struct {
	Node map[pulumi.StringOutput]pulumi.StringPtrOutput
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		linkerdProject := "Linkerd"
		cfg := config.New(ctx, "")

		result, err := metal.LookupProject(ctx, &metal.LookupProjectArgs{
			Name: &linkerdProject,
		}, nil)
		if err != nil {
			log.Panicf("Error %v", err)
			return err
		}

		if result == nil {
			log.Panicf("No project named %s", linkerdProject)
		}

		log.Printf("Got project %s, with id %s", result.Name, result.ProjectId)

		controlPlaneNode := createControlPlaneNode(
			NODE_PREFIX,
			VM_SIZE,
			DATACENTER,
			NODE_OS,
			NODE_PLAN,
			result.ProjectId,
			ctx,
		)
		// should this be in createNode?
		// make this configurable to a limit
		nodeMap := make(map[pulumi.StringOutput]pulumi.StringPtrOutput)

		workerNodeCount := cfg.GetInt("workerNodeCount")

		log.Printf("Found workerNodeCount %d", workerNodeCount)

		if workerNodeCount == 0 {
			workerNodeCount = 4
		}

		for i := 1; i <= workerNodeCount; i++ {

			nodeName := fmt.Sprintf("%s-%0d", NODE_PREFIX, i)
			da, err := createDeviceArgs(
				nodeName,        // make this configurable with a default
				"s3.xlarge.x86", // big bada boom, this can fail if no servers are available
				"sv",
				"ubuntu_20_10",
				"hourly",
				result.ProjectId)

			if err != nil {
				log.Printf("Error creating device args: %s", err)
			}

			node, err := createNode(ctx, nodeName, da)

			if err != nil {
				log.Printf("Error creating node: %s", err)
			}

			if node == nil {
				log.Print("No testDevice?")
			}

			log.Printf("device state: %+v", node)

			nodeMap[node.Hostname] = pulumi.StringPtrOutput(node.ID())
		}

		// ctx.Export("nodes", pulumi.ToStringMapOutput(nodeMap))

		// allAddresses := controlPlaneNode.IpAddresses.ApplyT(func(ipAddresses []metal.DeviceIpAddress) []string {
		// 	allAddresses := make([]string, 2)
		// 	for _, v := range ipAddresses {
		// 		allAddresses = append(allAddresses, v.ReservationIds...)
		// 	}
		// 	if err != nil {
		// 		log.Panicf("Error looking up devices %s", err)
		// 	}

		// 	// deviceResult
		// 	return allAddresses
		// }).(pulumi.StringArrayOutput)

		ctx.Export("control-plane-ip", controlPlaneNode.IpAddresses)
		ctx.Export("testing", pulumi.String("exported"))

		// dev, err := metal.GetDevice(ctx,
		// 	fmt.Sprintf("%s-control-plane", NODE_PREFIX),
		// 	controlPlaneNode.ID().ApplyT(func(id pulumi.ID) *string {
		// 		idString := fmt.Sprintf("%v", id)
		// 		return &idString
		// 	}).(pulumi.IDInput), nil)

		// if err != nil {
		// 	log.Panicf("Error looking up devices %s", err)
		// }

		// log.Printf("%s", &dev.HardwareReservationId)

		return nil
	})
}

func createControlPlaneNode(nodePrefix string, vmSize string, dataCenter string,
	os string, plan string, projectId string, ctx *pulumi.Context) *metal.Device {
	nodeName := fmt.Sprintf("%s-control-plane", nodePrefix)
	da, err := createDeviceArgs(
		nodeName, // make this configurable with a default
		vmSize,   // big bada boom, this can fail if no servers are available
		dataCenter,
		os,
		plan,
		projectId)

	if err != nil {
		log.Printf("Error creating device args: %s", err)
	}

	node, err := createNode(ctx, nodeName, da)

	if err != nil {
		log.Printf("Error creating node: %s", err)
	}

	if node == nil {
		log.Print("No testDevice?")
	}

	log.Printf("device state: %+v", node.State.OutputState)

	return node
}

func createNode(ctx *pulumi.Context, nodeName string,
	deviceArgs *metal.DeviceArgs) (*metal.Device, error) {

	node, err := metal.NewDevice(ctx,
		nodeName,
		deviceArgs,
		pulumi.DeleteBeforeReplace(true))

	if err != nil {
		return nil, err
	}

	ctx.Export(nodeName, node)

	return node, nil
}

func createDeviceArgs(hostname string, plan string, metro string, os string,
	billingCycle string, projectId string) (*metal.DeviceArgs, error) {
	return &metal.DeviceArgs{
		Hostname:        pulumi.String(hostname),
		Plan:            pulumi.String(plan),
		Metro:           pulumi.String(metro),
		OperatingSystem: pulumi.String(os),
		BillingCycle:    pulumi.String(billingCycle),
		ProjectId:       pulumi.String(projectId),
	}, nil
}
