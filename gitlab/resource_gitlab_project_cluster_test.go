package gitlab

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabProjectCluster_basic(t *testing.T) {
	var cluster gitlab.ProjectCluster
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabProjectClusterDestroy,
		Steps: []resource.TestStep{
			// Create a project and cluster with default options
			{
				Config: testAccGitlabProjectClusterConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectClusterExists("gitlab_project_cluster.foo", &cluster),
					testAccCheckGitlabProjectClusterAttributes(&cluster, &testAccGitlabProjectClusterExpectedAttributes{
						Name:                        fmt.Sprintf("foo-cluster-%d", rInt),
						EnvironmentScope:            "*",
						KubernetesApiURL:            "https://123.123.123",
						KubernetesAuthorizationType: "abac",
					}),
				),
			},
			// Update cluster
			{
				Config: testAccGitlabProjectClusterUpdateConfig(rInt, "abac"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectClusterExists("gitlab_project_cluster.foo", &cluster),
					testAccCheckGitlabProjectClusterAttributes(&cluster, &testAccGitlabProjectClusterExpectedAttributes{
						Name:                        fmt.Sprintf("foo-cluster-%d", rInt),
						EnvironmentScope:            "*",
						KubernetesApiURL:            "https://124.124.124",
						KubernetesCACert:            "some-cert",
						KubernetesNamespace:         "changed-namespace",
						KubernetesAuthorizationType: "abac",
					}),
				),
			},
			// Update authorization type cluster
			{
				Config: testAccGitlabProjectClusterUpdateConfig(rInt, "rbac"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectClusterExists("gitlab_project_cluster.foo", &cluster),
					testAccCheckGitlabProjectClusterAttributes(&cluster, &testAccGitlabProjectClusterExpectedAttributes{
						Name:                        fmt.Sprintf("foo-cluster-%d", rInt),
						EnvironmentScope:            "*",
						KubernetesApiURL:            "https://124.124.124",
						KubernetesCACert:            "some-cert",
						KubernetesNamespace:         "changed-namespace",
						KubernetesAuthorizationType: "rbac",
					}),
				),
			},
		},
	})
}

func TestAccGitlabProjectCluster_import(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabProjectClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabProjectClusterConfig(rInt),
			},
			{
				ResourceName:            "gitlab_project_cluster.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enabled", "kubernetes_token"},
			},
		},
	})
}

type testAccGitlabProjectClusterExpectedAttributes struct {
	Name                        string
	EnvironmentScope            string
	KubernetesApiURL            string
	KubernetesCACert            string
	KubernetesNamespace         string
	KubernetesAuthorizationType string
}

func testAccCheckGitlabProjectClusterExists(n string, cluster *gitlab.ProjectCluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %q", n)
		}

		project, clusterID, err := projectIdAndClusterIdFromId(rs.Primary.ID)
		if err != nil {
			return err
		}

		conn := testAccProvider.Meta().(*gitlab.Client)

		gotCluster, _, err := conn.ProjectCluster.GetCluster(project, clusterID)
		if err != nil {
			return err
		}

		*cluster = *gotCluster

		return nil
	}
}

func testAccCheckGitlabProjectClusterDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_project_cluster" {
			continue
		}

		project, clusterID, err := projectIdAndClusterIdFromId(rs.Primary.ID)
		if err != nil {
			return err
		}

		gotCluster, resp, err := conn.ProjectCluster.GetCluster(project, clusterID)
		if err == nil {
			if gotCluster != nil && fmt.Sprintf("%d", gotCluster.ID) == project {
				return fmt.Errorf("project cluster still exists")
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
	}

	return nil
}

func testAccCheckGitlabProjectClusterAttributes(cluster *gitlab.ProjectCluster, want *testAccGitlabProjectClusterExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cluster.Name != want.Name {
			return fmt.Errorf("got name %q; want %q", cluster.Name, want.Name)
		}

		if cluster.EnvironmentScope != want.EnvironmentScope {
			return fmt.Errorf("got environment scope %q; want %q", cluster.EnvironmentScope, want.EnvironmentScope)
		}

		if cluster.PlatformKubernetes.APIURL != want.KubernetesApiURL {
			return fmt.Errorf("got kubernetes api url %q; want %q", cluster.PlatformKubernetes.APIURL, want.KubernetesApiURL)
		}

		if cluster.PlatformKubernetes.CaCert != want.KubernetesCACert {
			return fmt.Errorf("got kubernetes ca cert %q; want %q", cluster.PlatformKubernetes.CaCert, want.KubernetesCACert)
		}

		if cluster.PlatformKubernetes.Namespace != want.KubernetesNamespace {
			return fmt.Errorf("got kubernetes namespace %q; want %q", cluster.PlatformKubernetes.Namespace, want.KubernetesNamespace)
		}

		if cluster.PlatformKubernetes.AuthorizationType != want.KubernetesAuthorizationType {
			return fmt.Errorf("got kubernetes authorization type %q; want %q", cluster.PlatformKubernetes.AuthorizationType, want.KubernetesAuthorizationType)
		}

		return nil
	}
}

func testAccGitlabProjectClusterConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-project-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource gitlab_project_cluster "foo" {
  project                       = "${gitlab_project.foo.id}"
  name                          = "foo-cluster-%d"
  kubernetes_api_url            = "https://123.123.123"
  kubernetes_token              = "some-token"
  kubernetes_authorization_type = "abac"
}
`, rInt, rInt)
}

func testAccGitlabProjectClusterUpdateConfig(rInt int, authType string) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-project-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource gitlab_project_cluster "foo" {
  project                       = "${gitlab_project.foo.id}"
  name                          = "foo-cluster-%d"
  kubernetes_api_url            = "https://124.124.124"
  kubernetes_token              = "some-token"
  kubernetes_ca_cert            = "some-cert"
  kubernetes_namespace          = "changed-namespace"
  kubernetes_authorization_type = "%s"
}
`, rInt, rInt, authType)
}
