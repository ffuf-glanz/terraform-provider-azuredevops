package azuredevops

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/tfhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/validate"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/release"
)

func resourceReleaseDefinition() *schema.Resource {
	variableGroups := &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type:         schema.TypeInt,
			ValidateFunc: validation.IntAtLeast(1),
		},
		Optional: true,
	}

	configurationVariableValue := &schema.Schema{
		Type: schema.TypeSet,
		// TODO : Abstract this because it is used by variable group and release definitions.
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
				},
				"allow_override": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"is_secret": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
		Optional: true,
		Set: func(i interface{}) int {
			item := i.(map[string]interface{})
			return schema.HashString(item["name"].(string))
		},
	}

	configurationVariableMap := &schema.Schema{
		Type: schema.TypeMap,
		Elem: configurationVariableValue,
	}

	taskInputValidation := map[string]*schema.Schema{
		"expression": {
			Type:     schema.TypeString,
			Required: true,
			Default:  "",
		},
		"message": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
	}

	workflowTasks := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: taskInputValidation,
		},
	}

	releaseDefinitionDeployStep := &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"tasks": workflowTasks,
			},
		},
	}

	rank := &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
		Default:  1,
	}

	approvalOptions := &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		MinItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"auto_triggered_and_previous_environment_approved_can_be_skipped": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"enforce_identity_revalidation": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"execution_order": {
					Type:         schema.TypeString,
					Required:     true,
					Default:      "beforeGates",
					ValidateFunc: validation.StringInSlice([]string{"afterGatesAlways", "afterSuccessfulGates", "beforeGates"}, false),
				},
				"release_creator_can_be_approver": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"required_approver_count": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"timeout_in_minutes": {
					Type:     schema.TypeInt,
					Optional: true,
				},
			},
		},
	}

	releaseDefinitionApprovalStep := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				// TODO : wire this up.
				"approver_id": {
					Type:     schema.TypeString,
					Optional: true,
					// TODO : validation - is this a UUID or int?
				},
				"rank": rank,
				"isAutomated": {
					Type:     schema.TypeBool,
					Required: true,
					Default:  true,
				},
				"isNotificationOn": {
					Type:     schema.TypeBool,
					Required: true,
					Default:  false,
				},
			},
		},
	}

	releaseDefinitionApprovals := &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"approvals":        releaseDefinitionApprovalStep,
				"approval_options": approvalOptions,
			},
		},
	}

	releaseDefinitionEnvironment := &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MinItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"rank": rank,
				// "owner": owner
				"variables":            configurationVariableMap,
				"variable_groups":      variableGroups,
				"pre_deploy_approvals": releaseDefinitionApprovals,
				"deploy_step":          releaseDefinitionDeployStep,
			},
		},
	}

	return &schema.Resource{
		Create: resourceReleaseDefinitionCreate,
		Read:   resourceReleaseDefinitionRead,
		Update: resourceReleaseDefinitionUpdate,
		Delete: resourceReleaseDefinitionDelete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"revision": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"path": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "\\",
				ValidateFunc: validate.FilePathOrEmpty,
			},
			"variable_groups": variableGroups,
			"source": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ibiza", "portalExtensionApi", "restApi", "undefined", "userInterface"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"variables": configurationVariableMap,
			"release_name_format": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Release-$(rev:r)",
			},
			"environments": releaseDefinitionEnvironment,

			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}

	/*
		"rename_me": {
			Type:     schema.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"rename_me": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
	*/

}

func resourceReleaseDefinitionCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	releaseDefinition, projectID, err := expandReleaseDefinition(d)
	if err != nil {
		return fmt.Errorf("error creating resource Build Definition: %+v", err)
	}

	createdReleaseDefinition, err := createReleaseDefinition(clients, releaseDefinition, projectID)
	if err != nil {
		return fmt.Errorf("error creating resource Build Definition: %+v", err)
	}

	flattenReleaseDefinition(d, createdReleaseDefinition, projectID)
	return nil
}

func flattenReleaseDefinition(d *schema.ResourceData, releaseDefinition *release.ReleaseDefinition, projectID string) {
	d.SetId(strconv.Itoa(*releaseDefinition.Id))

	d.Set("project_id", projectID)
	d.Set("name", *releaseDefinition.Name)
	d.Set("path", *releaseDefinition.Path)
	d.Set("variable_groups", *releaseDefinition.VariableGroups)
	d.Set("source", *releaseDefinition.Source)
	d.Set("description", *releaseDefinition.Description)
	d.Set("variables", flattenReleaseDefinitionVariables(releaseDefinition))
	d.Set("release_name_format", *releaseDefinition.ReleaseNameFormat)
	d.Set("url", *releaseDefinition.Url)
	d.Set("is_deleted", *releaseDefinition.IsDeleted)

	revision := 0
	if releaseDefinition.Revision != nil {
		revision = *releaseDefinition.Revision
	}

	d.Set("revision", revision)
}

