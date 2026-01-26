package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kinde-oss/kinde-go/kinde"
	"github.com/kinde-oss/kinde-go/kinde/management_api"
	"github.com/kinde-oss/kinde-go/oauth2/client_credentials"
)

// This is a test script to verify GetOrganizationConnections works correctly
// with the local SDK changes. It tests that connection properties are properly populated.
//
// Prerequisites:
// 1. Set up environment variables:
//    - KINDE_ISSUER_URL: Your Kinde issuer URL
//    - KINDE_CLIENT_ID: Your M2M client ID
//    - KINDE_CLIENT_SECRET: Your M2M client secret
//    - KINDE_ORGANIZATION_CODE: The organization code to test with
//
// 2. Use the test_local_sdk.sh script to set up the local SDK in your project
//
// 3. Run: go run test_get_organization_connections.go

func main() {
	// Get environment variables
	issuerURL := os.Getenv("KINDE_ISSUER_URL")
	clientID := os.Getenv("KINDE_CLIENT_ID")
	clientSecret := os.Getenv("KINDE_CLIENT_SECRET")
	organizationCode := os.Getenv("KINDE_ORGANIZATION_CODE")
	customAudience := os.Getenv("KINDE_AUDIENCE") // Optional: custom audience

	if issuerURL == "" || clientID == "" || clientSecret == "" || organizationCode == "" {
		log.Fatal(`
Missing required environment variables:
  - KINDE_ISSUER_URL: Your Kinde issuer URL
  - KINDE_CLIENT_ID: Your M2M client ID
  - KINDE_CLIENT_SECRET: Your M2M client secret
  - KINDE_ORGANIZATION_CODE: The organization code to test with

Example:
  export KINDE_ISSUER_URL="https://your-subdomain.kinde.com"
  export KINDE_CLIENT_ID="your_client_id"
  export KINDE_CLIENT_SECRET="your_client_secret"
  export KINDE_ORGANIZATION_CODE="org_abc123"
`)
	}

	ctx := context.Background()

	// Create client credentials flow
	// Note: The audience must be whitelisted in your Kinde M2M application settings
	var opts []client_credentials.Option
	
	if customAudience != "" {
		// Use custom audience if provided
		fmt.Printf("Using custom audience: %s\n", customAudience)
		opts = append(opts, client_credentials.WithAudience(customAudience))
	} else {
		// Use default Management API audience
		opts = append(opts, client_credentials.WithKindeManagementAPI(issuerURL))
	}
	
	opts = append(opts, client_credentials.WithTokenValidation(true))
	
	clientCredentialsFlow, err := client_credentials.NewClientCredentialsFlow(
		issuerURL,
		clientID,
		clientSecret,
		opts...,
	)
	if err != nil {
		log.Fatalf("Failed to create client credentials flow: %v", err)
	}

	// Create management API client
	kindeManagementAPI, err := kinde.NewManagementAPI(ctx, issuerURL, clientCredentialsFlow)
	if err != nil {
		log.Fatalf("Failed to create management API client: %v", err)
	}

	// Call GetOrganizationConnections
	fmt.Printf("Calling GetOrganizationConnections for organization: %s\n", organizationCode)
	resp, err := kindeManagementAPI.GetOrganizationConnections(ctx, management_api.GetOrganizationConnectionsParams{
		OrganizationCode: organizationCode,
	})
	if err != nil {
		// Check for OAuth audience error and provide helpful guidance
		errStr := err.Error()
		if strings.Contains(errStr, "audience") && strings.Contains(errStr, "whitelisted") {
			log.Fatalf(`
❌ OAuth Configuration Error: Audience not whitelisted

The error indicates that the Management API audience hasn't been whitelisted for your M2M application.

To fix this:
1. Go to your Kinde Dashboard: https://app.kinde.com
2. Navigate to: Settings > Applications > [Your M2M Application]
3. In the "Allowed audiences" section, add: %s/api
4. Save the changes
5. Try running the test again

Alternatively, if you have a different audience configured, you can:
- Set KINDE_AUDIENCE environment variable to use a custom audience
- Or modify the test to use client_credentials.WithAudience() instead of WithKindeManagementAPI()

Original error: %v
`, issuerURL, err)
		}
		log.Fatalf("Failed to get organization connections: %v", err)
	}

	// Handle error responses
	switch response := resp.(type) {
	case *management_api.GetOrganizationConnectionsBadRequest:
		fmt.Println("\n❌ API returned Bad Request (400)")
		// Cast to ErrorResponse to access GetErrors
		errResp := (*management_api.ErrorResponse)(response)
		errors := errResp.GetErrors()
		if len(errors) > 0 {
			for i, err := range errors {
				if code, ok := err.Code.Get(); ok {
					fmt.Printf("Error %d Code: %s\n", i+1, code)
				}
				if message, ok := err.Message.Get(); ok {
					fmt.Printf("Error %d Message: %s\n", i+1, message)
				}
			}
		}
		fmt.Println("\nPossible causes:")
		fmt.Println("  - Invalid organization code (check that it's correct)")
		fmt.Println("  - Organization doesn't exist")
		fmt.Println("  - Missing required permissions")
		fmt.Println("\nTo find valid organization codes:")
		fmt.Println("  1. Go to Kinde Dashboard: https://app.kinde.com")
		fmt.Println("  2. Navigate to: Organizations")
		fmt.Println("  3. Use the organization 'code' (not the name or ID)")
		log.Fatalf("Cannot proceed with invalid organization code: %s", organizationCode)
		
	case *management_api.GetOrganizationConnectionsForbidden:
		fmt.Println("\n❌ API returned Forbidden (403)")
		// Cast to ErrorResponse to access GetErrors
		errResp := (*management_api.ErrorResponse)(response)
		errors := errResp.GetErrors()
		if len(errors) > 0 {
			for i, err := range errors {
				if code, ok := err.Code.Get(); ok {
					fmt.Printf("Error %d Code: %s\n", i+1, code)
				}
				if message, ok := err.Message.Get(); ok {
					fmt.Printf("Error %d Message: %s\n", i+1, message)
				}
			}
		}
		fmt.Println("\nPossible causes:")
		fmt.Println("  - M2M application missing 'read:organization_connections' scope")
		fmt.Println("  - Insufficient permissions")
		fmt.Println("\nTo fix:")
		fmt.Println("  1. Go to Kinde Dashboard: https://app.kinde.com")
		fmt.Println("  2. Navigate to: Settings > Applications > [Your M2M Application]")
		fmt.Println("  3. Ensure 'read:organization_connections' scope is enabled")
		log.Fatalf("Access forbidden - check permissions")
		
	case *management_api.GetOrganizationConnectionsTooManyRequests:
		fmt.Println("\n⚠️  API returned Too Many Requests (429)")
		// Cast to ErrorResponse to access GetErrors
		errResp := (*management_api.ErrorResponse)(response)
		errors := errResp.GetErrors()
		if len(errors) > 0 {
			for i, err := range errors {
				if code, ok := err.Code.Get(); ok {
					fmt.Printf("Error %d Code: %s\n", i+1, code)
				}
				if message, ok := err.Message.Get(); ok {
					fmt.Printf("Error %d Message: %s\n", i+1, message)
				}
			}
		}
		fmt.Println("\nPlease wait a moment and try again.")
		log.Fatalf("Rate limit exceeded")
	}

	// Check response type - should be success response
	response, ok := resp.(*management_api.GetConnectionsResponse)
	if !ok {
		log.Fatalf("Unexpected response type: %T\n\nThis might indicate an API change or an unhandled error response.", resp)
	}

	// Verify the fix: Check that connection properties are populated
	fmt.Println("\n=== Response Analysis ===")
	
	if code, ok := response.Code.Get(); ok {
		fmt.Printf("Code: %s\n", code)
	}
	if message, ok := response.Message.Get(); ok {
		fmt.Printf("Message: %s\n", message)
	}

	fmt.Printf("\nNumber of connections: %d\n", len(response.Connections))

	if len(response.Connections) == 0 {
		fmt.Println("\n⚠️  No connections found. This is expected if the organization has no connections.")
		fmt.Println("   The fix ensures that when connections exist, their properties will be populated.")
		return
	}

	// Verify each connection has populated properties
	fmt.Println("\n=== Connection Details ===")
	allPropertiesPopulated := true

	for i, conn := range response.Connections {
		fmt.Printf("\nConnection %d:\n", i+1)
		
		hasID := conn.ID.IsSet()
		hasName := conn.Name.IsSet()
		hasDisplayName := conn.DisplayName.IsSet()
		hasStrategy := conn.Strategy.IsSet()

		if hasID {
			if id, ok := conn.ID.Get(); ok {
				fmt.Printf("  ✓ ID: %s\n", id)
			}
		} else {
			fmt.Printf("  ✗ ID: NOT SET (THIS INDICATES THE BUG)\n")
			allPropertiesPopulated = false
		}

		if hasName {
			if name, ok := conn.Name.Get(); ok {
				fmt.Printf("  ✓ Name: %s\n", name)
			}
		} else {
			fmt.Printf("  ✗ Name: NOT SET (THIS INDICATES THE BUG)\n")
			allPropertiesPopulated = false
		}

		if hasDisplayName {
			if displayName, ok := conn.DisplayName.Get(); ok {
				fmt.Printf("  ✓ Display Name: %s\n", displayName)
			}
		} else {
			fmt.Printf("  - Display Name: Not set (optional field)\n")
		}

		if hasStrategy {
			if strategy, ok := conn.Strategy.Get(); ok {
				fmt.Printf("  ✓ Strategy: %s\n", strategy)
			}
		} else {
			fmt.Printf("  ✗ Strategy: NOT SET (THIS INDICATES THE BUG)\n")
			allPropertiesPopulated = false
		}
	}

	fmt.Println("\n=== Test Result ===")
	if allPropertiesPopulated {
		fmt.Println("✅ SUCCESS: All connection properties are properly populated!")
		fmt.Println("   The fix is working correctly.")
	} else {
		fmt.Println("❌ FAILURE: Some connection properties are empty!")
		fmt.Println("   This indicates the bug still exists.")
		os.Exit(1)
	}
}
