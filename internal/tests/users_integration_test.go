package tests

import (
	"testing"
)

func TestUserRoutes(t *testing.T) {
  t.Log(`Integration tests for user routes...`)

  tests := []testCases{
    // POST
    {"POST - Create a new user - Success Case", userPostSuccess},
  }

  runTestCasesInParallel(t, tests)
}

func userPostSuccess(t *testing.T) {
  // Test case data
  tc := struct {
    route string
  }{
    route: "/users",
  }

  t.Error(tc)
}
