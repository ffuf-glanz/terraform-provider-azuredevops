package azuredevops

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/release"
	"github.com/microsoft/azure-devops-go-api/azuredevops/webapi"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
)

type Properties struct {
	DefinitionCreationSource *string
	IntegrateJiraWorkItems   *bool
	IntegrateBoardsWorkItems *bool
}

type ReleaseDeployPhase struct {
	// Dynamic based on PhaseType
	DeploymentInput interface{} `json:"deploymentInput,omitempty"`
	// WorkflowTasks
	WorkflowTasks *[]release.WorkflowTask `json:"workflowTasks,omitempty"`
	// Gets or sets the reference name of the task.
	RefName *string `json:"refName,omitempty"`
	// Name of the phase.
	Name *string `json:"name,omitempty"`
	// Type of the phase.
	PhaseType *release.DeployPhaseTypes `json:"phaseType,omitempty"`
	// Rank of the phase.
	Rank *int `json:"rank,omitempty"`

	// Deployment jobs of the phase.
	//DeploymentJobs *[]release.DeploymentJob `json:"deploymentJobs,omitempty"`

	// Phase execution error logs.
	//ErrorLog *string `json:"errorLog,omitempty"`

	// Deprecated:
	//Id *int `json:"id,omitempty"`

	// List of manual intervention tasks execution information in phase.
	//ManualInterventions *[]release.ManualIntervention `json:"manualInterventions,omitempty"`

	// ID of the phase.
	//PhaseId *string `json:"phaseId,omitempty"`

	// Run Plan ID of the phase.
	//RunPlanId *uuid.UUID `json:"runPlanId,omitempty"`

	// Phase start time.
	//StartedOn *azuredevops.Time `json:"startedOn,omitempty"`

	// Status of the phase.
	//Status *release.DeployPhaseStatus `json:"status,omitempty"`
}

type ArtifactDownloadModeType string

type artifactDownloadModeTypeValuesType struct {
	Skip      ArtifactDownloadModeType
	Selective ArtifactDownloadModeType
	All       ArtifactDownloadModeType
}

var ArtifactDownloadModeTypeValues = artifactDownloadModeTypeValuesType{
	Skip:      "Skip",
	Selective: "Selective",
	All:       "All",
}

type ReleaseDefinitionDemand struct {
	// Name of the demand.
	Name *string `json:"name,omitempty"`
	// The value of the demand.
	Value *string `json:"value,omitempty"`
}

// Takes the result of flatmap.Expand for an array of strings and returns a []*string
func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(string))
		}
	}
	return vs
}

// Takes the result of flatmap.Expand for an array of strings and returns a []*int
func expandIntList(configured []interface{}) []int {
	vs := make([]int, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(int))
		}
	}
	return vs
}

// Takes the result of schema.Set of strings and returns a []*string
func expandStringSet(configured *schema.Set) []*string {
	return expandStringList(configured.List())
}

// Takes list of pointers to strings. Expand to an array of raw strings and returns a []interface{} to keep compatibility w/ schema.NewSetschema.NewSet
func flattenStringList(list []*string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, *v)
	}
	return vs
}

func flattenStringSet(list []*string) *schema.Set {
	return schema.NewSet(schema.HashString, flattenStringList(list))
}

func expandReleaseDefinitionsProperties(d []interface{}) (interface{}, error) {
	if len(d) != 1 {
		return nil, fmt.Errorf("unexpectedly did not find a properties block in the env")
	}
	asMap := d[0].(map[string]interface{})
	return &Properties{
		DefinitionCreationSource: converter.String(asMap["definition_creation_source"].(string)),
		IntegrateJiraWorkItems:   converter.Bool(asMap["integrate_jira_work_items"].(bool)),
		IntegrateBoardsWorkItems: converter.Bool(asMap["integrate_boards_work_items"].(bool)),
	}, nil
}

