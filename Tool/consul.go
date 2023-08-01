// Tool/consul.go
package Tool

import (
	"github.com/hashicorp/consul/api"
)

var client *api.Client
var res api.AgentServiceRegistration

// RegService:服务注册
func RegService(address string, id string, name string, serviceIP string, servicePort int, Interval string, chickHttp string, tag ...string) error {
	// 用于客户端的配置
	config := api.DefaultConfig()
	config.Address = address
	// 用于服务注册
	res = api.AgentServiceRegistration{}
	res.ID = id
	res.Name = name
	res.Address = serviceIP
	res.Port = servicePort
	res.Tags = tag
	chicks := api.AgentServiceCheck{}
	// 间隔时间
	chicks.Interval = Interval
	chicks.HTTP = chickHttp
	// 赋值
	res.Check = &chicks

	// 创建客户端
	var err error
	client, err = api.NewClient(config)
	if err != nil {
		return err
	}
	err = client.Agent().ServiceRegister(&res)
	if err != nil {
		return err
	}
	return nil
}

// LogOutServer 程序关闭后解除注册
func LogOutServer() {
	_ = client.Agent().ServiceDeregister(res.ID)
}
