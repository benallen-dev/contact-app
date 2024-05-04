# Contact.App

The contacts app from Hypermedia Systems.

Currently building out the Web 1.0 application in Go, but not done yet.

[https://hypermedia.systems/a-web-1-0-application/](https://hypermedia.systems/a-web-1-0-application/)

Once that's done I'll tag the 1.0 commit and start adding the HTMX features. Perhaps it'll be a handy starting point for someone else who doesn't want to reimplement the old school app themselves.

## Acknowledgements

### Missing.css

Instead of pulling CSS from the internet, this project self hosts v0.2.0 of [missing.css](https://github.com/bigskysoftware/missing), which is governed by a [BSD 2-clause license](https://www.tldrlegal.com/license/bsd-2-clause-license-freebsd). The full license text has been prepended as a comment to the CSS file.

The benefit of doing this is I noticed 2 redirects and a total of around 60ms before the css was loaded when I was using the CDN version. Now it's ~5ms and 55ms is a lot when your actual document loads in 2ms.

## BTW

The structure of this program is pretty terrible with how things depend on each other but writing a nicely-factored application wasn't exactly the goal here. No, seriously. It's pretty bad, and the validation is a chaotic mess. Frankly I'm surprised it works at all
