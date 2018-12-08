# Test cases

## Definition

Test cases can be defined by creating a file inside the directory `simple` (for
the simple test suite) or by creating a new directory, and thus a new test
suite, in the `testdata` directory. This new directory should contain a set of
input files along with `BUILD.bazel` files (input designated as
`BUILD.bazel.in` and matching output files designated as `BUILD.bazel.out`).

NOTE: Do not make a `BUILD.bazel` or `BUILD` file in any directory or you will
introduce a boundary between that directory and testdata directory's glob. From
[globs
documentation](https://docs.bazel.build/versions/master/skylark/build-style.html#globs)
"Recursive globs make BUILD files difficult to reason about because they skip
subdirectories containing BUILD files."

