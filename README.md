# eventbridge-s3-parquet

This project provides an example of how to use [Amazon S3](https://aws.amazon.com/s3/) Event Notifications with [Amazon EventBridge](https://aws.amazon.com/eventbridge/) with the lambda written in [Go](https://go.dev).

# Overview

This project is deployed using cloudformation, with a stack managing the S3 data storage, and a stack holding the EventBridge rule and lambda which is triggered each time an s3 object is created. These stacks are decoupled which is a nice change from the original s3 events options and avoids complex cross service IAM policies. 

If the file is a [parquet](https://parquet.apache.org/documentation/latest/) format file the lambda receiving the event will log some information using [Apache Arrow](https://arrow.apache.org/).

I have also included an example of using [oapi-codegen](https://github.com/deepmap/oapi-codegen) to generate Go types for the s3 `Object Created` [Amazon EventBridge OpenAPI 3 Schema](https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-schema.html) provided by AWS.

# References

* [New – Use Amazon S3 Event Notifications with Amazon EventBridge](https://aws.amazon.com/blogs/aws/new-use-amazon-s3-event-notifications-with-amazon-eventbridge/)

# License

This application is released under Apache 2.0 license and is copyright [Mark Wolfe](https://www.wolfe.id.au).
