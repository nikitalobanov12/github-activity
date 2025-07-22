# GitHub Activity CLI

Sample solution for the [github-user-activity](https://roadmap.sh/projects/github-user-activity) challenge from roadmap.sh.

## How to run

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



