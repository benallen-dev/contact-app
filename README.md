# Contact.App

The contacts app from Hypermedia Systems.

Currently building out the Web 1.0 application in Go, but not done yet.

[https://hypermedia.systems/a-web-1-0-application/](https://hypermedia.systems/a-web-1-0-application/)

Once that's done I'll tag the 1.0 commit and start adding the HTMX features. Perhaps it'll be a handy starting point for someone else who doesn't want to reimplement the old school app themselves.

## TODO

- [ ] Delete doesn't persist properly because I suffer from pointer skill issues
- [ ] `flash` functionality to queue messages
- [ ] validation of field inputs

## BTW

The structure of this program is pretty terrible with how things depend on each other but 
writing a nicely-factored application wasn't exactly the goal here.
