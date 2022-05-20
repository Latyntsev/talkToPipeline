# talkToPipeline

CLI tool for integration other CLI tools with telegram

## Example

```bash
export ttp_bot_token=<telegram_bot_id>:<telegram_bot_token>
export ttp_chat_id=<chat_id>

talkToPipeline | fastlane spaceauth | talkToPipeline --send_only true
```
