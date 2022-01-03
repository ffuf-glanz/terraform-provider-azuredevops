package testutils

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops"
)

// initialized once, so it can be shared by each acceptance test
var provider = azuredevops.Provider()

// GetProvider returns the azuredevops provider
func GetProvider() *schema.Provider {
	return provider
}

// GetProviders returns a map of all providers needed for the project
func GetProviders() map[string]terraform.ResourceProvider {
	return map[string]terraform.ResourceProvider{
		"azuredevops": GetProvider(),
	}
}

// PreCheck checks that the requisite environment variables are set
func PreCheck(t *testing.T, additionalEnvVars *[]string) {
	requiredEnvVars := []string{
		"AZDO_ORG_SERVICE_URL",
		"AZDO_PERSONAL_ACCESS_TOKEN",
	}
	if additionalEnvVars != nil {
		requiredEnvVars = append(requiredEnvVars, *additionalEnvVars...)
	}
	missing := false
	for _, variable := range requiredEnvVars {
		if _, ok := os.LookupEnv(variable); !ok {
			missing = true
			t.Errorf("`%s` must be set for this acceptance test!", variable)
		}
	}
	if missing {
		t.Fatalf("Some environment variables missing.")
	}
}

// GenerateResourceName generates a random name with a constant prefix, useful for acceptance tests
func GenerateResourceName() string {
	return "test-acc-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
}
