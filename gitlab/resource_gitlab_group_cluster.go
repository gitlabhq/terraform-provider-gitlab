package gitlab

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Fourcast/go-gitlab"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceGitlabGroupCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabGroupClusterCreate,
		Read:   resourceGitlabGroupClusterRead,
		Update: resourceGitlabGroupClusterUpdate,
		Delete: resourceGitlabGroupClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
			"managed": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"platform_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"environment_scope": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "*",
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kubernetes_api_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"kubernetes_token": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"kubernetes_ca_cert": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"kubernetes_authorization_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "rbac",
				ValidateFunc: validation.StringInSlice([]string{"rbac", "abac", "unknown_authorization"}, false),
			},
			"management_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceGitlabGroupClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	group := d.Get("group").(string)

	pk := gitlab.AddGroupPlatformKubernetesOptions{
		APIURL: gitlab.String(d.Get("kubernetes_api_url").(string)),
		Token:  gitlab.String(d.Get("kubernetes_token").(string)),
	}

	if v, ok := d.GetOk("kubernetes_ca_cert"); ok {
		pk.CaCert = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("kubernetes_authorization_type"); ok {
		pk.AuthorizationType = gitlab.String(v.(string))
	}

	options := &gitlab.AddGroupClusterOptions{
		Name:               gitlab.String(d.Get("name").(string)),
		Enabled:            gitlab.Bool(d.Get("enabled").(bool)),
		Managed:            gitlab.Bool(d.Get("managed").(bool)),
		PlatformKubernetes: &pk,
	}

	if v, ok := d.GetOk("domain"); ok {
		options.Domain = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("environment_scope"); ok {
		options.EnvironmentScope = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("management_project_id"); ok {
		options.ManagementProjectID = gitlab.String(v.(string))
	}

	log.Printf("[DEBUG] create gitlab group cluster %q/%q", group, *options.Name)

	cluster, _, err := client.GroupCluster.AddCluster(group, options)

	if err != nil {
		return err
	}

	clusterIdString := fmt.Sprintf("%d", cluster.ID)
	d.SetId(buildTwoPartID(&group, &clusterIdString))

	return resourceGitlabGroupClusterRead(d, meta)
}

func resourceGitlabGroupClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	group, clusterId, err := groupIdAndClusterIdFromId(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] read gitlab group cluster %q/%d", group, clusterId)

	cluster, _, err := client.GroupCluster.GetCluster(group, clusterId)
	if err != nil {
		return err
	}

	d.Set("group", group)
	d.Set("name", cluster.Name)
	d.Set("domain", cluster.Domain)
	d.Set("created_at", cluster.CreatedAt.String())
	d.Set("provider_type", cluster.ProviderType)
	d.Set("platform_type", cluster.PlatformType)
	d.Set("environment_scope", cluster.EnvironmentScope)
	d.Set("cluster_type", cluster.ClusterType)

	d.Set("kubernetes_api_url", cluster.PlatformKubernetes.APIURL)
	d.Set("kubernetes_ca_cert", cluster.PlatformKubernetes.CaCert)
	d.Set("kubernetes_authorization_type", cluster.PlatformKubernetes.AuthorizationType)

	if cluster.ManagementProject == nil {
		d.Set("management_project_id", "")
	} else {
		d.Set("management_project_id", strconv.Itoa(cluster.ManagementProject.ID))
	}

	return nil
}

func resourceGitlabGroupClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	group, clusterId, err := groupIdAndClusterIdFromId(d.Id())
	if err != nil {
		return err
	}

	options := &gitlab.EditGroupClusterOptions{}

	if d.HasChange("name") {
		options.Name = gitlab.String(d.Get("name").(string))
	}

	if d.HasChange("domain") {
		options.Domain = gitlab.String(d.Get("domain").(string))
	}

	if d.HasChange("environment_scope") {
		options.EnvironmentScope = gitlab.String(d.Get("environment_scope").(string))
	}

	pk := &gitlab.EditGroupPlatformKubernetesOptions{}

	if d.HasChange("kubernetes_api_url") {
		pk.APIURL = gitlab.String(d.Get("kubernetes_api_url").(string))
	}

	if d.HasChange("kubernetes_token") {
		pk.Token = gitlab.String(d.Get("kubernetes_token").(string))
	}

	if d.HasChange("kubernetes_ca_cert") {
		pk.CaCert = gitlab.String(d.Get("kubernetes_ca_cert").(string))
	}

	if *pk != (gitlab.EditGroupPlatformKubernetesOptions{}) {
		options.PlatformKubernetes = pk
	}

	if d.HasChange("management_project_id") {
		options.ManagementProjectID = gitlab.String(d.Get("management_project_id").(string))
	}

	if *options != (gitlab.EditGroupClusterOptions{}) {
		log.Printf("[DEBUG] update gitlab group cluster %q/%d", group, clusterId)
		_, _, err := client.GroupCluster.EditCluster(group, clusterId, options)
		if err != nil {
			return err
		}
	}

	return resourceGitlabGroupClusterRead(d, meta)
}

func resourceGitlabGroupClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	group, clusterId, err := groupIdAndClusterIdFromId(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] delete gitlab group cluster %q/%d", group, clusterId)

	_, err = client.GroupCluster.DeleteCluster(group, clusterId)

	return err
}

func groupIdAndClusterIdFromId(id string) (string, int, error) {
	group, clusterIdString, err := parseTwoPartID(id)
	if err != nil {
		return "", 0, err
	}

	clusterId, err := strconv.Atoi(clusterIdString)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get clusterId: %v", err)
	}

	return group, clusterId, nil
}
