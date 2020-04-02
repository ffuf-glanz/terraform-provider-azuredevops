package azuredevops

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/release"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/tfhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/validate"
	"strconv"
)

// TagsSchema list of tags
var TagsSchema = &schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
}

func resourceReleaseDefinition() *schema.Resource {

	rank := &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Default:  1,
	}

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

	// TODO : import these from YAML
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

	demand := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"value": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	demands := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: demand,
		},
	}

	artifactItems := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}

	artifactDownloadInputBase := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"alias": {
					Type:     schema.TypeString,
					Required: true,
				},
				"artifact_download_mode": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(ArtifactDownloadModeTypeValues.All),
						string(ArtifactDownloadModeTypeValues.Selective),
						string(ArtifactDownloadModeTypeValues.Skip),
					}, false),
				},
				"artifact_items": artifactItems,
				"artifact_type": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}

	artifactsDownloadInput := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"artifact_download_input_base": artifactDownloadInputBase,
			},
		},
	}

	overrideInputs := &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	task := map[string]*schema.Schema{
		"always_run": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"condition": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "succeeded()",
		},
		"continue_on_error": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		// TODO : environment
		// "environment": environment,
		// TODO : inputs
		//"inputs": {
		//	Type:     schema.TypeString,
		//	Required: true,
		//},
		"display_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"override_input": overrideInputs,
		// TODO : Remove ref_name
		//"ref_name": {
		//	Type:     schema.TypeString,
		//	Required: true,
		//},
		// TODO : Remove task_id is going to be derived from task
		//"task_id": {
		//	Type:     schema.TypeString,
		//	Required: true,
		//},
		// TODO : Remove task_id is going to be derived from task
		//"version": {
		//	Type:     schema.TypeString,
		//	Required: true,
		//},
		"timeout_in_minutes": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
		"inputs": {
			Type:     schema.TypeMap,
			Optional: true,
		},
		"task": {
			Type:     schema.TypeString,
			Required: true,
			// ValidateFunc: // TODO check for pattern name@version /[a-zA-Z]+\@\d+/
		},
		"id": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validate.UUID,
		},
	}

	tasks := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: task,
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
				"task": tasks,
			},
		},
	}

	approvalOptions := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
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
		"task": tasks,
	}

	releaseDefinitionGates := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: releaseDefinitionGate,
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

	skipArtifactsDownload := &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
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
				"email_notification_type": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "OnlyOnFailure",
				},
				"email_recipients": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "release.environment.owner;release.creator",
				},
				"enable_access_token": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"publish_deployment_status": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
				},
				"pull_request_deployment_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"skip_artifacts_download": skipArtifactsDownload,
				"timeout_in_minutes": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  0,
				},
			},
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
			Optional: true,
			Default:  "",
		},
	}

	conditions := &schema.Schema{
		Type:     schema.TypeSet,
		MinItems: 1,
		MaxItems: 1,
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

	releaseDefinitionProperties := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"definition_creation_source": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "ReleaseNew",
				},
				"integrate_jira_work_items": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"integrate_boards_work_items": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}

	releaseDefinitionEnvironmentProperties := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.IntAtLeast(1),
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

	artifact := map[string]*schema.Schema{
		"alias": {
			Type:     schema.TypeString,
			Optional: true,
		},
		// TODO : definition_reference
		//"definition_reference": {
		//	Type:     schema.TypeInt,
		//	Optional: true,
		//},
		"is_primary": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"is_retained": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"type": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{ // NOTE : May need to use custom enum
				string(release.AgentArtifactTypeValues.GitHub),
				string(release.AgentArtifactTypeValues.Tfvc),
				string(release.AgentArtifactTypeValues.Build),
				string(release.AgentArtifactTypeValues.Custom),
				string(release.AgentArtifactTypeValues.ExternalTfsBuild),
				string(release.AgentArtifactTypeValues.FileShare),
				string(release.AgentArtifactTypeValues.Jenkins),
				string(release.AgentArtifactTypeValues.Nuget),
				string(release.AgentArtifactTypeValues.TfGit),
				string(release.AgentArtifactTypeValues.TfsOnPrem),
				string(release.AgentArtifactTypeValues.XamlBuild),
			}, false),
		},
	}

	artifacts := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: artifact,
		},
	}

	approval := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				// TODO : id
				//"id": {
				//	Type:     schema.TypeInt,
				//	Optional: true,
				//	Default:  0,
				//},
				// TODO : approver_id
				//"approver_id": {
				//	Type:         schema.TypeString,
				//	Optional:     true,
				//	ValidateFunc: validate.UUID,
				//},
				"rank": rank,
				"is_automated": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
				},
				// TODO : is_notification_on
				//"is_notification_on": {
				//	Type:     schema.TypeBool,
				//	Optional: true,
				//	Default:  false,
				//},
			},
		},
	}

	releaseDefinitionApproval := &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"approval":         approval,
				"approval_options": approvalOptions,
			},
		},
	}

	retentionPolicy := &schema.Schema{
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

	allowScriptsToAccessOauthToken := &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}

	artifactDownload := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				// TODO : create this :)
			},
		},
	}

	agentJob := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"override_input": overrideInputs,
				"demand":         demands,
				"rank":           rank,
				"agent_pool_hosted_azure_pipelines": {
					Type:     schema.TypeSet,
					Optional: true,
					MinItems: 1,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"agent_pool_id": {
								Type:     schema.TypeInt,
								Required: true,
							},
							"agent_specification": {
								Type:     schema.TypeString,
								Required: true,
								ValidateFunc: validation.StringInSlice([]string{
									"macOS-10.13",
									"macOS-10.14",
									"ubuntu-16.04",
									"ubuntu-18.04",
									"vs2015-win2012r2",
									"vs2017-win2016",
									"win1803",
									"windows-2019",
								}, false),
							},
						},
					},
				},
				"artifact_download": artifactDownload,
				"agent_pool_private": {
					Type:     schema.TypeSet,
					Optional: true,
					MinItems: 1,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"agent_pool_id": {
								Type:     schema.TypeString,
								Required: true,
							},
						},
					},
				},
				"timeout_in_minutes": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  0,
				},
				"max_execution_time_in_minutes": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  1,
				},
				"condition": {
					Type:     schema.TypeString,
					Required: true,
				},
				"multi_configuration": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"multipliers": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "A list of comma separated configuration variables to use. These are defined on the Variables tab. For example, OperatingSystem, Browser will run the tasks for both variables.",
							},
							"number_of_agents": {
								Type:         schema.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntAtLeast(1),
							},
							"continue_on_error": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  false,
							},
						},
					},
				},
				"multi_agent": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"max_number_of_agents": {
								Type:         schema.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntAtLeast(1),
							},
							"continue_on_error": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  false,
							},
						},
					},
				},
				"skip_artifacts_download":             skipArtifactsDownload,
				"allow_scripts_to_access_oauth_token": allowScriptsToAccessOauthToken,
				"task":                                tasks,
			},
		},
	}

	deploymentGroupJob := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"override_input": overrideInputs,
				"demand":         demands,
				"rank":           rank,
				"deployment_group_id": {
					Type:     schema.TypeInt,
					Required: true,
				},
				"tags": TagsSchema,
				"multiple": {
					Type:     schema.TypeSet,
					Optional: true,
					MinItems: 1,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"max_targets_in_parallel": {
								Type:     schema.TypeInt,
								Required: true,
							},
						},
					},
				},
				"allow_scripts_to_access_oauth_token": allowScriptsToAccessOauthToken,
				"timeout_in_minutes": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  0,
				},
				"max_execution_time_in_minutes": {
					Type:         schema.TypeInt,
					Optional:     true,
					Default:      1,
					ValidateFunc: validation.IntAtLeast(1),
				},
				"condition": {
					Type:     schema.TypeString,
					Required: true,
				},
				"task":                    tasks,
				"skip_artifacts_download": skipArtifactsDownload,
			},
		},
	}

	agentlessJob := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"override_input": overrideInputs,
				"rank":           rank,
				"timeout_in_minutes": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  0,
				},
				"max_execution_time_in_minutes": {
					Type:         schema.TypeInt,
					Optional:     true,
					Default:      1,
					ValidateFunc: validation.IntAtLeast(1),
				},
				"condition": {
					Type:     schema.TypeString,
					Required: true,
				},
				"multi_configuration": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"multipliers": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "A list of comma separated configuration variables to use. These are defined on the Variables tab. For example, OperatingSystem, Browser will run the tasks for both variables.",
							},
							"continue_on_error": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  false,
							},
						},
					},
				},
				"task":                    tasks,
				"skip_artifacts_download": skipArtifactsDownload,
			},
		},
	}

	stage := &schema.Schema{
		// TODO : can this be a TypeList and not require the user to supply rank?
		Type:     schema.TypeSet,
		Required: true,
		MinItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  -1,
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
				"variable":             configurationVariables,
				"variable_groups":      variableGroups,
				"pre_deploy_approval":  releaseDefinitionApproval,
				"post_deploy_approval": releaseDefinitionApproval,
				"deploy_step":          releaseDefinitionDeployStep,
				"agent_job":            agentJob,
				"deployment_group_job": deploymentGroupJob,
				"agentless_job":        agentlessJob,
				"retention_policy":     retentionPolicy,

				// TODO : This is missing from the docs
				// "runOptions": runOptions,
				// TODO : Rename this?
				"artifacts_download_input": artifactsDownloadInput,
				"environment_options":      environmentOptions,
				"demands": &schema.Schema{
					Type:       schema.TypeSet,
					Optional:   true,
					Deprecated: "Use DeploymentInput.Demands instead",
					Elem: &schema.Resource{
						Schema: demand,
					},
				},
				"conditions":            conditions,
				"execution_policy":      environmentExecutionPolicy,
				"schedules":             schedules,
				"properties":            releaseDefinitionEnvironmentProperties,
				"pre_deployment_gates":  releaseDefinitionGatesStep,
				"post_deployment_gates": releaseDefinitionGatesStep,
				"environment_triggers":  environmentTriggers,
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
				ValidateFunc: validate.Path,
			},
			"variable_groups": variableGroups,
			"source": {
				Type:     schema.TypeString,
				Computed: true,
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
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"tags":       TagsSchema,
			"properties": releaseDefinitionProperties,
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by terraform",
			},
			"created_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stage":     stage,
			"artifacts": artifacts,
		},
	}
}

