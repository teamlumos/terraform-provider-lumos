"""Dump the apiserver OpenAPI document including the flag-gated /v1 surface.

Every /v1/* route in the apiserver is registered with ``include_in_schema=False``
and gated at runtime by ``require_v1_enabled`` (LaunchDarkly flag
``external-v1-api-enabled``). ``app.openapi()`` therefore omits the whole V1 API,
which is why the checked-in specs contain no /v1 paths.

This script flips ``include_in_schema`` on the in-memory route objects only, so
the running server still advertises exactly what it did before. Nothing under
apiserver/ is modified.

Run from the provider repo:
    python3 scripts/dump-openapi.py --out original-openapi.json
"""

import argparse
import json
import os
import sys


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--out", default="original-openapi.json")
    parser.add_argument(
        "--apiserver",
        default=os.path.expanduser("~/develop/lumos/backend/python/services/apiserver"),
    )
    args = parser.parse_args()

    sys.path.insert(0, args.apiserver)
    from apiserver.main import app  # noqa: E402

    forced = 0
    for route in app.routes:
        if getattr(route, "path", "").startswith("/v1/") and not getattr(
            route, "include_in_schema", True
        ):
            route.include_in_schema = True
            forced += 1

    app.openapi_schema = None
    doc = app.openapi()

    with open(args.out, "w") as handle:
        handle.write(json.dumps(doc, sort_keys=True, indent=2) + "\n")

    v1 = [p for p in doc["paths"] if p.startswith("/v1/")]
    print(f"forced {forced} gated routes visible", file=sys.stderr)
    print(f"wrote {args.out}: {len(doc['paths'])} paths, {len(v1)} under /v1", file=sys.stderr)


if __name__ == "__main__":
    main()
