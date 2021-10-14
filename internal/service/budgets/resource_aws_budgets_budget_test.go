package budgets_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/budgets"
	"github.com/hashicorp/go-multierror"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	tfbudgets "github.com/hashicorp/terraform-provider-aws/internal/service/budgets"
	tfbudgets "github.com/hashicorp/terraform-provider-aws/internal/service/budgets"
	tfbudgets "github.com/hashicorp/terraform-provider-aws/internal/service/budgets"
	tfbudgets "github.com/hashicorp/terraform-provider-aws/internal/service/budgets"
	tfbudgets "github.com/hashicorp/terraform-provider-aws/internal/service/budgets"
	tfbudgets "github.com/hashicorp/terraform-provider-aws/internal/service/budgets"
	tfbudgets "github.com/hashicorp/terraform-provider-aws/internal/service/budgets"
	tfbudgets "github.com/hashicorp/terraform-provider-aws/internal/service/budgets"
	tfbudgets "github.com/hashicorp/terraform-provider-aws/internal/service/budgets"
	tfbudgets "github.com/hashicorp/terraform-provider-aws/internal/service/budgets"
	tfbudgets "github.com/hashicorp/terraform-provider-aws/internal/service/budgets"
	tfbudgets "github.com/hashicorp/terraform-provider-aws/internal/service/budgets"
	tfbudgets "github.com/hashicorp/terraform-provider-aws/internal/service/budgets"
)

func init() {
	resource.AddTestSweepers("aws_budgets_budget", &resource.Sweeper{
		Name: "aws_budgets_budget",
		F:    testSweepBudgetsBudgets,
	})
}

func testSweepBudgetsBudgets(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.(*conns.AWSClient).BudgetsConn
	accountID := client.(*conns.AWSClient).AccountID
	input := &budgets.DescribeBudgetsInput{
		AccountId: aws.String(accountID),
	}
	var sweeperErrs *multierror.Error

	for {
		output, err := conn.DescribeBudgets(input)
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping Budgets sweep for %s: %s", region, err)
			return sweeperErrs.ErrorOrNil() // In case we have completed some pages, but had errors
		}
		if err != nil {
			sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error retrieving Budgets: %w", err))
			return sweeperErrs
		}

		for _, budget := range output.Budgets {
			name := aws.StringValue(budget.BudgetName)

			log.Printf("[INFO] Deleting Budget: %s", name)
			_, err := conn.DeleteBudget(&budgets.DeleteBudgetInput{
				AccountId:  aws.String(accountID),
				BudgetName: aws.String(name),
			})
			if tfawserr.ErrMessageContains(err, budgets.ErrCodeNotFoundException, "") {
				continue
			}
			if err != nil {
				sweeperErr := fmt.Errorf("error deleting Budget (%s): %w", name, err)
				log.Printf("[ERROR] %s", sweeperErr)
				sweeperErrs = multierror.Append(sweeperErrs, sweeperErr)
				continue
			}
		}

		if aws.StringValue(output.NextToken) == "" {
			break
		}
		input.NextToken = output.NextToken
	}

	return sweeperErrs.ErrorOrNil()
}

