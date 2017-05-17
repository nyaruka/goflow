package config

// Mailroom is our top level configuration object
type Mailroom struct {
	Port  int    `default:"8080"`
	DB    string `default:"postgres://mailroom@localhost/mailroom?sslmode=disable"`
	Redis string `default:"redis://localhost:6379/0"`

	Workers int `default:"24"`

	S3_Region       string `default:"us-east-1"`
	S3_Media_Bucket string `default:"courier-media"`
	S3_Media_Prefix string `default:"/media/"`

	AWS_Access_Key_ID     string `default:"missing_aws_access_key_id"`
	AWS_Secret_Access_Key string `default:"missing_aws_secret_access_key"`
}
