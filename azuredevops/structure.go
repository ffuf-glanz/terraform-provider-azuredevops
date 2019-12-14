package azuredevops

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/release"
	"github.com/microsoft/azure-devops-go-api/azuredevops/webapi"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/crud/distributedtask"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"strings"
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

type DeploymentHealthOptionType string
type deploymentHealthOptionValuesType struct {
	OneTargetAtATime ArtifactDownloadModeType
	Custom           ArtifactDownloadModeType
}

var DeploymentHealthOptionTypeValues = deploymentHealthOptionValuesType{
	OneTargetAtATime: "OneTargetAtATime",
	Custom:           "Custom",
}

type ReleaseDefinitionDemand struct {
	// Name of the demand.
	Name *string `json:"name,omitempty"`
	// The value of the demand.
	Value *string `json:"value,omitempty"`
}

type MachineGroupDeploymentMultiple struct {
}

type ReleaseHostedAzurePipelines struct {
	AgentSpecification *release.AgentSpecification
	QueueId            *int
}

func expandStringList(d []interface{}) []string {
	vs := make([]string, 0, len(d))
	for _, v := range d {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(string))
		}
	}
	return vs
}
func expandStringSet(d *schema.Set) []string {
	return expandStringList(d.List())
}

func expandIntList(d []interface{}) []int {
	vs := make([]int, 0, len(d))
	for _, v := range d {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(int))
		}
	}
	return vs
}
func expandIntSet(configured *schema.Set) []int {
	return expandIntList(configured.List())
}

func expandReleaseDefinitionsProperties(d map[string]interface{}) Properties {
	return Properties{
		DefinitionCreationSource: converter.String(d["definition_creation_source"].(string)),
		IntegrateJiraWorkItems:   converter.Bool(d["integrate_jira_work_items"].(bool)),
		IntegrateBoardsWorkItems: converter.Bool(d["integrate_boards_work_items"].(bool)),
	}
}
func expandReleaseDefinitionsPropertiesList(d []interface{}) []Properties {
	vs := make([]Properties, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			vs = append(vs, expandReleaseDefinitionsProperties(val))
		}
	}
	return vs
}
func expandReleaseDefinitionsPropertiesSet(d *schema.Set) interface{} {
	d2 := expandReleaseDefinitionsPropertiesList(d.List())
	if len(d2) != 1 {
		return nil
	}
	return d2[0]
}

func expandReleaseCondition(d map[string]interface{}) release.Condition {
	conditionType := release.ConditionType(d["condition_type"].(string))
	return release.Condition{
		ConditionType: &conditionType,
		Name:          converter.String(d["name"].(string)),
		Value:         converter.String(d["value"].(string)),
	}
}
func expandReleaseConditionList(d []interface{}) []release.Condition {
	vs := make([]release.Condition, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			vs = append(vs, expandReleaseCondition(val))
		}
	}
	return vs
}
func expandReleaseConditionSet(d *schema.Set) []release.Condition {
	return expandReleaseConditionList(d.List())
}

