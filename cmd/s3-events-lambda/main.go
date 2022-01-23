package main

import (
	"github.com/alecthomas/kong"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/wolfeidau/eventbridge-s3-parquet/internal/flags"
	"github.com/wolfeidau/eventbridge-s3-parquet/internal/s3events"
	lmw "github.com/wolfeidau/lambda-go-extras/middleware"
	"github.com/wolfeidau/lambda-go-extras/middleware/raw"
	zlog "github.com/wolfeidau/lambda-go-extras/middleware/zerolog"
)

var (
	version = "unknown"

	cli flags.S3Events
)

func main() {
	kong.Parse(&cli,
		kong.Vars{"version": version}, // bind a var for version
	)

	// build up a list of fields which will be included in all log messages
	flds := lmw.FieldMap{"version": version}

	ch := lmw.New(
		zlog.New(zlog.Fields(flds)), // assign a logger and bind it in the context
	)

	if cli.RawEventLogging {
		ch.Use(raw.New(raw.Fields(flds))) // if raw event logging is enabled dump everything to the log in and out
	}

	h := s3events.NewHandler(cli)

	// register our lambda handler with the middleware configured
	lambda.StartHandler(ch.Then(h))
}
