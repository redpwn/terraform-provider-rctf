package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/redpwn/terraform-provider-rctf/internal/rctf"
	"github.com/segmentio/ksuid"
)

func resourceChallenge() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChallengePut,
		ReadContext:   resourceChallengeRead,
		UpdateContext: resourceChallengePut,
		DeleteContext: resourceChallengeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"category": {
				Type:     schema.TypeString,
				Required: true,
			},
			"author": {
				Type:     schema.TypeString,
				Required: true,
			},
			"file": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"url": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"max_points": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"min_points": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"flag": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tiebreak_eligible": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"sort_weight": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
		},
	}
}

func resourceChallengePut(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	r := m.(*rctf.Client)
	var diags diag.Diagnostics
	id := d.Id()
	if id == "" {
		id = ksuid.New().String()
		d.SetId(id)
	}
	c := rctf.Challenge{
		Id:          id,
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Category:    d.Get("category").(string),
		Author:      d.Get("author").(string),
		Files:       []rctf.ChallengeFile{},
		Points: rctf.ChallengePoints{
			Min: d.Get("min_points").(int),
			Max: d.Get("max_points").(int),
		},
		Flag:             d.Get("flag").(string),
		TiebreakEligible: d.Get("tiebreak_eligible").(bool),
		SortWeight:       d.Get("sort_weight").(int),
	}
	for _, v := range d.Get("file").([]interface{}) {
		f := v.(map[string]interface{})
		c.Files = append(c.Files, rctf.ChallengeFile{
			Name: f["name"].(string),
			Url:  f["url"].(string),
		})
	}
	if err := r.PutChallenge(ctx, c); err != nil {
		return diag.Errorf("put challenge: %s", err)
	}
	diags = append(diags, resourceChallengeRead(ctx, d, m)...)
	return diags
}

func flattenFiles(files []rctf.ChallengeFile) []map[string]interface{} {
	var f []map[string]interface{}
	for _, file := range files {
		f = append(f, map[string]interface{}{
			"name": file.Name,
			"url":  file.Url,
		})
	}
	return f
}

func resourceChallengeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	r := m.(*rctf.Client)
	var diags diag.Diagnostics
	id := d.Id()
	c, err := r.Challenge(ctx, id)
	if err != nil {
		return diag.Errorf("get challenge: %s", err)
	}
	if err := d.Set("id", c.Id); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", c.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", c.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("category", c.Category); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("author", c.Author); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("file", flattenFiles(c.Files)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("min_points", c.Points.Min); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("max_points", c.Points.Max); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("flag", c.Flag); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tiebreak_eligible", c.TiebreakEligible); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("sort_weight", c.SortWeight); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceChallengeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	r := m.(*rctf.Client)
	var diags diag.Diagnostics
	id := d.Id()
	if err := r.DeleteChallenge(ctx, id); err != nil {
		return diag.Errorf("delete challenge: %s", err)
	}
	return diags
}
