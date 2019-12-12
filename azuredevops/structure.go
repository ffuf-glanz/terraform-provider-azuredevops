package azuredevops

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/release"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
)

type AgentDeploymentInput struct {
	// Gets or sets the job condition.
	Condition *string `json:"condition,omitempty"`
	// Gets or sets the job cancel timeout in minutes for deployment which are cancelled by user for this release environment.
	JobCancelTimeoutInMinutes *int `json:"jobCancelTimeoutInMinutes,omitempty"`
	// Gets or sets the override inputs.
	OverrideInputs *map[string]string `json:"overrideInputs,omitempty"`
	// Gets or sets the job execution timeout in minutes for deployment which are queued against this release environment.
	TimeoutInMinutes *int `json:"timeoutInMinutes,omitempty"`
	// Artifacts that downloaded during job execution.
	ArtifactsDownloadInput *release.ArtifactsDownloadInput `json:"artifactsDownloadInput,omitempty"`
	// List demands that needs to meet to execute the job.
	Demands *[]interface{} `json:"demands,omitempty"`
	// Indicates whether to include access token in deployment job or not.
	EnableAccessToken *bool `json:"enableAccessToken,omitempty"`
	// Id of the pool on which job get executed.
	QueueId *int `json:"queueId,omitempty"`
	// Indicates whether artifacts downloaded while job execution or not.
	SkipArtifactsDownload *bool `json:"skipArtifactsDownload,omitempty"`
	// Specification for an agent on which a job gets executed.
	AgentSpecification *release.AgentSpecification `json:"agentSpecification,omitempty"`
	// Gets or sets the image ID.
	ImageId *int `json:"imageId,omitempty"`
	// Gets or sets the parallel execution input.
	ParallelExecution interface{} `json:"parallelExecution,omitempty"`
}

