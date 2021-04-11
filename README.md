### Resources
- Download the bank partner API from:
```
docker pull lendoab/interview-service:stable
```
- Run the bank partner API service:
```
docker run -p 8000:8000 lendoab/interview-service:stable
```
- Once you have the bank partner API container up and running you can access its API from 
```
http://localhost:8000/docs
```
or if you are running docker for windows from
```
http://<docker-machine-ip>:8000/docs
```

### Description
This exercise is a simplified version of Lendo’s domain, where a customer application is the
most important part of it. A customer application represents the intent of a customer to get
a loan. A job is used for retrieving the status of a certain application. Jobs are needed since
not all bank partners answer immediately.

In this exercise you will need to communicate towards a bank partner API. You have to send
a customer application to the bank partner which will assess it asynchronously and return it
in its initial pending status. Since the assessment of an application can take from 5 to 20
seconds, you will have to poll for updates on individual applications to make sure we keep
up to date records of all the applications in our persistent storage.

Your solution has to expose HTTP endpoints for creating a customer application, getting a
certain application and also getting all applications with a certain status. It will also be
responsible for polling the application status using the job endpoint described in the
documentation.

### Requirements
- We expect you to package your solution in a way that simplifies its execution.
- Feel free to decide which persistent datastore suits your solution.
- Feel free to implement your solution in a language that you are comfortable with.
- Provide good documentation for your API.
- Make sure the code is well tested.

### For experienced candidates we generally look for the points below
- Find a way to decouple the job polling mechanism by creating another microservice
to demonstrate asynchronous communication. This represents a simplified version
of Lendo’s domain and we believe that it is a good way for the candidate to show
their ability to separate concerns and design a small distributed system.
- Use a message broker for exchange messages between the microservices.
- Provide a K8s manifest or docker compose file.

### Further Notes
- The given assignment is a way for interviewers to evaluate the candidate’s skills and
engage in some discussions during the interview. For different experience levels,
different aspects on how the solution was built will be assessed. We expect the
candidate to have the exercise fresh in mind for a more fluent conversation.
- If there are any questions about the assignment, feel free to reach out to your HR
contact. They will put you in contact with the tech interviewer(s) assigned to your
application.
