package provider

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/redpwn/terraform-provider-rctf/internal/rctf"
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

var idReg = regexp.MustCompile(`[^\w-]+`)

func resourceChallengePut(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	r := m.(*rctf.Client)
	c := rctf.Challenge{
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
	c.Id = d.Id()
	if c.Id == "" {
		c.Id = idReg.ReplaceAllString(fmt.Sprintf("%s-%s", c.Category, c.Name), "-")
		d.SetId(c.Id)
		d.Set("id", c.Id)
	}
	if err := r.PutChallenge(ctx, c); err != nil {
		return diag.Errorf("put challenge: %s", err)
	}
	return diags
}

func resourceChallengeRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	r := m.(*rctf.Client)
	id := d.Id()
	c, err := r.Challenge(ctx, id)
	if err != nil {
		rctfErr := &rctf.Error{}
		if errors.As(err, rctfErr) && rctfErr.Kind == "badChallenge" {
			d.SetId("")
			return
		}
		return diag.Errorf("get challenge: %s", err)
	}
	d.Set("id", c.Id)
	d.Set("name", c.Name)
	d.Set("description", c.Description)
	d.Set("category", c.Category)
	d.Set("author", c.Author)
	var f []map[string]interface{}
	for _, file := range c.Files {
		f = append(f, map[string]interface{}{
			"name": file.Name,
			"url":  file.Url,
		})
	}
	d.Set("file", f)
	d.Set("min_points", c.Points.Min)
	d.Set("max_points", c.Points.Max)
	d.Set("flag", c.Flag)
	d.Set("tiebreak_eligible", c.TiebreakEligible)
	d.Set("sort_weight", c.SortWeight)
	return
}

func resourceChallengeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	r := m.(*rctf.Client)
	id := d.Id()
	if err := r.DeleteChallenge(ctx, id); err != nil {
		return diag.Errorf("delete challenge: %s", err)
	}
	return
}
