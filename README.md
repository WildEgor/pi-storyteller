## About
General idea use open-source models that can be used for story and/or for images generation (for example, Kandinsky and Ollama).

## Features
- [x] Generate images by user prompts;
- [x] Handle "job" (generate images) per "user" request (for telegram uses nickname); 
- [x] Priority queue for "jobs"; 
- [x] Limit for "bot" commands;

## Usage

- Download latest [release](https://github.com/WildEgor/pi-storyteller/releases) or build from source (change `.goreleaser.yml` and run `task build`);
- Put `bin` executable to any directory (default `/app`);
- Rename `config.example.yml` to `config.yml` and place near executable (also can specify path via `PI_STORYTELLER_CONFIG_PATH` env);
- Change `telegram.token` in `config.yml`. See [guide](https://core.telegram.org/bots/tutorial);
- Change `kandinsky.key` and `kandinsky.secret` in `config.yml`. See [guide](https://fusionbrain.ai/docs/en/doc/api-dokumentaciya/);
- Optional. Add telegram nicknames to `app.priority_list` in `config.yml` (unlimited image generations).
 
## References
- https://github.com/tvldz/storybook - using Python scripts and several models and show generated images at ink display;
- https://github.com/vitoplantamura/OnnxStream - lightweight wrapper for image generation;
- https://github.com/go-telegram-bot-api/telegram-bot-api - Telegram API client for Go;
- https://github.com/go-resty/resty - HTTP client for Go;
