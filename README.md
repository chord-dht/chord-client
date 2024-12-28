# chord-client

## Release

<https://github.com/chord-dht/chord-client/releases>

## Compile

Download and unzip the dist from: <https://github.com/chord-dht/chord-frontend/releases>

```shell
REPO="chord-dht/chord-frontend"
LATEST_RELEASE=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep "browser_download_url.*dist.tar.gz" | cut -d '"' -f 4)
curl -L -o dist.tar.gz $LATEST_RELEASE
tar -xzf dist.tar.gz

go build

./chord-client
```