func flattenReleaseDefinition(d *schema.ResourceData, releaseDefinition *release.ReleaseDefinition, projectID string) {
	d.SetId(strconv.Itoa(*releaseDefinition.Id))

	d.Set("project_id", projectID)
	// TODO : if any of the props below are optional, then the code could use converter.ToX
	d.Set("name", *releaseDefinition.Name)
	d.Set("path", *releaseDefinition.Path)
	d.Set("variable_groups", *releaseDefinition.VariableGroups)
	d.Set("source", *releaseDefinition.Source)
	d.Set("description", converter.ToString(releaseDefinition.Description, ""))
	d.Set("variable", flattenReleaseDefinitionVariables(releaseDefinition.Variables))
	d.Set("release_name_format", *releaseDefinition.ReleaseNameFormat)
	d.Set("url", *releaseDefinition.Url)
	d.Set("is_deleted", *releaseDefinition.IsDeleted)
	d.Set("tags", *releaseDefinition.Tags)
	d.Set("properties", flattenReleaseDefinitionProperties(releaseDefinition.Properties))
	d.Set("comment", *releaseDefinition.Comment)
	d.Set("created_on", *releaseDefinition.CreatedOn)
	d.Set("modified_on", *releaseDefinition.ModifiedOn)

	d.Set("stage", flattenReleaseDefinitionEnvironmentList(releaseDefinition.Environments))

	// TODO : build flattening for 3 items below.
	// d.Set("artifacts", flattenReleaseDefinitionArtifactsList(releaseDefinition.Environments))
	// d.Set("triggers", flattenReleaseDefinitionTriggersList(releaseDefinition.Environments))

	//d2 := d.Get("stage").(*schema.Set)
	//d4 := d.Get("name")
	//
	//d3, _ := d2.List()[0].(map[string]interface{})
	//fmt.Println(d3, d4)

	revision := 0
	if releaseDefinition.Revision != nil {
		revision = *releaseDefinition.Revision
	}
	d.Set("revision", revision)
}

func resourceReleaseDefinitionCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	releaseDefinition, projectID, err := expandReleaseDefinition(d)
	if err != nil {
		return fmt.Errorf("error creating resource Release Definition: %+v", err)
	}

	createdReleaseDefinition, err := createReleaseDefinition(clients, releaseDefinition, projectID)
	if err != nil {
		return fmt.Errorf("error creating resource Release Definition: %+v", err)
	}

	flattenReleaseDefinition(d, createdReleaseDefinition, projectID)
	return nil
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

func createReleaseDefinition(clients *config.AggregatedClient, releaseDefinition *release.ReleaseDefinition, project string) (*release.ReleaseDefinition, error) {
	createdBuild, err := clients.ReleaseClient.CreateReleaseDefinition(clients.Ctx, release.CreateReleaseDefinitionArgs{
		ReleaseDefinition: releaseDefinition,
		Project:           &project,
	})

	return createdBuild, err
}

func expandReleaseDefinition(d *schema.ResourceData) (*release.ReleaseDefinition, string, error) {
	projectID := d.Get("project_id").(string)

	// Look for the ID. This may not exist if we are within the context of a "create" operation,
	// so it is OK if it is missing.
	releaseDefinitionID, err := strconv.Atoi(d.Id())
	var releaseDefinitionReference *int = nil
	if err == nil {
		releaseDefinitionReference = &releaseDefinitionID
	}

	variableGroups := expandIntList(d.Get("variable_groups").([]interface{}))
	environments := expandReleaseDefinitionEnvironmentSet(d.Get("stage").(*schema.Set))
	variables := expandReleaseConfigurationVariableValueSet(d.Get("variable").(*schema.Set))
	properties := expandReleaseDefinitionsPropertiesSet(d.Get("properties").(*schema.Set))
	//artifacts := expandReleaseArtifactSet(d.Get("artifact").(*schema.Set))
	tags := expandStringList(d.Get("tags").([]interface{}))

	releaseDefinition := release.ReleaseDefinition{
		Id:                releaseDefinitionReference,
		Name:              converter.String(d.Get("name").(string)),
		Path:              converter.String(d.Get("path").(string)),
		Revision:          converter.Int(d.Get("revision").(int)),
		Description:       converter.String(d.Get("description").(string)),
		Environments:      &environments,
		Variables:         &variables,
		ReleaseNameFormat: converter.String(d.Get("release_name_format").(string)),
		VariableGroups:    &variableGroups,
		Properties:        &properties,
		//Artifacts:         &artifacts,
		Comment: converter.String(d.Get("comment").(string)),
		Tags:    &tags,
	}

	data, err := json.Marshal(releaseDefinition)
	fmt.Println(string(data))

	return &releaseDefinition, projectID, nil
}