// WIP
func expandReleaseDefinitionEnvironment(d map[string]interface{}) release.ReleaseDefinitionEnvironment {
	variableGroups := expandIntList(d["variable_groups"].([]interface{}))
	deployStep := expandReleaseDefinitionDeployStepSet(d["deploy_step"].(*schema.Set))
	variables := expandReleaseConfigurationVariableValueSet(d["variable"].(*schema.Set))
	conditions := expandReleaseConditionSet(d["conditions"].(*schema.Set))
	demands := expandReleaseDefinitionDemandSet(d["demands"].(*schema.Set))
	environmentOptions := expandReleaseEnvironmentOptionsSet(d["environment_options"].(*schema.Set))
	retentionPolicy := expandReleaseEnvironmentRetentionPolicySet(d["retention_policy"].(*schema.Set))
	preDeployApprovals := expandReleaseDefinitionApprovalsSet(d["pre_deploy_approval"].(*schema.Set))
	postDeployApprovals := expandReleaseDefinitionApprovalsSet(d["post_deploy_approval"].(*schema.Set))
	properties := expandReleaseDefinitionsPropertiesSet(d["properties"].(*schema.Set))
	agentJobs := expandReleaseDeployPhaseSet(d["agent_job"].(*schema.Set), release.DeployPhaseTypesValues.AgentBasedDeployment)
	deploymentGroupJobs := expandReleaseDeployPhaseSet(d["deployment_group_job"].(*schema.Set), release.DeployPhaseTypesValues.MachineGroupBasedDeployment)
	agentlessJobs := expandReleaseDeployPhaseSet(d["agentless_job"].(*schema.Set), release.DeployPhaseTypesValues.RunOnServer)

	deployPhases := append(append(agentJobs, deploymentGroupJobs...), agentlessJobs...)

	return release.ReleaseDefinitionEnvironment{
		Conditions:         &conditions,
		Demands:            &demands,
		DeployPhases:       &deployPhases,
		DeployStep:         deployStep,
		EnvironmentOptions: environmentOptions,
		//EnvironmentTriggers: nil,
		//ExecutionPolicy:     nil,
		//Id:                  converter.Int(d["id"].(int)),
		Name:                converter.String(d["name"].(string)),
		PostDeployApprovals: postDeployApprovals,
		//PostDeploymentGates: nil,
		PreDeployApprovals: preDeployApprovals,
		//PreDeploymentGates:  nil,
		//ProcessParameters:   nil,
		Properties:      &properties,
		QueueId:         nil,
		Rank:            converter.Int(d["rank"].(int)),
		RetentionPolicy: retentionPolicy,
		//RunOptions:      nil,
		//Schedules:       nil,
		VariableGroups: &variableGroups,
		Variables:      &variables,
	}
}
func expandReleaseDefinitionEnvironmentList(d []interface{}) []release.ReleaseDefinitionEnvironment {
	vs := make([]release.ReleaseDefinitionEnvironment, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			vs = append(vs, expandReleaseDefinitionEnvironment(val))
		}
	}
	return vs
}
func expandReleaseDefinitionEnvironmentSet(configured *schema.Set) []release.ReleaseDefinitionEnvironment {
	return expandReleaseDefinitionEnvironmentList(configured.List())
}

func expandReleaseArtifact(d map[string]interface{}) release.Artifact {
	// NOTE : Might have to build a custom struct because AgentArtifactType doesn't equal comments in code. See below.
	// It can have value as 'Build', 'Jenkins', 'GitHub', 'Nuget', 'Team Build (external)', 'ExternalTFSBuild', 'Git', 'TFVC', 'ExternalTfsXamlBuild'.
	artifactType := release.AgentArtifactType(d["type"].(string))
	return release.Artifact{
		Alias:      converter.String(d["Alias"].(string)),
		IsPrimary:  converter.Bool(d["IsPrimary"].(bool)),
		IsRetained: converter.Bool(d["IsRetained"].(bool)),
		SourceId:   converter.String(d["SourceId"].(string)),
		Type:       converter.String(string(artifactType)),
	}
}
func expandReleaseArtifactList(d []interface{}) []release.Artifact {
	vs := make([]release.Artifact, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			vs = append(vs, expandReleaseArtifact(val))
		}
	}
	return vs
}
func expandReleaseArtifactSet(d *schema.Set) []release.Artifact {
	return expandReleaseArtifactList(d.List())
}

func expandReleaseConfigurationVariableValue(d map[string]interface{}) (string, release.ConfigurationVariableValue) {
	return d["name"].(string), release.ConfigurationVariableValue{
		AllowOverride: converter.Bool(d["allow_override"].(bool)),
		Value:         converter.String(d["value"].(string)),
		IsSecret:      converter.Bool(d["is_secret"].(bool)),
	}
}
func expandReleaseConfigurationVariableValueList(d []interface{}) map[string]release.ConfigurationVariableValue {
	vs := make(map[string]release.ConfigurationVariableValue)
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			key, d2 := expandReleaseConfigurationVariableValue(val)
			vs[key] = d2
		}
	}
	return vs
}
func expandReleaseConfigurationVariableValueSet(d *schema.Set) map[string]release.ConfigurationVariableValue {
	return expandReleaseConfigurationVariableValueList(d.List())
}

