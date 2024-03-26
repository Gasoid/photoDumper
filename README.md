# photoDumper
[![CI](https://github.com/Gasoid/photoDumper/actions/workflows/ci.yml/badge.svg)](https://github.com/Gasoid/photoDumper/actions/workflows/ci.yml) 
[![codecov](https://codecov.io/gh/Gasoid/photoDumper/branch/main/graph/badge.svg?token=I5MSN7TKRL)](https://codecov.io/gh/Gasoid/photoDumper)


Tool downloads photos from VK albums (except from system albums)

![screen](screen.webp)


## Download page
[go to Download page](https://github.com/Gasoid/photoDumper/releases)


### Features:
- oauth2
- exif metadata: dateTime, GPS coordinates
- download all albums
- download a particular album

### Static files
- `tar xvfp <(curl -sL https://github.com/Gasoid/photoDumper/releases/download/1.1.0/build.zip)`
- or `go generate staticAssets.go`

### Run:
```bash
go run ./
```

## API Docs (swagger routines)
Regenerate docs:
```bash
swag init
```

Format swagger comments:
```bash
swag fmt
```

## Known issues
- albums with the same name will not be downloaded
- use `ulimit -n 1024` in order to fix 'too many open files' issue


### Tags
- Download all photos from vk account
- скачать все альбомы с вконтакте без регистрации и смс
- скачать фото с вк бесплатно
- приложение для скачивания всех фото из вк


**If you want to support my project you can buy something from this wishlist on amazon**
https://www.amazon.de/hz/wishlist/ls/H7GAFNFUC063?ref_=wl_share