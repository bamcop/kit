const config = {
    "autoload": {
        "classmap": [
            "database/seeds",
            "autoload.classmap.1"
        ],
        "psr-4": {
            "App\\": "autoload.psr-4.App\\"
        }
    },
    "autoload-dev": {
        "psr-4": {
            "Tests\\": "autoload-dev.psr-4.Tests\\"
        }
    },
    "config": {
        "optimize-autoloader": true,
        "preferred-install": "dist",
        "sort-packages": true
    },
    "description": "The Laravel Framework.",
    "extra": {
        "laravel": {
            "dont-discover": []
        }
    },
    "keywords": [
        "framework",
        "keywords.1"
    ],
    "license": `
jsonObj := gabs.New()

jsonObj.Array("foo", "array")
// Or .ArrayP("foo.array")

jsonObj.ArrayAppend(10, "foo", "array")
jsonObj.ArrayAppend(20, "foo", "array")
jsonObj.ArrayAppend(30, "foo", "array")

fmt.Println(jsonObj.String())
`,
    "minimum-stability": "dev",
    "name": "laravel/laravel",
    "prefer-stable": true,
    "require": {
        "codingyu/ueditor": "^3.0",
        "encore/laravel-admin": "1.*",
        "fabpot/goutte": "^4.0",
        "fideloper/proxy": "^4.0",
        "freyo/flysystem-qcloud-cos-v5": "^2.0",
        "fruitcake/laravel-cors": "^2.0",
        "guzzlehttp/guzzle": "~6.0",
        "jacobcyl/ali-oss-storage": "2.1",
        "jenssegers/agent": "^2.6",
        "jenssegers/mongodb": "3.6.x",
        "laravel/framework": "^6.0",
        "laravel/tinker": "^1.0",
        "magicalex/write-ini-file": "^2.0",
        "overtrue/wechat": "~4.0",
        "php": "^7.2",
        "symfony/translation-contracts": "require.symfony/translation-contracts"
    },
    "require-dev": {
        "filp/whoops": "^2.0",
        "fzaninotto/faker": "^1.4",
        "mockery/mockery": "^1.0",
        "nunomaduro/collision": "^3.0",
        "phpunit/phpunit": "require-dev.phpunit/phpunit"
    },
    "scripts": {
        "post-autoload-dump": [
            "Illuminate\Foundation\ComposerScripts::postAutoloadDump",
            "scripts.post-autoload-dump.1"
        ],
        "post-create-project-cmd": [
            "scripts.post-create-project-cmd.0"
        ],
        "post-root-package-install": [
            "scripts.post-root-package-install.0"
        ]
    },
    "type": "type"
}
