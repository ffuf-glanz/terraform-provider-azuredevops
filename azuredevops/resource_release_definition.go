package azuredevops

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strconv"

	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/tfhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/validate"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/microsoft/azure-devops-go-api/azuredevops/release"
	"github.com/microsoft/azure-devops-go-api/azuredevops/webapi"
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

	configurationVariableValue := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"value": {
			Type:     schema.TypeString,
			Required: true,
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
	}

	configurationVariables := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: configurationVariableValue,
		},
		Set: func(i interface{}) int {
			item := i.(map[string]interface{})
			return schema.HashString(item["name"].(string))
		},
	}

	//taskInputValidation := map[string]*schema.Schema{
	//	"expression": {
	//		Type:     schema.TypeString,
	//		Required: true,
	//	},
	//	"message": {
	//		Type:     schema.TypeString,
	//		Optional: true,
	//	},
	//}

	workFlowTask := map[string]*schema.Schema{
		"always_run": {
			Type:     schema.TypeBool,
			Required: true,
		},
		"condition": {
			Type:     schema.TypeString,
			Required: true,
		},
		"continue_on_error": {
			Type:     schema.TypeBool,
			Required: true,
		},
		"definition_type": {
			Type:     schema.TypeString,
			Required: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Required: true,
		},
		// TODO : Define obj
		"environment": {
			Type:     schema.TypeString,
			Required: true,
		},
		// TODO : Define obj
		"inputs": {
			Type:     schema.TypeString,
			Required: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		// TODO : Define obj
		"override_inputs": {
			Type:     schema.TypeString,
			Required: true,
		},
		"ref_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"task_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"timeout_in_minutes": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"version": {
			Type:     schema.TypeString,
			Required: true,
		},
	}

	workflowTasks := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: workFlowTask,
		},
	}

	releaseDefinitionDeployStep := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  0,
				},
				"tasks": workflowTasks,
			},
		},
	}

	rank := &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Default:  1,
	}

	approvalOptions := &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		MinItems: 1,
		MaxItems: 1,
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
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(release.ApprovalExecutionOrderValues.AfterGatesAlways),
						string(release.ApprovalExecutionOrderValues.AfterSuccessfulGates),
						string(release.ApprovalExecutionOrderValues.BeforeGates),
					}, false),
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

	releaseDefinitionApprovalStep := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
		"approver_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validate.UUID,
		},
		"rank": rank,
		"is_automated": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"is_notification_on": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}

	approvals := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: releaseDefinitionApprovalStep,
		},
	}

	releaseDefinitionGatesOptions := &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"is_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"minimum_success_duration": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"sampling_interval": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"stabilization_time": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"timeout": {
					Type:     schema.TypeInt,
					Optional: true,
				},
			},
		},
	}

	releaseDefinitionGate := map[string]*schema.Schema{
		"tasks": workflowTasks,
	}

	releaseDefinitionGates := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: releaseDefinitionGate,
		},
	}

	releaseDefinitionApprovals := &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"approvals":        approvals,
				"approval_options": approvalOptions,
			},
		},
	}

	releaseDefinitionGatesStep := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  0,
				},
				"gates":         releaseDefinitionGates,
				"gates_options": releaseDefinitionGatesOptions,
			},
		},
	}

	environmentRetentionPolicy := &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"days_to_keep": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  30,
				},
				"releases_to_keep": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  3,
				},
				"retain_build": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
				},
			},
		},
	}

	deployPhase := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"phase_type": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(release.DeployPhaseTypesValues.AgentBasedDeployment),
				string(release.DeployPhaseTypesValues.DeploymentGates),
				string(release.DeployPhaseTypesValues.MachineGroupBasedDeployment),
				string(release.DeployPhaseTypesValues.RunOnServer),
				string(release.DeployPhaseTypesValues.Undefined),
			}, false),
		},
		"rank": rank,
		"ref_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"workflow_tasks": workflowTasks,
		// TODO : This is missing from the docs
		// "deploymentInput": deploymentInput
	}

	deployPhases := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MinItems: 1,
		Elem: &schema.Resource{
			Schema: deployPhase,
		},
	}

	environmentOptions := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"auto_link_work_items": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"badge_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"publish_deployment_status": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"pull_request_deployment_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}

	demand := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"value": {
			Type:     schema.TypeString,
			Required: true,
		},
	}

	demands := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: demand,
		},
	}

	condition := map[string]*schema.Schema{
		"condition_type": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(release.ConditionTypeValues.Undefined),
				string(release.ConditionTypeValues.Artifact),
				string(release.ConditionTypeValues.EnvironmentState),
				string(release.ConditionTypeValues.Event),
			}, false),
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"value": {
			Type:     schema.TypeString,
			Required: true,
		},
	}

	conditions := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: condition,
		},
	}

	environmentExecutionPolicy := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"concurrency_count": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  1,
				},
				"queue_depth_count": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  0,
				},
			},
		},
	}

	schedule := map[string]*schema.Schema{
		"days_to_release": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(release.ScheduleDaysValues.All),
				string(release.ScheduleDaysValues.Friday),
				string(release.ScheduleDaysValues.Monday),
				string(release.ScheduleDaysValues.None),
				string(release.ScheduleDaysValues.Saturday),
				string(release.ScheduleDaysValues.Sunday),
				string(release.ScheduleDaysValues.Thursday),
				string(release.ScheduleDaysValues.Tuesday),
				string(release.ScheduleDaysValues.Wednesday),
			}, false),
		},
		"job_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"schedule_only_with_changes": {
			Type:     schema.TypeBool,
			Required: true,
		},
		"start_hours": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"start_minutes": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"time_zone_id": {
			Type:     schema.TypeString,
			Required: true,
		},
	}

	schedules := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: schedule,
		},
	}

	properties := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}

	environmentTrigger := map[string]*schema.Schema{
		"definition_environment_id": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"release_definition_id": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"trigger_content": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"trigger_type": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(release.EnvironmentTriggerTypeValues.Undefined),
				string(release.EnvironmentTriggerTypeValues.DeploymentGroupRedeploy),
				string(release.EnvironmentTriggerTypeValues.RollbackRedeploy),
			}, false),
		},
	}

	environmentTriggers := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: environmentTrigger,
		},
	}

	releaseDefinitionEnvironment := &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MinItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  0,
				},
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"rank": rank,
				// TODO : Is this something you would want to set
				"owner_id": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validate.UUID,
				},
				"variable":              configurationVariables,
				"variable_groups":       variableGroups,
				"pre_deploy_approvals":  releaseDefinitionApprovals,
				"deploy_step":           releaseDefinitionDeployStep,
				"post_deploy_approvals": releaseDefinitionApprovals,
				"deploy_phases":         deployPhases,
				// TODO : This is missing from the docs
				// "runOptions": runOptions
				"environment_options":   environmentOptions,
				"demands":               demands,
				"conditions":            conditions,
				"execution_policy":      environmentExecutionPolicy,
				"schedules":             schedules,
				"properties":            properties,
				"pre_deployment_gates":  releaseDefinitionGatesStep,
				"post_deployment_gates": releaseDefinitionGatesStep,
				"environment_triggers":  environmentTriggers,
				"retention_policy":      environmentRetentionPolicy,
				"badge_url": {
					Type:     schema.TypeString,
					Computed: true,
				},
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
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(release.ReleaseDefinitionSourceValues.Undefined),
					string(release.ReleaseDefinitionSourceValues.RestApi),
					string(release.ReleaseDefinitionSourceValues.PortalExtensionApi),
					string(release.ReleaseDefinitionSourceValues.Ibiza),
					string(release.ReleaseDefinitionSourceValues.UserInterface),
				}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"variable": configurationVariables,
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

			"created_on": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"modified_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
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
	d.Set("variable", flattenReleaseDefinitionVariables(releaseDefinition))
	d.Set("release_name_format", *releaseDefinition.ReleaseNameFormat)
	d.Set("url", *releaseDefinition.Url)
	d.Set("is_deleted", *releaseDefinition.IsDeleted)
	d.Set("created_on", *releaseDefinition.CreatedOn)
	d.Set("modified_on", *releaseDefinition.ModifiedOn)

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

	variableGroups := buildVariableGroups(d.Get("variable_groups").([]interface{}))
	environments, environmentsError := buildEnvironments(d.Get("environments").([]interface{}))
	if environmentsError != nil {
		return nil, "", environmentsError
	}
	variables, variablesError := buildVariables(d.Get("variable").(*schema.Set).List())
	if variablesError != nil {
		return nil, "", variablesError
	}

	releaseDefinition := release.ReleaseDefinition{
		Id:                releaseDefinitionReference,
		Name:              converter.String(d.Get("name").(string)),
		Path:              converter.String(d.Get("path").(string)),
		Revision:          converter.Int(d.Get("revision").(int)),
		Source:            &release.ReleaseDefinitionSourceValues.RestApi,
		Description:       converter.String(d.Get("description").(string)),
		Environments:      &environments,
		Variables:         &variables,
		ReleaseNameFormat: converter.String(d.Get("release_name_format").(string)),
		VariableGroups:    &variableGroups,
	}

	data, err := json.MarshalIndent(releaseDefinition, "", "\t")
	fmt.Println(string(data))

	return &releaseDefinition, projectID, nil
}

