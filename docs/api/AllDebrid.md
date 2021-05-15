# AllDebrid API

API docs: <https://docs.alldebrid.com>

## Get your API key

1. Login on <https://alldebrid.com/>
2. Visit <https://alldebrid.com/apikeys>

Store the key in a variable in the shell:

```bash
AD_APIKEY=...
```

## Get your user

```bash
curl --silent "https://api.alldebrid.com/v4/user?agent=myAppName&apikey=${AD_APIKEY}" | jq .
```

Should lead to something like:

```json
{
  "status": "success",
  "data": {
    "user": {
      "username": "yourusername",
      "email": "foo@example.com",
      "isPremium": true,
      "isTrial": false,
      "isSubscribed": false,
      "premiumUntil": 123,
      "lang": "en",
      "preferedDomain": "com",
      "fidelityPoints": 456,
      "limitedHostersQuotas": {
        "fileal": 2000,
        "filespace": 2000,
        "filefactory": 3000,
        "gigapeta": 10000,
        "videobin": 10000,
        "isra": 3000,
        "rapidgator": 50000,
        "brupload": 3000,
        "uploadcloud": 2000,
        "userscloud": 3000,
        "wipfiles": 3000,
        "wdupload": 5000,
        "ddl": 50000,
        "flashbit": 5000,
        "anzfile": 3000
      }
    }
  }
}
```

Using a bad API key doesn't lead to a non-"200 OK" response, but to a response body indicating the error:

```json
{
  "status": "error",
  "error": {
    "code": "AUTH_BAD_APIKEY",
    "message": "The auth apikey is invalid"
  }
}
```

## Check if a torrent is instantly available

We use the torrent of the 720p version of "Night of the Living Dead" from 1968 from YTS. The movie is in the public domain, so it's legal to download, stream and share.

Magnet: `magnet:?xt=urn:btih:50B7DAFB7137CBECF045F78E8EFBE4AC1A90D139&dn=Night+of+the+Living+Dead+%281968%29+%5B720p%5D+%5BYTS.MX%5D&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2F9.rarbg.to%3A2710%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.ch%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.cyberia.is%3A6969%2Fannounce&tr=http%3A%2F%2Fp4p.arenabg.com%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce`

=> Info hash: `50B7DAFB7137CBECF045F78E8EFBE4AC1A90D139`

Request:

```bash
curl --silent --form "magnets[]=50B7DAFB7137CBECF045F78E8EFBE4AC1A90D139" "https://api.alldebrid.com/v4/magnet/instant?agent=myAppName&apikey=${AD_APIKEY}" | jq .
```

Should lead to something like:

```json
{
  "status": "success",
  "data": {
    "magnets": [
      {
        "magnet": "50B7DAFB7137CBECF045F78E8EFBE4AC1A90D139",
        "hash": "50b7dafb7137cbecf045f78e8efbe4ac1a90d139",
        "instant": true
      }
    ]
  }
}

```

Note that there's no differentiation between the files in the torrent. If a hash is available, all files are available.

## Upload magnet URL

We're using the same magnet URL as above, for "Night of the Living Dead", which is in the public domain.

```bash
curl --silent --form "magnets[]=magnet:?xt=urn:btih:50B7DAFB7137CBECF045F78E8EFBE4AC1A90D139&dn=Night+of+the+Living+Dead+%281968%29+%5B720p%5D+%5BYTS.MX%5D&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2F9.rarbg.to%3A2710%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.ch%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.cyberia.is%3A6969%2Fannounce&tr=http%3A%2F%2Fp4p.arenabg.com%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce" "https://api.alldebrid.com/v4/magnet/upload?agent=myAppName&apikey=${AD_APIKEY}" | jq .
```

Should lead to something like:

```json
{
  "status": "success",
  "data": {
    "magnets": [
      {
        "magnet": "magnet:?xt=urn:btih:50B7DAFB7137CBECF045F78E8EFBE4AC1A90D139&dn=Night+of+the+Living+Dead+%281968%29+%5B720p%5D+%5BYTS.MX%5D&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2F9.rarbg.to%3A2710%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.ch%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.cyberia.is%3A6969%2Fannounce&tr=http%3A%2F%2Fp4p.arenabg.com%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce",
        "hash": "50b7dafb7137cbecf045f78e8efbe4ac1a90d139",
        "name": "Night of the Living Dead (1968) [720p] [YTS.MX]",
        "filename_original": "",
        "size": 828818888,
        "ready": true,
        "id": 456
      }
    ]
  }
}
```

