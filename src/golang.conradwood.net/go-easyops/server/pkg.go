/*
implements the server-side of gRPC.

As part of the https://www.yacloud.eu cloud the server registers itself upon startup and thus becomes accessible
for other servers and clients.

# Registration

Typically, servers register at the yacloud registry (registry.yacloud.eu).
It is also possible (albeit discouraged) to run a private cloud with the yacloud software.

The server needs at least a registry and authentication service(s).

Servers can also be run in standalone mode (see -ge_standalone option). If so, they need no external services, but are limited to local clients and no load-balancing, failover or authentication.

Multiple instances of the same server may register at the registry. Which server is chosen for any one rpc is handled in the client (see client package).

# Types

The Server supports 3 types of services:

  - gRPC (default), see NewServerDef()

  - TCP, see NewTCPServerDef(name)

  - HTML, deprecated.

  - Status, information about the status of this instance. Automatically added.

# gRPC server

this is the default and most commonly used service type. The server registers with the registry and claims to support the gRPC protocol. This implementation also automatically adds 'Status' to any gRPC servers.
gRPC calls are authenticated, that is, all calls require either a user or a service in the context. Calls with neither of the two will be rejected. Whilst this behaviour can be changed, it should only be required for low-level services. For example, the registry must accept rpc calls without authentication, because the first calls to the registry will be a lookup for the auth-service, which is required to create user objects in contexts. (chicken and egg...)

a server is typically started like this:

	sd := server.NewServerDef()
	sd.SetPort(4100)
	sd.SetRegister(server.Register(
	func(g *grpc.Server) error {
	  pb.RegisterEchoServiceServer(g, &echoServer{})
	  return nil
	},
	))
	err := server.ServerStartup(sd)

the "echoServer" is required to implement the service definition as defined by the .proto file.

# TCP server

this is a server that exposes a proprietary TCP connection on a port. The server instance registers the port at the registry, so that edge proxies (like h2gproxy) may proxy tcp connections to instances.

# Status server

this exposes an http api. Help can be found at [instance:port]/internal/help, for example https://10.1.1.1:4100/internal/help. Status is a bit of a misnomer. With the status api, one can:
  - shutdown the service
  - query the instance health
  - query the build information, including instance version
  - retrieve metrics (compatible with prometheus)
  - clear cache(s) of the instances
  - change commandline flags on-the-fly
  - view current outbound grpc Connections of this instance
  - view registered (as supposed to current) grpc connetions
  - download go debug information of this instance
*/
package server
