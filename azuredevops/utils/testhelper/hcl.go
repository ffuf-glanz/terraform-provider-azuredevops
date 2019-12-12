package testhelper

import (
	"fmt"
	"strings"
)

// TestAccAzureGitRepoResource HCL describing an AzDO GIT repository resource
func TestAccAzureGitRepoResource(projectName string, gitRepoName string, initType string) string {
	azureGitRepoResource := fmt.Sprintf(`
resource "azuredevops_azure_git_repository" "gitrepo" {
	project_id      = azuredevops_project.project.id
	name            = "%s"
	initialization {
		init_type = "%s"
	}
}`, gitRepoName, initType)

	projectResource := TestAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, azureGitRepoResource)
}

// TestAccGroupDataSource HCL describing an AzDO Group Data Source
func TestAccGroupDataSource(projectName string, groupName string) string {
	dataSource := fmt.Sprintf(`
data "azuredevops_group" "group" {
	project_id = azuredevops_project.project.id
	name       = "%s"
}`, groupName)

	projectResource := TestAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, dataSource)
}

// TestAccProjectResource HCL describing an AzDO project
func TestAccProjectResource(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
	project_name       = "%s"
	description        = "%s-description"
	visibility         = "private"
	version_control    = "Git"
	work_item_template = "Agile"
}`, projectName, projectName)
}

// TestAccUserEntitlementResource HCL describing an AzDO UserEntitlement
func TestAccUserEntitlementResource(principalName string) string {
	return fmt.Sprintf(`
resource "azuredevops_user_entitlement" "user" {
	principal_name     = "%s"
	account_license_type = "express"
}`, principalName)
}

// TestAccServiceEndpointGitHubResource HCL describing an AzDO service endpoint
func TestAccServiceEndpointGitHubResource(projectName string, serviceEndpointName string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_github" "serviceendpoint" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"

}`, serviceEndpointName)

	projectResource := TestAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

// TestAccServiceEndpointDockerHubResource HCL describing an AzDO service endpoint
func TestAccServiceEndpointDockerHubResource(projectName string, serviceEndpointName string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_dockerhub" "serviceendpoint" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"
}`, serviceEndpointName)

	projectResource := TestAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

// TestAccVariableGroupResource HCL describing an AzDO variable group
func TestAccVariableGroupResource(projectName string, variableGroupName string, allowAccess bool) string {
	variableGroupResource := fmt.Sprintf(`
resource "azuredevops_variable_group" "vg" {
	project_id  = azuredevops_project.project.id
	name        = "%s"
	description = "A sample variable group."
	allow_access = %t
	variable {
		name      = "key1"
		value     = "value1"
		is_secret = true
	}

	variable {
		name  = "key2"
		value = "value2"
	}

	variable {
		name = "key3"
	}
}`, variableGroupName, allowAccess)

	projectResource := TestAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, variableGroupResource)
}

// TestAccAgentPoolResource HCL describing an AzDO Agent Pool
func TestAccAgentPoolResource(poolName string) string {
	return fmt.Sprintf(`
resource "azuredevops_agent_pool" "pool" {
	name           = "%s"
	auto_provision = false
	pool_type      = "automation"
	}`, poolName)
}

// TestAccBuildDefinitionResource HCL describing an AzDO build definition
func TestAccBuildDefinitionResource(projectName string, buildDefinitionName string, buildPath string) string {
	buildDefinitionResource := fmt.Sprintf(`
resource "azuredevops_build_definition" "build" {
	project_id      = azuredevops_project.project.id
	name            = "%s"
	agent_pool_name = "Hosted Ubuntu 1604"
	path			= "%s"

	repository {
	  repo_type             = "GitHub"
	  repo_name             = "repoOrg/repoName"
	  branch_name           = "branch"
	  yml_path              = "path/to/yaml"
	}
}`, buildDefinitionName, strings.ReplaceAll(buildPath, `\`, `\\`))

	projectResource := TestAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, buildDefinitionResource)
}

// HCL describing an AzDO release definition
func TestAccReleaseDefinitionResource(projectName string, releaseDefinitionName string, releasePath string) string {
	releaseDefinitionResource := fmt.Sprintf(`
