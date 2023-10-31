# gAgents Registry Service

The gAgents Registry Service is a vital part of the distributed network architecture, serving as the hub for agent registration and event dissemination. Agents register with the registry, declaring their capabilities and services. When an agent seeks a particular service, the registry emits an event to the corresponding service providers through a Kafka topic, facilitating a direct communication between the requesting agent and the service providers.

Inorder to run the registry you will need to create a `.env` file in the root directory of the registry and provide an `OPENAI_API_KEY`

## Registration Protocol

Agents get registered via a `RegisterAgent` RPC call, passing an `AgentRegistration` message. Here's the `AgentRegistration` message structure:

```proto
message AgentRegistration {
    string unique_id = 1; 
    string grpc_address = 2; 
    string agent_type = 3; 
    repeated string capabilities = 4; 
    Metadata metadata = 5;
    Status status = 6;
    Authentication auth_data = 7;
    ContactInformation contact_info = 8;
    repeated string tags = 10;
}
```

Field descriptions:

- `unique_id`: A unique identifier for the agent.
- `grpc_address`: The address (including port) of the agent's gRPC server.
- `agent_type`: The role or type of the agent.
- `capabilities`: A list of services or capabilities provided by the agent.
- `metadata`: Additional metadata about the agent.
- `status`: The current status of the agent.
- `auth_data`: Authentication data for the agent.
- `contact_info`: Contact information for the agent.
- `tags`: A list of tags or labels associated with the agent.

## Example Registration

Here's an example `AgentRegistration` message formatted as a JSON object:

```json
{
    "unique_id": "agent-12345",
    "grpc_address": "localhost:50051",
    "agent_type": "Service Provider",
    "capabilities": ["negotiation", "delivery"],
    "metadata": {
        "software_version": "1.0.0",
        "location": {
            "latitude": 51.509865,
            "longitude": -0.118092
        },
        "region": "London",
        "organization": "ACME Corp"
    },
    "status": {
        "current_status": "ACTIVE"
    },
    "auth_data": {
        "token": "some-auth-token",
        "public_key": "some-public-key"
    },
    "contact_info": {
        "email": "agent@example.com",
        "secondary_channel": "http://example.com/contact"
    },
    "tags": ["urgent", "high-priority"]
}
```

## Event Dissemination

Post registration, when an agent requests a service, the registry emits an event to the Kafka topic corresponding to the requested service. Service providers subscribed to the topic receive the event, initiating a direct communication with the requesting agent to fulfill the service request. 

This event-driven mechanism ensures a real-time, decentralized interaction between agents, minimizing the latency in service provision and enhancing the network's overall efficiency.

## Using the Registry

1. **Registering an Agent**:
    - Agents register themselves with the registry by sending a `RegisterAgent` RPC call with their details.
    - The registry, upon successful registration, emits an event to the Kafka service.

2. **Requesting a Service**:
    - An agent sends a service request.
    - The registry emits an event to the Kafka topic corresponding to the requested service.

3. **Providing a Service**:
    - Service provider agents subscribed to the topic receive the event.
    - They initiate direct communication with the requesting agent to provide the requested service.
