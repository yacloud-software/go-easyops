testcases:

(old==v1, with context_builder==cb)

* serialise cb->cb
* serialise old->cb
* serialise old->old
* deserialise old->cb
* deserialise cb->cb
* fork() old->cb
* fork() cb->cb

* client->server cb->old
* client->server cb->cb

* check contexts with/without -token=xx
