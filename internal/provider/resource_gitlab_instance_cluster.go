package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/xanzy/go-gitlab"
)

func resourceGitlabInstanceCluster() *schema.Resource {
	return &schema.Resource{
		Description: "This resource allows you to create and manage instance clusters for your GitLab instances.\n" +
			"For further information on clusters, consult the [gitlab\n" +
			"documentation](https://docs.gitlab.com/ee/user/instance/clusters/).",

		CreateContext: resourceGitlabInstanceClusterCreate,
		ReadContext:   resourceGitlabInstanceClusterRead,
		UpdateContext: resourceGitlabInstanceClusterUpdate,
		DeleteContext: resourceGitlabInstanceClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of cluster.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"domain": {
				Description: "The base domain of the cluster.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enabled": {
				Description: "Determines if cluster is active or not. Defaults to `true`. This attribute cannot be read.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    true,
			},
			"managed": {
				Description: "Determines if cluster is managed by gitlab or not. Defaults to `true`. This attribute cannot be read.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    true,
			},
			"created_at": {
				Description: "Create time.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"provider_type": {
				Description: "Provider type.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"platform_type": {
				Description: "Platform type.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"environment_scope": {
				Description: "The associated environment to the cluster. Defaults to `*`.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "*",
			},
			"cluster_type": {
				Description: "Cluster type.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"kubernetes_api_url": {
				Description: "The URL to access the Kubernetes API.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"kubernetes_token": {
				Description: "The token to authenticate against Kubernetes. This attribute cannot be read.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"kubernetes_ca_cert": {
				Description: "TLS certificate (needed if API is using a self-signed TLS certificate).",
				Type:        schema.TypeString,
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.TrimSpace(old) == strings.TrimSpace(new)
				},
			},
			"kubernetes_namespace": {
				Description: "The unique namespace related to the instance.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"kubernetes_authorization_type": {
				Description:  "The cluster authorization type. Valid values are `rbac`, `abac`, `unknown_authorization`. Defaults to `rbac`.",
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "rbac",
				ValidateFunc: validation.StringInSlice([]string{"rbac", "abac", "unknown_authorization"}, false),
			},
			"management_project_id": {
				Description: "The ID of the management project for the cluster.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceGitlabInstanceClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	pk := gitlab.AddPlatformKubernetesOptions{
		APIURL: gitlab.String(d.Get("kubernetes_api_url").(string)),
		Token:  gitlab.String(d.Get("kubernetes_token").(string)),
	}

	if v, ok := d.GetOk("kubernetes_ca_cert"); ok {
		pk.CaCert = gitlab.String(v.(string))
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

	if v, ok := d.GetOk("management_project_id"); ok {
		options.ManagementProjectID = gitlab.String(v.(string))
	}

	log.Printf("[DEBUG] create gitlab instance cluster %q", *options.Name)

	cluster, _, err := client.InstanceCluster.AddCluster(options, gitlab.WithContext(ctx))

	if err != nil {
		return diag.FromErr(err)
	}

	clusterIdString := fmt.Sprintf("%d", cluster.ID)
	d.SetId(clusterIdString)

	return resourceGitlabInstanceClusterRead(ctx, d, meta)
}

func resourceGitlabInstanceClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	clusterId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] read gitlab instance cluster %d", clusterId)

	cluster, _, err := client.InstanceCluster.GetCluster(clusterId, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] gitlab instance cluster not found %d", clusterId)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

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

	if cluster.ManagementProject == nil {
		d.Set("management_project_id", "")
	} else {
		d.Set("management_project_id", strconv.Itoa(cluster.ManagementProject.ID))
	}

	return nil
}

func resourceGitlabInstanceClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	clusterId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
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

	if d.HasChange("namespace") {
		pk.Namespace = gitlab.String(d.Get("namespace").(string))
	}

	if *pk != (gitlab.EditPlatformKubernetesOptions{}) {
		options.PlatformKubernetes = pk
	}

	if d.HasChange("management_project_id") {
		options.ManagementProjectID = gitlab.String(d.Get("management_project_id").(string))
	}

	if *options != (gitlab.EditClusterOptions{}) {
		log.Printf("[DEBUG] update gitlab instance cluster %d", clusterId)
		_, _, err := client.InstanceCluster.EditCluster(clusterId, options, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGitlabInstanceClusterRead(ctx, d, meta)
}

func resourceGitlabInstanceClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	clusterId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] delete gitlab instance cluster %d", clusterId)

	_, err = client.InstanceCluster.DeleteCluster(clusterId, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
