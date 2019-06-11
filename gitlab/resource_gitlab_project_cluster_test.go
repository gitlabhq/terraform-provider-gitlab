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
				Config: testAccGitlabProjectClusterConfig(rInt, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectClusterExists("gitlab_project_cluster.foo", &cluster),
					testAccCheckGitlabProjectClusterAttributes(&cluster, &testAccGitlabProjectClusterExpectedAttributes{
						Name:                        fmt.Sprintf("foo-cluster-%d", rInt),
						Domain:                      "example.com",
						EnvironmentScope:            "*",
						KubernetesApiURL:            "https://123.123.123",
						KubernetesCACert:            projectClusterFakeCert,
						KubernetesAuthorizationType: "abac",
					}),
				),
			},
			// create an unmanaged cluster
			{
				Config: testAccGitlabProjectClusterConfig(rInt, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectClusterExists("gitlab_project_cluster.foo", &cluster),
					testAccCheckGitlabProjectClusterAttributes(&cluster, &testAccGitlabProjectClusterExpectedAttributes{
						Name:                        fmt.Sprintf("foo-cluster-%d", rInt),
						Domain:                      "example.com",
						EnvironmentScope:            "*",
						KubernetesApiURL:            "https://123.123.123",
						KubernetesCACert:            projectClusterFakeCert,
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
						Domain:                      "example-new.com",
						EnvironmentScope:            "*",
						KubernetesApiURL:            "https://124.124.124",
						KubernetesCACert:            projectClusterFakeCert,
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
						Domain:                      "example-new.com",
						EnvironmentScope:            "*",
						KubernetesApiURL:            "https://124.124.124",
						KubernetesCACert:            projectClusterFakeCert,
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
				Config: testAccGitlabProjectClusterConfig(rInt, true),
			},
			{
				ResourceName:            "gitlab_project_cluster.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enabled", "kubernetes_token", "managed"},
			},
		},
	})
}

type testAccGitlabProjectClusterExpectedAttributes struct {
	Name                        string
	Domain                      string
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

		if cluster.Domain != want.Domain {
			return fmt.Errorf("got domain %q; want %q", cluster.Domain, want.Domain)
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

func testAccGitlabProjectClusterConfig(rInt int, managed bool) string {
	m := "false"
	if managed {
		m = "true"
	}

	return fmt.Sprintf(`
variable "cert" {
  default = <<EOF
%s
EOF
}

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
  domain                        = "example.com"
  managed                       = "%s"
  kubernetes_api_url            = "https://123.123.123"
  kubernetes_token              = "some-token"
  kubernetes_ca_cert            = "${trimspace(var.cert)}"
  kubernetes_authorization_type = "abac"
}
`, projectClusterFakeCert, rInt, rInt, m)
}

func testAccGitlabProjectClusterUpdateConfig(rInt int, authType string) string {
	return fmt.Sprintf(`
variable "cert" {
  default = <<EOF
%s
EOF
}

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
  domain                        = "example-new.com"
  kubernetes_api_url            = "https://124.124.124"
  kubernetes_token              = "some-token"
  kubernetes_ca_cert            = "${trimspace(var.cert)}"
  kubernetes_namespace          = "changed-namespace"
  kubernetes_authorization_type = "%s"
}
`, projectClusterFakeCert, rInt, rInt, authType)
}

var projectClusterFakeCert = `-----BEGIN CERTIFICATE-----
MIICljCCAX4CCQDV7q2baHBlJjANBgkqhkiG9w0BAQsFADANMQswCQYDVQQGEwJV
SzAeFw0xOTAzMjUxMTUxNTZaFw0xOTA0MjQxMTUxNTZaMA0xCzAJBgNVBAYTAlVL
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA+bJAnlcGVMXjGdGcPFYf
aAyAlJLdef22vmjFQUgUw8HvblpUrkHYjVxdvVvg5tNFGLvnCIBRITJ6CQfl3f2i
ZL+SJCNZEWILt5TQTRQG09uab6An+ztm/XLJyHHUp0cEeI+aYifTuykB+cAxOLoA
+tWPq6i07Er+f/UcpntMxNi3b3LVpvdB5tcRvN6F2aXblLR3O7gvrmI4XA1u0Wba
LwDRgbS5bLy/AkrrZm7XMuy4zlbWEAMI3tgRInS+ENYKPuY0Hl3cxJBsI5EQ+mFx
aran+jYp7QQO/7VnrzRH7ZblWhcVKW9QoDFl78ZpT1sAwRif2ZFqMhw0sGiptLFb
AwIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQCgVxkvDwQPtpmx0WNriKsr5WeMvb6r
5DRhLOyyA7HncayAFCvAhk5M+x2wxWuuKzOPKsjSJpZDU0+2alVhZzzWbxSKoX7y
oW8+2ioodyfrW5vCPLEMfyqg2VGh+0F8PadVL96GZL20WYxCJ3eCuM7NFXG2ZciB
GJ48/0Tdc593QHg+19Jitq0xEL6V1dq5C5qhQxrikG3e3a+YYEZNCGwMj+2MhY2J
Up9FUfSTZR1MzQFi/7Dr2zffyuzFZk7IXrvA0foBe0GKPtiWQJ0/JHqkZfbfAEYw
c3fx6O/MhKijdlkbcpOanqD7PQEfUymTFLp2fZu2a1GRKIbfabGyyxGy
-----END CERTIFICATE-----`
