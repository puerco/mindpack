# mindpack

Mindpack is a CLI utility to manage minder bundles. 

A bundle is an package that groups profiles and rule types. Minder uses
bundles to ship profiles together with its rules and to keep them up to date.

## Install

To install mindpack, clone this repository and run `go build`

```bash

git clone git@github.com:puerco/mindpack.git
cd mindpack
go build

## Test

./mindpack

mindpack: manage minder bundles

Usage:
  mindpack [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        initializes a mindpack source directory
  pack        writes a mindpack bundle to a distributable archive
  version     Prints the version

Flags:
  -h, --help   help for mindpack


```

## Usage

To start a bundle, create a new directory and put profiles and rule type
definition files in it:

```bash 
# Create the bundle dir structure:
mkdir mybundle/profiles
mkdir mybundle/rule_types

# Add profile and rule_type data
curl -o mybundle/profiles/branch-protection.yaml \
     https://raw.githubusercontent.com/stacklok/minder-rules-and-profiles/main/profiles/github/branch-protection.yaml

curl -o mybundle/rule_types/branch_protection_enabled.yaml \
     https://raw.githubusercontent.com/stacklok/minder-rules-and-profiles/main/rule-types/github/branch_protection_enabled.yaml

# Use mindpack to initialize the new bundle. This writes the new 
# bundle manifest:

mindpack init --source=mybundle --name=my-bundle --version=v0.1.0

# Pack the bundle in a new pacakge ready to ship:

mindpack pack --source=mybundle/ -f my-bundle-0.0.1.mpk

```

## Bundle Structure

A minder bundle is an archive that packs together profiles, rule types and a 
signed manifest describing the data. We have a [full specification of the
minder bundles](specification.md) but here is a short summary:

### Directory Structure

Minder bundles are built from filesystem sources with a specific structure
that packs together a manifest file, profiles and rule types. Here is a simple
example:

```
./bundles/branch-protection
├── manifest.json
├── profiles
│   └── branch-protection.yaml
└── rule_types
    └── branch_protection_enabled.yaml
```

### Manifest Structure

The bundle manifest is a json file that packs together metadata about the bundle
and a listing of the files to verify them.

```json
{
  "metadata": {
    "name": "branch-protection",
    "version": "v0.0.1",
    "date": "2024-03-08T01:11:57-06:00"
  },
  "files": {
    "profiles": [
      {
        "name": "branch-protection.yaml",
        "hashes": {
          "sha-256": "f3682a1cb5ab92c0cc71dd913338bf40a89ec324024f8d3f500be0e2aa4a9ae1"
        }
      }
    ],
    "ruleTypes": [
      {
        "name": "branch_protection_enabled.yaml",
        "hashes": {
          "sha-256": "10198b8cac16cd1d983a0a6fbb950816448f65e8f1d7a7407e2ff94949b42ccb"
        }
      }
    ]
  }
}

```


