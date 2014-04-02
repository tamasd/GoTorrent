GoTorrent
=========

BitTorrent implementation written in Go (golang).

Implementation is in progress, a lot of stuff is missing.

Roadmap
=======

- peer list download from tracker
- download and seed a torrent
- upnp + nat-pmp
- bandwidth measuring/limiting
- tls
- dht
- magnet link
- uTP
- scraping
- http seeds
- ipv6
- merkle trees

Status
======

bencode
-------

Mostly complete implementation. Float handling is not implemented. The package has only smoke tests, failures and error
handling is not really tested. The biggest missing feature is to add support to struct tags override the struct names.

client
------

The torrent client itself. If you want to use this library just to download/seed torrents, this is what you are looking
for.

client/config
-------------

Stores configuration for the client to avoid circular dependencies with the tracker package

magnet
------

A struct is built on the top of the url, but retrieving the actual torrent won't be implemented until DHT is ready.

metainfo
--------

This package is actually just a struct which represents the torrent metainfo structure.

torrent
-------

Wrapper structure one the torrent file which is being downloaded/seeded.

tracker
-------

Torrent tracker manager.

util
----

Misc functions are here.
