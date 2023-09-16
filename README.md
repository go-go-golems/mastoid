# mastoid - download mastodon threads

![](https://img.shields.io/github/license/go-go-golems/mastoid)
![](https://img.shields.io/github/actions/workflow/status/go-go-golems/mastoid/push.yml?branch=main)

mastoid is a command line tool for downloading and rendering Mastodon threads. It allows you to easily archive threads for your notes or share them in a variety of formats.

![image](https://github.com/go-go-golems/mastoid/assets/128441/68e64a9b-b924-4c37-99c8-519fd735259c)


## Features

- Download a thread as structured data in JSON, CSV, etc using `mastoid thread`
- Render a thread in multiple formats like Markdown, HTML, JSON, etc using `mastoid render`
- Transform structured output using the glazed library

## Installation

To install the `mastoid` command line tool with homebrew, run:

  ```bash
  brew tap go-go-golems/go-go-go
  brew install go-go-golems/go-go-go/mastoid
  ```

To install the `mastoid` command using apt-get, run:

  ```bash
  echo "deb [trusted=yes] https://apt.fury.io/go-go-golems/ /" >> /etc/apt/sources.list.d/fury.list
  apt-get update
  apt-get install mastoid
  ```

To install using `yum`, run:

  ```bash
  echo "
  [fury]
  name=Gemfury Private Repo
  baseurl=https://yum.fury.io/go-go-golems/
  enabled=1
  gpgcheck=0
  " >> /etc/yum.repos.d/fury.repo
yum install mastoid
```

To install using `go get`, run:

```bash
go get -u github.com/go-go-golems/mastoid/cmd/mastoid
```

Finally, you can install by downloading the binaries straight from [github](https://github.com/go-go-golems/mastoid/releases).

```
go get github.com/go-go-golems/mastoid
```

## Usage

### Registering mastoid

Before using mastoid, you need to register it as an application on your Mastodon server. This allows mastoid to authenticate with your account and access private threads.

Register mastoid using the `register` command. Provide the Mastodon server URL with `--server` (default is https://hachyderm.io):

```
mastoid register --server https://example.com
2:19PM INF Registering app...
2:19PM INF App registration successful!

2:19PM INF Authorizing app...
Access Token: 
2:20PM INF Name Name=mastoid
2:20PM INF App authorization successful!
```

This will initiate an OAuth authorization flow. You will need to visit the given URL in your browser, log in to your Mastodon account, and paste the provided access code back into the mastoid CLI.

Once authorized, mastoid will save the OAuth credentials to `~/.mastoid/config.yaml`.

### Downloading a thread

Use the `thread` command to download a thread as structured data. Provide the status URL or ID of the initial post in the thread:

```
mastoid thread -s https://example.com/@user/12345
...
```

This will retrieve all posts in the thread, not just replies to the given post.

By default, the output is a human-readable table. To specify another format, use the `--output` option:

```
mastoid thread -s https://example.com/@user/12345 --output json
[
{
  "Author": "mnl",
  "Content": "\u003cp\u003ezum glueck kann ich den frust wegcoden...\u003c/p\u003e",
  "CreatedAt": "2023-09-16T10:18:58.203Z",
  "Depth": 0,
  "ID": "111074314695602779",
  "InReplyToID": null,
  "SiblingIndex": 0
}
]

```

To only output certain fields, use `--fields`:

```
mastoid thread -s https://example.com/@user/12345 --fields ID,Author,Content 
+--------------------+--------+--------------------------------------------------+
| ID                 | Author | Content                                          |
+--------------------+--------+--------------------------------------------------+
| 111074314695602779 | mnl    | <p>zum glueck kann ich den frust wegcoden...</p> |
+--------------------+--------+--------------------------------------------------+
```

### Rendering a thread

Instead of structured data, you can render the thread in a more readable format using the `render` command:
- Markdown (default)
- HTML
- Plain text
- JSON (tree structure, not tabular)

In Markdown, the thread structure is shown using nested blockquotes:

```
❯ mastoid render -s https://hachyderm.io/@mnl/111039612582661457

> Author: mnl (2023-09-10 07:13:45.988 +0000 UTC)
> URL: https://hachyderm.io/@mnl/111039612582661457
> Author URL: https://hachyderm.io/@mnl
> 
> I want to peek poke in a windows exe and see if I can find some boogers. I last did that 15 years ago with ida pro, what do the cool kids use these days? I’m on Linux, if that matters. 
> Besides the exe, interested in decoding a binary protocol. I’m pretty sure chatgpt will slay for that, but are there some cool tools I can look at there too? Always rolled my own.
> #reverseEngineering #security #crackThePlanet
> 
> > > > Author: XXX@mastodon.social (2023-09-10 08:53:24 +0000 UTC)
> > > > 
> > > > @mnl absolutely no experience with this but not mentioned yet was Frida
> > > > 
> > > Author: XXX@mastodon.social (2023-09-10 07:29:02 +0000 UTC)
> > > 
> > > @mnl IDA Pro ($$$) is still king. Ghidra (free) and Binary Ninja ($) are getting there.   
> > > Automated protocol reversing: This is still hard.
> > > https://github.com/techge/PRE-list
> > > 
> > Author: XXX@chaos.social (2023-09-10 07:18:04 +0000 UTC)
> > 
> > @mnl The main modern equivalent of IDA is called Ghidra and was released by NSA around 2019. It supports debugging on both Linux and Windows, and has WinDbg integration.
> > 
> Author: XXXn@pleroma.m68k.church (2023-09-10 09:46:00.14 +0000 UTC)
> @XXX @mnl I wouldn't call them equivalent, they are slightly different and complete each other.
> 
> Author: XXX@chaos.social (2023-09-10 09:48:39 +0000 UTC)
> 
> @XXXn @mnl Let's say "cultural equivalent". Ghidra beats IDA among hacktivists because it's free. In the old days we shared licenses - one company used a stuffed animal as the token.
> 
> Author: mnl (2023-09-10 07:15:34.646 +0000 UTC)
> 
> Things I might want to do to poke some boogers: break on network comm calls in a debugger, what’s the best debugging combo for binaries without symbols and wine ? #reverseEngineering
> 
> Author: XXX@merveilles.town (2023-09-15 23:58:26 +0000 UTC)
> 
> @mnl You might want to use Frida for that.  I know @XXX uses it for some cool protocol RE stuff.
> 
> Author: XXX@merveilles.town (2023-09-16 07:59:57 +0000 UTC)
> 
> @XXX @mnl not sure how well Frida and Wine work together, so far I've always used it natively
> 
```

The text output shows the thread structure using nested indentation:

```
❯ mastoid render -s https://hachyderm.io/@mnl/111039612582661457 --output text

+ Author: mnl (2023-09-10 07:13:45.988 +0000 UTC)
| URL: https://hachyderm.io/@mnl/111039612582661457
| Author URL: https://hachyderm.io/@mnl
| 
| I want to peek poke in a windows exe and see if I can find some boogers. I last did that 15 years ago with ida pro, what do the cool kids use these days? I’m on Linux, if that matters. 
| Besides the exe, interested in decoding a binary protocol. I’m pretty sure chatgpt will slay for that, but are there some cool tools I can look at there too? Always rolled my own.
| #reverseEngineering #security #crackThePlanet

      + Author: XXX@mastodon.social (2023-09-10 08:53:24 +0000 UTC)
      | 
      | @mnl absolutely no experience with this but not mentioned yet was Frida
      
    + Author: XXX@mastodon.social (2023-09-10 07:29:02 +0000 UTC)
    | 
    | @mnl IDA Pro ($$$) is still king. Ghidra (free) and Binary Ninja ($) are getting there.   
    | Automated protocol reversing: This is still hard.
    | https://github.com/techge/PRE-list
    
  + Author: XXX@chaos.social (2023-09-10 07:18:04 +0000 UTC)
  | 
  | @mnl The main modern equivalent of IDA is called Ghidra and was released by NSA around 2019. It supports debugging on both Linux and Windows, and has WinDbg integration.
  
+ Author: XXX@pleroma.m68k.church (2023-09-10 09:46:00.14 +0000 UTC)
| @XXX @mnl I wouldn't call them equivalent, they are slightly different and complete each other.

+ Author: XXX@chaos.social (2023-09-10 09:48:39 +0000 UTC)
| 
| @XXXn @mnl Let's say "cultural equivalent". Ghidra beats IDA among hacktivists because it's free. In the old days we shared licenses - one company used a stuffed animal as the token.

+ Author: mnl (2023-09-10 07:15:34.646 +0000 UTC)
| 
| Things I might want to do to poke some boogers: break on network comm calls in a debugger, what’s the best debugging combo for binaries without symbols and wine ? #reverseEngineering

+ Author: XXX@merveilles.town (2023-09-15 23:58:26 +0000 UTC)
| 
| @mnl You might want to use Frida for that.  I know @XXX uses it for some cool protocol RE stuff.

+ Author: XXX@merveilles.town (2023-09-16 07:59:57 +0000 UTC)
| 
| @XXX @mnl not sure how well Frida and Wine work together, so far I've always used it natively
```

Specify the output format with `--output`. For example:

```
mastoid render -s https://example.com/@user/12345 --output html > thread.html
```

By default, the JSON output contains every field returned by the Mastodon API. This can be quite verbose.

## Contributing

Contributions are welcome! mastoid is an experimental project so expect breaking changes.
