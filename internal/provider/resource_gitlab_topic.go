package provider

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_topic", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_topic`" + ` resource allows to manage the lifecycle of topics that are then assignable to projects.

-> Topics are the successors for project tags. Aside from avoiding terminology collisions with Git tags, they are more descriptive and better searchable.

~> Deleting a topic was implemented in GitLab 14.9. For older versions of GitLab set ` + "`soft_destroy = true`" + ` to empty out a topic instead of deleting it.

**Upstream API**: [GitLab REST API docs for topics](https://docs.gitlab.com/ee/api/topics.html)
`,

		CreateContext: resourceGitlabTopicCreate,
		ReadContext:   resourceGitlabTopicRead,
		UpdateContext: resourceGitlabTopicUpdate,
		DeleteContext: resourceGitlabTopicDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The topic's name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"soft_destroy": {
				Description: "Empty the topics fields instead of deleting it.",
				Type:        schema.TypeBool,
				Optional:    true,
				Deprecated:  "GitLab 14.9 introduced the proper deletion of topics. This field is no longer needed.",
			},
			"description": {
				Description: "A text describing the topic.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"avatar": {
				Description: "A local path to the avatar image to upload. **Note**: not available for imported resources",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"_avatar_hash": {
				Description: "The hash of the avatar image. **Note**: this is an internally used attribute to track the avatar image.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"avatar_url": {
				Description: "The URL of the avatar image.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
		CustomizeDiff: func(ctx context.Context, rd *schema.ResourceDiff, i interface{}) error {
			if v, ok := rd.GetOk("avatar"); ok {
				avatarPath := v.(string)
				avatarFile, err := os.Open(avatarPath)
				if err != nil {
					return fmt.Errorf("Unable to open avatar file %s: %s", avatarPath, err)
				}

				avatarHash, err := getSha256(avatarFile)
				if err != nil {
					return fmt.Errorf("Unable to get hash from avatar file %s: %s", avatarPath, err)
				}

				if rd.Get("_avatar_hash").(string) != avatarHash {
					if err := rd.SetNew("_avatar_hash", avatarHash); err != nil {
						return fmt.Errorf("Unable to set _avatar_hash: %s", err)
					}
				}
			}
			return nil
		},
	}
})

func resourceGitlabTopicCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	options := &gitlab.CreateTopicOptions{
		Name: gitlab.String(d.Get("name").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		options.Description = gitlab.String(v.(string))
	}

	avatarHash := ""
	if v, ok := d.GetOk("avatar"); ok {
		avatar, hash, err := resourceGitlabTopicGetAvatar(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		options.Avatar = avatar
		avatarHash = hash
	}

	log.Printf("[DEBUG] create gitlab topic %s", *options.Name)

	topic, _, err := client.Topics.CreateTopic(options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.Errorf("Failed to create topic %q: %s", *options.Name, err)
	}

	d.SetId(fmt.Sprintf("%d", topic.ID))
	if options.Avatar != nil {
		d.Set("_avatar_hash", avatarHash)
	}
	return resourceGitlabTopicRead(ctx, d, meta)
}

func resourceGitlabTopicRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	topicID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Failed to convert topic id %s to int: %s", d.Id(), err)
	}
	log.Printf("[DEBUG] read gitlab topic %d", topicID)

	topic, _, err := client.Topics.GetTopic(topicID, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] gitlab group %s not found so removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Failed to read topic %d: %s", topicID, err)
	}

	d.SetId(fmt.Sprintf("%d", topic.ID))
	d.Set("name", topic.Name)
	d.Set("description", topic.Description)
	d.Set("avatar_url", topic.AvatarURL)

	return nil
}

func resourceGitlabTopicUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	options := &gitlab.UpdateTopicOptions{}

	if d.HasChange("name") {
		options.Name = gitlab.String(d.Get("name").(string))
	}

	if d.HasChange("description") {
		options.Description = gitlab.String(d.Get("description").(string))
	}

	if d.HasChanges("avatar", "_avatar_hash") {
		avatarPath := d.Get("avatar").(string)
		var avatar *gitlab.TopicAvatar
		var avatarHash string
		// NOTE: the avatar should be removed
		if avatarPath == "" {
			avatar = &gitlab.TopicAvatar{}
			avatarHash = ""
		} else {
			changedAvatar, changedAvatarHash, err := resourceGitlabTopicGetAvatar(avatarPath)
			if err != nil {
				return diag.FromErr(err)
			}
			avatar = changedAvatar
			avatarHash = changedAvatarHash
		}
		options.Avatar = avatar
		d.Set("_avatar_hash", avatarHash)
	}

	log.Printf("[DEBUG] update gitlab topic %s", d.Id())

	topicID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Failed to convert topic id %s to int: %s", d.Id(), err)
	}
	_, _, err = client.Topics.UpdateTopic(topicID, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.Errorf("Failed to update topic %d: %s", topicID, err)
	}
	return resourceGitlabTopicRead(ctx, d, meta)
}

func resourceGitlabTopicDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	topicID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Failed to convert topic id %s to int: %s", d.Id(), err)
	}
	softDestroy := d.Get("soft_destroy").(bool)

	deleteNotSupported, err := isGitLabVersionLessThan(client, "14.9")()
	if err != nil {
		return diag.FromErr(err)
	}
	if !softDestroy && deleteNotSupported {
		return diag.Errorf("GitLab 14.9 introduced the proper deletion of topics. Set `soft_destroy = true` to empty out a topic instead of deleting it.")
	}

	// NOTE: the `soft_destroy` field is deprecated and will be removed in a future version.
	//       It was only introduced because GitLab prior to 14.9 didn't support topic deletion.
	if softDestroy {
		log.Printf("[WARN] Not deleting gitlab topic %s. Instead emptying its description", d.Id())

		options := &gitlab.UpdateTopicOptions{
			Description: gitlab.String(""),
		}

		_, _, err = client.Topics.UpdateTopic(topicID, options, gitlab.WithContext(ctx))
		if err != nil {
			return diag.Errorf("Failed to update topic %d: %s", topicID, err)
		}

		return nil
	}

	log.Printf("[DEBUG] delete gitlab topic %s", d.Id())

	if _, err = client.Topics.DeleteTopic(topicID, gitlab.WithContext(ctx)); err != nil {
		return diag.Errorf("Failed to delete topic %d: %s", topicID, err)
	}

	return nil
}

func resourceGitlabTopicGetAvatar(avatarPath string) (*gitlab.TopicAvatar, string, error) {
	avatarFile, err := os.Open(avatarPath)
	if err != nil {
		return nil, "", fmt.Errorf("Unable to open avatar file %s: %s", avatarPath, err)
	}

	avatarImageBuf := &bytes.Buffer{}
	teeReader := io.TeeReader(avatarFile, avatarImageBuf)
	avatarHash, err := getSha256(teeReader)
	if err != nil {
		return nil, "", fmt.Errorf("Unable to get hash from avatar file %s: %s", avatarPath, err)
	}

	avatarImageReader := bytes.NewReader(avatarImageBuf.Bytes())
	return &gitlab.TopicAvatar{
		Filename: avatarPath,
		Image:    avatarImageReader,
	}, avatarHash, nil
}

func getSha256(r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", fmt.Errorf("Unable to get hash %s", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
