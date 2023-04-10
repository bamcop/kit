const config = {
    "description": "The Laravel Framework.",
    "extra": {
        "laravel": {
            "classmap": [
                "database/seeds",
                "extra.laravel.classmap.1"
            ],
            "code": "extra.laravel.code"
        }
    },
    "keywords": [
        "framework",
        "keywords.1"
    ],
    "license": `Copyright (c) 2019 Ashley Jeffs

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
`,
    "name": "laravel/laravel",
    "shell": `#!/bin/sh

git rev-list -1 HEAD > /dev/null 2>&1
if [ $? -eq 0 ]; then
    export GIT_COMMIT=$(git describe --dirty --always)
    export BUILD_TIME=$(date "+%F %T")
    export GO_VERSION=$(go version)
else
    echo "fail"
fi

gox -ldflags                                                            \
    "                                                                   \
    -X 'github.com/lifeng1992/build_info.gitCommit=${GIT_COMMIT}'       \
    -X 'github.com/lifeng1992/build_info.buildTime=${BUILD_TIME}'       \
    -X 'github.com/lifeng1992/build_info.goVersion=${GO_VERSION}'       \
    "                                                                   \
    "$@"
`,
    "type": "type"
}
