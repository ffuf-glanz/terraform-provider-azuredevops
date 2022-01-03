package azuredevops

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	crud "github.com/microsoft/terraform-provider-azuredevops/azuredevops/crud/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/tfhelper"
)

// ResourceServiceEndpointSonarQube schema and implementation for SonarQube service endpoint resource
func ResourceServiceEndpointSonarQube() *schema.Resource {
	r := crud.GenBaseServiceEndpointResource(flattenServiceEndpointSonarQube, expandServiceEndpointSonarQube, parseImportedProjectIDAndServiceEndpointID)

	r.Schema["url"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		Description:  "Url for the SonarQube Server",
	}

	r.Schema["token"] = &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		Sensitive:        true,
		DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
		ValidateFunc:     validation.StringIsNotWhiteSpace,
		Description:      "Authentication Token generated through SonarQube (go to My Account > Security > Generate Tokens)",
	}
	// Add a spot in the schema to store the token secretly
	stSecretHashKey, stSecretHashSchema := tfhelper.GenerateSecreteMemoSchema("token")
	r.Schema[stSecretHashKey] = stSecretHashSchema

	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointSonarQube(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string) {
	serviceEndpoint, projectID := crud.DoBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Scheme: converter.String("UsernamePassword"),
		Parameters: &map[string]string{
			"username": d.Get("token").(string),
		},
	}
	serviceEndpoint.Type = converter.String("sonarqube")
	serviceEndpoint.Url = converter.String(d.Get("url").(string))
	return serviceEndpoint, projectID
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointSonarQube(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	crud.DoBaseFlattening(d, serviceEndpoint, projectID)

	tfhelper.HelpFlattenSecret(d, "token")

	d.Set("url", *serviceEndpoint.Url)
	d.Set("token", (*serviceEndpoint.Authorization.Parameters)["username"])
}
