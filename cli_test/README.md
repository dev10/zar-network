# Zar CLI Integration tests

The zar cli integration tests live in this folder. You can run the full suite by running:

```bash
go test -mod=readonly -p 4 `go list ./cli_test/...` -tags=cli_test
```

> NOTE: While the full suite runs in parallel, some of the tests can take up to a minute to complete

### Test Structure

This integration suite [uses a thin wrapper](https://godoc.org/github.com/cosmos/cosmos-sdk/tests) over the [`os/exec`](https://golang.org/pkg/os/exec/) package. This allows the integration test to run against built binaries (both `zard` and `zarcli` are used) while being written in golang. This allows tests to take advantage of the various golang code we have for operations like marshal/unmarshal, crypto, etc...

> NOTE: The tests will use whatever `zard` or `zarcli` binaries are available in your `$PATH`. You can check which binary will be run by the suite by running `which zard` or `which zarcli`. If you have your `$GOPATH` properly setup they should be in `$GOPATH/bin/zar*`. This will ensure that your test uses the latest binary you have built

Tests generally follow this structure:

```go
func TestMyNewCommand(t *testing.T) {
  t.Parallel()
	f := InitFixtures(t)

	// start zard server
	proc := f.GDStart()
	defer proc.Stop(false)

  // Your test code goes here...

	f.Cleanup()
}
```

This boilerplate above:

- Ensures the tests run in parallel. Because the tests are calling out to `os/exec` for many operations these tests can take a long time to run.
- Creates `.zard` and `.zarcli` folders in a new temp folder.
- Uses `zarcli` to create 2 accounts for use in testing: `foo` and `bar`
- Creates a genesis file with coins (`1000footoken,1000feetoken,150stake`) controlled by the `foo` key
- Generates an initial bonding transaction (`gentx`) to make the `foo` key a validator at genesis
- Starts `zard` and stops it once the test exits
- Cleans up test state on a successful run

### Notes when adding/running tests

- Because the tests run against a built binary, you should make sure you build every time the code changes and you want to test again, otherwise you will be testing against an older version. If you are adding new tests this can easily lead to confusing test results.
- The [`test_helpers.go`](./test_helpers.go) file is organized according to the format of `zarcli` and `zard` commands. There are comments with section headers describing the different areas. Helper functions to call CLI functionality are generally named after the command (e.g. `zarcli query staking validator` would be `QueryStakingValidator`). Try to keep functions grouped by their position in the command tree.
- Test state that is needed by `tx` and `query` commands (`home`, `chain_id`, etc...) is stored on the `Fixtures` object. This makes constructing your new tests almost trivial.
- Sometimes if you exit a test early there can be still running `zard` and `zarcli` processes that will interrupt subsequent runs. Still running `zarcli` processes will block access to the keybase while still running `zard` processes will block ports and prevent new tests from spinning up. You can ensure new tests spin up clean by running `pkill -9 zard && pkill -9 zarcli` before each test run.
- Most `query` and `tx` commands take a variadic `flags` argument. This pattern allows for the creation of a general function which is easily modified by adding flags. See the `TxSend` function and its use for a good example.
- `Tx*` functions follow a general pattern and return `(success bool, stdout string, stderr string)`. This allows for easy testing of multiple different flag configurations. See `TestZarCLICreateValidator` or `TestZarCLISubmitProposal` for a good example of the pattern.
