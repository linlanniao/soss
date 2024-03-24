package controller

type Controller struct {
	client IRemoveStorageClient
	config *Config
}

func NewController(client IRemoveStorageClient) *Controller {
	return &Controller{client: client}
}

func (c *Controller) List(prefix string) error {
	return c.client.List(c.config.Bucket, prefix)
}

func (c *Controller) Upload(prefix string, files ...string) error {
	return c.client.Upload(c.config.Bucket, prefix, files)
}

func (c *Controller) Download(outputDir string, files ...string) error {
	return c.client.Download(c.config.Bucket, files, outputDir)
}
