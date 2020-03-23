package azuredevops

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/tfhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/validate"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/build"
)

func resourceBuildDefinition() *schema.Resource {
	var filterSchema = map[string]*schema.Schema{
		"include": {
			Type:          schema.TypeString,
			ConflictsWith: []string{"exclude"},
		},
		"exclude": {
			Type:          schema.TypeString,
			ConflictsWith: []string{"exclude"},
		},
	}

	var branchFilter = &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		MinItems: 1,
		Elem: &schema.Resource{
			Schema: filterSchema,
		},
	}

	var pathFilter = &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: filterSchema,
		},
	}

	var scheduleSchema = map[string]*schema.Schema{
		"branch_filter": branchFilter,
		"schedule_job_id": {
			Type: schema.TypeString,
		},
		"only_on_changes": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"day": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"None", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday", "All"}, false),
		},
		"hour": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"minute": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"time_zone_id": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}

	var schedule = &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: scheduleSchema,
		},
	}

	// TODO : What is this? Why is it on PullRequest and ContinuousIntegration?
	var settingsSourceType = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
	}

	return &schema.Resource{
		Create: resourceBuildDefinitionCreate,
		Read:   resourceBuildDefinitionRead,
		Update: resourceBuildDefinitionUpdate,
		Delete: resourceBuildDefinitionDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, buildDefinitionID, err := ParseImportedProjectIDAndID(meta.(*config.AggregatedClient), d.Id())
				if err != nil {
					return nil, fmt.Errorf("Error parsing the build definition ID from the Terraform resource data: %v", err)
				}
				d.Set("project_id", projectID)
				d.SetId(fmt.Sprintf("%d", buildDefinitionID))

				return []*schema.ResourceData{d}, nil
			},
		},
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
				Default:      `\`,
				ValidateFunc: validate.Path,
			},
			"variable_groups": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntAtLeast(1),
				},
				MinItems: 1,
				Optional: true,
			},
			"agent_pool_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Hosted Ubuntu 1604",
			},
			"repository": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"yml_path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"repo_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"repo_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"GitHub", "TfsGit"}, false),
						},
						"branch_name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "master",
						},
						"service_connection_id": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
			},
			// BuildDefinition.triggers.
			// TODO: if all triggers below are empty create a "None" trigger if the SDK doesn't do it automatically.
			// TODO : convert triggers below to single Trigger array. Assign Enum Type to each from list below.
			// None = 1, ContinuousIntegration = 2, BatchedContinuousIntegration = 4, Schedule = 8, GatedCheckIn = 16,
			// BatchedGatedCheckIn = 32, PullRequest = 64, BuildCompletion = 128,
			// TODO : can you mix and match triggers or have more than 1? If not then add "conflicts_with" to every trigger.
			// TODO : convert "day" on schedule trigger into enum int. see below.
			// None = 0, Monday = 1, Tuesday = 2, Wednesday = 4, Thursday = 8, Friday = 16, Saturday = 32, Sunday = 64, All = 127

			"ci_trigger": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"batch": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"branch_filter": branchFilter,
						"max_concurrent_builds_per_branch": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"path_filter": pathFilter,
						"polling_interval": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"polling_job": {
							Type: schema.TypeInt,
							// TODO : is this required?
							Optional: true,
						},
						"settings_source_type": settingsSourceType,
					},
				},
			},
			"schedule_trigger": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schedule": schedule,
					},
				},
			},
			"gated_checkin_trigger": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path_filter": pathFilter,
						"run_ci": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"use_workspace_mappings": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"pull_request_trigger": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_cancel": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"forks": {
							Type:     schema.TypeSet,
							Optional: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"share_secrets": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"branch_filter": branchFilter,
						"path_filter":   pathFilter,
						// isCommentRequiredForPullRequest && requireCommentsForNonTeamMembersOnly
						"comment_required": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"All", "NonTeamMembers"}, false),
						},
						"settings_source_type": settingsSourceType,
					},
				},
			},
		},
	}
}

func resourceBuildDefinitionCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	buildDefinition, projectID, err := expandBuildDefinition(d)
	if err != nil {
		return fmt.Errorf("Error creating resource Build Definition: %+v", err)
	}

	createdBuildDefinition, err := createBuildDefinition(clients, buildDefinition, projectID)
	if err != nil {
		return fmt.Errorf("Error creating resource Build Definition: %+v", err)
	}

	flattenBuildDefinition(d, createdBuildDefinition, projectID)
	return resourceBuildDefinitionRead(d, m)
}

func flattenBuildDefinition(d *schema.ResourceData, buildDefinition *build.BuildDefinition, projectID string) {
	d.SetId(strconv.Itoa(*buildDefinition.Id))

	d.Set("project_id", projectID)
	d.Set("name", *buildDefinition.Name)
	d.Set("path", *buildDefinition.Path)
	d.Set("repository", flattenRepository(buildDefinition))
	d.Set("agent_pool_name", *buildDefinition.Queue.Pool.Name)

	d.Set("variable_groups", flattenVariableGroups(buildDefinition))

	revision := 0
	if buildDefinition.Revision != nil {
		revision = *buildDefinition.Revision
	}

	d.Set("revision", revision)
}

func flattenVariableGroups(buildDefinition *build.BuildDefinition) []int {
	if buildDefinition.VariableGroups == nil {
		return nil
	}

	variableGroups := make([]int, len(*buildDefinition.VariableGroups))

	for i, variableGroup := range *buildDefinition.VariableGroups {
		variableGroups[i] = *variableGroup.Id
	}

	return variableGroups
}

func createBuildDefinition(clients *config.AggregatedClient, buildDefinition *build.BuildDefinition, project string) (*build.BuildDefinition, error) {
	createdBuild, err := clients.BuildClient.CreateDefinition(clients.Ctx, build.CreateDefinitionArgs{
		Definition: buildDefinition,
		Project:    &project,
	})

	return createdBuild, err
}

func resourceBuildDefinitionRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	projectID, buildDefinitionID, err := tfhelper.ParseProjectIDAndResourceID(d)

	if err != nil {
		return err
	}

	buildDefinition, err := clients.BuildClient.GetDefinition(clients.Ctx, build.GetDefinitionArgs{
		Project:      &projectID,
		DefinitionId: &buildDefinitionID,
	})

	if err != nil {
		return err
	}

	flattenBuildDefinition(d, buildDefinition, projectID)
	return nil
}

func resourceBuildDefinitionDelete(d *schema.ResourceData, m interface{}) error {
	if d.Id() == "" {
		return nil
	}

	clients := m.(*config.AggregatedClient)
	projectID, buildDefinitionID, err := tfhelper.ParseProjectIDAndResourceID(d)
	if err != nil {
		return err
	}

	err = clients.BuildClient.DeleteDefinition(m.(*config.AggregatedClient).Ctx, build.DeleteDefinitionArgs{
		Project:      &projectID,
		DefinitionId: &buildDefinitionID,
	})

	return err
}

func resourceBuildDefinitionUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	buildDefinition, projectID, err := expandBuildDefinition(d)
	if err != nil {
		return err
	}

	updatedBuildDefinition, err := clients.BuildClient.UpdateDefinition(m.(*config.AggregatedClient).Ctx, build.UpdateDefinitionArgs{
		Definition:   buildDefinition,
		Project:      &projectID,
		DefinitionId: buildDefinition.Id,
	})

	if err != nil {
		return err
	}

	flattenBuildDefinition(d, updatedBuildDefinition, projectID)
	return nil
}

func flattenRepository(buildDefiniton *build.BuildDefinition) interface{} {
	yamlFilePath := ""

	// The process member can be of many types -- the only typing information
	// available from the compiler is `interface{}` so we can probe for known
	// implementations
	if processMap, ok := buildDefiniton.Process.(map[string]interface{}); ok {
		yamlFilePath = processMap["yamlFilename"].(string)
	}

	if yamlProcess, ok := buildDefiniton.Process.(*build.YamlProcess); ok {
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

func expandBuildDefinitionTrigger(d map[string]interface{}, triggerType build.DefinitionTriggerType) interface{} {
	switch triggerType {
	case build.DefinitionTriggerTypeValues.ContinuousIntegration:
		return build.ContinuousIntegrationTrigger{
			BatchChanges: converter.Bool(d["batch"].(bool)),
			// TODO : map values
		}
	case build.DefinitionTriggerTypeValues.Schedule:
		return build.ScheduleTrigger{
			// TODO : map values
		}
	case build.DefinitionTriggerTypeValues.GatedCheckIn:
		return build.GatedCheckInTrigger{
			// TODO : map values
		}
	case build.DefinitionTriggerTypeValues.PullRequest:
		return build.PullRequestTrigger{
			// TODO : map values
		}
	}
	return nil
}
func expandBuildDefinitionTriggerList(d []interface{}, triggerType build.DefinitionTriggerType) []interface{} {
	vs := make([]interface{}, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			vs = append(vs, expandBuildDefinitionTrigger(val, triggerType))
		}
	}
	return vs
}
func expandBuildDefinitionTriggerSet(configured *schema.Set, triggerType build.DefinitionTriggerType) []interface{} {
	return expandBuildDefinitionTriggerList(configured.List(), triggerType)
}

func expandBuildDefinition(d *schema.ResourceData) (*build.BuildDefinition, string, error) {
	projectID := d.Get("project_id").(string)
	repositories := d.Get("repository").(*schema.Set).List()

	variableGroupsInterface := d.Get("variable_groups").(*schema.Set).List()
	variableGroups := make([]build.VariableGroup, len(variableGroupsInterface))

	for i, variableGroup := range variableGroupsInterface {
		variableGroups[i] = *buildVariableGroup(variableGroup.(int))
	}

	// Note: If configured, this will be of length 1 based on the schema definition above.
	if len(repositories) != 1 {
		return nil, "", fmt.Errorf("Unexpectedly did not find repository metadata in the resource data")
	}

	repository := repositories[0].(map[string]interface{})

	repoName := repository["repo_name"].(string)
	repoType := repository["repo_type"].(string)
	repoURL := ""
	if strings.EqualFold(repoType, "github") {
		repoURL = fmt.Sprintf("https://github.com/%s.git", repoName)
	}

	ciTriggers := expandBuildDefinitionTriggerSet(d.Get("ci_trigger").(*schema.Set), build.DefinitionTriggerTypeValues.ContinuousIntegration)
	scheduleTriggers := expandBuildDefinitionTriggerSet(d.Get("schedule_trigger").(*schema.Set), build.DefinitionTriggerTypeValues.Schedule)
	gatedCheckinTriggers := expandBuildDefinitionTriggerSet(d.Get("gated_checkin_trigger").(*schema.Set), build.DefinitionTriggerTypeValues.GatedCheckIn)
	pullRequestTriggers := expandBuildDefinitionTriggerSet(d.Get("pull_request_trigger").(*schema.Set), build.DefinitionTriggerTypeValues.PullRequest)

	buildTriggers := append(append(append(ciTriggers, scheduleTriggers...), gatedCheckinTriggers...), pullRequestTriggers...)

	// Look for the ID. This may not exist if we are within the context of a "create" operation,
	// so it is OK if it is missing.
	buildDefinitionID, err := strconv.Atoi(d.Id())
	var buildDefinitionReference *int
	if err == nil {
		buildDefinitionReference = &buildDefinitionID
	} else {
		buildDefinitionReference = nil
	}

	agentPoolName := d.Get("agent_pool_name").(string)
	buildDefinition := build.BuildDefinition{
		Id:       buildDefinitionReference,
		Name:     converter.String(d.Get("name").(string)),
		Path:     converter.String(d.Get("path").(string)),
		Revision: converter.Int(d.Get("revision").(int)),
		Repository: &build.BuildRepository{
			Url:           &repoURL,
			Id:            &repoName,
			Name:          &repoName,
			DefaultBranch: converter.String(repository["branch_name"].(string)),
			Type:          &repoType,
			Properties: &map[string]string{
				"connectedServiceId": repository["service_connection_id"].(string),
			},
		},
		Process: &build.YamlProcess{
			YamlFilename: converter.String(repository["yml_path"].(string)),
		},
		Queue: &build.AgentPoolQueue{
			Name: &agentPoolName,
			Pool: &build.TaskAgentPoolReference{
				Name: &agentPoolName,
			},
		},
		QueueStatus:    &build.DefinitionQueueStatusValues.Enabled,
		Type:           &build.DefinitionTypeValues.Build,
		Quality:        &build.DefinitionQualityValues.Definition,
		VariableGroups: &variableGroups,
		Triggers:       &buildTriggers,
	}

	return &buildDefinition, projectID, nil
}

func buildVariableGroup(id int) *build.VariableGroup {
	return &build.VariableGroup{
		Id: &id,
	}
}
