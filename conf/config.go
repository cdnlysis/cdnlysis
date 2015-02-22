package conf

type Config struct {
	initalized bool

	Engine struct {
		Verbose bool
		Threads int
	}

	SyncProgress struct {
		Path string
	}

	S3 struct {
		Prefix    string
		AccessKey string
		SecretKey string
		Bucket    string
		Region    string
	}

	Logs struct {
		Prefix   string
		Location string
	}
}
