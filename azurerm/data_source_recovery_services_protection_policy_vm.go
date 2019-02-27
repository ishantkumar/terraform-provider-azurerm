package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmRecoveryServicesProtectionPolicyVm() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmRecoveryServicesProtectionPolicyVmRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"recovery_vault_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceArmRecoveryServicesProtectionPolicyVmRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).recoveryServicesProtectionPoliciesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)

	log.Printf("[DEBUG] Reading Recovery Service Protection Policy %q (resource group %q)", name, resourceGroup)

	protectionPolicy, err := client.Get(ctx, vaultName, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(protectionPolicy.Response) {
			return fmt.Errorf("Error: Recovery Services Protection Policy %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("Error making Read request on Recovery Service Protection Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*protectionPolicy.ID)
	d.Set("name", protectionPolicy.Name)
	d.Set("resource_group_name", resourceGroup)

	flattenAndSetTags(d, protectionPolicy.Tags)
	return nil
}
