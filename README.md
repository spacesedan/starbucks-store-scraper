# Starbucks Store Scraper

I am currently looking for a new startbucks to transfer to, and the current
process of getting a stores contact information is a little tedious, so to help
me out I wrote this program. It works by running a headless browser that scrapes
the start bucks store locator webpage and gathers store information and stores
that information in a json file. With this information available I am able to
contact store managers and inquire whether they are accepting transfers at this
time.

## How to use

Please use v2 located in the v2 directory, this initial build is a little naive
and I tried to brute force something to work. I'm just keeping it as a reminder
to put a little extra hard work to make something work the way I intended it
too.

How to build

```go
go build main.go
```

How to use

```go
./main -link=<starucks-store-loactor-link>
```

## Flags:

- link: This defines the link used to scrape store information from, without
  it the program will not run

- headless: This defines whether a a broswer will be visible when executing
  the program. This option is true by default, which means that the broswer will
  be running in the background and close when the scrape is finished.
