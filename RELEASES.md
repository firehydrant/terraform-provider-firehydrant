# Releases

Releases are typically run once we have a few user facing changes ready to be released. It is
generally good to batch up changes together but urgent bug fixes or important new features
that need to go out quickly are good reasons to release earlier. 

Things like test fixes or updates to the README do not require a release.

Any updates to documentation in the `/docs` folder require a release in order to show up
on registry.terraform.io.

## Preparing for a release

1. Compare the `main` branch with the previous release to see all the changes you will be
   releasing by looking at the following in GitHub (where vX.Y.Z represents the previous
   release version):
   ```
   https://github.com/firehydrant/terraform-provider-firehydrant/compare/vX.Y.Z...main
   ```
   Make sure any major customer-facing changes have an entry in the changelog.
2. Review the changelog to check for the following:
   - The changelog follows [the changelog format](./README.md#updating-the-changelog)
   - Information is accurate and free of typos 
   - All links work 
   - If there are any breaking changes, a "Notes" section has been added with 
     information on how to upgrade to the new version.
3. Do a final run-through to test all the resources and data sources that have changed.
   - Initialize and create resources and data sources using the previous provider version.
     You can check that you are using the previous version by running `terraform version`.
     The output should look something like this (where vX.Y.Z represents the previous
     release version):
     ```
     $ terraform version
     Terraform v1.2.1
     on darwin_arm64
     + provider registry.terraform.io/firehydrant/firehydrant vX.Y.Z
     ```
   - Upgrade your Terraform version to a local build of the provider.
     You can check that you are using the local build by looking at the output
     while running various Terraform commands. 
     The output should look something like this:
     ```
     $ terraform plan
     │ Warning: Provider development overrides are in effect
     │
     │ The following provider development overrides are set in the CLI configuration:
     │  - firehydrant/firehydrant in /PATH/TO/YOUR/OVERRIDE/DIRECTORY
     ```
   - Test that running `terraform plan` with your local provider build shows no changes and 
     doesn't introduce any unexpected breaking changes.
4. Decide on an appropriate version for your release. In general, we follow the [Terraform
   provider best practices for versioning](https://www.terraform.io/plugin/sdkv2/best-practices/versioning),
   which follow the guidelines of [Semantic Versioning](https://semver.org/).
   - Breaking changes should be very rare but if they are necessary, they should always be a minor version
     increment, as this provider is still at version 0.
   - New resources and data sources and other large features should generally be a minor version increment.
   - Small bug fixes, documentation fixes, and other changes that leave the provider functionally equivalent 
     to the previous version should generally be a patch version increment.
5. Open a PR that removes the "(Unreleased)" marker from the changelog version and introduces any other changes 
   you need to make to the changelog before release.
6. Merge this PR before starting the release.

## Creating a release

A new release is kicked off by adding a lightweight tag to the repository. 

1. Checkout the `main` branch and make sure you have pulled down the latest changes.
2. Make sure the changelog is ready to go and has the version you want to release 
   and the "(Unreleased)" marker has been removed.
3. Create a new [lightweight tag](https://git-scm.com/book/en/v2/Git-Basics-Tagging#_lightweight_tags)
   in the format `vX.Y.Z` that matches the version in the changelog by running:. 
   ```
   git tag vX.Y.Z
   ```
4. Check that your tag has been added to the latest commit by running:
   ```
   git show vX.Y.Z
   ```
5. Once you are certain you've tagged the correct commit, push your tag to the `main` branch
   to kick off the release by running:
   ```
   git push origin vX.Y.Z
   ```
6. You did it! A new run of the GitHub Actions release workflow should have been kicked off,
   which you can monitor at https://github.com/firehydrant/terraform-provider-firehydrant/actions/workflows/release.yml.
   - If any errors occur during the release, see the [Troubleshooting a release section](#troubleshooting-a-release) below. 
7. Once the release has succeeded, check that you can pull down the latest version and run `terraform apply` with
   a config that uses all the resources and data sources that have changed.
8. Check that the [documentation](https://registry.terraform.io/providers/firehydrant/firehydrant/latest/docs) has
   successfully updated. This may take a few minutes to show up.

## Cleaning up after a release

1. Open a PR adding a new version for the next release to the changelog with the 
   "(Unreleased)" marker. The new version should increment the patch number by one.
   The exact version you use doesn't really matter too much though, as it can be changed 
   before the next release. 

## Troubleshooting a release

It is rare but occasionally a release can fail. When that happens, there are a few different things
you can try to figure out what's going on and fix it.

Documentation on [why releases are set up this way can be found here](https://www.terraform.io/registry/providers/publishing#github-actions-preferred) 
and the [scaffolding repository this is based on can be found here](https://github.com/hashicorp/terraform-provider-scaffolding/blob/main/.github/workflows/release.yml). 

1. Check the error output. Often this will explain what went wrong and help you determine if you should
   rerun the release workflow or do something else to fix it. 
2. If the error message is unhelpful, identify which step of the release workflow is failing. Find this step
   in the [release workflow config file](./.github/workflows/release.yml) and search for the associated repo.
   Check through existing issues, search for the error you are getting in the source code, etc. 
3. If rerunning the release won't fix the problem you are seeing, you might need to 
   [delete the tag](https://git-scm.com/book/en/v2/Git-Basics-Tagging#_deleting_tags)
   while you work on fixing the release workflow by running `git push origin --delete vX.Y.Z`
   Once you have the workflow fixed, you can create and push the tag again to kick off a new release. 
