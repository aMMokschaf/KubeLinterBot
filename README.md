# KubeLinterBot

KubeLinterBot calls KubeLinter with one ore more .yaml or .yml, interprets KubeLinter's output and posts a comment to the relevant commit via the github-API if there is a security-problem.

## How to install
1. If there is no file named _kube-linter-bot-configuration.yaml_ in the _KubeLinterBot_-folder, copy the sample-file from the samples-folder.
Change _repository.reponame_, _repository.username_ and _bot.port_ according to your wishes.
Generate a safe secret and add it as _webhook.secret_. You will need this later while installing the webhook.
Important: Don't change repository.user.accessToken manually here.
2. Run **make build** in /KubeLinterBot
3. Run **./kube-linter-bot**
4. Authorize with github in your browser on http://localhost:7000  
You can remove authorization in your github-account-settings.
5. Install a webhook (will be automated in future versions) here:
https://github.com/[your-username]/[your-repository]/settings/hooks
and set these webhook-options:
>1. **Payload URL**: Your kubelinterbot-server address
>2. **Content type**: application/json
>3. **Secret**: The secret you generated for the configuration-file earlier.

>4. Select "Let me select individual events" and then choose "Pull Requests" and "Pushes".
>5. Make sure "Active" is activated. 
>6. Click "Add webhook". You're done!

## How to use
If there is a push- or pull-request-event in the watched repository, KubeLinterBot will automatically call KubeLinter and post the results. 

There are deployment-files for Kubernetes and a docker-file included. You can find them in the _deployment_-folder.
