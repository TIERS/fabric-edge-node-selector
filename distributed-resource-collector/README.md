# Distributed Resource Collector & Heartbeat
Simple webserver that collects relevant system data, and measures the latency between the DRC host and all the hostnames received.

# v0.2
#### Resource Collection
 - Updated the data structure to better fit requirements
#### Latency Measurement
 - Retrieve a list of hostnames to measure latency
 - Latency measurer connets to hostname via SSH, runs an output command and verifies that the output is correct
 - The measurement is the comparison between two timestamps, one before the SSH is stablished and the second after the output is verified
#### Fabric APP Integration
 - Connects via configurable API REST to the Fabric APP
 - Publishes Data Collection and Latency Measures to the distributed ledger
#### CRON
 - Configurable CRON timer for Resource Collection and Latency Measurement
 - Automatically runs the Resource Collector and the Latency Measurement and publishes the data to the Fabric APP

![dca12149aa1985525ca122ecc29082378a1739bd230244fa8c823d770357de50](https://user-images.githubusercontent.com/16642619/154840557-980054ec-5a73-475a-bc7c-a0e564a3631f.jpg)



# v0.1
 - Performs required (known) data collection tasks
 - Outputs data as JSON via a GET handler

Data collection takes a fraction of a second to complete (Linux environment)

2022-01-DD HH:26:17.305794417 +0200 EET

2022-01-DD HH:26:17.547917812 +0200 EET

Web server and handler add very little overhead to the operation
![image](https://user-images.githubusercontent.com/16642619/151067744-43a6913b-775a-4c7e-91fc-db7e30474bda.png)

# Planned
 - Post data automatically to Central Service (Heartbeat)
 - Perform latency checking tasks
    - Receive (and maybe Store) parameters from Central Service

# Possible expansions
 - Add authentication mechanism (CA & Central Service) to securize endpoints
