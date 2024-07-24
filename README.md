
## About
General idea use open-source models (like Ollama) that can be used for story generation and for images (for example, kandinsky).
Text generation must be generated using Raspberry Pi Zero 2 (need benchmarking) and using Go for image generation calls.
Also, Go server must render output html story files and send back to frontend.
Frontend simply show viewer (image carousel) with pagination and input field (if we want to send custom prompt) and "Generate" button (with cancel)

## Features
1) Handle "job" (generate images) per "user" request (for telegram uses nickname);
2) Priority queue for "jobs";
3) Limit for "bot" commands;
 
## References
- https://github.com/tvldz/storybook - using Python scripts and several models and show generated images at ink display;
- https://github.com/vitoplantamura/OnnxStream - lightweight wrapper for image generation;
- https://github.com/go-telegram-bot-api/telegram-bot-api - Telegram API client fo Go;
- https://github.com/go-resty/resty - HTTP client for Go;
