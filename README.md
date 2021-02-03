# KubeLinterBot

KubeLinterBot calls KubeLinter with one ore more .yaml or .yml, interprets KubeLinter's output and posts a comment to the relevant commit via the github-API if there is a security-problem. There is a Kube-Linter-binary (version 0.1.6.) included (_kubelinter_-folder). Check https://github.com/stackrox/kube-linter/releases for updates.

## How to install
1. If there is no file named _kube-linter-bot-configuration.yaml_ in the _KubeLinterBot_-folder, copy the sample-file from the samples-folder.
->Change _repository.reponame_, _repository.username_ and _bot.port_ according to your wishes.<-remove
Generate a safe secret and add it as _user.secret_. You will need this later while installing the webhook(s).

2. Generate a personal access token here: https://github.com/settings/tokens. If you don't want to generate a token now and your server has a browser you can use, you can skip step 2 now and later do step 5 instead.
You will need to check the following options:
>1. repo: If you want to lint private repositories, check _repo_. If you only want to lint public repositories, check public_repo
>2. Not yet implemented: admin:repo_hook: If you want your webhook installed automatically, check this.
>3. Click "Generate token". Github will display it right away. Copy said token to _user.accessToken_. 

3. Run **make build** in /KubeLinterBot/
4. Run **./kube-linter-bot**
5. (Skip if you did step 2): Authorize with github in your browser on http://localhost:7000  
You can remove authorization in your github-account-settings.
6. For every repository you want KubeLinterBot to watch, install a webhook (will be automated in future versions) here:
https://github.com/[owner-of-repository-name]/[your-repository]/settings/hooks
and set these webhook-options:
>1. **Payload URL**: Your kubelinterbot-server address
>2. **Content type**: application/json
>3. **Secret**: The secret you generated for the configuration-file earlier.
>4. Select "Let me select individual events" and then choose "Pull Requests" and "Pushes".
>5. Make sure "Active" is activated. 
>6. Click "Add webhook". You're done!

## How to use
If there is a push- or pull-request-event in the watched repository, KubeLinterBot will automatically call KubeLinter, process its output and post the results as a commit-comment (in case of a push-event) or a review-comment requesting changes (in case of a pull-request). 

There are deployment-files for Kubernetes and a docker-file included. You can find the Kubernetes-files in the _deployment_-folder and the Dockerfile in the KubeLinterBot-folder.
