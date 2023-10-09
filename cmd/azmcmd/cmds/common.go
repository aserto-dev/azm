package cmds

type Common struct {
	Host     string `name:"host" env:"ASERTO_DIR_SVC" default:"localhost:9292"`
	APIKey   string `name:"api-key" env:"ASERTO_DIR_KEY" default:""`
	TenantID string `name:"tenant-id" env:"ASERTO_TENANT_ID" default:""`
	Insecure bool   `name:"insecure" env:"ASERTO_SKIP_TLS_VERIFICATION" default:"false"`
}
