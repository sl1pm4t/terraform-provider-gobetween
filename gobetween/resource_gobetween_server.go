package gobetween

import (
	"github.com/hashicorp/terraform/helper/schema"
	gb "github.com/yyyar/gobetween/src/config"
)

func resourceGobetweenServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoBetweenServerCreate,
		// Update: resourceGoBetweenServerUpdate,
		Delete: resourceGoBetweenServerDelete,
		// Exists: resourceGoBetweenServerExists,
		Read: resourceGoBetweenServerRead,

		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"balance": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"bind": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"protocol": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "tcp",
			},

			"max_connections": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Optional: true,
				Default:  0,
			},

			"client_idle_timeout": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "1m",
			},

			"backend_idle_timeout": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "1m",
			},

			"backend_connection_timeout": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "1m",
			},

			"discovery": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "static",
						},

						"fail_policy": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "keeplast",
						},

						"interval": {
							Type:     schema.TypeInt,
							ForceNew: true,
							Optional: true,
							Default:  0,
						},

						"timeout": {
							Type:     schema.TypeInt,
							ForceNew: true,
							Optional: true,
							Default:  0,
						},

						"static_list": {
							Type:     schema.TypeList,
							ForceNew: true,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceGoBetweenServerCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(GBProvider).Client

	name := d.Get("name").(string)
	s := &gb.Server{}
	s.Balance = d.Get("balance").(string)
	s.Bind = d.Get("bind").(string)

	if v, ok := d.GetOk("max_connections"); ok {
		max := v.(int)
		s.MaxConnections = &max
	}

	if v, ok := d.GetOk("client_idle_timeout"); ok {
		cit := v.(string)
		s.ClientIdleTimeout = &cit
	}

	if v, ok := d.GetOk("backend_idle_timeout"); ok {
		bit := v.(string)
		s.BackendIdleTimeout = &bit
	}

	if v, ok := d.GetOk("backend_connection_timeout"); ok {
		bct := v.(string)
		s.BackendConnectionTimeout = &bct
	}

	// build static backend list
	staticList := make([]string, 0)
	if v, ok := d.GetOk("discovery.0.static_list"); ok {
		for _, static := range v.([]interface{}) {
			staticList = append(staticList, static.(string))
		}
	}

	s.Discovery = &gb.DiscoveryConfig{
		StaticDiscoveryConfig: &gb.StaticDiscoveryConfig{StaticList: staticList},
	}
	if v, ok := d.GetOk("discovery.0.kind"); ok {
		s.Discovery.Kind = v.(string)
	}
	if v, ok := d.GetOk("discovery.0.fail_policy"); ok {
		s.Discovery.Failpolicy = v.(string)
	}

	err := c.Api.AddServer(name, s)
	if err != nil {
		return err
	}

	d.SetId(name)

	return resourceGoBetweenServerRead(d, meta)
}

func resourceGoBetweenServerRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(GBProvider).Client

	name := d.Get("name").(string)

	s, err := c.Api.GetServer(name)
	if err != nil {
		return err
	}

	if s == nil {
		d.SetId("")
		return nil
	}

	d.Set("balance", s.Balance)
	d.Set("bind", s.Bind)
	d.Set("max_connections", s.MaxConnections)
	d.Set("client_idle_timeout", s.ClientIdleTimeout)
	d.Set("backend_idle_timeout", s.BackendIdleTimeout)
	d.Set("backend_connection_timeout", s.BackendConnectionTimeout)
	d.Set("discovery.0.kind", s.Discovery.Kind)
	d.Set("discovery.0.fail_policy", s.Discovery.Failpolicy)
	d.Set("discovery.0.static_list", s.Discovery.StaticList)

	return nil
}

func resourceGoBetweenServerDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(GBProvider).Client

	name := d.Get("name").(string)

	err := c.Api.DeleteServer(name)
	if err != nil {
		return err
	}

	return nil
}
