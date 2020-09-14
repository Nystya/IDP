# Scout
Repository for my IDP project

This service is a freelancing platform that allows employers to find suitable freelancers for their projects. 

The platform is composed of four services: 
1. Authentication & Reverse Proxy - Written using Django and uses PostgreSQL for user accounts storage. It uses GRPC to redirect requests to the appropiate service.
2. A service that handles all interactions with Jobs - Written in Go and uses Neo4j for data storage.
3. A service that handles all interactions with user Profiles - Written in Go and uses Neo4j for data storage.
4. A service that handles metrics - Written in Go, it listens on an MQTT message queue and writes the data in an InfluxDB instance.
