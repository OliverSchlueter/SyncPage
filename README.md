# SyncPage Server

SyncPage Server is a lightweight web server that automatically fetches and serves the latest release files from specified GitHub repositories. It is designed to update static content regularly for each site registered with it. This setup is ideal for serving multiple static sites and updating them directly from GitHub.

## Features

- **Automatic Updates**: SyncPage fetches and unpacks the latest GitHub release for each site in a fixed interval
- **Multiple Sites**: Easily serve multiple sites by defining each site with a GitHub repository.
- **Error Handling**: Automatically handles and reports errors if assets or files are missing.
- **Static File Serving**: Efficiently serves static files from each registered site.


## Usage

Start the application once and it will generate a sites.json in the ./data/sites/ directory.

Example config:

```json
[
  {
    "Name": "docs",
    "Repo": {
      "Owner": "fancymcplugins",
      "Repo": "docs",
      "AuthToken": "YOUR_TOKEN"
    },
    "WorkflowName": "Build documentation",
    "ArtifactName": "docs",
    "FileName": "\\.zip$"
  }
]
```

the docs will be avaiable under `http://localhost:8080/docs`.