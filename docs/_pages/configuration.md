---
title: "Configuration"
classes: wide
permalink: /configuration/
---

At the first execution, `Yo` will ask you to provide an [OpenAI API key](https://platform.openai.com/account/api-keys).

It will then generate your configuration in the file `~/.config/yo.json`, with the following structure:

```JS
{
  "openai_key": "sk-xxxxxxxxx",       // OpenAI API key
  "openai_proxy": "",                 // OpenAI API proxy (if needed)
  "openai_temperature": 0.2,          // OpenAI API temperature
  "user_default_prompt_mode": "exec", // user prefered prompt mode: "exec" (default) or "chat"
  "user_preferences": ""              // user preferences, expressed in natural language
}
```

In `REPL` mode, you can press `ctrl+s` to edit settings, they will be hot reloaded once you close your editor.
