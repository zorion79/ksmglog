# ksmglog - get audit log records from [Kaspersky Secure Mail Gateway](https://www.kaspersky.ru/small-to-medium-business-security/mail-security-appliance)
[![Build Status](https://travis-ci.com/zorion79/ksmglog.svg?branch=master)](https://travis-ci.com/zorion79/ksmglog)
[![Coverage Status](https://coveralls.io/repos/github/zorion79/ksmglog/badge.svg)](https://coveralls.io/github/zorion79/ksmglog)

This library get last 500 records from Kaspersky Secure Mail Gateway (KSMG).

- Collect records from many servers in one request
- Every POST request provide timeout before run, because KSMG can return empty response.
- Can run as service and return `chan *Record`

## Install

`go get -u github.com/zorion79/ksmglog`

## Usage

- define options `Opts` with url's like `https://ksmg01/ksmg/en-US/cgi-bin/klwi`
- make service `NewService(opts Opts)`
- grab logs `GetLogs` return `type Record`
- get `service.Channel()` and grab only latest Records