type ServerDeploymentInput struct {
	// Gets or sets the job condition.
	Condition *string `json:"condition,omitempty"`
	// Gets or sets the job cancel timeout in minutes for deployment which are cancelled by user for this release environment.
	JobCancelTimeoutInMinutes *int `json:"jobCancelTimeoutInMinutes,omitempty"`
	// Gets or sets the override inputs.
	OverrideInputs *map[string]string `json:"overrideInputs,omitempty"`
	// Gets or sets the job execution timeout in minutes for deployment which are queued against this release environment.
	TimeoutInMinutes *int `json:"timeoutInMinutes,omitempty"`
	// Gets or sets the parallel execution input.
	ParallelExecution interface{} `json:"parallelExecution,omitempty"`
}

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

	// TODO : Figure out if this is used for Deployment Job Group or all 3 Job Types.
	// Deployment jobs of the phase.
	//DeploymentJobs *[]release.DeploymentJob `json:"deploymentJobs,omitempty"`

	// TODO : Add manual_intervention {} (block) under agentless_job { } (block)
	// List of manual intervention tasks execution information in phase.
	ManualInterventions *[]release.ManualIntervention `json:"manualInterventions,omitempty"`

	// TODO : Remove below properties after a little R&D

	// TODO : Going to remove Id. As it is Deprecated.
	// Deprecated:
	//Id *int `json:"id,omitempty"`

	// TODO : Consider removing ID. As I do not believe you can change the value via the API.
	// TODO : If you can change via API then allow updating. Also explore if this cause a ForceNew/ForceReplace
	// ID of the phase.
	//PhaseId *string `json:"phaseId,omitempty"`

	// TODO : Consider removing RunPlanId. It is stateful data about the current state of the pipeline.
	// TODO : This does not seem like something one would want with terraform.
	// Run Plan ID of the phase.
	//RunPlanId *uuid.UUID `json:"runPlanId,omitempty"`

	// TODO : Consider removing StartedOn. It is stateful data about the current state of the pipeline.
	// TODO : This does not seem like something one would want with terraform.
	// Phase start time.
	//StartedOn *azuredevops.Time `json:"startedOn,omitempty"`

	// TODO : Consider removing Status. It is stateful data about the current state of the pipeline.
	// TODO : This does not seem like something one would want with terraform.
	// Status of the phase.
	//Status *release.DeployPhaseStatus `json:"status,omitempty"`

	// TODO : Consider removing ErrorLog. It is stateful data about the current state of the pipeline.
	// TODO : This does not seem like something one would want with terraform.
	// Phase execution error logs.
	//ErrorLog *string `json:"errorLog,omitempty"`
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
func expandStringSet(configured *schema.Set) []string {
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

func expandReleaseDefinitionEnvironmentList(d []interface{}) (*[]release.ReleaseDefinitionEnvironment, error) {
	vs := make([]release.ReleaseDefinitionEnvironment, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			v2, err := expandReleaseDefinitionEnvironment(val)
			if err != nil {
				return nil, err
			}
			vs = append(vs, *v2)
		}
	}
	return &vs, nil
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
	/*
		variableGroups := d["variable_groups"].([]interface{})
		variableGroupsMap := make([]int, len(variableGroups))
		for i, variableGroup := range variableGroups {
			variableGroupsMap[i] = variableGroup.(int)
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

		environmentOptions, environmentOptionsError := expandReleaseEnvironmentOptions(d["environment_options"].(*schema.Set).List())
		if environmentOptionsError != nil {
			return nil, environmentOptionsError
		}

	*/

	configurationRetentionPolicy := d["retention_policy"].(*schema.Set).List()
	retentionPolicy, err := expandReleaseEnvironmentRetentionPolicy(configurationRetentionPolicy[0].(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	configPreDeployApprovals := d["pre_deploy_approval"].(*schema.Set).List()
	preDeployApprovals, err2 := expandReleaseDefinitionApprovals(configPreDeployApprovals[0].(map[string]interface{}))
	if err2 != nil {
		return nil, err2
	}

	configPostDeployApprovals := d["post_deploy_approval"].(*schema.Set).List()
	postDeployApprovals, err3 := expandReleaseDefinitionApprovals(configPostDeployApprovals[0].(map[string]interface{}))
	if err3 != nil {
		return nil, err3
	}

	configAgentJobs := d["agent_job"].(*schema.Set).List()

	agentJobs, err4 := expandReleaseDeployPhaseList(configAgentJobs, release.DeployPhaseTypesValues.AgentBasedDeployment)
	if err4 != nil {
		return nil, err4
	}

	configDeployGroupJobs := d["deployment_group_job"].(*schema.Set).List()

	deploymentGroupJobs, err5 := expandReleaseDeployPhaseList(configDeployGroupJobs, release.DeployPhaseTypesValues.MachineGroupBasedDeployment)
	if err5 != nil {
		return nil, err5
	}

	configAgentlessJobs := d["agentless_job"].(*schema.Set).List()
	agentlessJobs, err6 := expandReleaseDeployPhaseList(configAgentlessJobs, release.DeployPhaseTypesValues.RunOnServer)
	if err6 != nil {
		return nil, err6
	}

	deployPhases := append(append(agentJobs, deploymentGroupJobs...), agentlessJobs...)

	releaseDefinitionEnvironment := release.ReleaseDefinitionEnvironment{
		//Conditions:          &conditionsMap,
		//Demands:             &demands,
		DeployPhases: &deployPhases,
		//DeployStep:          deployStepMap,
		//EnvironmentOptions:  environmentOptions,
		//EnvironmentTriggers: nil,
		//ExecutionPolicy:     nil,
		//Id:                  converter.Int(d["id"].(int)),
		Name: converter.String(d["name"].(string)),
		//Owner: &webapi.IdentityRef{
		//	Id:            converter.String(d["owner_id"].(string)),
		//	IsAadIdentity: converter.Bool(true),
		//	IsContainer:   converter.Bool(false),
		//},
		PostDeployApprovals: postDeployApprovals,
		//PostDeploymentGates: nil,
		PreDeployApprovals: preDeployApprovals,
		//PreDeploymentGates:  nil,
		//ProcessParameters:   nil,
		// Properties:          &releaseDefinitionEnvironmentProperties,
		QueueId:         nil,
		Rank:            converter.Int(d["rank"].(int)),
		RetentionPolicy: retentionPolicy,
		//RunOptions:      nil,
		//Schedules:       nil,
		//VariableGroups:  &variableGroupsMap,
		//Variables:       &variables,
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

func expandReleaseDeployPhaseList(d []interface{}, t release.DeployPhaseTypes) ([]interface{}, error) {
	vs := make([]interface{}, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			v2, err := expandReleaseDeployPhase(val, t)
			if err != nil {
				return nil, err
			}
			vs = append(vs, v2)
		}
	}
	return vs, nil
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

func expandReleaseDeployPhase(d map[string]interface{}, t release.DeployPhaseTypes) (interface{}, error) {
	var deploymentInput interface{}
	switch t {
	case release.DeployPhaseTypesValues.AgentBasedDeployment:
		{
			agentDeploymentInput, err := expandReleaseAgentDeploymentInput(d)
			if err != nil {
				return nil, err
			}
			deploymentInput = agentDeploymentInput
		}
	case release.DeployPhaseTypesValues.MachineGroupBasedDeployment:
		{
			machineGroupDeploymentInput, err := expandReleaseMachineGroupDeploymentInput(d)
			if err != nil {
				return nil, err
			}
			deploymentInput = machineGroupDeploymentInput
		}
	case release.DeployPhaseTypesValues.RunOnServer:
		{
			serverDeploymentInput, err := expandReleaseServerDeploymentInput(d)
			if err != nil {
				return nil, err
			}
			deploymentInput = serverDeploymentInput
		}
	}

	deployPhase := ReleaseDeployPhase{
		DeploymentInput: &deploymentInput, // release.AgentDeploymentInput || release.MachineGroupDeploymentInput{}
		Rank:            converter.Int(d["rank"].(int)),
		PhaseType:       &t,
		Name:            converter.String(d["name"].(string)),
		//RefName:         converter.String(d["ref_name"].(string)),
		//WorkflowTasks:   nil,
	}
	return deployPhase, nil
}

func expandReleaseMachineGroupDeploymentInput(d map[string]interface{}) (*release.MachineGroupDeploymentInput, error) {

	demands, demandsError := expandReleaseDefinitionDemandList(d["demand"].(*schema.Set).List())
	if demandsError != nil {
		return nil, demandsError
	}

	deploymentHealthOption := "OneTargetAtATime"
	configurationMultiple := d["multiple"].(*schema.Set).List()
	hasMultiple := len(configurationMultiple) > 0
	if hasMultiple {
		deploymentHealthOption = "Custom"
	}

	tags := expandStringList(d["tags"].([]interface{}))

	return &release.MachineGroupDeploymentInput{
		Condition:                 converter.String(d["condition"].(string)),
		JobCancelTimeoutInMinutes: converter.Int(d["max_execution_time_in_minutes"].(int)),
		OverrideInputs:            nil, // TODO : OverrideInputs
		TimeoutInMinutes:          converter.Int(d["timeout_in_minutes"].(int)),
		//ArtifactsDownloadInput:    artifactsDownloadInput,
		Demands:           &demands,
		EnableAccessToken: converter.Bool(d["allow_scripts_to_access_oauth_token"].(bool)), // TODO : enable_access_token
		QueueId:           converter.Int(d["deployment_group_id"].(int)),
		//SkipArtifactsDownload:     converter.Bool(d["skip_artifacts_download"].(bool)),
		DeploymentHealthOption: converter.String(deploymentHealthOption),
		//HealthPercent:
		Tags: &tags,
	}, nil
}

func expandReleaseAgentDeploymentInput(d map[string]interface{}) (*AgentDeploymentInput, error) {
	agentSpecification := &release.AgentSpecification{}

	//artifactsDownloadInput, err := expandReleaseArtifactsDownloadInput(d["artifacts_download_input"].(*schema.Set).List())
	//if err != nil {
	//	return nil, err
	//}
	//
	demands, demandsError := expandReleaseDefinitionDemandList(d["demand"].(*schema.Set).List())
	if demandsError != nil {
		return nil, demandsError
	}

	queueId := converter.Int(0)
	configAgentPoolHostedAzurePipelines := d["agent_pool_hosted_azure_pipelines"].(*schema.Set).List()
	configAgentPoolPrivate := d["agent_pool_private"].(*schema.Set).List()

	hasHosted, hasPrivate := len(configAgentPoolHostedAzurePipelines) > 0, len(configAgentPoolPrivate) > 0
	if hasHosted && hasPrivate {
		return nil, fmt.Errorf("conflit %s and %s specify only one", "agent_pool_hosted_azure_pipelines", "agent_pool_private")
	}

	if hasHosted {
		d2 := configAgentPoolHostedAzurePipelines[0].(map[string]interface{})
		agentSpecification.Identifier = converter.String(d2["agent_specification"].(string))
		queueId = converter.Int(d2["agent_pool_id"].(int))
	}

	var parallelExecution interface{}
	configurationMultiConfiguration := d["multi_configuration"].(*schema.Set).List()
	configurationMultiAgent := d["multi_agent"].(*schema.Set).List()

	hasMultiConfig, hasMultiAgent := len(configurationMultiConfiguration) > 0, len(configurationMultiAgent) > 0
	if hasMultiConfig && hasMultiAgent {
		return nil, fmt.Errorf("conflit %s and %s specify only one", "multi_configuration", "multi_agent")
	} else if hasMultiConfig {
		d2 := configurationMultiConfiguration[0].(map[string]interface{})
		parallelExecution = &release.MultiConfigInput{
			Multipliers:           converter.String(d2["multipliers"].(string)),
			MaxNumberOfAgents:     converter.Int(d2["number_of_agents"].(int)),
			ParallelExecutionType: &release.ParallelExecutionTypesValues.MultiConfiguration,
			ContinueOnError:       converter.Bool(d2["continue_on_error"].(bool)),
		}
	} else if hasMultiAgent {
		d2 := configurationMultiAgent[0].(map[string]interface{})
		parallelExecution = &release.ParallelExecutionInputBase{
			MaxNumberOfAgents:     converter.Int(d2["max_number_of_agents"].(int)),
			ParallelExecutionType: &release.ParallelExecutionTypesValues.MultiMachine,
			ContinueOnError:       converter.Bool(d2["continue_on_error"].(bool)),
		}
	} else {
		parallelExecution = &release.ExecutionInput{
			ParallelExecutionType: &release.ParallelExecutionTypesValues.None,
		}
	}

	return &AgentDeploymentInput{
		Condition:                 converter.String(d["condition"].(string)),
		JobCancelTimeoutInMinutes: converter.Int(d["max_execution_time_in_minutes"].(int)),
		OverrideInputs:            nil, // TODO : OverrideInputs
		TimeoutInMinutes:          converter.Int(d["timeout_in_minutes"].(int)),
		//ArtifactsDownloadInput:    artifactsDownloadInput,
		Demands: &demands,
		//EnableAccessToken:         converter.Bool(d["enable_access_token"].(bool)), // TODO : enable_access_token
		QueueId: queueId,
		//SkipArtifactsDownload:     converter.Bool(d["skip_artifacts_download"].(bool)),
		AgentSpecification: agentSpecification,
		// ImageId: ,
		ParallelExecution: &parallelExecution,
	}, nil
}

func expandReleaseServerDeploymentInput(d map[string]interface{}) (*ServerDeploymentInput, error) {

	var parallelExecution interface{}
	configurationMultiConfiguration := d["multi_configuration"].(*schema.Set).List()

	hasMultiConfig := len(configurationMultiConfiguration) > 0
	if hasMultiConfig {
		d2 := configurationMultiConfiguration[0].(map[string]interface{})
		parallelExecution = &release.MultiConfigInput{
			Multipliers:           converter.String(d2["multipliers"].(string)),
			ParallelExecutionType: &release.ParallelExecutionTypesValues.MultiConfiguration,
			ContinueOnError:       converter.Bool(d2["continue_on_error"].(bool)),
		}
	} else {
		parallelExecution = &release.ExecutionInput{
			ParallelExecutionType: &release.ParallelExecutionTypesValues.None,
		}
	}

	return &ServerDeploymentInput{
		Condition:                 converter.String(d["condition"].(string)),
		JobCancelTimeoutInMinutes: converter.Int(d["max_execution_time_in_minutes"].(int)),
		OverrideInputs:            nil, // TODO : OverrideInputs
		TimeoutInMinutes:          converter.Int(d["timeout_in_minutes"].(int)),
		ParallelExecution:         &parallelExecution,
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
	name := d["name"].(string)
	configValue := d["value"].(string)
	if len(configValue) > 0 {
		name += " -equals " + configValue
	}
	demand := ReleaseDefinitionDemand{
		Name: converter.String(name),
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
	return &release.EnvironmentRetentionPolicy{
		DaysToKeep: converter.Int(d["days_to_keep"].(int)),
		//RetainBuild:    converter.Bool(d["retain_build"].(bool)),
		ReleasesToKeep: converter.Int(d["releases_to_keep"].(int)),
	}, nil
}

func expandReleaseDefinitionApprovals(d map[string]interface{}) (*release.ReleaseDefinitionApprovals, error) {
	configurationApprovals := d["approval"].(*schema.Set).List()
	approvals, err := expandReleaseDefinitionApprovalStepList(configurationApprovals)
	if err != nil {
		return nil, err
	}

	//approvalOptions := d["approval_options"].(*schema.Set).List()
	//if len(approvalOptions) != 1 {
	//	return nil, fmt.Errorf("unexpectedly did not find a approval options in the approvals data")
	//}
	//d2 := approvalOptions[0].(map[string]interface{})
	//executionOrder := release.ApprovalExecutionOrder(d2["execution_order"].(string))

	releaseDefinitionApprovals := release.ReleaseDefinitionApprovals{
		Approvals: approvals,
		//ApprovalOptions: &release.ApprovalOptions{
		//	AutoTriggeredAndPreviousEnvironmentApprovedCanBeSkipped: converter.Bool(d2["auto_triggered_and_previous_environment_approved_can_be_skipped"].(bool)),
		//	EnforceIdentityRevalidation:                             converter.Bool(d2["enforce_identity_revalidation"].(bool)),
		//	ExecutionOrder:                                          &executionOrder,
		//	ReleaseCreatorCanBeApprover:                             converter.Bool(d2["release_creator_can_be_approver"].(bool)),
		//	RequiredApproverCount:                                   converter.Int(d2["required_approver_count"].(int)),
		//	TimeoutInMinutes:                                        converter.Int(d2["timeout_in_minutes"].(int)),
		//},
	}
	return &releaseDefinitionApprovals, nil
}

func expandReleaseDefinitionApprovalStepList(d []interface{}) (*[]release.ReleaseDefinitionApprovalStep, error) {
	vs := make([]release.ReleaseDefinitionApprovalStep, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			v2, err := expandReleaseDefinitionApprovalStep(val)
			if err != nil {
				return nil, err
			}
			vs = append(vs, v2)
		}
	}
	return &vs, nil
}

func expandReleaseDefinitionApprovalStep(d map[string]interface{}) (release.ReleaseDefinitionApprovalStep, error) {
	//configurationApprover := d["approver"]
	//var approverMap *webapi.IdentityRef
	//if configurationApprover != nil {
	//	approverMap = &webapi.IdentityRef{
	//		Id: converter.String(d["approver_id"].(string)),
	//	}
	//}

	return release.ReleaseDefinitionApprovalStep{
		//Id:               converter.Int(d["id"].(int)),
		//Approver:         approver,
		IsAutomated: converter.Bool(d["is_automated"].(bool)),
		//IsNotificationOn: converter.Bool(d["is_notification_on"].(bool)),
		Rank: converter.Int(d["rank"].(int)),
	}, nil
}
