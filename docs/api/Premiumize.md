# Premiumize API

API docs: <https://app.swaggerhub.com/apis/premiumize.me/api>

## Get your API key

1. Login on <https://www.premiumize.me/>
2. Visit <https://www.premiumize.me/account>

Store the key in a variable in the shell:

```bash
PM_APIKEY=...
```

## Get your account info

```bash
curl --silent "https://www.premiumize.me/api/account/info?apikey=${PM_APIKEY}" | jq .
```

Should lead to something like:

```json
{
  "status": "success",
  "customer_id": "123",
  "premium_until": 456,
  "limit_used": 0,
  "space_used": 789
}
```

Using a bad API key doesn't lead to a non-"200 OK" response, but to a response body indicating the error:

```json
{
  "status": "error",
  "message": "customer_id and pin parameter missing or not logged in "
}
```

## Check if a torrent is in the cache

We use the torrent of the 720p version of "Night of the Living Dead" from 1968 from YTS. The movie is in the public domain, so it's legal to download, stream and share.

Magnet: `magnet:?xt=urn:btih:50B7DAFB7137CBECF045F78E8EFBE4AC1A90D139&dn=Night+of+the+Living+Dead+%281968%29+%5B720p%5D+%5BYTS.MX%5D&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2F9.rarbg.to%3A2710%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.ch%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.cyberia.is%3A6969%2Fannounce&tr=http%3A%2F%2Fp4p.arenabg.com%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce`

=> Info hash: `50B7DAFB7137CBECF045F78E8EFBE4AC1A90D139`

Request:

```bash
curl --silent --get --data-urlencode "items[]=50B7DAFB7137CBECF045F78E8EFBE4AC1A90D139" "https://www.premiumize.me/api/cache/check?apikey=${PM_APIKEY}" | jq .
```

Should lead to something like:

```json
{
  "status": "success",
  "response": [
    true
  ],
  "transcoded": [
    true
  ],
  "filename": [
    "Night Of The Living Dead (1968) [BluRay] [720p] [YTS.AM]"
  ],
  "filesize": [
    "828818888"
  ]
}
```

Note that there's no differentiation between the files in the torrent. If a hash is available, all files are available.

## Create direct download link

We're using the same magnet URL as above, for "Night of the Living Dead", which is in the public domain.

Request:

```bash
curl --silent --data-urlencode "src=magnet:?xt=urn:btih:50B7DAFB7137CBECF045F78E8EFBE4AC1A90D139&dn=Night+of+the+Living+Dead+%281968%29+%5B720p%5D+%5BYTS.MX%5D&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2F9.rarbg.to%3A2710%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.ch%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.cyberia.is%3A6969%2Fannounce&tr=http%3A%2F%2Fp4p.arenabg.com%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce" "https://www.premiumize.me/api/transfer/directdl?apikey=${PM_APIKEY}" | jq .
```

Should lead to something like:

```json
{
  "status": "success",
  "content": [
    {
      "path": "Night Of The Living Dead (1968) [BluRay] [720p] [YTS.AM]/www.YTS.AM.jpg",
      "size": "58132",
      "link": "https://jenny.plusnet.club/dl/abc012/345/678/def901.234/www.YTS.AM.jpg",
      "stream_link": null,
      "transcode_status": "not_applicable"
    },
    {
      "path": "Night Of The Living Dead (1968) [BluRay] [720p] [YTS.AM]/Night.Of.The.Living.Dead.1968.720p.BluRay.x264-[YTS.AM].mp4",
      "size": "828760756",
      "link": "https://lyanna.lmklol.link/dl/ghi567-890jkl/345/678/mno123.321/Night.Of.The.Living.Dead.1968.720p.BluRay.x264-%5BYTS.AM%5D.mp4",
      "stream_link": "https://lyanna.lmklol.link/dl/ghi567-890jkl/345/678/mno123.321/Night.Of.The.Living.Dead.1968.720p.BluRay.x264-%5BYTS.AM%5D.mp4",
      "transcode_status": "good_as_is"
    }
  ]
}
```

## Add a transfer

When a magnet is not cached you can't use the request to create a direct download link. You first have to add a transfer.

We're using the magnet URL of the torrent of the 1080p version for the same movie as above, "Night of the Living Dead", which is in the public domain.

```bash
curl --silent --data-urlencode "src=magnet:?xt=urn:btih:11EA02584FA6351956F35671962AB46354D99060&dn=Night+of+the+Living+Dead+%281968%29+%5B1080p%5D+%5BYTS.MX%5D&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2F9.rarbg.to%3A2710%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.ch%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.cyberia.is%3A6969%2Fannounce&tr=http%3A%2F%2Fp4p.arenabg.com%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce" "https://www.premiumize.me/api/transfer/create?apikey=${PM_APIKEY}" | jq .
```

Should lead to something like:

```json
{
  "status": "success",
  "type": "savetocloud",
  "id": "xyz321",
  "name": "Night of the Living Dead (1968) [1080p] [YTS.MX]"
}
```

## Get info about the transfer

```bash
curl --silent "https://www.premiumize.me/api/transfer/list?apikey=${PM_APIKEY}" | jq .
```

Should lead to something like:

```json
{
  "status": "success",
  "transfers": [
    {
      "id": "xyz321-VREg",
      "name": "Night Of The Living Dead (1968) [BluRay] [1080p] [YTS.AM]",
      "message": null,
      "status": "finished",
      "progress": 0,
      "folder_id": "-uvw654",
      "file_id": null,
      "src": "magnet:?xt=urn:btih:11ea02584fa6351956f35671962ab46354d99060&dn=Night+Of+The+Living+Dead+(1968)+BluRay+1080p+YTS+%5BYIFY%5D&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp://tracker.internetwarriors.net:1337&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969"
    }
  ]
}
```

Note that the response contains *all* past transfers (except for the ones that were directly turned into a direct download link via the `/directdl` endpoint).
