# KubeLinterBot

KubeLinterBot calls KubeLinter with one ore more .yaml, interprets KubeLinter's output and posts a comment to the relevant commit via the github-API if there is a security-problem.

Installation:
1. If there is no file named kube-linter-bot-configuration.yaml in the KubeLinterBot-folder, copy the sample-file from the samples-folder.
Change reponame, username and bot.port according to your wishes.
Generate a safe secret and type it to hookSecret. You will need this later while installing the webhook.
Important: Don't change access-token manually here.
2. Run make build in /KubeLinterBot
3. Run ./kube-linter-bot
4. Authorize with github in your browser on http://localhost:7000
You can remove authorization in your github-account-settings.
5. Install a webhook (will be automated in future versions) here:
https://github.com/[your-username]/[your-repository]/settings/hooks
Set these options:
Payload URL: Your kubelinterbot-server address
Content type: application/json
Secret: The secret you generated for the configuration-file earlier.

Select "Let me select individual events" and then choose "Pull Requests" and "Pushes".

Make sure you "Active" is activated. 
Click "Add webhook". You're done!

Use:
If there is a push- or pull-request-event in the watched repository, kubelinterbot will automatically call kubelinter and post the results. 