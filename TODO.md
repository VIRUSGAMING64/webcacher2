Todo maybe


[] Add a config to replace arguments
ConfigExample:
t=w2_time
search="candela"
Execution example
if requests contains
https://www.google.com?t=12 then argument t is changed to actual time ->   https://www.google.com?t=1790514602794827369

https://www.google.com?search=hola+soy+german   -->   https://www.google.com?search=candela

[] Add web gui to show stats
[] Add Graph visualization to gui
[] Change memcache in GlobalCache to LRU cache with max size