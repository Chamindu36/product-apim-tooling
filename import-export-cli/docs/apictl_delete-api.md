## apictl delete-api

Delete API

### Synopsis

Delete an API from an environment

```
apictl delete-api (--name <name-of-the-api> --version <version-of-the-api> --provider <provider-of-the-api> --environment <environment-from-which-the-api-should-be-deleted>) [flags]
```

### Examples

```
apictl delete-api -n TwitterAPI -v 1.0.0 -r admin -e dev
apictl delete-api -n FacebookAPI -v 2.1.0 -e production
NOTE: All the 3 flags (--name (-n), --version (-v), and --environment (-e)) are mandatory.
If the --provider (-r) is not specified, the logged-in user will be considered as the provider.
```

### Options

```
  -e, --environment string   Environment from which the API should be deleted
  -h, --help                 help for delete-api
  -n, --name string          Name of the API to be deleted
  -r, --provider string      Provider of the API
  -v, --version string       Version of the API to be deleted
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications

