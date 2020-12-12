# YASRP - Yet Another Simple Reverse Proxy

### Goals:
* Reverse HTTP Proxy
* Caching feature

### Timeline
- v0.1
Prototype to work with the ideas considering the goals and restrictions.

- v0.2
Recfactory to create some useful stuffs like config files, reorganize the main work on a package(__ReverseProxy__).  Pending work on Caching features.  

- v0.2.1
Created the CacheEngine model to allow many cache implementations like memory, filesystem, redis and etc.
    - Memory Cache implemented and working.
    - Dummy Cache implemented. This Cache Engine, does nothing. Work as example.

#### Known issues:
- Doesn't work(yet) with SSL/TLS targets;
- Can't handle chunked body;
- Doesn't have any timeout(Read/Write/Connect);
- Memory cache doesn't have entries limit. This could be dangerous.;
- Need to make the content-type filter to work.(cache only some types);
- There some issues when doing reverse proxy. Maybe improving in same way as **ProxyPass** and **ProxyPassReverse** on Apache httpd;

