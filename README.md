## Under construction

#### What is highLoadParser?
It is high load version of [postParser](https://github.com/vynovikov/postParser). Kubernetes is used for self-healing and autoscaling. Kafka is for highly performant communication. ClickHouse is for storing large amounts of metadata and logs.
![highLoad group](forManual/highLoad.png)

#### Adaptation
Previously built [postParser](https://github.com/vynovikov/postParser) is changed according to following steps:
- Multi-threading is removed. Processing is carried out in single thread way for simplicity and kubernetes resource allocation clearance. Tuning for incoming load is provided by Kubernetes autoscaling
- gRPC transmitters/receivers are replaced with more performant Apache Kafka producers/consumers
- Parser is obtained metric handling module necessary for autoscaling purposes