func expandReleaseDefinitionDeployStep(d map[string]interface{}) release.ReleaseDefinitionDeployStep {
	tasks := expandReleaseWorkFlowTaskSet(d["task"].(*schema.Set))
	return release.ReleaseDefinitionDeployStep{
		Id:    converter.Int(d["id"].(int)),
		Tasks: &tasks,
	}
}
func expandReleaseDefinitionDeployStepList(d []interface{}) []release.ReleaseDefinitionDeployStep {
	vs := make([]release.ReleaseDefinitionDeployStep, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			vs = append(vs, expandReleaseDefinitionDeployStep(val))
		}
	}
	return vs
}
func expandReleaseDefinitionDeployStepSet(d *schema.Set) *release.ReleaseDefinitionDeployStep {
	d2 := expandReleaseDefinitionDeployStepList(d.List())
	if len(d2) != 1 {
		return nil
	}
	return &d2[0]
}

func expandReleaseDeployPhase(d map[string]interface{}, t release.DeployPhaseTypes) ReleaseDeployPhase {
	workflowTasks := expandReleaseWorkFlowTaskList(d["task"].([]interface{}))
	var deploymentInput interface{}
	switch t {
	case release.DeployPhaseTypesValues.AgentBasedDeployment:
		deploymentInput = expandReleaseAgentDeploymentInput(d)
	case release.DeployPhaseTypesValues.MachineGroupBasedDeployment:
		deploymentInput = expandReleaseMachineGroupDeploymentInput(d)
	case release.DeployPhaseTypesValues.RunOnServer:
		deploymentInput = expandReleaseServerDeploymentInput(d)
	}
	return ReleaseDeployPhase{
		DeploymentInput: &deploymentInput,
		Rank:            converter.Int(d["rank"].(int)),
		PhaseType:       &t,
		Name:            converter.String(d["name"].(string)),
		//RefName:         converter.String(d["ref_name"].(string)),
		WorkflowTasks: &workflowTasks,
	}
}
func expandReleaseDeployPhaseList(d []interface{}, t release.DeployPhaseTypes) []interface{} {
	vs := make([]interface{}, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			vs = append(vs, expandReleaseDeployPhase(val, t))
		}
	}
	return vs
}
func expandReleaseDeployPhaseSet(d *schema.Set, t release.DeployPhaseTypes) []interface{} {
	return expandReleaseDeployPhaseList(d.List(), t)
}

func expandReleaseEnvironmentOptions(d map[string]interface{}) release.EnvironmentOptions {
	return release.EnvironmentOptions{
		AutoLinkWorkItems:            converter.Bool(d["auto_link_work_items"].(bool)),
		BadgeEnabled:                 converter.Bool(d["badge_enabled"].(bool)),
		EmailNotificationType:        converter.String(d["email_notification_type"].(string)),
		EmailRecipients:              converter.String(d["email_recipients"].(string)),
		EnableAccessToken:            converter.Bool(d["enable_access_token"].(bool)),
		PublishDeploymentStatus:      converter.Bool(d["publish_deployment_status"].(bool)),
		PullRequestDeploymentEnabled: converter.Bool(d["pull_request_deployment_enabled"].(bool)),
		SkipArtifactsDownload:        converter.Bool(d["skip_artifacts_download"].(bool)),
		TimeoutInMinutes:             converter.Int(d["timeout_in_minutes"].(int)),
	}
}
func expandReleaseEnvironmentOptionsList(d []interface{}) []release.EnvironmentOptions {
	vs := make([]release.EnvironmentOptions, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			vs = append(vs, expandReleaseEnvironmentOptions(val))
		}
	}
	return vs
}
func expandReleaseEnvironmentOptionsSet(d *schema.Set) *release.EnvironmentOptions {
	d2 := expandReleaseEnvironmentOptionsList(d.List())
	if len(d2) != 1 {
		return nil
	}
	return &d2[0]
}

func expandReleaseMachineGroupDeploymentInputMultiple(d map[string]interface{}) MachineGroupDeploymentMultiple {
	return MachineGroupDeploymentMultiple{}
}
func expandReleaseMachineGroupDeploymentInputMultipleList(d []interface{}) []MachineGroupDeploymentMultiple {
	vs := make([]MachineGroupDeploymentMultiple, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			vs = append(vs, expandReleaseMachineGroupDeploymentInputMultiple(val))
		}
	}
	return vs
}
func expandReleaseMachineGroupDeploymentInputMultipleSet(d *schema.Set) *MachineGroupDeploymentMultiple {
	d2 := expandReleaseMachineGroupDeploymentInputMultipleList(d.List())
	if len(d2) != 1 {
		return nil
	}
	return &d2[0]
}

