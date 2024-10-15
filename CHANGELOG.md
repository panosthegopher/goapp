# 2024/03/29

v0.1.0: Cloning the original codebase from pfcafe/goapp repo, to work on the PF's knowledge assessment (Go).

v0.1.1: Problem 1 -> Fixing printing issue and also providing the minimum unit tests to provide a basic validation in the updated code.

v1.1.1: 
    a. Feature A -> Modifing the random string generator to generate only hex values and verify its accuracy and resource usage by creating a test and a benchmark run.
    b. Feature B -> Extent the API to also return the Hex value in WS connection.

v1.1.2: Providing an enhancement where reading values from a configuration file. This will be needed later for the pproof and client, as well as if/when this goApp was moved to the Cloud.

v2.0.0: Adding a fix for the memory usage issue which has been observed during many WS sessions. Providing a solution for better performance on concurrent executions and without leaving resources uncleaned.
