package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

type CustomAttributeGetter func(int, string, ...gitlab.RequestOptionFunc) (*gitlab.CustomAttribute, *gitlab.Response, error)
type CustomAttributeSetter func(int, gitlab.CustomAttribute, ...gitlab.RequestOptionFunc) (*gitlab.CustomAttribute, *gitlab.Response, error)
type CustomAttributeDeleter func(int, string, ...gitlab.RequestOptionFunc) (*gitlab.Response, error)

type CreateGetter func(*gitlab.Client) CustomAttributeGetter
type CreateSetter func(*gitlab.Client) CustomAttributeSetter
type CreateDeleter func(*gitlab.Client) CustomAttributeDeleter

func CreateCustomAttributeResource(idName string, createGetter CreateGetter, createSetter CreateSetter, createDeleter CreateDeleter, description string) *schema.Resource {
	setToState := func(d *schema.ResourceData, userId int, customAttribute *gitlab.CustomAttribute) {
		// lintignore:R001
		d.Set(idName, userId)
		d.Set("key", customAttribute.Key)
		d.Set("value", customAttribute.Value)
	}

	readFunc := func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		client := meta.(*gitlab.Client)
		getter := createGetter(client)
		log.Printf("[DEBUG] read Custom Attribute %s", d.Id())

		id, key, err := parseId(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		customAttribute, _, err := getter(id, key)
		if err != nil {
			return diag.FromErr(err)
		}

		setToState(d, id, customAttribute)
		return nil
	}

	setFunc := func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		client := meta.(*gitlab.Client)
		setter := createSetter(client)

		id := d.Get(idName).(int)
		options := &gitlab.CustomAttribute{
			Key:   d.Get("key").(string),
			Value: d.Get("value").(string),
		}

		log.Printf("[DEBUG] set (create or update) Custom Attribute %s with value %s for %s %d", options.Key, options.Value, idName, id)

		customAttribute, _, err := setter(id, *options, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(buildId(id, customAttribute.Key))
		return readFunc(ctx, d, meta)
	}

	deleteFunc := func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		client := meta.(*gitlab.Client)
		deleter := createDeleter(client)
		log.Printf("[DEBUG] delete Custom Attribute %s", d.Id())

		id, key, err := parseId(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = deleter(id, key, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}

		return nil
	}

	return &schema.Resource{
		Description:   description,
		CreateContext: setFunc,
		ReadContext:   readFunc,
		UpdateContext: setFunc,
		DeleteContext: deleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			idName: {
				Description: fmt.Sprintf("The id of the %s.", idName),
				Type:        schema.TypeInt,
				Required:    true,
			},
			"key": {
				Description: "Key for the Custom Attribute.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"value": {
				Description: "Value for the Custom Attribute.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func parseId(id string) (int, string, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return -1, "", fmt.Errorf("unexpected ID format (%q). Expected id:key", id)
	}

	subjectId, err := strconv.Atoi(parts[0])
	if err != nil {
		return -1, "", fmt.Errorf("unexpected ID format (%q). Expected id:key whereas `id` must be an integer", id)
	}

	return subjectId, parts[1], nil
}

func buildId(id int, key string) string {
	return fmt.Sprintf("%d:%s", id, key)
}
