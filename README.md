# Go-ToolbOX
Chinese Documentation - [中文文档](./README-ZH.md)

> Go-toolbOx is characterized by daily work requirements and extension development, encapsulated generic tool classes

[![stable](https://img.shields.io/badge/stable-stable-green.svg)](https://github.com/kamalyes/go-toolbox)
[![license](https://img.shields.io/github/license/kamalyes/go-toolbox)]()
[![download](https://img.shields.io/github/downloads/kamalyes/go-toolbox/total)]()
[![release](https://img.shields.io/github/v/release/kamalyes/go-toolbox)]()
[![commit](https://img.shields.io/github/last-commit/kamalyes/go-toolbox)]()
[![issues](https://img.shields.io/github/issues/kamalyes/go-toolbox)]()
[![pull](https://img.shields.io/github/issues-pr/kamalyes/go-toolbox)]()
[![fork](https://img.shields.io/github/forks/kamalyes/go-toolbox)]()
[![star](https://img.shields.io/github/stars/kamalyes/go-toolbox)]()
[![go](https://img.shields.io/github/go-mod/go-version/kamalyes/go-toolbox)]()
[![size](https://img.shields.io/github/repo-size/kamalyes/go-toolbox)]()
[![contributors](https://img.shields.io/github/contributors/kamalyes/go-toolbox)]()
[![codecov](https://codecov.io/gh/kamalyes/go-toolbox/branch/master/graph/badge.svg)](https://codecov.io/gh/kamalyes/go-toolbox)
[![Go Report Card](https://goreportcard.com/badge/github.com/kamalyes/go-toolbox)](https://goreportcard.com/report/github.com/kamalyes/go-toolbox)
[![Go Reference](https://pkg.go.dev/badge/github.com/kamalyes/go-toolbox?status.svg)](https://pkg.go.dev/github.com/kamalyes/go-toolbox?tab=doc)
[![Sourcegraph](https://sourcegraph.com/github.com/kamalyes/go-toolbox/-/badge.svg)](https://sourcegraph.com/github.com/kamalyes/go-toolbox?badge)

**Go-toolbOx's key features are:**

- **Convert**: Conversion between data types, such as converting strings to integers or changing date formats from one form to another.

- **Desensitize**: Remove or obfuscate sensitive information to prevent data leakage, such as removing Personally Identifiable Information (PII) or encrypting data.

- **CRC**: Cyclic Redundancy Check for error detection in data transmission.

- **Error Handling**: Enhanced error handling capabilities to simplify error management.

- **HTTP Extensions**: Auxiliary tools for HTTP requests and responses.

- **Image Processing**: Tools for image processing and manipulation.

- **JSON Handling**: Lightweight data interchange format handling.

- **Location Services**: Information related to IP regions and more.

- **Math Extensions**: Extended functionalities for numerical computations.

- **Time Handling**: Parsing, validating, manipulating, and displaying dates and times to simplify date and time handling.

- **OS Interface**: Programming interfaces for interacting with the operating system.

- **Queue**: Implementation of queue data structures.

- **Random Numbers**: Random number generators suitable for various applications.

- **Retry Mechanism**: The process of retrying an operation when it fails, commonly used for network requests and database operations to enhance system reliability.

- **Scheduling**: Task scheduling tools that support the execution of timed tasks.

- **Signature**: Verification of data integrity and origin, used for validating data integrity and origin, such as word signatures and message signatures.

- **SQL Builder**: Tools for constructing SQL queries.

- **String Handling**: Extended functionalities for string processing, providing features like formatting, splitting, and concatenation.

- **Synchronization Tools**: Tools for synchronization in concurrent programming.

- **Types**: Definitions and operations for various types.

- **Unit Conversion**: Tools for converting between units.

- **User Agent**: Tools for handling user agent strings.

- **UUID**: Generation of universally unique identifiers (UUID).

- **Validator**: Tools for validating data integrity, such as form validation and data format validation to ensure input data meets expected formats and rules.

- **Compression Tools**: Tools related to data compression and decompression.

## Getting started

### Prerequisites

requires [Go](https://go.dev/) version [1.20](https://go.dev/doc/devel/release#go1.20.0) or above.

### Getting

With [Go's module support](https://go.dev/wiki/Modules#how-to-use-modules), `go [build|run|test]` automatically fetches the necessary dependencies when you add the import in your code:

```sh
import "github.com/kamalyes/go-toolbox"
```

Alternatively, use `go get`:

```sh
go get -u github.com/kamalyes/go-toolbox
```