func expandReleaseMultiConfigInput(d map[string]interface{}) release.MultiConfigInput {
	return release.MultiConfigInput{
		Multipliers:           converter.String(d["multipliers"].(string)),
		MaxNumberOfAgents:     converter.Int(d["number_of_agents"].(int)),
		ParallelExecutionType: &release.ParallelExecutionTypesValues.MultiConfiguration,
		ContinueOnError:       converter.Bool(d["continue_on_error"].(bool)),
	}
}
func expandReleaseMultiConfigInputList(d []interface{}) []release.MultiConfigInput {
	vs := make([]release.MultiConfigInput, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			vs = append(vs, expandReleaseMultiConfigInput(val))
		}
	}
	return vs
}
func expandReleaseMultiConfigInputSet(d *schema.Set) *release.MultiConfigInput {
	d2 := expandReleaseMultiConfigInputList(d.List())
	if len(d2) != 1 {
		return nil
	}
	return &d2[0]
}

func expandReleaseParallelExecutionInputBase(d map[string]interface{}) release.ParallelExecutionInputBase {
	return release.ParallelExecutionInputBase{
		MaxNumberOfAgents:     converter.Int(d["max_number_of_agents"].(int)),
		ParallelExecutionType: &release.ParallelExecutionTypesValues.MultiMachine,
		ContinueOnError:       converter.Bool(d["continue_on_error"].(bool)),
	}
}
func expandReleaseParallelExecutionInputBaseList(d []interface{}) []release.ParallelExecutionInputBase {
	vs := make([]release.ParallelExecutionInputBase, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			vs = append(vs, expandReleaseParallelExecutionInputBase(val))
		}
	}
	return vs
}
func expandReleaseParallelExecutionInputBaseSet(d *schema.Set) *release.ParallelExecutionInputBase {
	d2 := expandReleaseParallelExecutionInputBaseList(d.List())
	if len(d2) != 1 {
		return nil
	}
	return &d2[0]
}

func expandReleaseAgentSpecification(d map[string]interface{}) release.AgentSpecification {
	return release.AgentSpecification{
		Identifier: converter.String(d["agent_specification"].(string)),
	}
}
func expandReleaseAgentSpecificationList(d []interface{}) []release.AgentSpecification {
	vs := make([]release.AgentSpecification, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			vs = append(vs, expandReleaseAgentSpecification(val))
		}
	}
	return vs
}
func expandReleaseAgentSpecificationSet(d *schema.Set) *release.AgentSpecification {
	d2 := expandReleaseAgentSpecificationList(d.List())
	if len(d2) != 1 {
		return nil
	}
	return &d2[0]
}

func expandReleaseHostedAzurePipelines(d map[string]interface{}) ReleaseHostedAzurePipelines {
	agentSpecification := expandReleaseAgentSpecification(d)
	return ReleaseHostedAzurePipelines{
		AgentSpecification: &agentSpecification,
		QueueId:            converter.Int(d["agent_pool_id"].(int)),
	}
}
func expandReleaseHostedAzurePipelinesList(d []interface{}) []ReleaseHostedAzurePipelines {
	vs := make([]ReleaseHostedAzurePipelines, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			vs = append(vs, expandReleaseHostedAzurePipelines(val))
		}
	}
	return vs
}
func expandReleaseHostedAzurePipelinesSet(d *schema.Set) (*release.AgentSpecification, int) {
	d2 := expandReleaseHostedAzurePipelinesList(d.List())
	if len(d2) != 1 {
		return nil, 0
	}
	return d2[0].AgentSpecification, *d2[0].QueueId
}

func expandReleaseArtifactsDownloadInputSet(d *schema.Set) *release.ArtifactsDownloadInput {
	downloadInputs := expandReleaseArtifactDownloadInputBaseSet(d)
	return &release.ArtifactsDownloadInput{
		DownloadInputs: &downloadInputs,
	}
}

