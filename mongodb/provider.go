package mongodb

import (
	"terraform-provider-mongodb/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"publickey": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"privatekey": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"orgid": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"mongodb_user": resourceUser(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"mongodb_user": dataSourceUser(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	publickey := d.Get("publickey").(string)
	privatekey := d.Get("privatekey").(string)
	orgid := d.Get("orgid").(string)
	return client.NewClient(publickey, privatekey, orgid), nil

}
