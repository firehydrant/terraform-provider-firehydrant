# Tests

Running the tests for this provider requires access to FireHydrant. Some resources or data sources
may require access to the Enterprise tier of FireHydrant to run successfully.

## 1. Set up necessary environment variables

1. `FIREHYDRANT_API_KEY` - (Required) A [bot token](https://support.firehydrant.com/hc/en-us/articles/360057722832-Creating-a-Bot-User)
   to use for testing in FireHydrant.
2. `FIREHYDRANT_BASE_URL` - (Optional) The FireHydrant API URL to connect to for testing.
   Defaults to `https://api.firehydrant.io/v1/`

You can set your environment variables using whatever method you'd like.
The following are instructions for setting up environment variables using [envchain](https://github.com/sorah/envchain).

1. Make sure you have envchain installed.
   [Instructions for this can be found in the envchain README](https://github.com/sorah/envchain#installation).
2. Pick a namespace for storing your environment variables. I suggest `terraform-provider-firehydrant`.
3. For each environment variable you need to set, run the following command:
   ```sh
   envchain --set YOUR_NAMESPACE_HERE ENVIRONMENT_VARIABLE_HERE
   ```
   **OR**

   Set all of the environment variables at once with the following command:
   ```sh
   envchain --set YOUR_NAMESPACE_HERE FIREHYDRANT_BASE_URL FIREHYDRANT_API_KEY
   ```

## 2. Run the tests

### Running all acceptance tests

#### With envchain:
```sh
$ envchain YOUR_NAMESPACE_HERE make testacc
```

#### Without envchain:
```sh
$ make testacc
```

### Running specific acceptance tests

The commands below use task lists as an example.

#### With envchain:
```sh
$ TESTARGS="-run TestAccTaskList" envchain YOUR_NAMESPACE_HERE make testacc
```

#### Without envchain:
```sh
$ TESTARGS="-run TestAccTaskList" make testacc
```

### Running the all non-acceptance tests

#### With envchain:
```sh
$ envchain YOUR_NAMESPACE_HERE make test
```

#### Without envchain:
```sh
$ make test
```
