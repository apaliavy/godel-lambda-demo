# Demo application 

This application was built as part of Godel Serverless in Go presentation. 
The logic in the following: 

![alt text](https://github.com/apaliavy/godel-lambda-demo/blob/main/assets/img/diagram.png?raw=true)

We want to handle two types of files upload into S3 bucket: 
- XLS files;
- CSV files.

When upload happens we trigger a lambda function which reads it content and saves as json file in another bucket. 

JSON-handler lambda picks it up, parses and stores into RDS postgres instance. 

## Notice 

Application structure is super-simplified (no abstractions, problems with single responsibility for `model` package, etc).
We want to concentrate on lambda features, not on abstractions in the code. 

## Deployment: 

At the moment app doesn't support auto-deploy, please use the following command: 

```shell
make zip
```

to build application zip file. 

This zip file should be uploaded as app code using AWS UI. 
