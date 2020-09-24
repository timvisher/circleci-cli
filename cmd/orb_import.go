package cmd

import (
"fmt"
  "strings"
	"github.com/CircleCI-Public/circleci-cli/api"
	"github.com/CircleCI-Public/circleci-cli/api/graphql"
  )

func refsToImport(opts orbOptions) []string, error {

}

type OrbImportPlan struct {
    newNamespaces []string
    newOrbs []string
    newVersions []string
    alreadyExistingVersions []string
}

// $ circleci orb import aengelberg

// Would you like to:
// Create namespace `aengelberg`
// Create orb `aengelberg/my-orb`
// Create orb `aengelberg/my-orb2`
// Import version `aengelberg/my-orb@1.2.3`
// Import version `aengelberg/my-orb2@4.5.6`
// Import version `aengelberg/my-orb3@5.6.7`
// `aengelberg/my-orb4@5.6.7` already exists

func importPlan(opts orbOptions, refs []string) OrbImportPlan, error {

}

func applyPlan(opts orbOptions, plan OrbImportPlan) error {

}


func importOrb(opts orbOptions) error {
	ref := opts.args[0] // TODO: might just be a namespace
    cloudClient := client.NewClient("https://circleci.com", "graphql-unstable", "", opts.config.Debug)
	cloudVersion, err := api.OrbInfo(cloudClient, ref)
	if err != nil {
	  return err
	}
	refWithoutVersion := strings.split(ref, "@")[0]
    orbComponents := strings.split(ref, "/")
    orbNs, orbName := orbComponents[0], orbComponents[1]
	latestServerVersion, err := api.OrbInfo(opts.cl, ref)
	if err != nil {
	  if ok := err.(api.ErrOrbVersionNotExists); !ok {
	    return err
	  }
	}
	if latestServerVersion == nil {
	  err = api.CreateOrb(opts.cl, orbNs, orbName)
	  if err != nil {
	    return err
	  }
	}
}