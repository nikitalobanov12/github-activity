# GitHub Activity CLI

A command-line tool built in Go that fetches and displays the recent public activity of any GitHub user. This tool provides a clean, formatted view of user activities including commits, issues, pull requests, and repository events.

## Features

- **Real-time Activity**: Fetches live data from GitHub's public API
- **Smart Grouping**: Automatically groups consecutive push events to the same repository
- **Clickable Links**: Terminal hyperlinks for repositories and user profiles (in supported terminals)
- **Activity Types**: Supports various GitHub events including:
  - Push events (commits)
  - Issues (opened/closed)
  - Pull requests (opened/closed)
  - Repository creation/deletion
  - Stars and forks
  - Branch creation/deletion
  - Repository visibility changes
- **Error Handling**: Graceful handling of invalid usernames and API errors
- **Limit Control**: Shows the 10 most recent activities

## Installation

### Quick Install (Recommended)
If you have Go installed, you can install the tool directly:

```bash
go install github.com/nikitalobanov12/github-activity@latest
```

After installation, you can run it from anywhere:
```bash
github-activity <username>
```

### Alternative: Build from source
```bash
git clone https://github.com/nikitalobanov12/github-activity
cd github-activity
go build -o github-activity .
./github-activity <username>
```

### Prerequisites
- Go 1.24.5 or later

## Usage

```bash
github-activity <username>
```

### Examples

```bash
# Display activity for a specific user
github-activity nikitalobanov12

# Show help
github-activity -h
github-activity --help
```

### Sample Output

```
Fetching activity for user: nikitalobanov12

Recent activity for nikitalobanov12:
- Pushed 3 commits to user/repository-name
- Opened a new issue in user/another-repo
- Starred user/awesome-project
- Created repository user/new-project
- Closed a pull request in user/existing-repo
```

## How It Works

1. **API Integration**: Uses GitHub's public events API (`https://api.github.com/users/{username}/events`)
2. **Data Processing**: Parses JSON response and structures event data
3. **Smart Grouping**: Combines consecutive push events to the same repository
4. **Formatting**: Converts event data into human-readable descriptions
5. **Terminal Output**: Displays formatted results with clickable hyperlinks

## Error Handling

The tool handles various error scenarios:
- **Invalid Username**: Shows error message for non-existent users
- **Network Issues**: Reports connection failures
- **API Errors**: Handles GitHub API rate limits and other HTTP errors
- **Empty Activity**: Notifies when no recent activity is found

## Project Context

This project is part of the [roadmap.sh Backend Developer roadmap](https://roadmap.sh/projects/github-user-activity). It's designed to help developers practice:
- API integration
- JSON parsing
- Command-line tool development
- Error handling
- Data formatting and presentation