func expandReleaseMachineGroupDeploymentInput(d map[string]interface{}) *release.MachineGroupDeploymentInput {
	tags := expandStringList(d["tags"].([]interface{}))
	artifactsDownloadInput := expandReleaseArtifactsDownloadInputSet(d["artifact_download"].(*schema.Set))
	demands := expandReleaseDefinitionDemandSet(d["demand"].(*schema.Set))
	multiple := expandReleaseMachineGroupDeploymentInputMultipleSet(d["multiple"].(*schema.Set))
	deploymentHealthOption := DeploymentHealthOptionTypeValues.OneTargetAtATime
	if multiple != nil {
		deploymentHealthOption = DeploymentHealthOptionTypeValues.Custom
	}
	return &release.MachineGroupDeploymentInput{
		Condition:                 converter.String(d["condition"].(string)),
		JobCancelTimeoutInMinutes: converter.Int(d["max_execution_time_in_minutes"].(int)),
		OverrideInputs:            nil, // TODO : OverrideInputs
		TimeoutInMinutes:          converter.Int(d["timeout_in_minutes"].(int)),
		ArtifactsDownloadInput:    artifactsDownloadInput,
		Demands:                   &demands,
		EnableAccessToken:         converter.Bool(d["allow_scripts_to_access_oauth_token"].(bool)), // TODO : enable_access_token
		QueueId:                   converter.Int(d["deployment_group_id"].(int)),
		SkipArtifactsDownload:     converter.Bool(d["skip_artifacts_download"].(bool)),
		DeploymentHealthOption:    converter.String(string(deploymentHealthOption)),
		//HealthPercent:
		Tags: &tags,
	}
}
func expandReleaseAgentDeploymentInput(d map[string]interface{}) AgentDeploymentInput {
	artifactsDownloadInput := expandReleaseArtifactsDownloadInputSet(d["artifact_download"].(*schema.Set))
	demands := expandReleaseDefinitionDemandSet(d["demand"].(*schema.Set))
	agentPoolPrivate := expandReleaseAgentSpecificationSet(d["agent_pool_private"].(*schema.Set))

	agentPoolHostedAzurePipelines, queueId := expandReleaseHostedAzurePipelinesSet(d["agent_pool_hosted_azure_pipelines"].(*schema.Set))
	//if agentPoolPrivate != nil && agentPoolHostedAzurePipelines != nil { // TODO : how to solve
	//	return nil, fmt.Errorf("conflit %s and %s specify only one", "agent_pool_hosted_azure_pipelines", "agent_pool_private")
	//}
	var agentSpecification *release.AgentSpecification
	if agentPoolHostedAzurePipelines != nil {
		agentSpecification = agentPoolHostedAzurePipelines
	} else {
		agentSpecification = agentPoolPrivate
	}

	var parallelExecution interface{} = &release.ExecutionInput{
		ParallelExecutionType: &release.ParallelExecutionTypesValues.None,
	}
	multiConfiguration := expandReleaseMultiConfigInputSet(d["multi_configuration"].(*schema.Set))
	multiAgent := expandReleaseParallelExecutionInputBaseSet(d["multi_agent"].(*schema.Set))
	//if multiConfiguration != nil && multiAgent != nil { // TODO : how to solve
	//	return nil, fmt.Errorf("conflit %s and %s specify only one", "multi_configuration", "multi_agent")
	//}
	if multiConfiguration != nil {
		parallelExecution = multiConfiguration
	} else if multiAgent != nil {
		parallelExecution = multiAgent
	}

	return AgentDeploymentInput{
		Condition:                 converter.String(d["condition"].(string)),
		JobCancelTimeoutInMinutes: converter.Int(d["max_execution_time_in_minutes"].(int)),
		OverrideInputs:            nil, // TODO : OverrideInputs
		TimeoutInMinutes:          converter.Int(d["timeout_in_minutes"].(int)),
		ArtifactsDownloadInput:    artifactsDownloadInput,
		Demands:                   &demands,
		EnableAccessToken:         converter.Bool(d["allow_scripts_to_access_oauth_token"].(bool)), // TODO : enable_access_token
		QueueId:                   &queueId,
		SkipArtifactsDownload:     converter.Bool(d["skip_artifacts_download"].(bool)),
		AgentSpecification:        agentSpecification,
		// ImageId: ,
		ParallelExecution: &parallelExecution,
	}
}
func expandReleaseServerDeploymentInput(d map[string]interface{}) *ServerDeploymentInput {
	var parallelExecution interface{} = &release.ExecutionInput{
		ParallelExecutionType: &release.ParallelExecutionTypesValues.None,
	}
	multiConfiguration := expandReleaseMultiConfigInputSet(d["multi_configuration"].(*schema.Set))
	if multiConfiguration != nil {
		parallelExecution = multiConfiguration
	}
	return &ServerDeploymentInput{
		Condition:                 converter.String(d["condition"].(string)),
		JobCancelTimeoutInMinutes: converter.Int(d["max_execution_time_in_minutes"].(int)),
		OverrideInputs:            nil, // TODO : OverrideInputs
		TimeoutInMinutes:          converter.Int(d["timeout_in_minutes"].(int)),
		ParallelExecution:         &parallelExecution,
	}
}