func TestAccAWSBudgetsBudget_basic(t *testing.T) {
	var budget budgets.Budget
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_budgets_budget.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(budgets.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, budgets.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccAWSBudgetsBudgetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSBudgetsBudgetConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccAWSBudgetsBudgetExists(resourceName, &budget),
					acctest.CheckResourceAttrAccountID(resourceName, "account_id"),
					acctest.CheckResourceAttrGlobalARN(resourceName, "arn", "budgets", fmt.Sprintf(`budget/%s`, rName)),
					resource.TestCheckResourceAttr(resourceName, "budget_type", "RI_UTILIZATION"),
					resource.TestCheckResourceAttr(resourceName, "cost_filter.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "cost_filter.*", map[string]string{
						"name":     "Service",
						"values.#": "1",
						"values.0": "Amazon Elasticsearch Service",
					}),
					resource.TestCheckResourceAttr(resourceName, "cost_filters.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "cost_filters.Service", "Amazon Elasticsearch Service"),
					resource.TestCheckResourceAttr(resourceName, "limit_amount", "100.0"),
					resource.TestCheckResourceAttr(resourceName, "limit_unit", "PERCENTAGE"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "time_period_end"),
					resource.TestCheckResourceAttrSet(resourceName, "time_period_start"),
					resource.TestCheckResourceAttr(resourceName, "time_unit", "QUARTERLY"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSBudgetsBudget_Name_Generated(t *testing.T) {
	var budget budgets.Budget
	resourceName := "aws_budgets_budget.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(budgets.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, budgets.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccAWSBudgetsBudgetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSBudgetsBudgetConfigNameGenerated(),
				Check: resource.ComposeTestCheckFunc(
					testAccAWSBudgetsBudgetExists(resourceName, &budget),
					resource.TestCheckResourceAttr(resourceName, "budget_type", "RI_COVERAGE"),
					resource.TestCheckResourceAttr(resourceName, "cost_filter.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "cost_filter.*", map[string]string{
						"name":     "Service",
						"values.#": "1",
						"values.0": "Amazon Redshift",
					}),
					resource.TestCheckResourceAttr(resourceName, "cost_filters.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "cost_filters.Service", "Amazon Redshift"),
					resource.TestCheckResourceAttr(resourceName, "limit_amount", "100.0"),
					resource.TestCheckResourceAttr(resourceName, "limit_unit", "PERCENTAGE"),
					create.TestCheckResourceAttrNameGenerated(resourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", "terraform-"),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "time_period_end"),
					resource.TestCheckResourceAttrSet(resourceName, "time_period_start"),
					resource.TestCheckResourceAttr(resourceName, "time_unit", "ANNUALLY"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSBudgetsBudget_NamePrefix(t *testing.T) {
	var budget budgets.Budget
	resourceName := "aws_budgets_budget.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(budgets.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, budgets.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccAWSBudgetsBudgetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSBudgetsBudgetConfigNamePrefix("tf-acc-test-prefix-"),
				Check: resource.ComposeTestCheckFunc(
					testAccAWSBudgetsBudgetExists(resourceName, &budget),
					resource.TestCheckResourceAttr(resourceName, "budget_type", "SAVINGS_PLANS_UTILIZATION"),
					resource.TestCheckResourceAttr(resourceName, "cost_filter.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "cost_filters.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "limit_amount", "100.0"),
					resource.TestCheckResourceAttr(resourceName, "limit_unit", "PERCENTAGE"),
					create.TestCheckResourceAttrNameFromPrefix(resourceName, "name", "tf-acc-test-prefix-"),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", "tf-acc-test-prefix-"),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "time_period_end"),
					resource.TestCheckResourceAttrSet(resourceName, "time_period_start"),
					resource.TestCheckResourceAttr(resourceName, "time_unit", "MONTHLY"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSBudgetsBudget_disappears(t *testing.T) {
	var budget budgets.Budget
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_budgets_budget.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(budgets.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, budgets.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccAWSBudgetsBudgetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSBudgetsBudgetConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccAWSBudgetsBudgetExists(resourceName, &budget),
					acctest.CheckResourceDisappears(acctest.Provider, ResourceBudget(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAWSBudgetsBudget_CostTypes(t *testing.T) {
	var budget budgets.Budget
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_budgets_budget.test"

	now := time.Now().UTC()
	ts1 := now.AddDate(0, 0, -14)
	ts2 := time.Date(2050, 1, 1, 00, 0, 0, 0, time.UTC)
	ts3 := now.AddDate(0, 0, -28)
	ts4 := time.Date(2060, 7, 1, 00, 0, 0, 0, time.UTC)
	startDate1 := tfbudgets.TimePeriodTimestampToString(&ts1)
	endDate1 := tfbudgets.TimePeriodTimestampToString(&ts2)
	startDate2 := tfbudgets.TimePeriodTimestampToString(&ts3)
	endDate2 := tfbudgets.TimePeriodTimestampToString(&ts4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(budgets.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, budgets.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccAWSBudgetsBudgetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSBudgetsBudgetConfigCostTypes(rName, startDate1, endDate1),
				Check: resource.ComposeTestCheckFunc(
					testAccAWSBudgetsBudgetExists(resourceName, &budget),
					resource.TestCheckResourceAttr(resourceName, "budget_type", "COST"),
					resource.TestCheckResourceAttr(resourceName, "cost_filter.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "cost_filter.*", map[string]string{
						"name":     "AZ",
						"values.#": "2",
						"values.0": acctest.Region(),
						"values.1": acctest.AlternateRegion(),
					}),
					resource.TestCheckResourceAttr(resourceName, "cost_filters.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "cost_filters.AZ", strings.Join([]string{acctest.Region(), acctest.AlternateRegion()}, ",")),
					resource.TestCheckResourceAttr(resourceName, "cost_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cost_types.0.include_credit", "true"),
					resource.TestCheckResourceAttr(resourceName, "cost_types.0.include_discount", "false"),
					resource.TestCheckResourceAttr(resourceName, "cost_types.0.include_other_subscription", "true"),
					resource.TestCheckResourceAttr(resourceName, "cost_types.0.include_recurring", "true"),
					resource.TestCheckResourceAttr(resourceName, "cost_types.0.include_refund", "true"),
					resource.TestCheckResourceAttr(resourceName, "cost_types.0.include_subscription", "true"),
					resource.TestCheckResourceAttr(resourceName, "cost_types.0.include_support", "true"),
					resource.TestCheckResourceAttr(resourceName, "cost_types.0.include_tax", "false"),
					resource.TestCheckResourceAttr(resourceName, "cost_types.0.include_upfront", "true"),
					resource.TestCheckResourceAttr(resourceName, "cost_types.0.use_amortized", "false"),
					resource.TestCheckResourceAttr(resourceName, "cost_types.0.use_blended", "true"),
					resource.TestCheckResourceAttr(resourceName, "limit_amount", "456.78"),
					resource.TestCheckResourceAttr(resourceName, "limit_unit", "USD"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "time_period_end", endDate1),
					resource.TestCheckResourceAttr(resourceName, "time_period_start", startDate1),
					resource.TestCheckResourceAttr(resourceName, "time_unit", "DAILY"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSBudgetsBudgetConfigCostTypesUpdated(rName, startDate2, endDate2),
				Check: resource.ComposeTestCheckFunc(
					testAccAWSBudgetsBudgetExists(resourceName, &budget),
					resource.TestCheckResourceAttr(resourceName, "budget_type", "COST"),
					resource.TestCheckResourceAttr(resourceName, "cost_filter.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "cost_filter.*", map[string]string{
						"name":     "AZ",
						"values.#": "2",
						"values.0": acctest.AlternateRegion(),
						"values.1": acctest.ThirdRegion(),
					}),
					resource.TestCheckResourceAttr(resourceName, "cost_filters.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "cost_filters.AZ", strings.Join([]string{acctest.AlternateRegion(), acctest.ThirdRegion()}, ",")),
					resource.TestCheckResourceAttr(resourceName, "cost_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cost_types.0.include_credit", "false"),
					resource.TestCheckResourceAttr(resourceName, "cost_types.0.include_discount", "true"),
					resource.TestCheckResourceAttr(resourceName, "cost_types.0.include_other_subscription", "true"),
					resource.TestCheckResourceAttr(resourceName, "cost_types.0.include_recurring", "true"),
					resource.TestCheckResourceAttr(resourceName, "cost_types.0.include_refund", "false"),
					resource.TestCheckResourceAttr(resourceName, "cost_types.0.include_subscription", "true"),
					resource.TestCheckResourceAttr(resourceName, "cost_types.0.include_support", "true"),
					resource.TestCheckResourceAttr(resourceName, "cost_types.0.include_tax", "true"),
					resource.TestCheckResourceAttr(resourceName, "cost_types.0.include_upfront", "true"),
					resource.TestCheckResourceAttr(resourceName, "cost_types.0.use_amortized", "false"),
					resource.TestCheckResourceAttr(resourceName, "cost_types.0.use_blended", "false"),
					resource.TestCheckResourceAttr(resourceName, "limit_amount", "567.89"),
					resource.TestCheckResourceAttr(resourceName, "limit_unit", "USD"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "time_period_end", endDate2),
					resource.TestCheckResourceAttr(resourceName, "time_period_start", startDate2),
					resource.TestCheckResourceAttr(resourceName, "time_unit", "DAILY"),
				),
			},
		},
	})
}

func TestAccAWSBudgetsBudget_Notifications(t *testing.T) {
	var budget budgets.Budget
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_budgets_budget.test"
	snsTopicResourceName := "aws_sns_topic.test"

	domain := acctest.RandomDomainName()
	emailAddress1 := acctest.RandomEmailAddress(domain)
	emailAddress2 := acctest.RandomEmailAddress(domain)
	emailAddress3 := acctest.RandomEmailAddress(domain)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(budgets.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, budgets.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccAWSBudgetsBudgetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSBudgetsBudgetConfigNotifications(rName, emailAddress1, emailAddress2),
				Check: resource.ComposeTestCheckFunc(
					testAccAWSBudgetsBudgetExists(resourceName, &budget),
					resource.TestCheckResourceAttr(resourceName, "budget_type", "USAGE"),
					resource.TestCheckResourceAttr(resourceName, "cost_filter.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "cost_filters.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "limit_amount", "432.1"),
					resource.TestCheckResourceAttr(resourceName, "limit_unit", "GBP"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "notification.*", map[string]string{
						"comparison_operator":          "GREATER_THAN",
						"notification_type":            "ACTUAL",
						"subscriber_email_addresses.#": "0",
						"subscriber_sns_topic_arns.#":  "1",
						"threshold":                    "150",
						"threshold_type":               "PERCENTAGE",
					}),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "notification.*.subscriber_sns_topic_arns.*", snsTopicResourceName, "arn"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "notification.*", map[string]string{
						"comparison_operator":          "EQUAL_TO",
						"notification_type":            "FORECASTED",
						"subscriber_email_addresses.#": "2",
						"subscriber_sns_topic_arns.#":  "0",
						"threshold":                    "200.1",
						"threshold_type":               "ABSOLUTE_VALUE",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "notification.*.subscriber_email_addresses.*", emailAddress1),
					resource.TestCheckTypeSetElemAttr(resourceName, "notification.*.subscriber_email_addresses.*", emailAddress2),
					resource.TestCheckResourceAttr(resourceName, "time_unit", "ANNUALLY"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSBudgetsBudgetConfigNotificationsUpdated(rName, emailAddress3),
				Check: resource.ComposeTestCheckFunc(
					testAccAWSBudgetsBudgetExists(resourceName, &budget),
					resource.TestCheckResourceAttr(resourceName, "budget_type", "USAGE"),
					resource.TestCheckResourceAttr(resourceName, "cost_filter.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "cost_filters.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "limit_amount", "432.1"),
					resource.TestCheckResourceAttr(resourceName, "limit_unit", "GBP"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "notification.*", map[string]string{
						"comparison_operator":          "LESS_THAN",
						"notification_type":            "ACTUAL",
						"subscriber_email_addresses.#": "1",
						"subscriber_sns_topic_arns.#":  "0",
						"threshold":                    "123.45",
						"threshold_type":               "ABSOLUTE_VALUE",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "notification.*.subscriber_email_addresses.*", emailAddress3),
					resource.TestCheckResourceAttr(resourceName, "time_unit", "ANNUALLY"),
				),
			},
		},
	})
}

func testAccAWSBudgetsBudgetExists(resourceName string, v *budgets.Budget) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Budget ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).BudgetsConn

		accountID, budgetName, err := tfbudgets.BudgetParseResourceID(rs.Primary.ID)

		if err != nil {
			return err
		}

		output, err := tfbudgets.FindBudgetByAccountIDAndBudgetName(conn, accountID, budgetName)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccAWSBudgetsBudgetDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).BudgetsConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_budgets_budget" {
			continue
		}

		accountID, budgetName, err := tfbudgets.BudgetParseResourceID(rs.Primary.ID)

		if err != nil {
			return err
		}

		_, err = tfbudgets.FindBudgetByAccountIDAndBudgetName(conn, accountID, budgetName)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("Budget Action %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccAWSBudgetsBudgetConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_budgets_budget" "test" {
  name         = %[1]q
  budget_type  = "RI_UTILIZATION"
  limit_amount = "100.0"
  limit_unit   = "PERCENTAGE"
  time_unit    = "QUARTERLY"

  cost_filters = {
    Service = "Amazon Elasticsearch Service"
  }
}
`, rName)
}

func testAccAWSBudgetsBudgetConfigNameGenerated() string {
	return `
resource "aws_budgets_budget" "test" {
  budget_type  = "RI_COVERAGE"
  limit_amount = "100.00"
  limit_unit   = "PERCENTAGE"
  time_unit    = "ANNUALLY"

  cost_filter {
    name   = "Service"
    values = ["Amazon Redshift"]
  }
}
`
}

func testAccAWSBudgetsBudgetConfigNamePrefix(namePrefix string) string {
	return fmt.Sprintf(`
resource "aws_budgets_budget" "test" {
  name_prefix  = %[1]q
  budget_type  = "SAVINGS_PLANS_UTILIZATION"
  limit_amount = "100"
  limit_unit   = "PERCENTAGE"
  time_unit    = "MONTHLY"
}
`, namePrefix)
}

func testAccAWSBudgetsBudgetConfigCostTypes(rName, startDate, endDate string) string {
	return fmt.Sprintf(`
resource "aws_budgets_budget" "test" {
  name         = %[1]q
  budget_type  = "COST"
  limit_amount = "456.78"
  limit_unit   = "USD"

  time_period_start = %[2]q
  time_period_end   = %[3]q
  time_unit         = "DAILY"

  cost_filter {
    name   = "AZ"
    values = [%[4]q, %[5]q]
  }

  cost_types {
    include_discount     = false
    include_subscription = true
    include_tax          = false
    use_blended          = true
  }
}
`, rName, startDate, endDate, acctest.Region(), acctest.AlternateRegion())
}

func testAccAWSBudgetsBudgetConfigCostTypesUpdated(rName, startDate, endDate string) string {
	return fmt.Sprintf(`
resource "aws_budgets_budget" "test" {
  name         = %[1]q
  budget_type  = "COST"
  limit_amount = "567.89"
  limit_unit   = "USD"

  time_period_start = %[2]q
  time_period_end   = %[3]q
  time_unit         = "DAILY"

  cost_filter {
    name   = "AZ"
    values = [%[4]q, %[5]q]
  }

  cost_types {
    include_credit       = false
    include_discount     = true
    include_refund       = false
    include_subscription = true
    include_tax          = true
    use_blended          = false
  }
}
`, rName, startDate, endDate, acctest.AlternateRegion(), acctest.ThirdRegion())
}

func testAccAWSBudgetsBudgetConfigNotifications(rName, emailAddress1, emailAddress2 string) string {
	return fmt.Sprintf(`
resource "aws_sns_topic" "test" {
  name = %[1]q
}

resource "aws_budgets_budget" "test" {
  name         = %[1]q
  budget_type  = "USAGE"
  limit_amount = "432.10"
  limit_unit   = "GBP"
  time_unit    = "ANNUALLY"

  notification {
    comparison_operator       = "GREATER_THAN"
    notification_type         = "ACTUAL"
    threshold                 = 150
    threshold_type            = "PERCENTAGE"
    subscriber_sns_topic_arns = [aws_sns_topic.test.arn]
  }

  notification {
    comparison_operator        = "EQUAL_TO"
    notification_type          = "FORECASTED"
    threshold                  = 200.10
    threshold_type             = "ABSOLUTE_VALUE"
    subscriber_email_addresses = [%[2]q, %[3]q]
  }
}
`, rName, emailAddress1, emailAddress2)
}

func testAccAWSBudgetsBudgetConfigNotificationsUpdated(rName, emailAddress1 string) string {
	return fmt.Sprintf(`
resource "aws_budgets_budget" "test" {
  name         = %[1]q
  budget_type  = "USAGE"
  limit_amount = "432.10"
  limit_unit   = "GBP"
  time_unit    = "ANNUALLY"

  notification {
    comparison_operator        = "LESS_THAN"
    notification_type          = "ACTUAL"
    threshold                  = 123.45
    threshold_type             = "ABSOLUTE_VALUE"
    subscriber_email_addresses = [%[2]q]
  }
}
`, rName, emailAddress1)
}
