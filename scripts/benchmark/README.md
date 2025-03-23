# Benchmark Scripts

This directory contains scripts for benchmarking and load testing the Code Controller (CC).

## Scripts

- **benchmark.sh** - Basic benchmark for CC functionality
- **cc_benchmark.sh** - Benchmarks CC commands with timing information
- **cc_scale_test.sh** - Tests CC with multiple implementations and features
- **git_benchmark.sh** - Pure Git operations benchmark
- **git_scale_benchmark.sh** - Large-scale Git operations testing

## Usage

Run any benchmark script from the project root:

=== CC Benchmark Test ===
Testing Git operations performance with 10 implementations
Setting up test environment...
Initializing project...
Initializing project: benchmark-test
Project benchmark-test initialized successfully
Generating 10 implementations (this will take some time)...
Generating implementations for: Create a simple Hello World application
Frameworks: [nodejs python go java ruby php rust cpp typescript kotlin]
Count: 10
Parallel: true
Current branch is: master
Creating branch impl-nodejs-1742708623...
Generating implementation for nodejs... This may take a while.
Creating branch impl-python-1742708623...
Generating implementation for python... This may take a while.
Creating branch impl-go-1742708623...
Generating implementation for go... This may take a while.
Creating branch impl-java-1742708623...
Generating implementation for java... This may take a while.
Creating branch impl-ruby-1742708623...
Generating implementation for ruby... This may take a while.
Creating branch impl-php-1742708623...
Generating implementation for php... This may take a while.
Creating branch impl-rust-1742708623...
Generating implementation for rust... This may take a while.
Creating branch impl-cpp-1742708623...
Generating implementation for cpp... This may take a while.
Creating branch impl-typescript-1742708623...
Generating implementation for typescript... This may take a while.
Creating branch impl-kotlin-1742708623...
Generating implementation for kotlin... This may take a while.
Implementations generated successfully
Generation completed in 0 seconds
Number of branches created: 21
Testing branch switching performance...
accepts 1 arg(s), received 0

Benchmark results will be printed to stdout. Some scripts will ask if you want to
clean up test artifacts when finished.
