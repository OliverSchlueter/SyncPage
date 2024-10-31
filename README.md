# SyncPage Server

SyncPage Server is a lightweight web server that automatically fetches and serves the latest release files from specified GitHub repositories. It is designed to update static content regularly for each site registered with it. This setup is ideal for serving multiple static sites and updating them directly from GitHub.

## Features

- **Automatic Updates**: SyncPage fetches and unpacks the latest GitHub release for each site in a fixed interval
- **Multiple Sites**: Easily serve multiple sites by defining each site with a GitHub repository.
- **Error Handling**: Automatically handles and reports errors if assets or files are missing.
- **Static File Serving**: Efficiently serves static files from each registered site.
