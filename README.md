# Indexer

This is a project written in Go. Its main purpose is to index a source of data to a Full Text Search engine.

## Getting Started

This project will index data with the following format:
```
Message-ID: [some-id]
Date: [some-date]
From: [some-email-address]
To: [some-email-addresses]
Subject: [some-subject]
Mime-Version: [some-mime-version]
Content-Type: [some-content-type]
Content-Transfer-Encoding: [some-transfer-encoding]
X-From: [some-x-from]
X-To: [some-x-to]
X-cc: [some-x-cc]
X-bcc: [some-x-bcc]
X-Folder: [some-x-folder]
X-Origin: [some-x-origin]
X-FileName: [some-x-fileName]

[some-body-content]
```

To a `ZincSearch` engine.

## Prerequisites

ZincSearch server should be running. More info [here](https://docs.zincsearch.com/quickstart/#installation). 

## Usage
Follow the next steps:
- Go to `src` folder
- Execute the command:
    ```
    go build .
    ```
- Run the executable file:
    ```
    indexer.go.exe [path_to_database]
    ```
    for example:
    ```
    indexer.go.exe ./sample
    ```
*Note.- If your data is large, it will take a while to index all that data.*

- Once finished you can go to your ZincSearch server and look for your information under the `email_index` index.

*Note.- A `data.json` will also be created containing all your data in json format.*