## Get the status of a torrent

Use the same ID we got from the previous step where we added the magnet URL, `"456"` in our example.

```bash
curl --silent --get --data-urlencode "id=456" "https://api.alldebrid.com/v4/magnet/status?agent=myAppName&apikey=${AD_APIKEY}" | jq .
```

Should lead to something like:

```json
{
  "status": "success",
  "data": {
    "magnets": {
      "id": 456,
      "filename": "Night Of The Living Dead (1968) [BluRay] [720p] [YTS.AM]",
      "size": 828818888,
      "hash": "50b7dafb7137cbecf045f78e8efbe4ac1a90d139",
      "status": "Ready",
      "statusCode": 4,
      "downloaded": 828818888,
      "uploaded": 828818888,
      "seeders": 0,
      "downloadSpeed": 0,
      "processingPerc": 0,
      "uploadSpeed": 0,
      "uploadDate": 1613926249,
      "completionDate": 1613926249,
      "links": [
        {
          "link": "https://uptobox.com/abc789",
          "filename": "Night.Of.The.Living.Dead.1968.720p.BluRay.x264-[YTS.AM].mp4",
          "size": 828760756,
          "files": [
            "Night.Of.The.Living.Dead.1968.720p.BluRay.x264-[YTS.AM].mp4"
          ]
        },
        {
          "link": "https://uptobox.com/def012",
          "filename": "www.YTS.AM.jpg",
          "size": 58132,
          "files": [
            "www.YTS.AM.jpg"
          ]
        }
      ],
      "type": "m",
      "notified": true,
      "version": 1
    }
  }
}
```

Different from RealDebrid, where after adding a magnet the status is "waiting_files_selection" at first, here we directly get the `"status": "Ready"` as well as links for each file. The links are no direct download links though, they're like OCH hoster links (like uploaded.net, Rapidgator, Mega.nz), so we take another step with AllDebrid to get to the final direct download.

## Generate HTTP download link

We're using the URL form the `"links"` from the previous response, so `"https://uptobox.com/abc789"` in our example.

```bash
curl --silent --get --data-urlencode "link=https://uptobox.com/abc789" "https://api.alldebrid.com/v4/link/unlock?agent=myAppName&apikey=${AD_APIKEY}" | jq .
```

Should lead to something like:

```json
{
  "status": "success",
  "data": {
    "link": "https://ghi.debrid.it/dl/345jkl/Night.Of.The.Living.Dead.1968.720p.BluRay.x264-%5BYTS.AM%5D.mp4",
    "host": "uptobox",
    "filename": "Night.Of.The.Living.Dead.1968.720p.BluRay.x264-[YTS.AM].mp4",
    "streaming": [],
    "paws": false,
    "filesize": 828760756,
    "streams": [
      {
        "quality": 480,
        "ext": "mp4",
        "filesize": 718147584,
        "name": "und",
        "link": "https://www12.uptostream.com/901pqr/480/0/video.mp4",
        "id": "480-und"
      },
      {
        "quality": 360,
        "ext": "mp4",
        "filesize": 451067904,
        "name": "und",
        "link": "https://www12.uptostream.com/901pqr/360/0/video.mp4",
        "id": "360-und"
      }
    ],
    "id": "678mno",
    "hostDomain": "uptobox.com"
  }
}
```

The `"link"` URL is streamable and downloadable.

## Delete a magnet

Use the same ID we got when we added the magnet URL, `"456"` in our example.

```bash
curl --silent --get --data-urlencode "id=456" "https://api.alldebrid.com/v4/magnet/delete?agent=myAppName&apikey=${AD_APIKEY}" | jq .
```

Should lead to something like:

```json
{
  "status": "success",
  "data": {
    "message": "Magnet was successfully deleted"
  }
}
```
