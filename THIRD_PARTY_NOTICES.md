# Third-party notices

`christiangeorgelucas/html-markdown-tools` is deployed as a compiled Go
binary that statically links the open-source libraries listed below. Every
one is permissively licensed (MIT, BSD-2-Clause, BSD-3-Clause, or
Apache-2.0). There is no copyleft-licensed code (GPL, LGPL, AGPL, MPL,
CDDL, EPL, SSPL) anywhere in the deployed build closure — the full package
set `go list -deps ./...` reports for this module, reproduced below.

## Summary

| Package | License |
|---|---|
| [github.com/JohannesKaufmann/html-to-markdown](https://github.com/JohannesKaufmann/html-to-markdown) v1.6.0 | MIT |
| [codeberg.org/readeck/go-readability/v2](https://codeberg.org/readeck/go-readability) v2.1.2 (MIT-licensed maintained fork of [go-shiori/go-readability](https://github.com/go-shiori/go-readability), itself a line-by-line port of Mozilla's Readability.js) | MIT |
| [github.com/PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery) | BSD-3-Clause |
| [github.com/andybalholm/cascadia](https://github.com/andybalholm/cascadia) | BSD-2-Clause |
| [github.com/go-shiori/dom](https://github.com/go-shiori/dom) | MIT |
| [github.com/gogs/chardet](https://github.com/gogs/chardet) | MIT |
| [github.com/itlightning/dateparse](https://github.com/itlightning/dateparse) | MIT |
| [gopkg.in/yaml.v2](https://github.com/go-yaml/yaml) (pulled in only by html-to-markdown's unused `frontmatter` plugin file — this package never calls it) | Apache-2.0 |
| [golang.org/x/net](https://pkg.go.dev/golang.org/x/net), [golang.org/x/text](https://pkg.go.dev/golang.org/x/text) (Go team) | BSD-3-Clause |
| [google.golang.org/protobuf](https://pkg.go.dev/google.golang.org/protobuf) (Axiom's own generated bindings dependency) | BSD-3-Clause |

`golang.org/x/net`/`golang.org/x/text`/`google.golang.org/protobuf` each ship
an additional patent grant alongside the BSD-3-Clause license text; that grant
is an additional permission, not a condition of redistribution, and is not
separately reproduced below.

## Why v1, not v2, of html-to-markdown

JohannesKaufmann/html-to-markdown's `main` branch is a from-scratch v2
rewrite; its `LinkStyle` option (inline vs. referenced links) is present in
the `Options` struct but its setter is explicitly commented out in the v2
source with `// TODO: allow changing the link style once the render logic
is implemented`. Since configurable link style is one of this package's
documented options, ConvertToMarkdown wraps the frozen, stable v1 line
(tag `v1.6.0`) instead, where `LinkStyle` ("inlined"/"referenced") is fully
implemented, verified directly against the library's source.

## License texts

### MIT (JohannesKaufmann/html-to-markdown, go-readability, go-shiori/dom, gogs/chardet, itlightning/dateparse)

```
Copyright (c) contributors of the respective project

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

(Each project's actual LICENSE file carries its own copyright-holder name;
the permission text above is byte-for-byte identical across all five.)

### BSD-3-Clause (goquery, golang.org/x/net, golang.org/x/text, google.golang.org/protobuf)

```
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

  * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
  * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
  * Neither the name of the copyright holder nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
```

### BSD-2-Clause (andybalholm/cascadia)

```
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
```

### Apache License 2.0 (gopkg.in/yaml.v2)

See https://www.apache.org/licenses/LICENSE-2.0 for the full text.
