# goload

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Overview

This Go modules is designed to watch for changes in a specific Go file and restart the application upon detecting changes. It utilizes the `fsnotify` package to monitor file system events and `os/exec` for restarting the application.

## Features

- **File Monitoring:** Monitors changes in the specified Go file.
- **Automatic Restart:** Automatically restarts the application upon detecting changes in the file.
- **Simple Usage:** Easy setup and usage with command-line arguments.

