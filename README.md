# c_tlse_challenge
## Goal
Build a serverless API that, given an IP, returns geolocation data

## Components
- AWS Lambda to do the work
- AWS Gateway to serve as an interface between our lambda function and external requests
- The solution is desployed in AWS (free-tier) by using Terraform


request --> AWS Gateway --> AWS Lambda (golang) --> External API to enrich data

## TODO
- Add unit tests
- Implement a fallback option in case the external api is down. One option could be MindMax's database
- Restrict the calls so it only accepts GET /queries/{ip}