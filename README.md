# courseware
A courseware web server with noms backend

To run:

1. install courseware-server:

cd courseware-server
go install

2. start noms (~/bin is where go installs the compiled executable):

~/bin/noms serve serve localhost:8000::courseware
Listening on port 8000...

3. start web server:

cd courseware-server
~/bin/courseware-server

The noms DB is listening in port 8000.
The web server is listening on port 8080.
