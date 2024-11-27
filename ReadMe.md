MIT License

Copyright (c) 2024 Jonas Lambert

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

***

# Fenix Inception

## ExecutionWorker
ExecutionWorker has the responsibility to feed TestInstructions to be executed to the correct Connectors. For all Connector that is not belonging to Fenix Inception itself then PubSub is used as transportation method. For Connectors belonging to Fenix Inception then gRPC is used.
Another responsibility is to receive published available TestInstructions, TestInstructionContainers, Allowed Users and Template-address.
![Fenix Inception - Worker](./Documentation/FenixInception-Overview-NonDetailed-Worker.png "Fenix Inception - Worker")

The following environment variable is needed for ExecutionWorker to be able to run locally.

| Environment variable                          | Example value                                                            | comment                                     |
|-----------------------------------------------|--------------------------------------------------------------------------|---------------------------------------------|
| AuthClientId                                  | 944682210385-jokmr7b6fdllr6k76kfo2hagic7kfvnt.apps.googleusercontent.com |                                             |
| AuthClientSecret                              | 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                                      |                                             |
| ExecutionLocationForFenixGuiBuilderServer     | LOCALHOST_NODOCKER                                                       | LOCALHOST_NODOCKER, LOCALHOST_DOCKER or GCP |
| ExecutionLocationForFenixTestExecutionServer  | LOCALHOST_NODOCKER                                                       | LOCALHOST_NODOCKER, LOCALHOST_DOCKER or GCP |
| ExecutionLocationForWorker                    | LOCALHOST_NODOCKER                                                       | LOCALHOST_NODOCKER, LOCALHOST_DOCKER or GCP |
| ExecutionWorkerPort                           | 6671                                                                     |                                             |
| FenixExecutionServerAddress                   | 127.0.0.1                                                                |                                             |
| FenixExecutionServerPort                      | 6670                                                                     |                                             |
| FenixGuiBuilderServerAddress                  | fenixguitestcasebuilderserver-must-be-logged-in-ffwegrgwrg-lz.a.run.app  |                                             |
| FenixGuiBuilderServerPort                     | 443                                                                      |                                             |
| GcpProject                                    | mycloud-run-project                                                      |                                             |
| LocalServiceAccountPath                       | #                                                                        |                                             |
| LoggingLevel                                  | DebugLevel                                                               | DebugLevel, InfoLevel                       |
| TestInstructionExecutionPubSubTopicBase       | ProcessTestInstructionExecutionRequest                                   |                                             |
| TestInstructionExecutionPubSubTopicSchema     | ProcessTestInstructionExecutionRequestSchema                             |                                             |
| UsePubSubWhenSendingTestInstructionExecutions | true                                                                     |                                             |





