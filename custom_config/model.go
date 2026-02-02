package custom_config

import (
	"github.com/New-Tatthep/microservice/util/convutil"
	"github.com/New-Tatthep/microservice/util/stringutil"
)

type CustomConfiguration struct {
	Test string `mapstructure:"test" json:"test"`
	Api  map[string]struct {
		BaseUrl  string            `mapstructure:"base-url" json:"base_url"`
		Endpoint map[string]string `mapstructure:"endpoint" json:"endpoint"`
	} `mapstructure:"api" json:"api"`
	CloudLoginEndpoint         string `mapstructure:"cloud-login-endpoint" json:"cloud_login_endpoint"`
	CloudExtendEndpoint        string `mapstructure:"cloud-extend-endpoint" json:"cloud_extend_endpoint"`
	CloudResetPasswordEndpoint string `mapstructure:"cloud-reset-password-endpoint" json:"cloud_reset_password_endpoint"`
	DialTimeout                int64  `mapstructure:"dial-timeout" json:"dial_time_out"`
	DialToIp                   string `mapstructure:"dial-to-ip" json:"dial_to_ip"`
	SocketSessionNamespace     string `mapstructure:"socket-session-namespace" json:"socket_session_namespace"`
	SocketServerEndpoint       string `mapstructure:"socket-server-endpoint" json:"socket_server_endpoint"`
	SocketApiKey               string `mapstructure:"socket-api-key" json:"socket_api_key"`

	ProduceCloudAllBranchUrl      string `mapstructure:"produce-cloud-all-branch-url" json:"produce_cloud_all_branch_url"`
	ProduceCloudSpecificBranchUrl string `mapstructure:"produce-cloud-specific-branch-url" json:"produce_cloud_specific_branch_url"`
}

func (cfg *CustomConfiguration) String() string {
	return stringutil.Json(*cfg)
}

func (cfg *CustomConfiguration) ToMap() map[string]any {
	return convutil.Obj2Map(*cfg)
}
