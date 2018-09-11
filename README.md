# Demo for https://channel9.msdn.com/Shows/Azure-Friday/Go-on-Azure-Part-5-Build-apps-with-the-Azure-SDK-for-Go

This demo app uses the Azure SDK for Go to create a container group.  It demonstrates the following concepts.

1. Using `auth.NewAuthorizerFromEnvironment()` for authorization.
2. Blocking on an asynchronous operation until it completes.
3. Optionally polling on the asynchronous operation while performing other work.
4. How to traverse paged results, either by page or via a page iterator.

# How to Build

1. Execute `go get github.com/jhendrixMSFT/c9demo1`
2. Execute `dep ensure`
3. Execute `go build`

# How to Run

Set the appropriate environment variables depending on the type of authentication you wish to use.
For authorization with a service principal set the following environment variables.
```
AZURE_CLIENT_ID
AZURE_CLIENT_SECRET
AZURE_TENANT_ID
```
More information about authorization can be found [here](https://github.com/azure/azure-sdk-for-go#more-authentication-details).

# License

This code is provided under the MIT license. See [LICENSE](./LICENSE) for details.