func expandReleaseDefinitionEnvironmentList(environments []interface{}) ([]release.ReleaseDefinitionEnvironment, error) {
	environmentsMap := make([]release.ReleaseDefinitionEnvironment, len(environments))
	for i, environment := range environments {
		env, err := expandReleaseDefinitionEnvironment(environment.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		environmentsMap[i] = *env
	}
	return environmentsMap, nil
}

func expandReleaseArtifactList(d interface{}) ([]release.Artifact, error) {
	if d != nil {
		return make([]release.Artifact, 0), nil
	}
	d2 := d.([]interface{})
	artifacts := make([]release.Artifact, len(d2))
	for i, d3 := range d2 {
		artifact, err := expandReleaseArtifact(d3.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		artifacts[i] = artifact
	}
	return artifacts, nil
}

func expandReleaseArtifact(d map[string]interface{}) (release.Artifact, error) {
	// NOTE : Might have to build a custom struct because AgentArtifactType doesn't equal comments in code. See below.
	// It can have value as 'Build', 'Jenkins', 'GitHub', 'Nuget', 'Team Build (external)', 'ExternalTFSBuild', 'Git', 'TFVC', 'ExternalTfsXamlBuild'.
	artifactType := release.AgentArtifactType(d["type"].(string))
	return release.Artifact{
		Alias: converter.String(d["Alias"].(string)),
		// DefinitionReference: converter.Bool(d["DefinitionReference"].(bool)),
		IsPrimary:  converter.Bool(d["IsPrimary"].(bool)),
		IsRetained: converter.Bool(d["IsRetained"].(bool)),
		SourceId:   converter.String(d["SourceId"].(string)),
		Type:       converter.String(string(artifactType)),
	}, nil
}

func expandReleaseConfigurationVariableValueSet(variables []interface{}) (map[string]release.ConfigurationVariableValue, error) {
	variablesMap := make(map[string]release.ConfigurationVariableValue)
	for _, variable := range variables {
		key, variableMap, err := expandReleaseConfigurationVariableValue(variable.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		variablesMap[key] = variableMap
	}
	return variablesMap, nil
}

func expandReleaseConfigurationVariableValue(d map[string]interface{}) (string, release.ConfigurationVariableValue, error) {
	return d["name"].(string), release.ConfigurationVariableValue{
		AllowOverride: converter.Bool(d["allow_override"].(bool)),
		Value:         converter.String(d["value"].(string)),
		IsSecret:      converter.Bool(d["is_secret"].(bool)),
	}, nil
}

func expandReleaseDefinitionEnvironment(d map[string]interface{}) (*release.ReleaseDefinitionEnvironment, error) {
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
		environmentRetentionPolicy, err := expandReleaseEnvironmentRetentionPolicy(retentionPolicy[0].(map[string]interface{}))
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
		environmentRetentionPolicy, err := expandReleaseDefinitionApprovals(preDeployApprovals[0].(map[string]interface{}))
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
		environmentRetentionPolicy, err := expandReleaseDefinitionApprovals(postDeployApprovals[0].(map[string]interface{}))
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
		releaseDefinitionDeployStep, err := expandReleaseDefinitionDeployStep(deployStep[0].(map[string]interface{}))
		deployStepMap = releaseDefinitionDeployStep
		if err != nil {
			return nil, err
		}
	}

	variables, variablesError := expandReleaseConfigurationVariableValueSet(d["variable"].(*schema.Set).List())
	if variablesError != nil {
		return nil, variablesError
	}

	conditions := d["conditions"].(*schema.Set).List()
	conditionsMap := make([]release.Condition, len(conditions))
	for i, condition := range conditions {
		asMap := condition.(map[string]interface{})
		conditionType := release.ConditionType(asMap["condition_type"].(string))
		value := ""
		if asMap["value"] != nil {
			value = asMap["value"].(string)
		}
		conditionsMap[i] = release.Condition{
			ConditionType: &conditionType,
			Name:          converter.String(asMap["name"].(string)),
			Value:         converter.String(value),
		}
	}

	demands, demandsError := expandReleaseDefinitionDemandList(d["demands"].(*schema.Set).List())
	if demandsError != nil {
		return nil, demandsError
	}

	deployPhases, deployPhasesError := expandReleaseDeployPhaseList(d["deploy_phases"].([]interface{}))
	if deployPhasesError != nil {
		return nil, deployPhasesError
	}

	environmentOptions, environmentOptionsError := expandReleaseEnvironmentOptions(d["environment_options"].(*schema.Set).List())
	if environmentOptionsError != nil {
		return nil, environmentOptionsError
	}

	releaseDefinitionEnvironment := release.ReleaseDefinitionEnvironment{
		Conditions:          &conditionsMap,
		Demands:             &demands,
		DeployPhases:        &deployPhases,
		DeployStep:          deployStepMap,
		EnvironmentOptions:  environmentOptions,
		EnvironmentTriggers: nil,
		ExecutionPolicy:     nil,
		Id:                  converter.Int(d["id"].(int)),
		Name:                converter.String(d["name"].(string)),
		Owner: &webapi.IdentityRef{
			Id:            converter.String(d["owner_id"].(string)),
			IsAadIdentity: converter.Bool(true),
			IsContainer:   converter.Bool(false),
		},
		PostDeployApprovals: postDeployApprovalsMap,
		PostDeploymentGates: nil,
		PreDeployApprovals:  preDeployApprovalsMap,
		PreDeploymentGates:  nil,
		ProcessParameters:   nil,
		// Properties:          &releaseDefinitionEnvironmentProperties,
		QueueId:         nil,
		Rank:            converter.Int(d["rank"].(int)),
		RetentionPolicy: retentionPolicyMap,
		RunOptions:      nil,
		Schedules:       nil,
		VariableGroups:  &variableGroupsMap,
		Variables:       &variables,
	}

	return &releaseDefinitionEnvironment, nil
}

func expandReleaseDefinitionDeployStep(d map[string]interface{}) (*release.ReleaseDefinitionDeployStep, error) {
	if d["tasks"] == nil {
		return nil, nil
	}

	tasks, err := expandReleaseWorkFlowTaskList(d["tasks"].([]interface{}))
	if err != nil {
		return nil, err
	}

	return &release.ReleaseDefinitionDeployStep{
		Id:    converter.Int(d["id"].(int)),
		Tasks: &tasks,
	}, nil
}

func expandReleaseDeployPhaseList(deployPhases []interface{}) ([]interface{}, error) {
	deployPhasesMap := make([]interface{}, len(deployPhases))
	for i, deployPhase := range deployPhases {
		demandMap, err := expandReleaseDeployPhase(deployPhase.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		deployPhasesMap[i] = demandMap
	}
	return deployPhasesMap, nil
}

func expandReleaseEnvironmentOptions(d []interface{}) (*release.EnvironmentOptions, error) {
	if len(d) != 1 {
		return nil, fmt.Errorf("unexpectedly did not build environment options")
	}
	asMap := d[0].(map[string]interface{})

	deployPhase := release.EnvironmentOptions{
		AutoLinkWorkItems:            converter.Bool(asMap["auto_link_work_items"].(bool)),
		BadgeEnabled:                 converter.Bool(asMap["badge_enabled"].(bool)),
		EmailNotificationType:        converter.String(asMap["email_notification_type"].(string)),
		EmailRecipients:              converter.String(asMap["email_recipients"].(string)),
		EnableAccessToken:            converter.Bool(asMap["enable_access_token"].(bool)),
		PublishDeploymentStatus:      converter.Bool(asMap["publish_deployment_status"].(bool)),
		PullRequestDeploymentEnabled: converter.Bool(asMap["pull_request_deployment_enabled"].(bool)),
		SkipArtifactsDownload:        converter.Bool(asMap["skip_artifacts_download"].(bool)),
		TimeoutInMinutes:             converter.Int(asMap["timeout_in_minutes"].(int)),
	}
	return &deployPhase, nil
}

func expandReleaseDeployPhase(d map[string]interface{}) (interface{}, error) {
	var deploymentInput interface{}

	phaseType := release.DeployPhaseTypes(d["phase_type"].(string))
	switch phaseType {
	case release.DeployPhaseTypesValues.AgentBasedDeployment:

		agentDeploymentInput := d["agent_deployment_input"].(*schema.Set).List()
		if len(agentDeploymentInput) != 1 {
			return nil, fmt.Errorf("unexpectedly did not find a agent deployment input in the deploy phases data")
		}
		agentDeploymentInputMap, deploymentInputErrors := expandReleaseAgentDeploymentInput(agentDeploymentInput[0].(map[string]interface{}))
		if deploymentInputErrors != nil {
			return nil, deploymentInputErrors
		}
		deploymentInput = agentDeploymentInputMap
	}

	deployPhase := ReleaseDeployPhase{
		DeploymentInput: &deploymentInput,
		Rank:            converter.Int(d["rank"].(int)),
		PhaseType:       &phaseType,
		Name:            converter.String(d["name"].(string)),
		RefName:         converter.String(d["ref_name"].(string)),
		WorkflowTasks:   nil,
	}
	return deployPhase, nil
}

func expandReleaseAgentDeploymentInput(d map[string]interface{}) (interface{}, error) {
	artifactsDownloadInput, err := expandReleaseArtifactsDownloadInput(d["artifacts_download_input"].(*schema.Set).List())
	if err != nil {
		return nil, err
	}

	demands, demandsError := expandReleaseDefinitionDemandList(d["demands"].(*schema.Set).List())
	if demandsError != nil {
		return nil, demandsError
	}

	parallelExecutionType := release.ParallelExecutionTypes(d["parallel_execution_type"].(string))

	return release.AgentDeploymentInput{
		Condition:                 converter.String(d["condition"].(string)),
		JobCancelTimeoutInMinutes: converter.Int(d["job_cancel_timeout_in_minutes"].(int)),
		OverrideInputs:            nil, // TODO : OverrideInputs
		TimeoutInMinutes:          converter.Int(d["timeout_in_minutes"].(int)),
		ArtifactsDownloadInput:    artifactsDownloadInput,
		Demands:                   &demands,
		EnableAccessToken:         converter.Bool(d["enable_access_token"].(bool)),
		QueueId:                   converter.Int(d["queue_id"].(int)),
		SkipArtifactsDownload:     converter.Bool(d["skip_artifacts_download"].(bool)),
		AgentSpecification: &release.AgentSpecification{
			Identifier: converter.String(d["agent_specification_identifier"].(string)),
		},
		// ImageId: converter.Int(d["image_id"].(int)),
		ParallelExecution: &release.ExecutionInput{
			ParallelExecutionType: &parallelExecutionType,
		},
	}, nil
}

func expandReleaseArtifactsDownloadInput(d []interface{}) (*release.ArtifactsDownloadInput, error) {
	if len(d) == 0 {
		return nil, nil
	}
	if len(d) != 1 {
		return nil, fmt.Errorf("unexpectedly did not find an artifacts download input")
	}
	downloadInputs, err := expandReleaseArtifactDownloadInputBases(d[0].([]interface{}))
	if err != nil {
		return nil, err
	}
	return &release.ArtifactsDownloadInput{
		DownloadInputs: downloadInputs,
	}, nil
}

func expandReleaseArtifactDownloadInputBases(d []interface{}) (*[]release.ArtifactDownloadInputBase, error) {
	artifactDownloadInputBases := make([]release.ArtifactDownloadInputBase, len(d))
	for i, data := range d {
		artifactDownloadInputBase, err := expandReleaseArtifactDownloadInputBase(data.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		artifactDownloadInputBases[i] = *artifactDownloadInputBase
	}
	return &artifactDownloadInputBases, nil
}

func expandReleaseArtifactDownloadInputBase(d map[string]interface{}) (*release.ArtifactDownloadInputBase, error) {
	dataArtifactItems := d["artifact_items"].([]interface{})
	artifactItems := make([]string, len(dataArtifactItems))
	for i, dataArtifactItem := range dataArtifactItems {
		artifactItems[i] = dataArtifactItem.(string)
	}

	return &release.ArtifactDownloadInputBase{
		Alias:                converter.String(d["alias"].(string)),
		ArtifactDownloadMode: converter.String(d["artifact_download_mode"].(string)),
		ArtifactItems:        &artifactItems,
		ArtifactType:         converter.String(d["artifact_type"].(string)),
	}, nil
}

func expandReleaseDefinitionDemandList(demands []interface{}) ([]interface{}, error) {
	demandsMap := make([]interface{}, len(demands))
	for i, demand := range demands {
		demandMap, err := expandReleaseDefinitionDemand(demand.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		demandsMap[i] = demandMap
	}
	return demandsMap, nil
}

func expandReleaseDefinitionDemand(d map[string]interface{}) (interface{}, error) {
	demand := ReleaseDefinitionDemand{
		Name:  converter.String(d["name"].(string)),
		Value: converter.String(d["value"].(string)),
	}
	return demand, nil
}

func expandReleaseWorkFlowTaskList(tasks []interface{}) ([]release.WorkflowTask, error) {
	tasksMap := make([]release.WorkflowTask, len(tasks))
	for i, approval := range tasks {
		releaseApproval, err := expandReleaseWorkFlowTask(approval.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		tasksMap[i] = *releaseApproval
	}
	return tasksMap, nil
}

func expandReleaseWorkFlowTask(d map[string]interface{}) (*release.WorkflowTask, error) {
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

func expandReleaseEnvironmentRetentionPolicy(d map[string]interface{}) (*release.EnvironmentRetentionPolicy, error) {
	environmentRetentionPolicy := release.EnvironmentRetentionPolicy{
		DaysToKeep:     converter.Int(d["days_to_keep"].(int)),
		RetainBuild:    converter.Bool(d["retain_build"].(bool)),
		ReleasesToKeep: converter.Int(d["releases_to_keep"].(int)),
	}
	return &environmentRetentionPolicy, nil
}

func expandReleaseDefinitionApprovals(d map[string]interface{}) (*release.ReleaseDefinitionApprovals, error) {
	approvals := d["approvals"].(*schema.Set).List()
	approvalsMap := make([]release.ReleaseDefinitionApprovalStep, len(approvals))

	for i, approval := range approvals {
		releaseApproval, err := expandReleaseDefinitionApprovalStep(approval.(map[string]interface{}))
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

func expandReleaseDefinitionApprovalStep(d map[string]interface{}) (*release.ReleaseDefinitionApprovalStep, error) {
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
