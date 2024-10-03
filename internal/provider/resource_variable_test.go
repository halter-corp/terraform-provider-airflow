package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAirflowVariable_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	rNameUpdated := acctest.RandomWithPrefix("tf-acc-test")

	resourceName := "airflow_variable.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAirflowVariableCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAirflowVariableConfigBasic(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "key", rName),
					resource.TestCheckResourceAttr(resourceName, "value", rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAirflowVariableConfigBasic(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "key", rName),
					resource.TestCheckResourceAttr(resourceName, "value", rNameUpdated),
				),
			},
		},
	})
}

func TestAccAirflowVariable_desc(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	rNameUpdated := acctest.RandomWithPrefix("tf-acc-test")
	rDesc := acctest.RandomWithPrefix("tf-acc-test")

	resourceName := "airflow_variable.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAirflowVariableCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAirflowVariableConfigDesc(rName, rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "key", rName),
					resource.TestCheckResourceAttr(resourceName, "value", rName),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAirflowVariableConfigDesc(rName, rNameUpdated, rDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "key", rName),
					resource.TestCheckResourceAttr(resourceName, "value", rNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
				),
			},
			{
				Config: testAccAirflowVariableConfigBasic(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "key", rName),
					resource.TestCheckResourceAttr(resourceName, "value", rName),
					// bug? Airflow API seems unhappy with trying to set description to null.
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
		},
	})
}

func testAccCheckAirflowVariableCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "airflow_variable" {
			continue
		}

		variable, res, err := client.ApiClient.VariableApi.GetVariable(client.AuthContext, rs.Primary.ID).Execute()
		if err == nil {
			if *variable.Key == rs.Primary.ID {
				return fmt.Errorf("Airflow Variable (%s) still exists.", rs.Primary.ID)
			}
		}

		if res != nil && res.StatusCode == 404 {
			continue
		}
	}

	return nil
}

func testAccAirflowVariableConfigBasic(rName, value string) string {
	return fmt.Sprintf(`
resource "airflow_variable" "test" {
  key    = %[1]q
  value  = %[2]q
}
`, rName, value)
}

func testAccAirflowVariableConfigDesc(rName, value, desc string) string {
	return fmt.Sprintf(`
resource "airflow_variable" "test" {
  key          = %[1]q
  value        = %[2]q
  description  = %[3]q
}
`, rName, value, desc)
}
