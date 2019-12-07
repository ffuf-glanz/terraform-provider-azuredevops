package azuredevops

import (
	//"context"
	//"errors"
	"fmt"

	//"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/testhelper"
	"strconv"
	"strings"
	"testing"

	//"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	//"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/release"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	//"github.com/stretchr/testify/require"
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

// validates that all supported repo types are allowed by the schema
/*
func TestAzureDevOpsReleaseDefinition_RepoTypeListIsCorrect(t *testing.T) {
	expectedRepoTypes := []string{"GitHub", "TfsGit"}
	repoSchema := resourceReleaseDefinition().Schema["repository"]
	repoTypeSchema := repoSchema.Elem.(*schema.Resource).Schema["repo_type"]

	for _, repoType := range expectedRepoTypes {
		_, errors := repoTypeSchema.ValidateFunc(repoType, "")
		require.Equal(t, 0, len(errors), "Repo type unexpectedly did not pass validation")
	}
}

// validates that and error is thrown if any of the un-supported file path characters are used
func TestAzureDevOpsReleaseDefinition_PathInvalidCharacterListIsError(t *testing.T) {
	expectedInvalidPathCharacters := []string{"<", ">", "|", ":", "$", "@", "\"", "/", "%", "+", "*", "?"}
	pathSchema := resourceReleaseDefinition().Schema["path"]

	for _, repoType := range expectedInvalidPathCharacters {
		_, errors := pathSchema.ValidateFunc(repoType, "")
		require.Equal(t, "<>|:$@\"/%+*? are not allowed", errors[0].Error())
	}
}

// verifies that the flatten/expand round trip yields the same release definition
func TestAzureDevOpsReleaseDefinition_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceReleaseDefinition().Schema, nil)
	flattenReleaseDefinition(resourceData, &testReleaseDefinition, testProjectID)

	releaseDefinitionAfterRoundTrip, projectID, err := expandReleaseDefinition(resourceData)

	require.Nil(t, err)
	require.Equal(t, testReleaseDefinition, *releaseDefinitionAfterRoundTrip)
	require.Equal(t, testProjectID, projectID)
}

// verifies that an expand will fail if there is insufficient configuration data found in the resource
func TestAzureDevOpsReleaseDefinition_Expand_FailsIfNotEnoughData(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceReleaseDefinition().Schema, nil)
	_, _, err := expandReleaseDefinition(resourceData)
	require.NotNil(t, err)
}

// verifies that if an error is produced on create, the error is not swallowed
func TestAzureDevOpsReleaseDefinition_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceData := schema.TestResourceDataRaw(t, resourceReleaseDefinition().Schema, nil)
	flattenReleaseDefinition(resourceData, &testReleaseDefinition, testProjectID)

	releaseClient := azdosdkmocks.NewMockReleaseClient(ctrl)
	clients := &config.AggregatedClient{ReleaseClient: releaseClient, Ctx: context.Background()}

	expectedArgs := release.CreateDefinitionArgs{Definition: &testReleaseDefinition, Project: &testProjectID}
	releaseClient.
		EXPECT().
		CreateDefinition(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateDefinition() Failed")).
		Times(1)

	err := resourceReleaseDefinitionCreate(resourceData, clients)
	require.Contains(t, err.Error(), "CreateDefinition() Failed")
}

// verifies that if an error is produced on a read, it is not swallowed
func TestAzureDevOpsReleaseDefinition_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceData := schema.TestResourceDataRaw(t, resourceReleaseDefinition().Schema, nil)
	flattenReleaseDefinition(resourceData, &testReleaseDefinition, testProjectID)

	releaseClient := azdosdkmocks.NewMockReleaseClient(ctrl)
	clients := &config.AggregatedClient{ReleaseClient: releaseClient, Ctx: context.Background()}

	expectedArgs := release.GetDefinitionArgs{DefinitionId: testReleaseDefinition.Id, Project: &testProjectID}
	releaseClient.
		EXPECT().
		GetDefinition(clients.Ctx, expectedArgs).
		Return(nil, errors.New("GetDefinition() Failed")).
		Times(1)

	err := resourceReleaseDefinitionRead(resourceData, clients)
	require.Equal(t, "GetDefinition() Failed", err.Error())
}

// verifies that if an error is produced on a delete, it is not swallowed
func TestAzureDevOpsReleaseDefinition_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceData := schema.TestResourceDataRaw(t, resourceReleaseDefinition().Schema, nil)
	flattenReleaseDefinition(resourceData, &testReleaseDefinition, testProjectID)

	releaseClient := azdosdkmocks.NewMockReleaseClient(ctrl)
	clients := &config.AggregatedClient{ReleaseClient: releaseClient, Ctx: context.Background()}

	expectedArgs := release.DeleteDefinitionArgs{DefinitionId: testReleaseDefinition.Id, Project: &testProjectID}
	releaseClient.
		EXPECT().
		DeleteDefinition(clients.Ctx, expectedArgs).
		Return(errors.New("DeleteDefinition() Failed")).
		Times(1)

	err := resourceReleaseDefinitionDelete(resourceData, clients)
	require.Equal(t, "DeleteDefinition() Failed", err.Error())
}

// verifies that if an error is produced on an update, it is not swallowed
func TestAzureDevOpsReleaseDefinition_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resourceData := schema.TestResourceDataRaw(t, resourceReleaseDefinition().Schema, nil)
	flattenReleaseDefinition(resourceData, &testReleaseDefinition, testProjectID)

	releaseClient := azdosdkmocks.NewMockReleaseClient(ctrl)
	clients := &config.AggregatedClient{ReleaseClient: releaseClient, Ctx: context.Background()}

	expectedArgs := release.UpdateDefinitionArgs{
		Definition:   &testReleaseDefinition,
		DefinitionId: testReleaseDefinition.Id,
		Project:      &testProjectID,
	}

	releaseClient.
		EXPECT().
		UpdateDefinition(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateDefinition() Failed")).
		Times(1)

	err := resourceReleaseDefinitionUpdate(resourceData, clients)
	require.Equal(t, "UpdateDefinition() Failed", err.Error())
}
*/

/**
 * Begin acceptance tests
 */

// validates that an apply followed by another apply (i.e., resource update) will be reflected in AzDO and the
// underlying terraform state.
func TestAccAzureDevOpsReleaseDefinition_CreateAndUpdate(t *testing.T) {
	projectName := testAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	releaseDefinitionPathEmpty := ""
	releaseDefinitionNameFirst := testAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	//releaseDefinitionNameSecond := testAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	//releaseDefinitionPathFirst := `\` + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	//releaseDefinitionPathSecond := `\` + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	//
	//releaseDefinitionPathThird := releaseDefinitionNameFirst + `\` + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	//releaseDefinitionPathFourth := releaseDefinitionNameSecond + `\` + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	tfReleaseDefNode := "azuredevops_release_definition.build"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testhelper.TestAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccReleaseDefinitionCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccReleaseDefinitionResource(projectName, releaseDefinitionNameFirst, releaseDefinitionPathEmpty),
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

// HCL describing an AzDO release definition
func testAccReleaseDefinitionResource(projectName string, releaseDefinitionName string, releasePath string) string {
	releaseDefinitionResource := fmt.Sprintf(`
resource "azuredevops_release_definition" "release" {
	project_id      = azuredevops_project.project.id
	name            = "%s"
	path			= "%s"
}`, releaseDefinitionName, strings.ReplaceAll(releasePath, `\`, `\\`))

	projectResource := testhelper.TestAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, releaseDefinitionResource)
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
