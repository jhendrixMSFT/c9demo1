package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2018-06-01/containerinstance"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
)

const (
	resourceGroup  = "demoresgroup1"
	containerGroup = "democontainergroup1"
	location       = "WestUS"
	containerName  = "democontainer1"
	containerImage = "appsvc/sample-hello-world:latest"
)

func main() {
	// create an authorizer from the following environment variables
	// AZURE_CLIENT_ID
	// AZURE_CLIENT_SECRET
	// AZURE_TENANT_ID
	// other types of authorization are supported, see the docs for auth.NewAuthorizerFromEnvironment()
	a, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		panic(err)
	}
	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")

	// create a resource group
	rgClient := resources.NewGroupsClient(subscriptionID)
	rgClient.Authorizer = a
	_, err = rgClient.CreateOrUpdate(context.Background(), resourceGroup, resources.Group{
		Location: to.StringPtr(location),
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("created resource group %s\n", resourceGroup)

	// delete the resource group at the end of the demo, fire-and-forget async operation
	defer rgClient.Delete(context.Background(), resourceGroup)

	// create the client for managing container groups and attach the authorizer
	cgClient := containerinstance.NewContainerGroupsClient(subscriptionID)
	cgClient.Authorizer = a

	// start async operation
	future, err := cgClient.CreateOrUpdate(context.Background(), resourceGroup, containerGroup, containerinstance.ContainerGroup{
		Location: to.StringPtr(location),
		ContainerGroupProperties: &containerinstance.ContainerGroupProperties{
			OsType: containerinstance.Linux,
			Containers: &[]containerinstance.Container{
				containerinstance.Container{
					Name: to.StringPtr(containerName),
					ContainerProperties: &containerinstance.ContainerProperties{
						Image: to.StringPtr(containerImage),
						Resources: &containerinstance.ResourceRequirements{
							Requests: &containerinstance.ResourceRequests{
								MemoryInGB: to.Float64Ptr(4.0),
								CPU:        to.Float64Ptr(2),
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("begin deployment of container group %s, waiting for deployment to complete...", containerGroup)

	// block until async operation is complete
	// default wait is 15 minutes, adjust as required
	// client.PollingDuration = 20 * time.Minute
	err = future.WaitForCompletionRef(context.Background(), cgClient.Client)
	if err != nil {
		panic(err)
	}
	fmt.Println("done!")

	// or perform custom polling on the async operation
	/*for done, err := future.Done(client); !done; done, err = future.Done(client) {
		if err != nil {
			panic(err)
		}
		// do some other stuff while waiting...
		fmt.Println(future.Status())
	}*/

	// list by page
	for page, err := cgClient.List(context.Background()); page.NotDone(); err = page.Next() {
		if err != nil {
			panic(err)
		}
		for _, cg := range page.Values() {
			fmt.Printf("found container group %s\n", *cg.Name)
		}
	}

	// or use the page iterator to seamlessly transition across pages
	/*for iter, err := cgClient.ListComplete(context.Background()); iter.NotDone(); err = iter.Next() {
		if err != nil {
			panic(err)
		}
		fmt.Println(*iter.Value().Name)
	}*/
}