// Convert AzDO Variables data structure to Terraform TypeSet
func flattenReleaseDefinitionVariables(variableGroup *release.ReleaseDefinition) interface{} {
	// Preallocate list of variable prop maps
	variables := make([]map[string]interface{}, len(*variableGroup.Variables))

	index := 0
	for k, v := range *variableGroup.Variables {
		variables[index] = map[string]interface{}{
			"name":      k,
			"value":     converter.ToString(v.Value, ""),
			"is_secret": converter.ToBool(v.IsSecret, false),
		}
		index = index + 1
	}

	return variables
}

func createReleaseDefinition(clients *config.AggregatedClient, releaseDefinition *release.ReleaseDefinition, project string) (*release.ReleaseDefinition, error) {
	createdBuild, err := clients.ReleaseClient.CreateReleaseDefinition(clients.Ctx, release.CreateReleaseDefinitionArgs{
		ReleaseDefinition: releaseDefinition,
		Project:           &project,
	})

	return createdBuild, err
}

func resourceReleaseDefinitionRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	projectID, releaseDefinitionID, err := tfhelper.ParseProjectIDAndResourceID(d)

	if err != nil {
		return err
	}

	releaseDefinition, err := clients.ReleaseClient.GetReleaseDefinition(clients.Ctx, release.GetReleaseDefinitionArgs{
		Project:      &projectID,
		DefinitionId: &releaseDefinitionID,
	})

	if err != nil {
		return err
	}

	flattenReleaseDefinition(d, releaseDefinition, projectID)
	return nil
}

func resourceReleaseDefinitionDelete(d *schema.ResourceData, m interface{}) error {
	if d.Id() == "" {
		return nil
	}

	clients := m.(*config.AggregatedClient)
	projectID, releaseDefinitionID, err := tfhelper.ParseProjectIDAndResourceID(d)
	if err != nil {
		return err
	}

	err = clients.ReleaseClient.DeleteReleaseDefinition(m.(*config.AggregatedClient).Ctx, release.DeleteReleaseDefinitionArgs{
		Project:      &projectID,
		DefinitionId: &releaseDefinitionID,
	})

	return err
}

func resourceReleaseDefinitionUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	releaseDefinition, projectID, err := expandReleaseDefinition(d)
	if err != nil {
		return err
	}

	updatedReleaseDefinition, err := clients.ReleaseClient.UpdateReleaseDefinition(m.(*config.AggregatedClient).Ctx, release.UpdateReleaseDefinitionArgs{
		ReleaseDefinition: releaseDefinition,
		Project:           &projectID,
	})

	if err != nil {
		return err
	}

	flattenReleaseDefinition(d, updatedReleaseDefinition, projectID)
	return nil
}

func expandReleaseDefinition(d *schema.ResourceData) (*release.ReleaseDefinition, string, error) {
	projectID := d.Get("project_id").(string)

	// Look for the ID. This may not exist if we are within the context of a "create" operation,
	// so it is OK if it is missing.
	releaseDefinitionID, err := strconv.Atoi(d.Id())
	var releaseDefinitionReference *int
	if err == nil {
		releaseDefinitionReference = &releaseDefinitionID
	} else {
		releaseDefinitionReference = nil
	}

	variableGroups := d.Get("variable_groups").(*schema.Set).List()
	variableGroupsMap := make([]int, len(variableGroups))
	for i, variableGroup := range variableGroups {
		variableGroupsMap[i] = variableGroup.(int)
	}

	releaseDefinition := release.ReleaseDefinition{
		Id:          releaseDefinitionReference,
		Name:        converter.String(d.Get("name").(string)),
		Path:        converter.String(d.Get("path").(string)),
		Revision:    converter.Int(d.Get("revision").(int)),
		Source:      &release.ReleaseDefinitionSourceValues.RestApi,
		Description: converter.String(d.Get("description").(string)),
		// Variables:
		ReleaseNameFormat: converter.String(d.Get("release_name_format").(string)),
		VariableGroups:    &variableGroupsMap,
	}

	data, err := json.MarshalIndent(releaseDefinition, "", "\t")
	fmt.Println(string(data))

	return &releaseDefinition, projectID, nil
}
