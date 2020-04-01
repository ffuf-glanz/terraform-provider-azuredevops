package azuredevops

import (
	"fmt"
	"github.com/stretchr/testify/require"

	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/testhelper"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/release"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
)

var testReleaseDefinition = release.ReleaseDefinition{
	Id:       converter.Int(100),
	Revision: converter.Int(1),
	Name:     converter.String("Name"),
	Path:     converter.String("\\"),
	//VariableGroups: &[]int()
}

/**
 * Begin unit tests
 */

// verifies that the flatten/expand round trip yields the same release definition
func TestAzureDevOpsReleaseDefinition_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceReleaseDefinition().Schema, nil)
	flattenReleaseDefinition(resourceData, &testReleaseDefinition, testProjectID)

	releaseDefinitionAfterRoundTrip, projectID, err := expandReleaseDefinition(resourceData)

	require.Nil(t, err)
	require.Equal(t, testReleaseDefinition, *releaseDefinitionAfterRoundTrip)
	require.Equal(t, testProjectID, projectID)
}

/**
 * Begin acceptance tests
 */

// validates that an apply followed by another apply (i.e., resource update) will be reflected in AzDO and the
// underlying terraform state.
func TestAccAzureDevOpsReleaseDefinition_CreateAndUpdate(t *testing.T) {
	projectName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	releaseDefinitionPathEmpty := ""
	releaseDefinitionNameFirst := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	//releaseDefinitionNameSecond := testAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	//releaseDefinitionPathFirst := `\` + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	//releaseDefinitionPathSecond := `\` + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	//
	//releaseDefinitionPathThird := releaseDefinitionNameFirst + `\` + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	//releaseDefinitionPathFourth := releaseDefinitionNameSecond + `\` + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	tfReleaseDefNode := "azuredevops_release_definition.release"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testhelper.TestAccPreCheck(t, nil) },
		Providers:    testAccProviders,
		CheckDestroy: testAccReleaseDefinitionCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testhelper.TestAccReleaseDefinitionResource(projectName, releaseDefinitionNameFirst, releaseDefinitionPathEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfReleaseDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfReleaseDefNode, "revision"),
					resource.TestCheckResourceAttr(tfReleaseDefNode, "name", releaseDefinitionNameFirst),
					resource.TestCheckResourceAttr(tfReleaseDefNode, "path", releaseDefinitionPathEmpty),
					testAccCheckReleaseDefinitionResourceExists(releaseDefinitionNameFirst),
				),
			},
			//, {
			//	Config: testAccReleaseDefinitionResource(projectName, releaseDefinitionNameSecond, releaseDefinitionPathEmpty),
			//	Check: resource.ComposeTestCheckFunc(
			//		resource.TestCheckResourceAttrSet(tfReleaseDefNode, "project_id"),
			//		resource.TestCheckResourceAttrSet(tfReleaseDefNode, "revision"),
			//		resource.TestCheckResourceAttr(tfReleaseDefNode, "name", releaseDefinitionNameSecond),
			//		resource.TestCheckResourceAttr(tfReleaseDefNode, "path", releaseDefinitionPathEmpty),
			//		testAccCheckReleaseDefinitionResourceExists(releaseDefinitionNameSecond),
			//	),
			//}, {
			//	Config: testAccReleaseDefinitionResource(projectName, releaseDefinitionNameFirst, releaseDefinitionPathFirst),
			//	Check: resource.ComposeTestCheckFunc(
			//		resource.TestCheckResourceAttrSet(tfReleaseDefNode, "project_id"),
			//		resource.TestCheckResourceAttrSet(tfReleaseDefNode, "revision"),
			//		resource.TestCheckResourceAttr(tfReleaseDefNode, "name", releaseDefinitionNameFirst),
			//		resource.TestCheckResourceAttr(tfReleaseDefNode, "path", releaseDefinitionPathFirst),
			//		testAccCheckReleaseDefinitionResourceExists(releaseDefinitionNameFirst),
			//	),
			//}, {
			//	Config: testAccReleaseDefinitionResource(projectName, releaseDefinitionNameFirst,
			//		releaseDefinitionPathSecond),
			//	Check: resource.ComposeTestCheckFunc(
			//		resource.TestCheckResourceAttrSet(tfReleaseDefNode, "project_id"),
			//		resource.TestCheckResourceAttrSet(tfReleaseDefNode, "revision"),
			//		resource.TestCheckResourceAttr(tfReleaseDefNode, "name", releaseDefinitionNameFirst),
			//		resource.TestCheckResourceAttr(tfReleaseDefNode, "path", releaseDefinitionPathSecond),
			//		testAccCheckReleaseDefinitionResourceExists(releaseDefinitionNameFirst),
			//	),
			//}, {
			//	Config: testAccReleaseDefinitionResource(projectName, releaseDefinitionNameFirst, releaseDefinitionPathThird),
			//	Check: resource.ComposeTestCheckFunc(
			//		resource.TestCheckResourceAttrSet(tfReleaseDefNode, "project_id"),
			//		resource.TestCheckResourceAttrSet(tfReleaseDefNode, "revision"),
			//		resource.TestCheckResourceAttr(tfReleaseDefNode, "name", releaseDefinitionNameFirst),
			//		resource.TestCheckResourceAttr(tfReleaseDefNode, "path", releaseDefinitionPathThird),
			//		testAccCheckReleaseDefinitionResourceExists(releaseDefinitionNameFirst),
			//	),
			//}, {
			//	Config: testAccReleaseDefinitionResource(projectName, releaseDefinitionNameFirst, releaseDefinitionPathFourth),
			//	Check: resource.ComposeTestCheckFunc(
			//		resource.TestCheckResourceAttrSet(tfReleaseDefNode, "project_id"),
			//		resource.TestCheckResourceAttrSet(tfReleaseDefNode, "revision"),
			//		resource.TestCheckResourceAttr(tfReleaseDefNode, "name", releaseDefinitionNameFirst),
			//		resource.TestCheckResourceAttr(tfReleaseDefNode, "path", releaseDefinitionPathFourth),
			//		testAccCheckReleaseDefinitionResourceExists(releaseDefinitionNameFirst),
			//	),
			//},
		},
	})
}

