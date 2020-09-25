package cmd

import (
	"fmt"
	"strings"

	"github.com/CircleCI-Public/circleci-cli/api"
	"github.com/CircleCI-Public/circleci-cli/api/graphql"
)

type orbImportPlan struct {
	NewNamespaces           []string
	NewOrbs                 []api.Orb
	NewVersions             []api.OrbVersion
	AlreadyExistingVersions []api.OrbVersion
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

func importPlan(opts orbOptions, orbVersions []api.OrbVersion) (orbImportPlan, error) {
	uniqueNamespaces := map[string]bool{}
	uniqueOrbs := map[string]api.Orb{}

	// Dedupe namespaces and orbs.
	for _, o := range orbVersions {
		ns, orb := o.Orb.Namespace, o.Orb.Name
		uniqueNamespaces[ns] = true
		uniqueOrbs[orb] = o.Orb
	}

	var plan orbImportPlan
	for ns := range uniqueNamespaces {
		_, err := api.NamespaceExists(opts.cl, ns) // TODO: this implementation will change to include the notion of "imported"
		if err != nil {
			return orbImportPlan{}, fmt.Errorf("namespace check failed: %s", err.Error())
		}

		plan.NewNamespaces = append(plan.NewNamespaces, ns)
	}

	for _, orb := range uniqueOrbs {
		_, err := api.OrbInfo(opts.cl, fmt.Sprintf("%s/%s", orb.Namespace, orb.Name))
		if _, ok := err.(*api.ErrOrbVersionNotExists); ok {
			plan.NewOrbs = append(plan.NewOrbs, orb)
			continue
		}
		if err != nil {
			return orbImportPlan{}, fmt.Errorf("orb info check failed: %s", err.Error())
		}
	}

	for _, o := range orbVersions {
		_, err := api.OrbInfo(opts.cl, fmt.Sprintf("%s/%s@%s", o.Orb.Namespace, o.Orb.Name, o.Version))
		if _, ok := err.(*api.ErrOrbVersionNotExists); ok {
			plan.NewVersions = append(plan.NewVersions, o)
			continue
		}
		if err != nil {
			return orbImportPlan{}, fmt.Errorf("orb info check failed: %s", err.Error())
		}

		plan.AlreadyExistingVersions = append(plan.AlreadyExistingVersions, o)
	}

	return plan, nil
}

func applyPlan(opts orbOptions, plan orbImportPlan) error {

	return nil
}

func importOrb(opts orbOptions) error {
	// 1. versionsToImport
	// 2. generateImportPlan
	// 3. applyImportPlan

	return nil
}

func isNamespace(ref string) bool {
	if len(strings.Split(ref, "/")) > 1 {
		return false
	}
	return true
}
