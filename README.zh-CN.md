![](https://github.com/smilingpoplar/rss-translate/actions/workflows/build.yml/badge.svg)

将 rss 翻译为中文，翻译后的 rss [见这里](https://github.com/smilingpoplar/rss-translate/tree/rss/)

# 配置

在 `config.yaml` 配置 rss 源，github workflow 将定时翻译

## github workflow

（在 repo 的 Settings > Security > Secrets and variables > Actions > Repository secrets）新建 secrets：

- GH_TOKEN：（在 github 的 Settings > Developer Settings > Personal access tokens > Tokens(classic)）[新建 token](https://github.com/settings/tokens/new)，新建时勾选 workflow 权限，拷贝 token 值到 GH_TOKEN
- U_NAME 和 U_EMAIL：git 用户名和邮箱
