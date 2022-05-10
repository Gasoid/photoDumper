# photoDumper
Application downloads album photos from VK

features:
- oauth2
- exif metadata: dateTime, GPS coordinates


Run:
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
