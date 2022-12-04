# TODO
* [ ] use log.* functions instead of fmt.println where appropriate
* [ ] only open the database for write when we do imports. in any other case `?mode=readonly`
* [ ] add tags and words deletion capabilities
* [ ] add exclude tags filtering (ie `-tags a,b,c -etags d` will list words from a,b,c tags and that dont have d tag)
* [ ] Make a proper module and ensure the dependencies get installed from `go.mod`