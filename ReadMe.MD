# loan-process-orchestrator

### Temporal Workflow Engine
Temporal is an open source platform that enables enterprises to maintain execution resiliency for business processes in the event of an infrastructure failure. 

Temporal does the heavy lifting of ensuring execution in the face of danger, such as a momentary network outage or an unexpected slowdown in disk access. 

![WithTemporalAndWihtout.png](common%2FWithTemporalAndWihtout.png) 
All the operational stuff: Retries, recovery, persistent state, message queues, scaling, work distribution, tracing and monitoring are all handled by the platform itself.\


### About the Loan Origination Example
For this example using only three steps in the workflow for the Loan Origination (LO) process. \

Using 3 out of 7 Loan Origination Process for the reference implementation
https://allcloud.in/all-cloud-blog/7-stages-in-loan-origination-process

1. Loan Application
2. Credit Decision
3. Loan Funding


### Setup temporal
Following the guide here https://learn.temporal.io/getting_started/go/dev_environment/
Install temporal cli
```
curl -sSf https://temporal.download/cli.sh | sh
```

```
temporal server start-dev
```

Temporal WebUI can be accessed here http://localhost:8233/

### Build and Start

Build
```
go build
```

Start Worker and Web App in different terminals
```
./loan-process-orchestrator workerapp
```
```
./loan-process-orchestrator webapp
```
WebApp can be accessed here:
http://localhost:8080/swagger/index.html
![webapp.png](common%2Fapis.png) 

#### How to use the app

1. Submit a Loan Application using the payload
```
{
  "aadhaar_number": "8989090987876767",
  "account_no": "90908765434321",
  "applicant_name": "John Doe",
  "pan_number": "AEERWQ1231L"
}
```
Successful request will return the response 200 and response data would be as follows:
```
{
  "data": {
    "workflow_id": "24136d9c-685a-41a5-94b5-4e95d244b605",
    "run_id": "109d4315-e726-47ac-9048-a752598fbb9a"
  }
}
```

2. Check the applicaiton status (state) by using the workflow_id and run_id.
![loan_app_status.png](common%2Floan_app_status.png)

3. Approval Process using the following payload.
```
{
  
   "workflow_id": "24136d9c-685a-41a5-94b5-4e95d244b605",
   "run_id": "109d4315-e726-47ac-9048-a752598fbb9a"
   "approved": true
}
```
![approval_process_using_singals.png](common%2Fapproval_process_using_singals.png)

Same can be validated in the temporal webui as well
![loan_app_approved_workflow.png](common%2Floan_app_approved_workflow.png)

4. Check the status for the loan application after the approval
![loan_app_status_approved.png](common%2Floan_app_status_approved.png)


### More about Temporal
From architecting monoliths with a single database and centralized state — to microservices where everything is distributed across multiple containers, servers, data centers, and even continents. Distributing things solves scaling concerns, but introduces a whole new world of problems, many of which were previously solved by monoliths.
Temporal would be able to maintains the high level of reliability various microservice processes or taks while also providing tremendous amounts of visibility into running processes.

#### Reliability and Fault Tolerance
Temporal is designed for long-running and resilient workflows. It can handle failures, retries, and recoveries gracefully, ensuring that your workflows remain reliable even in the face of system failures or interruptions.

#### Scalability
Temporal is built for scalability. It can handle a large number of concurrent workflows and can scale horizontally as your application's needs grow. This makes it suitable for both small and large-scale distributed systems.

#### Stateful Workflows
Temporal supports stateful workflows, which means the state of a workflow can be maintained across different activities and over time. This is useful for scenarios where the outcome of one activity affects the subsequent steps in the workflow.

Check more here https://learn.temporal.io/ to learn about Temporal.
