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

var _ = registerResource("gitlab_group_cluster", func() *schema.Resource {
	return &schema.Resource{
		Description: "This resource allows you to create and manage group clusters for your GitLab groups.\n" +
			"For further information on clusters, consult the [gitlab\n" +
			"documentation](https://docs.gitlab.com/ce/user/group/clusters/index.html).",

		CreateContext: resourceGitlabGroupClusterCreate,
		ReadContext:   resourceGitlabGroupClusterRead,
		UpdateContext: resourceGitlabGroupClusterUpdate,
		DeleteContext: resourceGitlabGroupClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"group": {
				Description: "The id of the group to add the cluster to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
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
				Description: "The token to authenticate against Kubernetes.",
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
})

func resourceGitlabGroupClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	cluster, _, err := client.GroupCluster.AddCluster(group, options, gitlab.WithContext(ctx))

	if err != nil {
		return diag.FromErr(err)
	}

	clusterIdString := fmt.Sprintf("%d", cluster.ID)
	d.SetId(buildTwoPartID(&group, &clusterIdString))

	return resourceGitlabGroupClusterRead(ctx, d, meta)
}

func resourceGitlabGroupClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	group, clusterId, err := groupIdAndClusterIdFromId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] read gitlab group cluster %q/%d", group, clusterId)

	cluster, _, err := client.GroupCluster.GetCluster(group, clusterId, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] gitlab group cluster not found %s/%d", group, clusterId)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
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

func resourceGitlabGroupClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	group, clusterId, err := groupIdAndClusterIdFromId(d.Id())
	if err != nil {
		return diag.FromErr(err)
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
		_, _, err := client.GroupCluster.EditCluster(group, clusterId, options, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGitlabGroupClusterRead(ctx, d, meta)
}

func resourceGitlabGroupClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	group, clusterId, err := groupIdAndClusterIdFromId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] delete gitlab group cluster %q/%d", group, clusterId)

	_, err = client.GroupCluster.DeleteCluster(group, clusterId, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
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
