package cmd

import (
	"github.com/CircleCI-Public/circleci-cli/api"
	"github.com/CircleCI-Public/circleci-cli/api/graphql"
  )

type OrbImportPlan struct {
    newNamespaces []string
    newOrbs []*api.Orb
    newVersions []*api.OrbVersion
    alreadyExistingVersions []*api.OrbVersion
}

// $ circleci orb import aengelberg

// Would you like to:
//  Create namespace `aengelberg`
//  Create orb `aengelberg/my-orb`
//  Create orb `aengelberg/my-orb2`
//  Import version `aengelberg/my-orb@1.2.3`
//  Import version `aengelberg/my-orb2@4.5.6`
//  Import version `aengelberg/my-orb3@5.6.7`
//  (`aengelberg/my-orb4@5.6.7` already exists)

func versionsToImport(opts orbOptions) ([]*api.OrbVersion, error) {
	ref := opts.args[0]
	// TODO: if only a namespace is passed in, fetch all orbs within namespace
	// TODO: support an `--all-versions` flag that gets all versions instead of latest version per orb?
    cloudClient := graphql.NewClient("https://circleci.com", "graphql-unstable", "", opts.cfg.Debug)
	version, err := api.OrbInfo(cloudClient, ref)
	if err != nil {
	    return nil, err
	}
	return []*api.OrbVersion{version}, nil
}

func importPlan(opts orbOptions, refs []string) (OrbImportPlan, error) {
    // TODO
    return OrbImportPlan{}, nil
}

func applyPlan(opts orbOptions, plan OrbImportPlan) error {
    // TODO
    return nil
}


func importOrb(opts orbOptions) error {
    return nil
}