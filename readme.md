# goso
a golang stack overflow notifier using growlnotify. mainly for my own use, but if someone else finds it useful, great. :P

### usage

* Assumes you have golang installed.
* Assumes you have `$GOPATH/bin/` on your `$PATH`. 
* Assumes you have [growlnotify](http://growl.info/downloads). 


The below script will search for SO questions tagged as "go" every 5 minutes, and notify you using growlnotify.

```
go get github.com/wfreeman/goso
goso -tags go
```

If you want to search for multiple tags, separate them by semicolon (and put them in quotes).

