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

func resourceReleaseDefinition() *schema.Resource {
	/*
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
			"override_inputs": overrideInputs,
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


		agentDeploymentInput := &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"condition": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"job_cancel_timeout_in_minutes": {
						Type:     schema.TypeInt,
						Optional: true,
						Default:  1,
					},
					"override_inputs": overrideInputs,
					"timeout_in_minutes": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"artifacts_download_input": artifactsDownloadInput,
					"demands":                  demands,
					"enable_access_token": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},
					"queue_id": {
						Type:     schema.TypeInt,
						Required: true,
					},
					"skip_artifacts_download": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},
					"agent_specification_identifier": {
						Type:     schema.TypeString,
						Required: true,
					},
					"image_id": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"parallel_execution_type": {
						Type:     schema.TypeString,
						Optional: true,
						Default:  release.ParallelExecutionTypesValues.None,
						ValidateFunc: validation.StringInSlice([]string{
							string(release.ParallelExecutionTypesValues.None),
							string(release.ParallelExecutionTypesValues.MultiConfiguration),
							string(release.ParallelExecutionTypesValues.MultiMachine),
						}, false),
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

			"agent_deployment_input": agentDeploymentInput,
			// TODO : GatesDeployPhase, MachineGroupBasedDeployPhase, RunOnServerDeployPhase
			// TODO : How to do Validation based on phase_type?
			//"gates_deployment_input" : gatesDeploymentInput,
			//"machine_group_deployment_input" : machineGroupDeploymentInput,
			//"run_on_server_deploy_phase" : runOnServerDeployPhase,
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
			Required: true,
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
					"skip_artifacts_download": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},
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
			Required: true,
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
			Required: true,
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
			Type:     schema.TypeList,
			Required: true,
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
	*/
	rank := &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Default:  1,
	}

	approval := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				//"id": {
				//	Type:     schema.TypeInt,
				//	Optional: true,
				//	Default:  0,
				//},
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
				"approval": approval,
				//"approval_options": approvalOptions,
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
				//"retain_build": {
				//	Type:     schema.TypeBool,
				//	Optional: true,
				//	Default:  true,
				//},
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
				"rank": rank,
				"hosted_azure_pipelines_agent_pool": {
					Type:     schema.TypeSet,
					Optional: true,
					MinItems: 1,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
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
					ConflictsWith: []string{"private_agent_pool"},
				},
				//"private_agent_pool": {
				//	Type:     schema.TypeSet,
				//	Optional: true,
				//	MinItems: 1,
				//	MaxItems: 1,
				//	Elem: &schema.Resource{
				//		Schema: map[string]*schema.Schema{
				//			"agent_pool_id": {
				//				Type:     schema.TypeString,
				//				Required: true,
				//			},
				//		},
				//	},
				//	ConflictsWith: []string{"hosted_azure_pipelines_agent_pool"},
				//},
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
			},
		},
	}

	stage := &schema.Schema{
		// TODO: can this be a TypeList and not require the user to supply rank?
		Type:     schema.TypeSet,
		Required: true,
		MinItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				//"id": {
				//	Type:     schema.TypeInt,
				//	Optional: true,
				//	Default:  -1,
				//},
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"rank": rank,
				// TODO : Is this something you would want to set
				//"owner_id": {
				//	Type:         schema.TypeString,
				//	Optional:     true,
				//	ValidateFunc: validate.UUID,
				//},
				//"variable":              configurationVariables,
				//"variable_groups":       variableGroups,
				"pre_deploy_approvals":  releaseDefinitionApproval,
				"post_deploy_approvals": releaseDefinitionApproval,
				//"deploy_step":           releaseDefinitionDeployStep,
				//"deploy_phases":         deployPhases,
				"agent_job":        agentJob,
				"retention_policy": retentionPolicy,

				// TODO : This is missing from the docs
				// "runOptions": runOptions
				//"environment_options": environmentOptions,
				//"demands": &schema.Schema{
				//	Type:       schema.TypeSet,
				//	Optional:   true,
				//	Deprecated: "Use DeploymentInput.Demands instead",
				//	Elem: &schema.Resource{
				//		Schema: demand,
				//	},
				//},
				//"conditions":            conditions,
				//"execution_policy":      environmentExecutionPolicy,
				//"schedules":             schedules,
				//"properties":            releaseDefinitionEnvironmentProperties,
				//"pre_deployment_gates":  releaseDefinitionGatesStep,
				//"post_deployment_gates": releaseDefinitionGatesStep,
				//"environment_triggers":  environmentTriggers,
				//"badge_url": {
				//	Type:     schema.TypeString,
				//	Computed: true,
				//},
			},
		},
	}

	return &schema.Resource{
		Create: resourceReleaseDefinitionCreate, // expand
		Read:   resourceReleaseDefinitionRead,   // flatten
		Update: resourceReleaseDefinitionUpdate, // expand
		Delete: resourceReleaseDefinitionDelete, // flatten

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			//"revision": {
			//	Type:     schema.TypeInt,
			//	Computed: true,
			//},
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
			//"variable_groups": variableGroups,
			//"source": {
			//	Type:     schema.TypeString,
			//	Optional: true,
			//	ValidateFunc: validation.StringInSlice([]string{
			//		string(release.ReleaseDefinitionSourceValues.Undefined),
			//		string(release.ReleaseDefinitionSourceValues.RestApi),
			//		string(release.ReleaseDefinitionSourceValues.PortalExtensionApi),
			//		string(release.ReleaseDefinitionSourceValues.Ibiza),
			//		string(release.ReleaseDefinitionSourceValues.UserInterface),
			//	}, false),
			//},
			//"description": {
			//	Type:     schema.TypeString,
			//	Optional: true,
			//	Default:  "",
			//},
			//"variable": configurationVariables,
			//"release_name_format": {
			//	Type:     schema.TypeString,
			//	Optional: true,
			//	Default:  "Release-$(rev:r)",
			//},
			"stage": stage,

			//"url": {
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},
			//"is_deleted": {
			//	Type:     schema.TypeBool,
			//	Computed: true,
			//},
			//
			//"created_on": {
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},
			//
			//"modified_on": {
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},
			//"properties": releaseDefinitionProperties,
			//"artifacts":  artifacts,
			//"comment": {
			//	Type:     schema.TypeString,
			//	Optional: true,
			//	Default:  "Managed by terraform",
			//},
		},
	}
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
	var releaseDefinitionReference *int
	if err == nil {
		releaseDefinitionReference = &releaseDefinitionID
	} else {
		releaseDefinitionReference = nil
	}

	//variableGroups := expandIntList(d.Get("variable_groups").([]interface{}))

	environments, environmentsError := expandReleaseDefinitionEnvironmentList(d.Get("stage").([]interface{}))
	if environmentsError != nil {
		return nil, "", environmentsError
	}

	/*
		variables, variablesError := expandReleaseConfigurationVariableValueSet(d.Get("variable").(*schema.Set).List())
		if variablesError != nil {
			return nil, "", variablesError
		}

		properties, propertiesErrors := expandReleaseDefinitionsProperties(d.Get("properties").(*schema.Set).List())
		if propertiesErrors != nil {
			return nil, "", propertiesErrors
		}

		artifacts, artifactsErrors := expandReleaseArtifactList(d.Get("artifacts"))
		if artifactsErrors != nil {
			return nil, "", artifactsErrors
		}

		now := azuredevops.Time{Time: time.Now().UTC()}
	*/

	releaseDefinition := release.ReleaseDefinition{
		Id:   releaseDefinitionReference,
		Name: converter.String(d.Get("name").(string)),
		Path: converter.String(d.Get("path").(string)),
		//Revision:          converter.Int(d.Get("revision").(int)),
		//Source:            &release.ReleaseDefinitionSourceValues.RestApi,
		//Description:       converter.String(d.Get("description").(string)),
		Environments: environments,
		//Variables:         &variables,
		//ReleaseNameFormat: converter.String(d.Get("release_name_format").(string)),
		//VariableGroups:    &variableGroups,
		//Properties:        &properties,
		//Artifacts:         &artifacts,
		//Comment:           converter.String(d.Get("comment").(string)),

		//CreatedBy:         &webapi.IdentityRef{},
		//CreatedOn: &now,
	}

	data, err := json.Marshal(releaseDefinition)
	fmt.Println(string(data))

	return &releaseDefinition, projectID, nil
}
