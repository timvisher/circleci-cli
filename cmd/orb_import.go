package cmd

import (
	"fmt"
	"strings"

	"github.com/CircleCI-Public/circleci-cli/api"
	"github.com/CircleCI-Public/circleci-cli/api/graphql"
)

type OrbImportPlan struct {
	newNamespaces           []string
	newOrbs                 []*api.Orb
	newVersions             []*api.OrbVersion
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

func versionsToImport(opts orbOptions) ([]api.OrbVersion, error) {
	cloudClient := graphql.NewClient("https://circleci.com", "graphql-unstable", "", opts.cfg.Debug)

	var orbVersions []api.OrbVersion
	for _, ref := range opts.args {
		// If its not a namespace, fetch using api.OrbInfo -> append to list
		if !isNamespace(ref) {
			version, err := api.OrbInfo(cloudClient, ref)
			if err != nil {
				return nil, fmt.Errorf("orb info: %s", err.Error())
			}

			orbVersions = append(orbVersions, *version)
			continue
		}

		// TODO: support an `--all-versions` flag that gets all versions instead of latest version per orb?
		// Note: fetching all orb versions may not be possible. The best we could do is fetch an arbitrarily large number.
		// Otherwise, do some other operation that grabs orb source data from a single namespace.
		obv, err := api.ListNamespaceOrbVersions(cloudClient, ref)
		if err != nil {
			return nil, fmt.Errorf("list namespace orb versions: %s", err.Error())
		}

		orbVersions = append(orbVersions, obv...)
	}

	return orbVersions, nil
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

func isNamespace(ref) bool {
	if len(strings.Split(ref, "/")) > 1 {
		return false
	}
	return true
}
