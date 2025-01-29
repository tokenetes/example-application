# Alpha Stocks

Alpha Stocks is a sample application that implements transaction tokens(TraTs) using Tokenetes. It runs on a Kubernetes cluster and can serve as a reference for integrating Tokenetes into other projects. The application has the following architecture:

~~~
                                    ╔════════════════════════╗                                                              
                                    ║                        ║                                                              
                                    ║                        ║                                                              
                                    ║                        ║                                                              
                                    ║ Tokenetes (Transaction ║                                                              
                                    ║    Tokens Service)     ║                                                              
                                    ║                        ║                                                              
                                    ║                        ║                                                              
                                    ║                        ║                                                              
                                    ║                        ║                                                              
                                    ╚════════════════════════╝                                                              
                                                 ▲                                                                          
                      ┌ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─│─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ 
                                                 │                                         ┌────────────────────────┐      │
                      │                          │                                         │                        │       
                                                 │                                         │                        │      │
                      │                          │                                         │                        │       
                                                 │                                         │                        │      │
                      │                          │                          ┌─────────────▶│     Stocks Service     │       
                                                 │                          │              │                        │      │
                      │                          │                          │              │                        │       
                                                 │                          │              │                        │      │
┌────────────┐        │                          │                          │              │                        │       
│            │                      ┌────────────────────────┐              │              └────────────────────────┘      │
│            │        │             │                        │              │                           ▲                   
│            │                      │                        │              │                           │                  │
│            │        │             │                        │              │                           │                   
│            │                      │                        │              │                           │                  │
│    User    │────────┼────────────▶│      API Gateway       │──────────────┤                           │                   
│            │                      │                        │              │                           │                  │
│            │        │             │                        │              │                           │                   
│            │                      │                        │              │                           │                  │
│            │        │             │                        │              │                           │                   
│            │                      └────────────────────────┘              │                           │                  │
└────────────┘        │                          │                          │              ┌────────────────────────┐       
       │                                         │                          │              │                        │      │
       │              │                          │                          │              │                        │       
       │                                         │                          │              │                        │      │
       │              │                          │                          │              │                        │       
       │                                         │                          └─────────────▶│     Order Service      │      │
       │              │                          │                                         │                        │       
       │                                         │                                         │                        │      │
       │              │                          │                                         │                        │       
       │                                         │                                         │                        │      │
       │              │                          │                                         └────────────────────────┘       
       │                                         │                                                                         │
       │              │                          │                                                                
       │               ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ┼ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ┘
       │                                         ▼                                                                          
       │                            ┌────────────────────────┐                                                              
       │                            │                        │                                                              
       │                            │                        │                                                              
       │                            │                        │                                                              
       │                            │   Dex OpenID Connect   │                                                              
       └───────────────────────────▶│   Identity Provider    │                                                              
                                    │                        │                                                              
                                    │                        │                                                              
                                    │                        │                                                              
                                    │                        │                                                              
                                    └────────────────────────┘                                                              
~~~

As shown in the diagram above, the API Gateway integrates with the Tokenetes service to obtain TraTs that it can use to assure identity and context in its calls downstream, to the Order and Stocks services. The Order Service also calls the Stocks Service and passes the TraT it received from the API Gateway to the Stocks Service. Because TraTs can be passed between downstream services, they can assure identity and call context in arbitrarily deep call chains. The short-lived nature of TraTs makes them relatively immune to replay attacks (unless the replay happens really quickly, and the replay is exactly the same as the information in the TraT).

## How to Run

## Backend

### Prerequisites

Ensure Kubernetes is installed and correctly configured on your machine before executing these commands.

The application SPIRE installation is set up with "docker-desktop" as the cluster name. This is specific to the Kubernetes setup on Docker Desktop. If you are using a different Kubernetes cluster, adjust the SPIRE configurations to match your specific environment:

- Navigate to the SPIRE configuration directory at `deploy/spire`.

- Edit the configuration files to replace "docker-desktop" with your actual cluster name.

- Review and modify other settings as needed to align with your cluster requirements.

### Deploying Services

Navigate to the `deploy` directory and run the below command to deploy the services:

```bash
./deploy.sh
```

### Cleaning Up

You can remove all generated Kubernetes resources using the command below:

```bash
./destroy.sh
```

### OIDC Authentication via Dex

The application uses Dex as its OIDC provider, configured at `deploy/alpha-stocks-dev/configs/dex-config.yaml`. If you need to add clients, update secrets, or manage users, please update this file as necessary.

### SPIRE Identity Management

The application incorporates SPIRE(the SPIFFE Runtime Environment) for workload identity management, with configurations located at `deploy/spire/`. To adjust service identities, modify configurations, or manage workload registrations, please refer to and update the appropriate files within the directory.


## Client(Frontend)

To start the client, navigate to the `frontend` directory and follow these steps:

### Install dependencies:

```bash
npm install
```

### Start the frontend server:

```bash
npm start
```

&nbsp;

For more detailed instructions refer to the service-specific README files in their respective directories.

## Tokenetes Documentation
For detailed documentation and setup guides of tokenetes please visit tokenetes official documentation page: [tokenetes.io](https://tokenetes.io)

## Contribute to Tokenetes
Contributions to the project are welcome, including feature enhancements, bug fixes, and documentation improvements.