package mongodb

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"terraform-provider-mongodb/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func validateEmail(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	value := v.(string)

	var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

	if !(emailRegex.MatchString(value)) {
		errs = append(errs, fmt.Errorf("Expected EmailId is not valid  %s", k))
		return warns, errs
	}
	return
}

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateEmail,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"invitationid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"userid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"inviterusername": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"orgname": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"country": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"email_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"first_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"roles": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"org_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"group_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"role_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	userRoles := d.Get("roles").([]interface{})
	ois := make([]string, 0)

	for _, role := range userRoles {
		i := role.(map[string]interface{})
		orgid := i["org_id"].(string)
		role := i["role_name"].(string)
		groupid := i["group_id"].(string)
		if groupid != "" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unless he accept the invitation,he cannot br given a project",
			})
			return diags
		}
		if orgid != "" {
			ois = append(ois, role)
		}

	}
	user := client.NewUser{
		Username: d.Get("username").(string),
		Roles:    ois,
	}
	newuser, err := apiClient.AddNewUser(&user)
	if err != nil {
		log.Println("[ERROR]: ", err)
		return diag.FromErr(err)
	}
	d.SetId(newuser.Username)
	d.Set("inviterusername", newuser.InviterUsername)
	d.Set("orgname", newuser.OrgName)
	d.Set("username", newuser.Username)
	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	if os.Getenv("ORGID") != "" {
		return diags
	}
	apiClient := m.(*client.Client)
	userId := d.Id()
	user, err := apiClient.GetUser(userId)
	if err != nil {
		log.Println("[ERROR]: ", err)
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return diag.FromErr(err)
		}
	}
	if len(user.ID) > 0 {
		fmt.Println("Set all values")
		d.SetId(user.Username)
		d.Set("userid", user.ID)
		d.Set("country", user.Country)
		d.Set("email_address", user.EmailAddress)
		d.Set("first_name", user.FirstName)
		d.Set("last_name", user.LastName)
		roles := make([]map[string]interface{}, 0)
		for _, v := range user.Roles {
			role := make(map[string]interface{})
			role["org_id"] = v.OrgID
			role["role_name"] = v.RoleName
			role["group_id"] = v.GroupID
			roles = append(roles, role)
		}
		d.Set("roles", roles)
		d.Set("username", user.Username)

	}
	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var _ diag.Diagnostics
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	if d.HasChange("username") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "User not allowed to change username",
			Detail:   "User not allowed to change username",
		})

		return diags
	}
	userId := d.Get("username").(string)
	userInfo, error := apiClient.GetUser(userId)
	if error != nil {
		log.Println("[ERROR]: ", error)
		if strings.Contains(error.Error(), "not found") {
			d.SetId("")
		} else {
			return diag.FromErr(error)
		}
	}
	userId = userInfo.ID
	roles := d.Get("roles").([]interface{})

	ois := []client.Role{}
	for _, role := range roles {
		i := role.(map[string]interface{})
		oi := client.Role{
			OrgID:    i["org_id"].(string),
			RoleName: i["role_name"].(string),
			GroupID:  i["group_id"].(string),
		}
		ois = append(ois, oi)
	}
	updatevalue := client.UpdateUser{
		Roles: ois,
	}
	_, err := apiClient.UpdateUser(&updatevalue, userId)
	if err != nil {
		log.Printf("[Error] Error updating user :%s", err)
		return diag.FromErr(err)
	}
	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")
	return diags
}