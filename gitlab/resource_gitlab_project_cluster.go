package gitlab

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/xanzy/go-gitlab"
)

func resourceGitlabProjectCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabProjectClusterCreate,
		Read:   resourceGitlabProjectClusterRead,
		Update: resourceGitlabProjectClusterUpdate,
		Delete: resourceGitlabProjectClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project": {
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
			"kubernetes_namespace": {
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
		},
	}
}

func resourceGitlabProjectClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	pk := gitlab.AddPlatformKubernetesOptions{
		APIURL: gitlab.String(d.Get("kubernetes_api_url").(string)),
		Token:  gitlab.String(d.Get("kubernetes_token").(string)),
	}

	if v, ok := d.GetOk("kubernetes_ca_cert"); ok {
		pk.CaCert = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("kubernetes_namespace"); ok {
		pk.Namespace = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("kubernetes_authorization_type"); ok {
		pk.AuthorizationType = gitlab.String(v.(string))
	}

	options := &gitlab.AddClusterOptions{
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

	log.Printf("[DEBUG] create gitlab project cluster %q/%q", project, *options.Name)

	cluster, _, err := client.ProjectCluster.AddCluster(project, options)

	if err != nil {
		return err
	}

	clusterIdString := fmt.Sprintf("%d", cluster.ID)
	d.SetId(buildTwoPartID(&project, &clusterIdString))

	return resourceGitlabProjectClusterRead(d, meta)
}

func resourceGitlabProjectClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	project, clusterId, err := projectIdAndClusterIdFromId(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] read gitlab project cluster %q/%d", project, clusterId)

	cluster, _, err := client.ProjectCluster.GetCluster(project, clusterId)
	if err != nil {
		return err
	}

	d.Set("project", project)
	d.Set("name", cluster.Name)
	d.Set("domain", cluster.Domain)
	d.Set("created_at", cluster.CreatedAt.String())
	d.Set("provider_type", cluster.ProviderType)
	d.Set("platform_type", cluster.PlatformType)
	d.Set("environment_scope", cluster.EnvironmentScope)
	d.Set("cluster_type", cluster.ClusterType)

	d.Set("kubernetes_api_url", cluster.PlatformKubernetes.APIURL)
	d.Set("kubernetes_ca_cert", cluster.PlatformKubernetes.CaCert)
	d.Set("kubernetes_namespace", cluster.PlatformKubernetes.Namespace)
	d.Set("kubernetes_authorization_type", cluster.PlatformKubernetes.AuthorizationType)

	return nil
}

func resourceGitlabProjectClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	project, clusterId, err := projectIdAndClusterIdFromId(d.Id())
	if err != nil {
		return err
	}

	options := &gitlab.EditClusterOptions{}

	if d.HasChange("name") {
		options.Name = gitlab.String(d.Get("name").(string))
	}

	if d.HasChange("domain") {
		options.Domain = gitlab.String(d.Get("domain").(string))
	}

	if d.HasChange("environment_scope") {
		options.EnvironmentScope = gitlab.String(d.Get("environment_scope").(string))
	}

	pk := &gitlab.EditPlatformKubernetesOptions{}

	if d.HasChange("kubernetes_api_url") {
		pk.APIURL = gitlab.String(d.Get("kubernetes_api_url").(string))
	}

	if d.HasChange("kubernetes_token") {
		pk.Token = gitlab.String(d.Get("kubernetes_token").(string))
	}

	if d.HasChange("kubernetes_ca_cert") {
		pk.CaCert = gitlab.String(d.Get("kubernetes_ca_cert").(string))
	}

	if d.HasChange("kubernetes_namespace") {
		pk.Namespace = gitlab.String(d.Get("kubernetes_namespace").(string))
	}

	if *pk != (gitlab.EditPlatformKubernetesOptions{}) {
		options.PlatformKubernetes = pk
	}

	if *options != (gitlab.EditClusterOptions{}) {
		log.Printf("[DEBUG] update gitlab project cluster %q/%d", project, clusterId)
		_, _, err := client.ProjectCluster.EditCluster(project, clusterId, options)
		if err != nil {
			return err
		}
	}

	return resourceGitlabProjectClusterRead(d, meta)
}

func resourceGitlabProjectClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project, clusterId, err := projectIdAndClusterIdFromId(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] delete gitlab project cluster %q/%d", project, clusterId)

	_, err = client.ProjectCluster.DeleteCluster(project, clusterId)

	return err
}

func projectIdAndClusterIdFromId(id string) (string, int, error) {
	project, clusterIdString, err := parseTwoPartID(id)
	if err != nil {
		return "", 0, err
	}

	clusterId, err := strconv.Atoi(clusterIdString)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get clusterId: %v", err)
	}

	return project, clusterId, nil
}
