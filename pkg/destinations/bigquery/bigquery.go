package bigquery

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/rs/zerolog/log"
	"github.com/scratchdata/scratchdata/pkg/util"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type BigQueryServer struct {
	CredentialsJsonString string `mapstructure:"credentials_json" schema:"credentials_json" form:"label:Credentials JSON String,type:textarea"`
	Location              string `mapstructure:"location" schema:"location" form:"label:Location,type:text"`

	// this will be implemented during runtime
	GCSBucketName string `mapstructure:"gcs_bucket_name" schema:"gcs_bucket_name" form:"label:GCS Bucket Name,type:text"`
	GCSFilePrefix string `mapstructure:"gcs_file_prefix" schema:"gcs_file_prefix" form:"label:GCS File Prefix,type:text"`

	DeleteFromGCS bool `mapstructure:"delete_from_gcs" schema:"delete_from_gcs" form:"label:Delete From GCS,type:bool"`

	Credentials *google.Credentials
	conn        *bigquery.Client
}

func openConn(s *BigQueryServer) (*bigquery.Client, error) {
	ctx := context.Background()

	credentials, err := google.CredentialsFromJSON(ctx, []byte(s.CredentialsJsonString), bigquery.Scope)
	if err != nil {
		log.Error().Err(err).Msg("bigquery credentials error")
		return nil, err
	}

	client, err := bigquery.NewClient(ctx, credentials.ProjectID, option.WithCredentials(credentials))
	if err != nil {
		log.Error().Err(err).Msg("bigquery conn error")
		return nil, err
	}

	s.Credentials = credentials

	log.Info().Msg("Connected to BigQuery")
	return client, nil
}

func OpenServer(settings map[string]interface{}) (*BigQueryServer, error) {
	srv := util.ConfigToStruct[BigQueryServer](settings)

	conn, err := openConn(srv)
	if err != nil {
		log.Error().Err(err).Msg("BigQuery OpenServer Error")
		return nil, err
	}
	srv.conn = conn
	return srv, nil
}

func (s *BigQueryServer) Close() error {
	return s.conn.Close()
}
