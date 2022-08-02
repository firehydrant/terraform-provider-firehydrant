---
name: Bug report
about: Let us know about an unexpected error, a crash, or an incorrect behavior.
labels: bug, new

---

<!--
Thank you for opening an issue! Please note that we try to keep this issue tracker reserved for bug reports 
and feature requests for the FireHydrant provider. For issues related to the FireHydrant platform, please 
visit [FireHydrant Support](https://support.firehydrant.com/hc/en-us).
-->

### Terraform Version
<!---
Run `terraform version` to show the version, and paste the result between the ``` marks below.
If you are not running the latest version of Terraform, please try upgrading because your 
issue may have already been fixed.
-->

```
...
```

### Terraform Configuration Files
<!--
Paste the relevant parts of your Terraform configuration between the ``` marks below.
-->

```terraform
...
```

### Debug Output
<!--
Full debug output can be obtained by running Terraform with the environment variable `TF_LOG=trace`. 
Please create a GitHub Gist containing the debug output. Please do _not_ paste the debug output in 
the issue, since debug output is long.

Debug output may contain sensitive information. Please review it before posting publicly.
-->

### Expected Behavior
<!--
What should have happened?
-->

### Actual Behavior
<!--
What actually happened?
-->

### Steps to Reproduce
<!--
Please list the full steps required to reproduce the issue.
-->

1. _Describe how to replicate the conditions_
1. _under which you experienced your issue_
1. _including example Terraform configs where necessary._

### Additional Context
<!--
Is there anything atypical about your situation that we should know? For example: 
is Terraform running in a wrapper script or in a CI system? Are you passing any 
unusual command line options or environment variables to opt-in to non-default behavior?
-->

### References
<!--
Are there any other GitHub issues (open or closed) or pull requests that should be linked here?
For example:
- #6017
-->