resource "azuredevops_release_definition" "release" {
  project_id = "merlin"
  name = "%s"
  path = "\\"

  stage {
    name = "Stage 1"
    rank = 1

    agent_job {
      name = "Agent job 1"
      rank = 1
 	  demand {
		name =  "equals_condition_name"
		value = "x"
	  }
      demand {
		name =  "exists_condition_name"
	  }
      agent_pool_hosted_azure_pipelines {
        agent_pool_id = 52
        agent_specification = "ubuntu-18.04"
      }
      timeout_in_minutes = 0
      max_execution_time_in_minutes = 1
      condition = "succeeded()"
		// overrideInput {} // TODO
		// enable_access_token ? Do we need this on this level?
    }

	agent_job {
      name = "Agent job 3"
      rank = 3
      agent_pool_hosted_azure_pipelines {
        agent_pool_id = 52
        agent_specification = "ubuntu-18.04"
      }
	
	  multi_configuration {
		multipliers = "OperatingSystem"
		number_of_agents = 1
	  }

      condition = "succeeded()"
		// overrideInput {} // TODO
		// enable_access_token ? Do we need this on this level?
    }

	agent_job {
      name = "Agent job 2"
      rank = 2
      agent_pool_hosted_azure_pipelines {
        agent_pool_id = 52
        agent_specification = "ubuntu-18.04"
      }
	  multi_agent {
		max_number_of_agents = 1
	  }
      condition = "succeeded()"
		// overrideInput {} // TODO
		// enable_access_token ? Do we need this on this level?
    }

    pre_deploy_approval {
      approval {
        is_automated = true
        rank = 1
      }
    }

    post_deploy_approval {
      approval {
        is_automated = true
        rank = 1
      }
    }

    retention_policy {
      days_to_keep = 1
      releases_to_keep = 1
    }
  }
}`, releaseDefinitionName)

	projectResource := TestAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, releaseDefinitionResource)
}

func TestAccReleaseDefinitionResourceTemp(projectName string, releaseDefinitionName string, releasePath string) string {
	releaseDefinitionResource := fmt.Sprintf(`
resource "azuredevops_release_definition" "release" {
  project_id = "merlin" // TODO: revert this back to azuredevops_project.project.id
  name = "%s"
  path = "\\"

  stage {
    name = "Stage 1"
    rank = 1

    deployment_group_job {
      name = "Deployment group job 1"
      rank = 1
      deployment_group_id = 1619 // QueueId
      //DeploymentHealthOption: OneAtATime
      tags = ["deployment_group_job_1"]
      timeout_in_minutes = 0
      max_execution_time_in_minutes = 1
      // download_artifact = {}
      allow_scripts_to_access_oauth_token = true
      condition = "succeedeed()"
      // tasks = file('foo.yml')
    }

	deployment_group_job {
      name = "Deployment group job 2"
      rank = 2
      deployment_group_id = 1619
      tags = ["deployment_group_job_2"]
      multiple { //DeploymentHealthOption: Custom
        max_targets_in_parallel = 0 // HealthPercent
      }
      timeout_in_minutes = 0
      max_execution_time_in_minutes = 1
      // download_artifact = {}
      allow_scripts_to_access_oauth_token = true
      condition = "succeedeed()"
      // tasks = file('foo.yml')
    }

    pre_deploy_approval {
      approval {
        is_automated = true
        rank = 1
      }
    }

    post_deploy_approval {
      approval {
        is_automated = true
        rank = 1
      }
    }

    retention_policy {
      days_to_keep = 1
      releases_to_keep = 1
    }
  }
}`, releaseDefinitionName)

	projectResource := TestAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, releaseDefinitionResource)
}

func TestAccReleaseDefinitionResourceAgentless(projectName string, releaseDefinitionName string, releasePath string) string {
	releaseDefinitionResource := fmt.Sprintf(`
resource "azuredevops_release_definition" "release" {
  project_id = "merlin"
  name = "%s"
  path = "\\"

  stage {
    name = "Stage 1"
    rank = 1

    agentless_job {
      name = "Agentless job 1"
      rank = 1

      timeout_in_minutes = 0
      condition = "succeeded()"
    }

    agentless_job {
      name = "Agentless job 2"
      rank = 2

      multi_configuration {
        multipliers = "OperatingSystem"
      }

      condition = "succeeded()"
    }

    pre_deploy_approval {
      approval {
        is_automated = true
        rank = 1
      }
    }

    post_deploy_approval {
      approval {
        is_automated = true
        rank = 1
      }
    }

    retention_policy {
      days_to_keep = 1
      releases_to_keep = 1
    }
  }
}`, releaseDefinitionName)

	projectResource := TestAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, releaseDefinitionResource)
}

// TestAccGroupMembershipResource full terraform stanza to standup a group membership
func TestAccGroupMembershipResource(projectName, groupName, userPrincipalName string) string {
	membershipDependenciesStanza := TestAccGroupMembershipDependencies(projectName, groupName, userPrincipalName)
	membershipStanza := `
resource "azuredevops_group_membership" "membership" {
	group = data.azuredevops_group.group.descriptor
	members = [azuredevops_user_entitlement.user.descriptor]
}`

	return membershipDependenciesStanza + "\n" + membershipStanza
}

// TestAccGroupMembershipDependencies all the dependencies needed to configure a group membership
func TestAccGroupMembershipDependencies(projectName, groupName, userPrincipalName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
	project_name = "%s"
}
data "azuredevops_group" "group" {
	project_id = azuredevops_project.project.id
	name       = "%s"
}
resource "azuredevops_user_entitlement" "user" {
	principal_name       = "%s"
	account_license_type = "express"
}

output "group_descriptor" {
	value = data.azuredevops_group.group.descriptor
}
output "user_descriptor" {
	value = azuredevops_user_entitlement.user.descriptor
}
`, projectName, groupName, userPrincipalName)
}
