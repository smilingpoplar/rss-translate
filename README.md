![](https://github.com/smilingpoplar/rss-translate/actions/workflows/build.yml/badge.svg) [![](https://shields.io/badge/README-中文-blue.svg)](README.zh-CN.md)

Translate RSS into another language, view [the translated RSS](https://github.com/smilingpoplar/rss-translate/tree/rss/) here.

# Configuration

Configure RSS feeds in `config.yaml`, and the GitHub workflow will perform periodic translations.

## GitHub Workflow

(In the repository's Settings > Security > Secrets and variables > Actions > Repository secrets) Create the following secrets:

- GH_TOKEN: (In GitHub's Settings > Developer Settings > Personal access tokens > Tokens(classic)) [Create a token](https://github.com/settings/tokens/new), select the workflow permission, and then copy the token to GH_TOKEN.
- U_NAME and U_EMAIL: Git username and email.