func expandReleaseArtifactDownloadInputBase(d map[string]interface{}) release.ArtifactDownloadInputBase {
	artifactItems := expandStringList(d["items"].([]interface{}))
	return release.ArtifactDownloadInputBase{
		Alias:                converter.String(d["alias"].(string)),
		ArtifactDownloadMode: converter.String(d["artifact_download_mode"].(string)),
		ArtifactItems:        &artifactItems,
		ArtifactType:         converter.String(d["artifact_type"].(string)),
	}
}
func expandReleaseArtifactDownloadInputBaseList(d []interface{}) []release.ArtifactDownloadInputBase {
	vs := make([]release.ArtifactDownloadInputBase, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			vs = append(vs, expandReleaseArtifactDownloadInputBase(val))
		}
	}
	return vs
}
func expandReleaseArtifactDownloadInputBaseSet(d *schema.Set) []release.ArtifactDownloadInputBase {
	return expandReleaseArtifactDownloadInputBaseList(d.List())
}

func expandReleaseDefinitionDemand(d map[string]interface{}) ReleaseDefinitionDemand {
	name := d["name"].(string)
	configValue := d["value"].(string)
	if len(configValue) > 0 {
		name += " -equals " + configValue
	}
	return ReleaseDefinitionDemand{
		Name: converter.String(name),
	}
}
func expandReleaseDefinitionDemandList(d []interface{}) []interface{} {
	vs := make([]interface{}, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			vs = append(vs, expandReleaseDefinitionDemand(val))
		}
	}
	return vs
}
func expandReleaseDefinitionDemandSet(d *schema.Set) []interface{} {
	return expandReleaseDefinitionDemandList(d.List())
}

func expandReleaseWorkFlowTask(d map[string]interface{}) release.WorkflowTask {
	task := strings.Split(d["task"].(string), "@")
	taskName, version := task[0], task[1]
	configInputs := d["inputs"].(map[string]interface{})
	configId := d["id"].(string)
	taskId, err := uuid.Parse(configId)
	if err != nil {
		taskId = distributedtask.TasksMap[taskName]
	}

	inputs := make(map[string]string)
	for k, v := range configInputs {
		inputs[k] = v.(string)
	}

	return release.WorkflowTask{
		TaskId:          &taskId,
		AlwaysRun:       converter.Bool(d["always_run"].(bool)),
		Condition:       converter.String(d["condition"].(string)),
		ContinueOnError: converter.Bool(d["continue_on_error"].(bool)),
		DefinitionType:  converter.String("task"),
		Enabled:         converter.Bool(d["enabled"].(bool)),
		//Environment:      converter.String(d["environment"].(string)),
		//Inputs:           converter.Int(d["inputs"].(int)),
		Name: converter.String(d["display_name"].(string)),
		//OverrideInputs:   converter.Int(d["override_inputs"].(int)),
		TimeoutInMinutes: converter.Int(d["timeout_in_minutes"].(int)),
		Version:          &version,
		Inputs:           &inputs,
	}
}
func expandReleaseWorkFlowTaskList(d []interface{}) []release.WorkflowTask {
	vs := make([]release.WorkflowTask, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			vs = append(vs, expandReleaseWorkFlowTask(val))
		}
	}
	return vs
}
func expandReleaseWorkFlowTaskSet(d *schema.Set) []release.WorkflowTask {
	return expandReleaseWorkFlowTaskList(d.List())
}

