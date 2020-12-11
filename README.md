# YASRP - Yet Another Simple Reverse Proxy

### Goals:
* Reverse HTTP Proxy
* Caching feature

### Timeline
- v0.1
Prototype to work with the ideas considering the goals and restrictions.

- v0.2
Recfactory to create some useful stuffs like config files, reorganize the main work on a package(__ReverseProxy__).  Pending work on Caching features.  
Known issues:
    - Doesn't work(yet) with SSL/TLS targets;
    - Can't handle chunked body;
    - Doesn't have any timeout(Read/Write/Connect)