func TestAccAzureDevOpsReleaseDefinition_CreateAndUpdate_Temp(t *testing.T) {
	projectName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	releaseDefinitionPathEmpty := ""
	releaseDefinitionNameFirst := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	tfReleaseDefNode := "azuredevops_release_definition.release"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testhelper.TestAccPreCheck(t, nil) },
		Providers:    testAccProviders,
		CheckDestroy: testAccReleaseDefinitionCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testhelper.TestAccReleaseDefinitionResourceTemp(projectName, releaseDefinitionNameFirst, releaseDefinitionPathEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfReleaseDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfReleaseDefNode, "revision"),
					resource.TestCheckResourceAttr(tfReleaseDefNode, "name", releaseDefinitionNameFirst),
					resource.TestCheckResourceAttr(tfReleaseDefNode, "path", releaseDefinitionPathEmpty),
					testAccCheckReleaseDefinitionResourceExists(releaseDefinitionNameFirst),
				),
			},
		},
	})
}

func TestAccAzureDevOpsReleaseDefinition_CreateAndUpdate_Agentless(t *testing.T) {
	projectName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	releaseDefinitionPathEmpty := ""
	releaseDefinitionNameFirst := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	tfReleaseDefNode := "azuredevops_release_definition.release"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testhelper.TestAccPreCheck(t, nil) },
		Providers:    testAccProviders,
		CheckDestroy: testAccReleaseDefinitionCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testhelper.TestAccReleaseDefinitionResourceAgentless(projectName, releaseDefinitionNameFirst, releaseDefinitionPathEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfReleaseDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfReleaseDefNode, "revision"),
					resource.TestCheckResourceAttr(tfReleaseDefNode, "name", releaseDefinitionNameFirst),
					resource.TestCheckResourceAttr(tfReleaseDefNode, "path", releaseDefinitionPathEmpty),
					testAccCheckReleaseDefinitionResourceExists(releaseDefinitionNameFirst),
				),
			},
		},
	})
}

// Given the name of an AzDO release definition, this will return a function that will check whether
// or not the definition (1) exists in the state and (2) exist in AzDO and (3) has the correct name
func testAccCheckReleaseDefinitionResourceExists(expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		releaseDef, ok := s.RootModule().Resources["azuredevops_release_definition.release"]
		if !ok {
			return fmt.Errorf("Did not find a release definition in the TF state")
		}

		releaseDefinition, err := getReleaseDefinitionFromResource(releaseDef)
		if err != nil {
			return err
		}

		if *releaseDefinition.Name != expectedName {
			return fmt.Errorf("Build Definition has Name=%s, but expected Name=%s", *releaseDefinition.Name, expectedName)
		}

		return nil
	}
}

// verifies that all release definitions referenced in the state are destroyed. This will be invoked
// *after* terrafform destroys the resource but *before* the state is wiped clean.
func testAccReleaseDefinitionCheckDestroy(s *terraform.State) error {
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_release_definition" {
			continue
		}

		// indicates the release definition still exists - this should fail the test
		if _, err := getReleaseDefinitionFromResource(resource); err == nil {
			return fmt.Errorf("Unexpectedly found a release definition that should be deleted")
		}
	}

	return nil
}

// given a resource from the state, return a release definition (and error)
func getReleaseDefinitionFromResource(resource *terraform.ResourceState) (*release.ReleaseDefinition, error) {
	releaseDefID, err := strconv.Atoi(resource.Primary.ID)
	if err != nil {
		return nil, err
	}

	projectID := resource.Primary.Attributes["project_id"]
	clients := testAccProvider.Meta().(*config.AggregatedClient)
	return clients.ReleaseClient.GetReleaseDefinition(clients.Ctx, release.GetReleaseDefinitionArgs{
		Project:      &projectID,
		DefinitionId: &releaseDefID,
	})
}