func buildEnvironments(environments []interface{}) ([]release.ReleaseDefinitionEnvironment, error) {
	environmentsMap := make([]release.ReleaseDefinitionEnvironment, len(environments))
	for i, environment := range environments {
		env, err := buildReleaseDefinitionEnvironment(environment.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		environmentsMap[i] = *env
	}
	return environmentsMap, nil
}

func buildVariableGroups(variableGroups []interface{}) []int {
	variableGroupsMap := make([]int, len(variableGroups))
	for i, variableGroup := range variableGroups {
		variableGroupsMap[i] = variableGroup.(int)
	}
	return variableGroupsMap
}

func buildVariables(variables []interface{}) (map[string]release.ConfigurationVariableValue, error) {
	variablesMap := make(map[string]release.ConfigurationVariableValue)
	for _, variable := range variables {
		key, variableMap, err := buildVariable(variable.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		variablesMap[key] = variableMap
	}
	return variablesMap, nil
}

func buildVariable(d map[string]interface{}) (string, release.ConfigurationVariableValue, error) {
	return d["name"].(string), release.ConfigurationVariableValue{
		AllowOverride: converter.Bool(d["allow_override"].(bool)),
		Value:         converter.String(d["value"].(string)),
		IsSecret:      converter.Bool(d["is_secret"].(bool)),
	}, nil
}

func buildReleaseDefinitionEnvironment(d map[string]interface{}) (*release.ReleaseDefinitionEnvironment, error) {
	variableGroups := d["variable_groups"].([]interface{})
	variableGroupsMap := make([]int, len(variableGroups))
	for i, variableGroup := range variableGroups {
		variableGroupsMap[i] = variableGroup.(int)
	}

	var retentionPolicyMap *release.EnvironmentRetentionPolicy
	if d["retention_policy"] != nil {
		retentionPolicy := d["retention_policy"].(*schema.Set).List()
		if len(retentionPolicy) != 1 {
			return nil, fmt.Errorf("unexpectedly did not find a retention policy in the environment data")
		}
		environmentRetentionPolicy, err := buildEnvironmentRetentionPolicy(retentionPolicy[0].(map[string]interface{}))
		retentionPolicyMap = environmentRetentionPolicy
		if err != nil {
			return nil, err
		}
	}

	var preDeployApprovalsMap *release.ReleaseDefinitionApprovals
	if d["pre_deploy_approvals"] != nil {
		preDeployApprovals := d["pre_deploy_approvals"].(*schema.Set).List()
		if len(preDeployApprovals) != 1 {
			return nil, fmt.Errorf("unexpectedly did not find a pre deploy approval in the environment data")
		}
		environmentRetentionPolicy, err := buildReleaseDefinitionApprovals(preDeployApprovals[0].(map[string]interface{}))
		preDeployApprovalsMap = environmentRetentionPolicy
		if err != nil {
			return nil, err
		}
	}

	var postDeployApprovalsMap *release.ReleaseDefinitionApprovals
	if d["post_deploy_approvals"] != nil {
		postDeployApprovals := d["post_deploy_approvals"].(*schema.Set).List()
		if len(postDeployApprovals) != 1 {
			return nil, fmt.Errorf("unexpectedly did not find a post deploy approval in the environment data")
		}
		environmentRetentionPolicy, err := buildReleaseDefinitionApprovals(postDeployApprovals[0].(map[string]interface{}))
		postDeployApprovalsMap = environmentRetentionPolicy
		if err != nil {
			return nil, err
		}
	}

	var deployStepMap *release.ReleaseDefinitionDeployStep
	if d["deploy_step"] != nil {
		deployStep := d["deploy_step"].(*schema.Set).List()
		if len(deployStep) != 1 {
			return nil, fmt.Errorf("unexpectedly did not find a deploy step in the environment data")
		}
		releaseDefinitionDeployStep, err := buildReleaseDefinitionDeployStep(deployStep[0].(map[string]interface{}))
		deployStepMap = releaseDefinitionDeployStep
		if err != nil {
			return nil, err
		}
	}

	variables, variablesError := buildVariables(d["variable"].(*schema.Set).List())
	if variablesError != nil {
		return nil, variablesError
	}

	conditions := d["conditions"].(*schema.Set).List()
	conditionsMap := make([]release.Condition, len(conditions))
	for i, condition := range conditions {
		asMap := condition.(map[string]interface{})
		conditionType := release.ConditionType(asMap["condition_type"].(string))
		conditionsMap[i] = release.Condition{
			ConditionType: &conditionType,
			Name:          converter.String(d["name"].(string)),
			Value:         converter.String(d["value"].(string)),
		}
	}

	releaseDefinitionEnvironment := release.ReleaseDefinitionEnvironment{
		Conditions:          &conditionsMap,
		CurrentRelease:      nil,
		Demands:             nil,
		DeployPhases:        nil,
		DeployStep:          deployStepMap,
		EnvironmentOptions:  nil,
		EnvironmentTriggers: nil,
		ExecutionPolicy:     nil,
		Id:                  converter.Int(d["id"].(int)),
		Name:                converter.String(d["name"].(string)),
		Owner: &webapi.IdentityRef{
			Id: converter.String(d["owner_id"].(string)),
		},
		PostDeployApprovals: postDeployApprovalsMap,
		PostDeploymentGates: nil,
		PreDeployApprovals:  preDeployApprovalsMap,
		PreDeploymentGates:  nil,
		ProcessParameters:   nil,
		Properties:          nil,
		QueueId:             nil,
		Rank:                converter.Int(d["rank"].(int)),
		RetentionPolicy:     retentionPolicyMap,
		RunOptions:          nil,
		Schedules:           nil,
		VariableGroups:      &variableGroupsMap,
		Variables:           &variables,
	}

	return &releaseDefinitionEnvironment, nil
}

func buildReleaseDefinitionDeployStep(d map[string]interface{}) (*release.ReleaseDefinitionDeployStep, error) {
	if d["tasks"] == nil {
		return nil, nil
	}

	tasks, err := buildWorkFlowTasks(d["tasks"].([]interface{}))
	if err != nil {
		return nil, err
	}

	return &release.ReleaseDefinitionDeployStep{
		Id:    converter.Int(d["id"].(int)),
		Tasks: &tasks,
	}, nil
}

func buildWorkFlowTasks(tasks []interface{}) ([]release.WorkflowTask, error) {
	tasksMap := make([]release.WorkflowTask, len(tasks))
	for i, approval := range tasks {
		releaseApproval, err := buildWorkflowTask(approval.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		tasksMap[i] = *releaseApproval
	}
	return tasksMap, nil
}

func buildWorkflowTask(d map[string]interface{}) (*release.WorkflowTask, error) {
	taskId, err := uuid.Parse(d["task_id"].(string))
	if err != nil {
		return nil, fmt.Errorf("error parsing task_id: %s. %v", d["task_id"].(string), err)
	}

	return &release.WorkflowTask{
		AlwaysRun:       converter.Bool(d["always_run"].(bool)),
		Condition:       converter.String(d["condition"].(string)),
		ContinueOnError: converter.Bool(d["continue_on_error"].(bool)),
		DefinitionType:  converter.String(d["definition_type"].(string)),
		Enabled:         converter.Bool(d["enabled"].(bool)),
		// Environment:      converter.String(d["environment"].(string)),
		//Inputs:           converter.Int(d["inputs"].(int)),
		Name: converter.String(d["name"].(string)),
		//OverrideInputs:   converter.Int(d["override_inputs"].(int)),
		RefName:          converter.String(d["ref_name"].(string)),
		TaskId:           &taskId,
		TimeoutInMinutes: converter.Int(d["timeout_in_minutes"].(int)),
		Version:          converter.String(d["version"].(string)),
	}, nil
}

func buildEnvironmentRetentionPolicy(d map[string]interface{}) (*release.EnvironmentRetentionPolicy, error) {
	environmentRetentionPolicy := release.EnvironmentRetentionPolicy{
		DaysToKeep:     converter.Int(d["days_to_keep"].(int)),
		RetainBuild:    converter.Bool(d["retain_build"].(bool)),
		ReleasesToKeep: converter.Int(d["releases_to_keep"].(int)),
	}
	return &environmentRetentionPolicy, nil
}

func buildReleaseDefinitionApprovals(d map[string]interface{}) (*release.ReleaseDefinitionApprovals, error) {
	approvals := d["approvals"].(*schema.Set).List()
	approvalsMap := make([]release.ReleaseDefinitionApprovalStep, len(approvals))

	for i, approval := range approvals {
		releaseApproval, err := buildReleaseApproval(approval.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		approvalsMap[i] = *releaseApproval
	}

	approvalOptions := d["approval_options"].(*schema.Set).List()
	if len(approvalOptions) != 1 {
		return nil, fmt.Errorf("unexpectedly did not find a approval options in the approvals data")
	}
	d2 := approvalOptions[0].(map[string]interface{})
	executionOrder := release.ApprovalExecutionOrder(d2["execution_order"].(string))

	releaseDefinitionApprovals := release.ReleaseDefinitionApprovals{
		Approvals: &approvalsMap,
		ApprovalOptions: &release.ApprovalOptions{
			AutoTriggeredAndPreviousEnvironmentApprovedCanBeSkipped: converter.Bool(d2["auto_triggered_and_previous_environment_approved_can_be_skipped"].(bool)),
			EnforceIdentityRevalidation:                             converter.Bool(d2["enforce_identity_revalidation"].(bool)),
			ExecutionOrder:                                          &executionOrder,
			ReleaseCreatorCanBeApprover:                             converter.Bool(d2["release_creator_can_be_approver"].(bool)),
			RequiredApproverCount:                                   converter.Int(d2["required_approver_count"].(int)),
			TimeoutInMinutes:                                        converter.Int(d2["timeout_in_minutes"].(int)),
		},
	}
	return &releaseDefinitionApprovals, nil
}

func buildReleaseApproval(d map[string]interface{}) (*release.ReleaseDefinitionApprovalStep, error) {
	approver := d["approver"]
	var approverMap *webapi.IdentityRef
	if approver != nil {
		approverMap = &webapi.IdentityRef{
			Id: converter.String(d["approver_id"].(string)),
		}
	}

	releaseDefinitionApprovalStep := release.ReleaseDefinitionApprovalStep{
		Id:               converter.Int(d["id"].(int)),
		Approver:         approverMap,
		IsAutomated:      converter.Bool(d["is_automated"].(bool)),
		IsNotificationOn: converter.Bool(d["is_notification_on"].(bool)),
		Rank:             converter.Int(d["rank"].(int)),
	}
	return &releaseDefinitionApprovalStep, nil
}
