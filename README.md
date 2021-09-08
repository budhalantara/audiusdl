# Golang Audius Downloader

## Build
```
go build
```

## Usage
Save to file with track name
```
./audiusdl https://audius.co/naskoedm/where-i-belong-239874
```
Save to file with custom file name
```
./audiusdl https://audius.co/naskoedm/where-i-belong-239874 > custom
```
Donwload and convert to mp3 with ffmpeg
```
./audiusdl https://audius.co/naskoedm/where-i-belong-239874 | ffmpeg -i pipe:0 -f mp3 -b:a 320K converted.mp3
```