func expandReleaseEnvironmentRetentionPolicy(d map[string]interface{}) release.EnvironmentRetentionPolicy {
	return release.EnvironmentRetentionPolicy{
		DaysToKeep: converter.Int(d["days_to_keep"].(int)),
		//RetainBuild:    converter.Bool(d["retain_build"].(bool)),
		ReleasesToKeep: converter.Int(d["releases_to_keep"].(int)),
	}
}
func expandReleaseEnvironmentRetentionPolicyList(d []interface{}) []release.EnvironmentRetentionPolicy {
	vs := make([]release.EnvironmentRetentionPolicy, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			vs = append(vs, expandReleaseEnvironmentRetentionPolicy(val))
		}
	}
	return vs
}
func expandReleaseEnvironmentRetentionPolicySet(d *schema.Set) *release.EnvironmentRetentionPolicy {
	d2 := expandReleaseEnvironmentRetentionPolicyList(d.List())
	if len(d2) != 1 {
		return nil
	}
	return &d2[0]
}

func expandReleaseDefinitionApprovalStep(d map[string]interface{}) release.ReleaseDefinitionApprovalStep {
	configurationApprover := d["approver"]
	var approver *webapi.IdentityRef
	if configurationApprover != nil {
		approver = &webapi.IdentityRef{
			Id: converter.String(d["approver_id"].(string)),
		}
	}
	return release.ReleaseDefinitionApprovalStep{
		//Id:               converter.Int(d["id"].(int)),
		Approver:    approver,
		IsAutomated: converter.Bool(d["is_automated"].(bool)),
		//IsNotificationOn: converter.Bool(d["is_notification_on"].(bool)),
		Rank: converter.Int(d["rank"].(int)),
	}
}
func expandReleaseDefinitionApprovalStepList(d []interface{}) []release.ReleaseDefinitionApprovalStep {
	vs := make([]release.ReleaseDefinitionApprovalStep, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			vs = append(vs, expandReleaseDefinitionApprovalStep(val))
		}
	}
	return vs
}
func expandReleaseDefinitionApprovalStepSet(d *schema.Set) []release.ReleaseDefinitionApprovalStep {
	return expandReleaseDefinitionApprovalStepList(d.List())
}

func expandReleaseApprovalOptions(d map[string]interface{}) release.ApprovalOptions {
	executionOrder := release.ApprovalExecutionOrder(d["execution_order"].(string))
	return release.ApprovalOptions{
		AutoTriggeredAndPreviousEnvironmentApprovedCanBeSkipped: converter.Bool(d["auto_triggered_and_previous_environment_approved_can_be_skipped"].(bool)),
		EnforceIdentityRevalidation:                             converter.Bool(d["enforce_identity_revalidation"].(bool)),
		ExecutionOrder:                                          &executionOrder,
		ReleaseCreatorCanBeApprover:                             converter.Bool(d["release_creator_can_be_approver"].(bool)),
		RequiredApproverCount:                                   converter.Int(d["required_approver_count"].(int)),
		TimeoutInMinutes:                                        converter.Int(d["timeout_in_minutes"].(int)),
	}
}
func expandReleaseApprovalOptionsList(d []interface{}) []release.ApprovalOptions {
	vs := make([]release.ApprovalOptions, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			vs = append(vs, expandReleaseApprovalOptions(val))
		}
	}
	return vs
}
func expandReleaseApprovalOptionsSet(d *schema.Set) *release.ApprovalOptions {
	d2 := expandReleaseApprovalOptionsList(d.List())
	if len(d2) != 1 {
		return nil
	}
	return &d2[0]
}

func expandReleaseDefinitionApprovals(d map[string]interface{}) release.ReleaseDefinitionApprovals {
	approvals := expandReleaseDefinitionApprovalStepSet(d["approval"].(*schema.Set))
	approvalOptions := expandReleaseApprovalOptionsSet(d["approval_options"].(*schema.Set))
	return release.ReleaseDefinitionApprovals{
		Approvals:       &approvals,
		ApprovalOptions: approvalOptions,
	}
}
func expandReleaseDefinitionApprovalsList(d []interface{}) []release.ReleaseDefinitionApprovals {
	vs := make([]release.ReleaseDefinitionApprovals, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			vs = append(vs, expandReleaseDefinitionApprovals(val))
		}
	}
	return vs
}
func expandReleaseDefinitionApprovalsSet(d *schema.Set) *release.ReleaseDefinitionApprovals {
	d2 := expandReleaseDefinitionApprovalsList(d.List())
	if len(d2) != 1 {
		return nil
	}
	return &d2[0]
}

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
