# A persistent storage server in Go

## Configuration
Each server instance is configured via a JSON file that's fed to it at startup.\
JSON configuration can have the following fields:
```json 
{
    "port": number,
    "internal_port": number,
    "log_file": string,
    "persistent": boolean,
    "slaves": number[],
    "methods": string[]
}
```
- `port` - a port to open the server on (see *API*)
- `internal_port` - a port for internal communication (see *Internal API*)
- `log_file` - a file for persistent logging (see *Persistent Logging*)
- `persistent` - if true, writes messages to `log_file` (see *Persistent Logging*)
- `slaves` - list of ports that will get requests resent to their `internal_port`
- `methods` - list of methods allowed on `port`

## API
This section describes API implemented on `port`. To read about API implemented on `internal_port`, see *Internal API*.\
These are endpoints for user interaction.

```http
GET /get?KEY HTTP/1.1
```
Gets a value by key `KEY`.\
If `methods` doesn't contain neither `"GET"` nor `"all"`, returns `403 Forbidden`.\
If that value exists, `200 OK` is returned and value is returned as content.\
Otherwise, `404 Not Found` is returned and key is returned as content.

```http
PUT /put?KEY HTTP/1.1
```
Puts a value in request content in key `KEY`.\
If `methods` doesn't contain neither `"PUT"` nor `"all"`, returns `403 Forbidden`.\
`200 OK` is always returned and the value is modified/created.
Then a /put_internal request is sent to every `internal_port` in `slaves`.

```http
DELETE /delete?KEY HTTP/1.1
```
Deletes a value by key `KEY`.\
If `methods` doesn't contain neither `"DELETE"` nor `"all"`, returns `403 Forbidden`.\
`200 OK` is always returned and the value is deleted if it exists.
Then a /delete_internal request is sent to every `internal_port` in `slaves`.

## Internal API
This section describes API implemented on `internal_port`. To read about API implemented on `port`, see *API*.\
These are endpoints for internal interaction of different servers, so they should not be exposed to the end user.

```http
GET /get_internal?KEY HTTP/1.1
PUT /put_internal?KEY HTTP/1.1
DELETE /delete_internal?KEY HTTP/1.1
```
These endpoints differ from their non-internal counterparts in 2 ways:

- They always do their storage interaction, i.e. they don't check `methods` config.
- They do not resend requests to `slaves`.

## Persistent logging
Persistent logging guarantees that on restart of the server data will be restored.
On startup, server instance reads data from `log_file`, that has new-line separated entries of the following format:
```
PUT {"key":"KEY","value":"VALUE"}
DELETE {"key":"KEY"}
```
> Note
> VALUE here is Base64 encoded.

On receiving a mutating request (i.e. PUT or DELETE), provided `persistent` is *true*, server writes new entries to `log_file`.
> Note
> Given a master and a slave server (master has `slaves` set to slave's `internal_port`), that means that mutating requests will result in duplicating `log_file` values.
> Since PUT and DELETE are idempotent, that shouldn't corrupt the data, but will result in redundant lines in `log_file`.

## Hot-reloading
Upon startup, server starts watching it's config file changes, and if it notices them, it will read it and have it's config values changed.\
If config file can't be read, hot-reloading feature gets disabled with no ability to reenable it (other than restaring the server altogether).\
Following values can be changed with hot-reloading:
- `slaves`
- `methods`
- `persistent`
