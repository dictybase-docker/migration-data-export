# Migration data export
This is a source repository for [docker](http://docker.io) image to run
all data export tasks of dictybase overhaul project.

## Exporting data 
### Prerequisites
* [Docker](https://www.docker.com/products/docker-app) and [docker-compose](https://docs.docker.com/compose/).
* A working connection to NU VPN.
### Running the exports
* Connect to NU VPN.
#### Stockcenter 
* Run compose
```
docker-compose -f dsc.yml up
```
It should take around few hours to complete. The script generally stay in the
terminal even after it got finished. So, close it by pressing Ctrl+c.

* Run make task
```
make create-tarball
```
Will export data from containers to a folder named `data` in the current directory.


## Command line (for understanding purpose only)
```
docker run --rm dictybase/migration-data-export -h
```
It will print out all the available scripts.

