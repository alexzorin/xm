# xm

This is a client implementation for the XMail Admin protocol.

[Documentation for this package](http://godoc.org/github.com/alexzorin/xm).

[Documentation for the protocol](http://www.xmailserver.org/Readme.html#description).

## Status

Basically, only `mluseradd` and `mluserdel` are properly implemented, but its trivial to add more commands
using `Cmd`. 

Responses that have multi-line, non-`RESSTRING` output are not implemented and not supported.

TLS authentication as per the XMail spec is not supported (though MD5 is used), but you can always use stunnel and pass
a TLS connection to `NewClient`.

Most of tests are done via integration, `go test -tags integration`.

## License

```
The MIT License (MIT)

Copyright (c) 2014 Alex Zorin

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```