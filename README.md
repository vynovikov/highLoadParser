## Changelog 
- operation process refactored from being asyncronous to synchronous
- architecture switched from hexagonal to clean
- TLS pair generation moved to Dockerfile

## Under construction

#### What is highLoadParser?
It is version of [postParser](https://github.com/vynovikov/postParser) group adapted for high load.  Kubernetes is used for self-healing and autoscaling. Kafka is for highly performant communication. 

#### Architecture
HighLoadParser has clean architecture. Such architecture is redundant for a service of this scale and is chosen for study purposes

![architecture](forManual/architecture.png)

#### Structured logging
Every architecture layer records log element for every request.  Log chain is stored by logger service in json format for convenient    display using ELK

![logChain](forManual/logChain.png)

#### Adaptation
Previously build [postParser](https://github.com/vynovikov/postParser) goes through adaptation steps:
- Multi-threading is removed. Processing uses single thread. For simplicity and kubernetes resource allocation clearance. Tuning for incoming load is provided by Kubernetes autoscaling
- gRPC transmitters/receivers are replaced with more performant Apache Kafka producers/consumers
- Parser obtains metric handling module necessary for autoscaling purposes
