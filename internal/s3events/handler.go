package s3events

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/apache/arrow/go/v7/parquet/file"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/eventbridge-s3-parquet/internal/events"
	"github.com/wolfeidau/eventbridge-s3-parquet/internal/events/s3created"
	"github.com/wolfeidau/eventbridge-s3-parquet/internal/flags"
)

type Handler struct {
	cfg flags.S3Events
}

func NewHandler(cfg flags.S3Events) *Handler {
	return &Handler{
		cfg: cfg,
	}
}

func (h *Handler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {

	event, err := events.ParseEvent(payload)
	if err != nil {
		return nil, err
	}

	switch v := event.Detail.(type) {
	case *s3created.ObjectCreated:
		return h.processCreated(ctx, v)
	}

	return nil, errors.New("failed to process event, unknown type")
}

func (h *Handler) processCreated(ctx context.Context, created *s3created.ObjectCreated) ([]byte, error) {

	// does the key have a suffix of .parquet

	if filepath.Ext(created.Object.Key) != ".parquet" {
		log.Ctx(ctx).Info().Str("key", created.Object.Key).Str("bucket", created.Bucket.Name).Msg("skipping non parquet file")

		return []byte(`{"msg": "ok"}`), nil
	}

	config, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	s3client := s3.NewFromConfig(config)

	res, err := s3client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(created.Bucket.Name),
		Key:    aws.String(created.Object.Key),
	})
	if err != nil {
		return nil, err
	}

	f, err := ioutil.TempFile("", "*.parquet")
	if err != nil {
		return nil, err
	}

	written, err := io.Copy(f, res.Body)
	if err != nil {
		return nil, err
	}

	// cleanup temp file
	defer os.Remove(f.Name())

	log.Ctx(ctx).Info().Int64("written", written).Msg("tempfile written")

	rdr, err := file.OpenParquetFile(f.Name(), true)
	if err != nil {
		return nil, err
	}

	log.Ctx(ctx).Info().Int("rowgroups", rdr.NumRowGroups()).Int("realcolumns", rdr.MetaData().Schema.Root().NumFields()).Msg("file info")

	return []byte(`{"msg": "ok"}`), nil
}
