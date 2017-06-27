package gobetween

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"gb_host": "GoBetween Server API Host or IP",
	}
}

// Provider returns a terraform.ResourceProvider
func Provider() terraform.ResourceProvider {
	// The provider definition
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"gobetween_server": resourceGobetweenServer(),
		},

		ConfigureFunc: providerConfigure,
	}
}

type GBProvider struct {
	Client *GbClient
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	host := d.Get("host").(string)
	port := d.Get("port").(int)

	c := &GbClient{}
	c.Init(fmt.Sprintf("http://%s:%d", host, port), "", "")

	err := c.Api.GetSystemInfo()
	if err != nil {
		return GBProvider{}, err
	}

	return GBProvider{Client: c}, nil
}
