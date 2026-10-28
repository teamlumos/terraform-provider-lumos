# Snapshot tests

Snapshot tests for the provider, run with the real `terraform` CLI against a
real Lumos API. Each case creates its own objects, walks them through a series
of configs, and destroys them. Every plan along the way is compared with a
snapshot, so a plan that changes more than it should, shows values as
`(known after apply)` when nothing is changing, or isn't empty right after an
apply shows up as a diff.

## Running

Requirements: Terraform 1.12.x (the snapshots were recorded with it, and other
versions may word plans differently), Perl and curl, which macOS and most
Linux distributions include, and jq.

```sh
export LUMOS_API_TOKEN=...                  # an API token for the org
export LUMOS_API_URL=https://api.lumos.com
./snapshot_tests/run.sh                     # every case
./snapshot_tests/run.sh 'lumos_app/*'       # one resource's cases
./snapshot_tests/run.sh lumos_app/custom_app_lifecycle
```

Optional variables, for cases that can't create what they need themselves.
A case that needs one that isn't set is skipped.

| Variable | What it is |
|---|---|
| `GROUP_ID` | A group in the org, for requester groups, app admins and assignees. |
| `LUMOS_TEST_TBA_LABEL` | The label of a time-based access option defined in the org. |

Everything a case creates has `bb-<run id>` in its name. When a case fails, the
runner still destroys what it created, keeps the case's workdir (with the full
`terraform.log`) and prints its path. If the destroy itself fails, objects
named with that run id are left in Lumos.

### Which provider build runs

`shared/versions.tf` pins the prerelease under test, and `terraform init`
downloads it from the registry.

To run a local build instead, build it and point Terraform at it with
`dev_overrides` in a CLI config file of its own, so the override doesn't apply
to other Terraform workspaces:

```sh
go build -o ~/terraform-providers/terraform-provider-lumos .
```

```hcl
# ~/terraform-providers/local-provider.tfrc
provider_installation {
  dev_overrides {
    "registry.terraform.io/teamlumos/lumos" = "/Users/<you>/terraform-providers"
  }
  direct {}
}
```

```sh
TF_CLI_CONFIG_FILE=~/terraform-providers/local-provider.tfrc ./snapshot_tests/run.sh
```

The runner sees the override and skips `terraform init`.

## Cases

A case is a folder under `cases/<resource>/`:

```
cases/lumos_app/custom_app_lifecycle/
  01-create.tf
  02-update-label.tf
  ...
  snapshots/
```

Each numbered file is the whole config for that step, not a change to the
previous one. The runner copies `shared/` into a fresh directory, then for each
step:

1. plans and compares the plan with `snapshots/<step>.plan.txt`;
2. applies that plan;
3. plans again, and fails if that plan has any changes;
4. fails if any boolean output isn't `true`, so a case can check what a data
   source reads back with an output such as
   `value = data.lumos_app.this.label == lumos_app.this.label`.

A step that is expected to leave a diff behind, such as a value written in a
form Lumos reads back differently, has `snapshots/<step>.replan.txt` with that
diff, and the plan after its apply is compared with it instead. `--update`
records such a diff whenever it sees one and warns about each, so delete any
that aren't a known, intended diff: otherwise a bug is recorded as expected.

A step whose file has a `# expect: error` line must fail at plan or apply. If
its plan succeeds, the plan is compared with `snapshots/<step>.plan.txt` as
usual; the error, from Terraform's first `Error:` line on, is compared with
`snapshots/<step>.error.txt`. The case fails if the step succeeds. The next
step starts from the state the error left.

A step can also be a shell script, `<step>.sh`, which changes objects in Lumos
outside Terraform, as someone editing in the Lumos UI would. It runs in the
case's workdir after sourcing `drift.sh`, whose `state_id` prints a resource's
id from state and whose `api` calls the API. The step fails if any command in
it fails. It has no snapshot of its own: the next step's plan shows the drift.

After the last step it plans the destroy from the last config that applied,
compares it with `snapshots/destroy.plan.txt`, and destroys.

In snapshots, UUIDs read as `<uuid-1>`, `<uuid-2>`, ... in the order they first
appear in the case, the run id reads as `<run>`, and the HTTP version in an API
error reads as `HTTP/<version>`. Terraform prints a set's elements sorted by
value, which for elements holding new UUIDs differs from run to run, so the
elements of the provider's set attributes are sorted again by their text with
the UUIDs numbered. The runner reads which attributes are sets from the
provider's schema.
