# Simple

## Explanation

This test demonstrates the parser not getting confused by comment sequences.

## Why do we care about this test case

With the trivial regexp based parser that we had before, it got tripped up
by multiline comments that had @import statements on the inside.
