# RealDebrid API

API docs: <https://api.real-debrid.com>

## Get your API key

1. Login on <https://real-debrid.com>
2. Visit <https://real-debrid.com/apitoken>

Store the token in a variable in the shell:

```bash
RD_APITOKEN=...
```

## Get your user

```bash
curl --silent -H "Authorization: Bearer ${RD_APITOKEN}" "https://api.real-debrid.com/rest/1.0/user" | jq .
```

Should lead to something like:

```json
{
  "id": 123,
  "username": "yourusername",
  "email": "foo*****@example.com",
  "points": 456,
  "locale": "en",
  "avatar": "https://fcdn.real-debrid.com/images/forum/empty.png",
  "type": "premium",
  "premium": 789,
  "expiration": "2021-02-20T12:34:56.000Z"
}
```

## Check if a torrent is instantly available

We use the torrent of the 720p version of "Night of the Living Dead" from 1968 from YTS. The movie is in the public domain, so it's legal to download, stream and share.

Magnet: `magnet:?xt=urn:btih:50B7DAFB7137CBECF045F78E8EFBE4AC1A90D139&dn=Night+of+the+Living+Dead+%281968%29+%5B720p%5D+%5BYTS.MX%5D&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2F9.rarbg.to%3A2710%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.ch%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.cyberia.is%3A6969%2Fannounce&tr=http%3A%2F%2Fp4p.arenabg.com%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce`

=> Info hash: `50B7DAFB7137CBECF045F78E8EFBE4AC1A90D139`

Request:

```bash
curl --silent -H "Authorization: Bearer ${RD_APITOKEN}" "https://api.real-debrid.com/rest/1.0/torrents/instantAvailability/50B7DAFB7137CBECF045F78E8EFBE4AC1A90D139" | jq .
```

Should lead to something like:

```json
{
  "50b7dafb7137cbecf045f78e8efbe4ac1a90d139": {
    "rd": [
      {
        "1": {
          "filename": "Night.Of.The.Living.Dead.1968.720p.BluRay.x264-[YTS.AM].mp4",
          "filesize": 828760756
        }
      }
    ]
  }
}
```

Note that there's only one file instantly available, while the torrent contains two files (the other one is the movie poster as jpg file).

## Add magnet URL

We're using the same magnet URL as above, for "Night of the Living Dead", which is in the public domain.

```bash
curl --silent -H "Authorization: Bearer ${RD_APITOKEN}" --data-urlencode magnet="magnet:?xt=urn:btih:50B7DAFB7137CBECF045F78E8EFBE4AC1A90D139&dn=Night+of+the+Living+Dead+%281968%29+%5B720p%5D+%5BYTS.MX%5D&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2F9.rarbg.to%3A2710%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.ch%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.cyberia.is%3A6969%2Fannounce&tr=http%3A%2F%2Fp4p.arenabg.com%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce" "https://api.real-debrid.com/rest/1.0/torrents/addMagnet" | jq .
```

Should lead to something like:

```json
{
  "id": "ABC123",
  "uri": "https://api.real-debrid.com/rest/1.0/torrents/info/ABC123"
}
```

## Get info about the torrent

Use the same ID we got from the previous step where we added the magnet URL, "ABC123" in our example.

```bash
curl --silent -H "Authorization: Bearer ${RD_APITOKEN}" "https://api.real-debrid.com/rest/1.0/torrents/info/ABC123" | jq .
```

Should lead to something like:

```json
{
  "id": "DEF456",
  "filename": "Night Of The Living Dead (1968) [BluRay] [720p] [YTS.AM]",
  "original_filename": "Night Of The Living Dead (1968) [BluRay] [720p] [YTS.AM]",
  "hash": "50b7dafb7137cbecf045f78e8efbe4ac1a90d139",
  "bytes": 828818888,
  "original_bytes": 828818888,
  "host": "real-debrid.com",
  "split": 2000,
  "progress": 0,
  "status": "waiting_files_selection",
  "added": "2021-02-20T12:33:31.000Z",
  "files": [
    {
      "id": 1,
      "path": "/Night.Of.The.Living.Dead.1968.720p.BluRay.x264-[YTS.AM].mp4",
      "bytes": 828760756,
      "selected": 0
    },
    {
      "id": 2,
      "path": "/www.YTS.AM.jpg",
      "bytes": 58132,
      "selected": 0
    }
  ],
  "links": []
}
```

Note the `"status": "waiting_files_selection"`, and also that the torrent contains two files, while in the instant availability check we got the info that only the first one is instantly available.

## Add one of the files to the RealDebrid downloads

The `files=1` selects the first file, and we use the same ID from the previous step, `DEF456` in our example.

```bash
curl --silent -H "Authorization: Bearer ${RD_APITOKEN}" --data-urlencode files=1 "https://api.real-debrid.com/rest/1.0/torrents/selectFiles/DEF456" | jq .
```

## Get info about the torrent *again*

```bash
curl --silent -H "Authorization: Bearer ${RD_APITOKEN}" "https://api.real-debrid.com/rest/1.0/torrents/info/DEF456" | jq .
```

Should lead to something like:

```json
{
  "id": "DEF456",
  "filename": "Night.Of.The.Living.Dead.1968.720p.BluRay.x264-[YTS.AM].mp4",
  "original_filename": "Night Of The Living Dead (1968) [BluRay] [720p] [YTS.AM]",
  "hash": "50b7dafb7137cbecf045f78e8efbe4ac1a90d139",
  "bytes": 828760756,
  "original_bytes": 828818888,
  "host": "real-debrid.com",
  "split": 2000,
  "progress": 100,
  "status": "downloaded",
  "added": "2021-02-20T12:33:31.000Z",
  "files": [
    {
      "id": 1,
      "path": "/Night.Of.The.Living.Dead.1968.720p.BluRay.x264-[YTS.AM].mp4",
      "bytes": 828760756,
      "selected": 1
    },
    {
      "id": 2,
      "path": "/www.YTS.AM.jpg",
      "bytes": 58132,
      "selected": 0
    }
  ],
  "links": [
    "https://real-debrid.com/d/GHI789"
  ],
  "ended": "2019-10-28T21:18:16.000Z"
}

```

Now it's `"status": "downloaded"` and there's a URL in `"links"`. Also, `"selected"` marks which file we previously selected for the download.

## Generate HTTP download link

We're using the URL form the `"links"` from the previous response, so `"https://real-debrid.com/d/GHI789"` in our example.

```bash
curl --silent -H "Authorization: Bearer ${RD_APITOKEN}" --data-urlencode "link=https://real-debrid.com/d/GHI789" "https://api.real-debrid.com/rest/1.0/unrestrict/link" | jq .
```

Should lead to something like:

```json
{
  "id": "JKL012",
  "filename": "Night.Of.The.Living.Dead.1968.720p.BluRay.x264-[YTS.AM].mp4",
  "mimeType": "video/mp4",
  "filesize": 828760756,
  "link": "https://real-debrid.com/d/GHI789",
  "host": "real-debrid.com",
  "host_icon": "https://fcdn.real-debrid.com/123/images/hosters/realdebrid.png",
  "chunks": 32,
  "crc": 1,
  "download": "https://123.download.real-debrid.com/d/JKL012/Night.Of.The.Living.Dead.1968.720p.BluRay.x264-%5BYTS.AM%5D.mp4",
  "streamable": 1
}
```

The `"download"` URL is streamable and downloadable.
