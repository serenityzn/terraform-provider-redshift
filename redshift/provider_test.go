package redshift

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	testAccProviders map[string]*schema.Provider
	testAccProvider  *schema.Provider
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"redshift": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	var host string
	if host = os.Getenv("REDSHIFT_HOST"); host == "" {
		t.Fatal("REDSHIFT_HOST must be set for acceptance tests")
	}
	if v := os.Getenv("REDSHIFT_USER"); v == "" {
		t.Fatal("REDSHIFT_USER must be set for acceptance tests")
	}
}

func initTemporaryCredentialsProvider(t *testing.T, provider *schema.Provider) {
	clusterIdentifier := getEnvOrSkip("REDSHIFT_TEMPORARY_CREDENTIALS_CLUSTER_IDENTIFIER", t)

	sdkClient, err := stsClient(t)
	if err != nil {
		t.Skip(fmt.Sprintf("Unable to load STS client due to: %s", err))
	}

	response, err := sdkClient.GetCallerIdentity(context.TODO(), nil)
	if err != nil {
		t.Skip(fmt.Sprintf("Unable to get current STS identity due to: %s", err))
	}
	if response == nil {
		t.Skip("Unable to get current STS identity. Empty response.")
	}

	config := map[string]interface{}{
		"temporary_credentials": []interface{}{
			map[string]interface{}{
				"cluster_identifier": clusterIdentifier,
			},
		},
	}
	if arn, ok := os.LookupEnv("REDSHIFT_TEMPORARY_CREDENTIALS_ASSUME_ROLE_ARN"); ok {
		config["temporary_credentials"].([]interface{})[0].(map[string]interface{})["assume_role"] = []interface{}{
			map[string]interface{}{
				"arn": arn,
			},
		}
	}
	diagnostics := provider.Configure(context.Background(), terraform.NewResourceConfigRaw(config))
	if diagnostics != nil {
		if diagnostics.HasError() {
			t.Fatalf("Failed to configure temporary credentials provider: %v", diagnostics)
		}
	}
}

func stsClient(t *testing.T) (*sts.Client, error) {
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	return sts.NewFromConfig(config), nil
}

func TestAccRedshiftTemporaryCredentials(t *testing.T) {
	provider := Provider()
	assume_role_arn := os.Getenv("REDSHIFT_TEMPORARY_CREDENTIALS_ASSUME_ROLE_ARN")
	defer os.Setenv("REDSHIFT_TEMPORARY_CREDENTIALS_ASSUME_ROLE_ARN", assume_role_arn)
	os.Unsetenv("REDSHIFT_TEMPORARY_CREDENTIALS_ASSUME_ROLE_ARN")
	prepareRedshiftTemporaryCredentialsTestCases(t, provider)
	client, ok := provider.Meta().(*Client)
	if !ok {
		t.Fatal("Unable to initialize client")
	}
	db, err := client.Connect()
	if err != nil {
		t.Fatalf("Unable to connect to database: %s", err)
	}
	defer db.Close()
}

func TestAccRedshiftTemporaryCredentialsAssumeRole(t *testing.T) {
	_ = getEnvOrSkip("REDSHIFT_TEMPORARY_CREDENTIALS_ASSUME_ROLE_ARN", t)
	provider := Provider()
	prepareRedshiftTemporaryCredentialsTestCases(t, provider)
	client, ok := provider.Meta().(*Client)
	if !ok {
		t.Fatal("Unable to initialize client")
	}
	db, err := client.Connect()
	if err != nil {
		t.Fatalf("Unable to connect to database: %s", err)
	}
	defer db.Close()
}

func prepareRedshiftTemporaryCredentialsTestCases(t *testing.T, provider *schema.Provider) {
	redshift_password := os.Getenv("REDSHIFT_PASSWORD")
	defer os.Setenv("REDSHIFT_PASSWORD", redshift_password)
	os.Unsetenv("REDSHIFT_PASSWORD")
	rawUsername := os.Getenv("REDSHIFT_USER")
	defer os.Setenv("REDSHIFT_USER", rawUsername)
	username := strings.ToLower(permanentUsername(rawUsername))
	os.Setenv("REDSHIFT_USER", username)
	initTemporaryCredentialsProvider(t, provider)
}

// TestAccRedshiftServerlessTemporaryCredentials tests the full connection flow
// for Redshift Serverless using GetCredentials, and validates that all SQL
// queries used by the provider actually execute successfully against a real
// Serverless cluster. Requires:
//
//	REDSHIFT_HOST                                      - Serverless endpoint
//	REDSHIFT_DATABASE                                  - Database name
//	REDSHIFT_TEMPORARY_CREDENTIALS_SERVERLESS_WORKGROUP - Workgroup name
//	REDSHIFT_TEMPORARY_CREDENTIALS_SERVERLESS_REGION   - AWS region (optional)
func TestAccRedshiftServerlessTemporaryCredentials(t *testing.T) {
	workgroupName := getEnvOrSkip("REDSHIFT_TEMPORARY_CREDENTIALS_SERVERLESS_WORKGROUP", t)
	region := os.Getenv("REDSHIFT_TEMPORARY_CREDENTIALS_SERVERLESS_REGION")

	provider := Provider()

	cfg := map[string]interface{}{
		"host":     os.Getenv("REDSHIFT_HOST"),
		"database": os.Getenv("REDSHIFT_DATABASE"),
		"type":     "serverless",
		"username": "unused", // ignored by serverless GetCredentials
		"temporary_credentials_serverless": []interface{}{
			map[string]interface{}{
				"workgroup_name": workgroupName,
				"region":         region,
			},
		},
	}

	diagnostics := provider.Configure(context.Background(), terraform.NewResourceConfigRaw(cfg))
	if diagnostics != nil && diagnostics.HasError() {
		t.Fatalf("Failed to configure serverless provider: %v", diagnostics)
	}

	client, ok := provider.Meta().(*Client)
	if !ok {
		t.Fatal("Unable to initialize client")
	}

	db, err := client.Connect()
	if err != nil {
		t.Fatalf("Unable to connect to Redshift Serverless: %s", err)
	}
	defer db.Close()

	// Validate the exact queries the provider runs on the serverless path.
	// These catch column-not-found errors before they reach production.
	queries := []struct {
		name string
		sql  string
	}{
		{
			name: "pg_user basic columns",
			sql:  "SELECT usesysid, usecreatedb, usesuper FROM pg_user LIMIT 1",
		},
		{
			name: "pg_user valuntil",
			sql:  "SELECT COALESCE(valuntil, 'infinity') FROM pg_user LIMIT 1",
		},
		{
			name: "pg_user current user lookup",
			sql:  "SELECT usesysid, usecreatedb, usesuper FROM pg_user WHERE usename = CURRENT_USER",
		},
		{
			name: "schema owner join",
			sql: `SELECT trim(svv_all_schemas.schema_name), trim(pg_user.usename)
				  FROM svv_all_schemas
				  LEFT JOIN pg_user ON pg_user.usesysid = svv_all_schemas.schema_owner
				  WHERE svv_all_schemas.database_name = current_database()
				  LIMIT 1`,
		},
	}

	for _, q := range queries {
		t.Run(q.name, func(t *testing.T) {
			rows, err := db.Query(q.sql)
			if err != nil {
				t.Fatalf("Query %q failed: %s\nSQL: %s", q.name, err, q.sql)
			}
			rows.Close()
		})
	}
}
