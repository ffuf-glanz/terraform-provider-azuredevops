package azuredevops

import (
	"fmt"
	"github.com/microsoft/azure-devops-go-api/azuredevops/taskagent"
	"strconv"
	"strings"

	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/tfhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/validate"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/release"
)

func resourceReleaseDefinition() *schema.Resource {
	return &schema.Resource{
		Create: resourceReleaseDefinitionCreate,
		Read:   resourceReleaseDefinitionRead,
		Update: resourceReleaseDefinitionUpdate,
		Delete: resourceReleaseDefinitionDelete,

		Schema: map[string]*schema.Schema{
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

			// TODO : Abstract this because it is used by build definitions and release definitions.
			"variable_groups": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntAtLeast(1),
				},
				MinItems: 1,
				Optional: true,
			},
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
			// TODO : Abstract this because it is used by variable group and release definitions.
			"variables": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"is_secret": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
				Required: true,
				MinItems: 1,
				Set: func(i interface{}) int {
					item := i.(map[string]interface{})
					return schema.HashString(item["name"].(string))
				},
			},
			"releaseNameFormat": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Release-$(rev:r)",
			},

			//"repository": {
			//	Type:     schema.TypeSet,
			//	Required: true,
			//	MinItems: 1,
			//	MaxItems: 1,
			//	Elem: &schema.Resource{
			//		Schema: map[string]*schema.Schema{
			//			"yml_path": {
			//				Type:     schema.TypeString,
			//				Required: true,
			//			},
			//			"repo_name": {
			//				Type:     schema.TypeString,
			//				Required: true,
			//			},
			//			"repo_type": {
			//				Type:         schema.TypeString,
			//				Required:     true,
			//				ValidateFunc: validation.StringInSlice([]string{"GitHub", "TfsGit"}, false),
			//			},
			//			"branch_name": {
			//				Type:     schema.TypeString,
			//				Optional: true,
			//				Default:  "master",
			//			},
			//			"service_connection_id": {
			//				Type:     schema.TypeString,
			//				Optional: true,
			//				Default:  "",
			//			},
			//		},
			//	},
			//},

			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"isDeleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceReleaseDefinitionCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	releaseDefinition, err := expandReleaseDefinition(d)
	if err != nil {
		return fmt.Errorf("Error creating resource Build Definition: %+v", err)
	}

	createdReleaseDefinition, err := createReleaseDefinition(clients, releaseDefinition)
	if err != nil {
		return fmt.Errorf("Error creating resource Build Definition: %+v", err)
	}

	flattenReleaseDefinition(d, createdReleaseDefinition)
	return nil
}

func flattenReleaseDefinition(d *schema.ResourceData, releaseDefinition *release.ReleaseDefinition) {
	d.SetId(strconv.Itoa(*releaseDefinition.Id))

	d.Set("name", *releaseDefinition.Name)
	d.Set("path", *releaseDefinition.Path)
	d.Set("variable_groups", *releaseDefinition.VariableGroups)
	d.Set("source", *releaseDefinition.Source)
	d.Set("description", *releaseDefinition.Description)
	d.Set("variables", flattenReleaseDefinitionVariables(releaseDefinition))
	d.Set("releaseNameFormat", *releaseDefinition.ReleaseNameFormat)
	d.Set("url", *releaseDefinition.Url)
	d.Set("isDeleted", *releaseDefinition.IsDeleted)

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

func createReleaseDefinition(clients *config.AggregatedClient, releaseDefinition *release.ReleaseDefinition) (*release.ReleaseDefinition, error) {
	createdBuild, err := clients.BuildClient.CreateDefinition(clients.Ctx, release.CreateDefinitionArgs{
		Definition: releaseDefinition,
		Project:    &project,
	})

	return createdBuild, err
}

func resourceReleaseDefinitionRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	projectID, releaseDefinitionID, err := tfhelper.ParseProjectIDAndResourceID(d)

	if err != nil {
		return err
	}

	releaseDefinition, err := clients.BuildClient.GetDefinition(clients.Ctx, release.GetDefinitionArgs{
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

	err = clients.BuildClient.DeleteDefinition(m.(*config.AggregatedClient).Ctx, release.DeleteReleaseDefinitionArgs{
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

	updatedReleaseDefinition, err := clients.BuildClient.UpdateDefinition(m.(*config.AggregatedClient).Ctx, release.UpdateDefinitionArgs{
		Definition:   releaseDefinition,
		Project:      &projectID,
		DefinitionId: releaseDefinition.Id,
	})

	if err != nil {
		return err
	}

	flattenReleaseDefinition(d, updatedReleaseDefinition, projectID)
	return nil
}

func flattenRepository(buildDefiniton *release.ReleaseDefinition) interface{} {
	yamlFilePath := ""

	// The process member can be of many types -- the only typing information
	// available from the compiler is `interface{}` so we can probe for known
	// implementations
	if processMap, ok := buildDefiniton.Process.(map[string]interface{}); ok {
		yamlFilePath = processMap["yamlFilename"].(string)
	}

	if yamlProcess, ok := buildDefiniton.Process.(*release.YamlProcess); ok {
		yamlFilePath = *yamlProcess.YamlFilename
	}

	return []map[string]interface{}{{
		"yml_path":              yamlFilePath,
		"repo_name":             *buildDefiniton.Repository.Name,
		"repo_type":             *buildDefiniton.Repository.Type,
		"branch_name":           *buildDefiniton.Repository.DefaultBranch,
		"service_connection_id": (*buildDefiniton.Repository.Properties)["connectedServiceId"],
	}}
}

func expandReleaseDefinition(d *schema.ResourceData) (*release.ReleaseDefinition, error) {
	repositories := d.Get("repository").(*schema.Set).List()

	variableGroupsInterface := d.Get("variable_groups").(*schema.Set).List()
	variableGroups := make([]int, len(variableGroupsInterface))

	// Note: If configured, this will be of length 1 based on the schema definition above.
	//if len(repositories) != 1 {
	//	return nil, fmt.Errorf("unexpectedly did not find repository metadata in the resource data")
	//}

	// Look for the ID. This may not exist if we are within the context of a "create" operation,
	// so it is OK if it is missing.
	releaseDefinitionID, err := strconv.Atoi(d.Id())
	var releaseDefinitionReference *int
	if err == nil {
		releaseDefinitionReference = &releaseDefinitionID
	} else {
		releaseDefinitionReference = nil
	}

	releaseDefinition := release.ReleaseDefinition{
		Id:       releaseDefinitionReference,
		Name:     converter.String(d.Get("name").(string)),
		Path:     converter.String(d.Get("path").(string)),
		Revision: converter.Int(d.Get("revision").(int)),
		// Source: release.ReleaseDefinitionSourceValues.RestApi,
		Description: converter.String(d.Get("description").(string)),
		// Variables:
		ReleaseNameFormat: converter.String(d.Get("releaseNameFormat").(string)),
		VariableGroups:    d.Get("variableGroups").(*[]int),
	}

	return &releaseDefinition, nil
}

func buildVariableGroup(id int) *release.VariableGroup {
	return &release.VariableGroup{
		Id: &id,
	}
}
