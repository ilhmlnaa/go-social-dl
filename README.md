# 🐦 Twitter Downloader (Golang)

**Twitter Downloader** is a simple media downloader built with **Golang** to fetch videos and images directly from Twitter. The project is designed with flexibility in mind, making it easy to extend support to other platforms like **Instagram** in the future.

## Features

- ✅ Download media (video/image) directly from Twitter.
- ✅ Simple and clean API endpoint.
- ✅ Modular structure for easy future expansion (e.g., Instagram).
- ✅ Lightweight and fast — run as a single binary or with Docker.

## Installation

### Run Locally

Clone the repository and start the server using Go:

```bash
git clone https://github.com/username/twitter-downloader.git
cd twitter-downloader
go run main.go
```

### Build as Binary

```bash
go build -o twitter-dl
./twitter-dl
```

### Run with Docker

Build the Docker image:

```bash
docker build -t twitter-downloader .
```

Run the Docker container:

```bash
docker run -d -p 3000:3000 --name twitter-dl \
  -e TWITTER_AUTH_TOKEN=your_twitter_auth_token \
  -e TWITTER_CSRF_TOKEN=your_twitter_csrf_token \
  twitter-downloader
```

Once running, your API will be available at:

```
http://localhost:3000/dl?url=twitter_url
```

## Environment Variables

You can optionally create a `.env` file in the project root:

```env
TWITTER_AUTH_TOKEN=your_twitter_auth_token
TWITTER_CSRF_TOKEN=your_twitter_csrf_token
PORT=3000
```

| Variable             | Description                                                             |
| -------------------- | ----------------------------------------------------------------------- |
| `TWITTER_AUTH_TOKEN` | Optional, used for authenticated requests to the Twitter API.           |
| `TWITTER_CSRF_TOKEN` | Optional, used with the auth token if needed.                           |
| `PORT`               | Optional, defaults to 3000. You can change this to any port you prefer. |

## API Endpoint

Currently, the project provides a single API endpoint:

**GET** `/dl?url={twitter_url}`

This endpoint allows you to download media from a Twitter link.

### Example usage with curl:

```bash
curl "http://localhost:3000/dl?url=https://twitter.com/username/status/1234567890"
```

If you prefer to get the direct media URL, the API will return a JSON response:

```json
{
  "status": "success",
  "urls": ["https://pbs.twimg.com/media/Gr35T-DWMAAuVfL.jpg"]
}
```

## Roadmap

- [x] Instagram media downloader support
- [ ] Instagram media downloader support

## License

This project is licensed under the MIT License — feel free to use, modify, and contribute.

---

**💡 Note:** Twitter authentication tokens can expire or change frequently. Make sure to use a valid token from your current Twitter session.
