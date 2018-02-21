/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"log"
	"net/http"
)

func resourceNsxtLogicalRouterLinkPortOnTier1() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtLogicalRouterLinkPortOnTier1Create,
		Read:   resourceNsxtLogicalRouterLinkPortOnTier1Read,
		Update: resourceNsxtLogicalRouterLinkPortOnTier1Update,
		Delete: resourceNsxtLogicalRouterLinkPortOnTier1Delete,

		Schema: map[string]*schema.Schema{
			"revision": getRevisionSchema(),
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
			},
			"display_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The display name of this resource. Defaults to ID if not set",
				Optional:    true,
				Computed:    true,
			},
			"tag": getTagsSchema(),
			"logical_router_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Identifier for logical router on which this port is created",
				Required:    true,
				ForceNew:    true,
			},
			"linked_logical_router_port_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Identifier for port on logical router to connect to",
				Required:    true,
				ForceNew:    true,
			},
			"service_binding": getResourceReferencesSchema(false, false, []string{"LogicalService"}, "Service Bindings"),
		},
	}
}

func resourceNsxtLogicalRouterLinkPortOnTier1Create(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	logicalRouterID := d.Get("logical_router_id").(string)
	linkedLogicalRouterPortID := d.Get("linked_logical_router_port_id").(string)
	serviceBinding := getServiceBindingsFromSchema(d, "service_binding")
	logicalRouterLinkPort := manager.LogicalRouterLinkPortOnTier1{
		Description:               description,
		DisplayName:               displayName,
		Tags:                      tags,
		LogicalRouterId:           logicalRouterID,
		LinkedLogicalRouterPortId: makeResourceReference("LogicalPort", linkedLogicalRouterPortID),
		ServiceBindings:           serviceBinding,
	}

	logicalRouterLinkPort, resp, err := nsxClient.LogicalRoutingAndServicesApi.CreateLogicalRouterLinkPortOnTier1(nsxClient.Context, logicalRouterLinkPort)

	if err != nil {
		return fmt.Errorf("Error during LogicalRouterLinkPortOnTier1 create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during LogicalRouterLinkPortOnTier1 create: %v", resp.StatusCode)
	}
	d.SetId(logicalRouterLinkPort.Id)

	return resourceNsxtLogicalRouterLinkPortOnTier1Read(d, m)
}

func resourceNsxtLogicalRouterLinkPortOnTier1Read(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical router link port on tier1 id")
	}

	logicalRouterLinkPort, resp, err := nsxClient.LogicalRoutingAndServicesApi.ReadLogicalRouterLinkPortOnTier1(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] LogicalRouterLinkPortOnTier1 %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during LogicalRouterLinkPortOnTier1 read: %v", err)
	}

	d.Set("revision", logicalRouterLinkPort.Revision)
	d.Set("description", logicalRouterLinkPort.Description)
	d.Set("display_name", logicalRouterLinkPort.DisplayName)
	setTagsInSchema(d, logicalRouterLinkPort.Tags)
	d.Set("logical_router_id", logicalRouterLinkPort.LogicalRouterId)
	d.Set("linked_logical_router_port_id", logicalRouterLinkPort.LinkedLogicalRouterPortId)
	setServiceBindingsInSchema(d, logicalRouterLinkPort.ServiceBindings, "service_binding")

	return nil
}

func resourceNsxtLogicalRouterLinkPortOnTier1Update(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical router link port on tier1 id")
	}

	revision := int64(d.Get("revision").(int))
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	logicalRouterID := d.Get("logical_router_id").(string)
	linkedLogicalRouterPortID := d.Get("linked_logical_router_port_id").(string)
	serviceBinding := getServiceBindingsFromSchema(d, "service_binding")
	logicalRouterLinkPort := manager.LogicalRouterLinkPortOnTier1{
		Revision:                  revision,
		Description:               description,
		DisplayName:               displayName,
		Tags:                      tags,
		LogicalRouterId:           logicalRouterID,
		LinkedLogicalRouterPortId: makeResourceReference("LogicalPort", linkedLogicalRouterPortID),
		ServiceBindings:           serviceBinding,
		ResourceType:              "LogicalRouterLinkPortOnTIER1",
	}

	logicalRouterLinkPort, resp, err := nsxClient.LogicalRoutingAndServicesApi.UpdateLogicalRouterLinkPortOnTier1(nsxClient.Context, id, logicalRouterLinkPort)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during LogicalRouterLinkPortOnTier1 %v update: %v (%+v)", id, err, resp)
	}

	return resourceNsxtLogicalRouterLinkPortOnTier1Read(d, m)
}

func resourceNsxtLogicalRouterLinkPortOnTier1Delete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical router link port on tier1 id")
	}

	localVarOptionals := make(map[string]interface{})
	resp, err := nsxClient.LogicalRoutingAndServicesApi.DeleteLogicalRouterPort(nsxClient.Context, id, localVarOptionals)
	if err != nil {
		return fmt.Errorf("Error during LogicalRouterLinkPortOnTier1 delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] LogicalRouterLinkPortOnTier1 %s not found", id)
		d.SetId("")
	}

	return nil
}
