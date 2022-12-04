# TODO
* [ ] use log.* functions instead of fmt.println where appropriate
* [x] only open the database for write when we do imports. in any other case `?mode=readonly`
* [ ] add tags and words deletion capabilities
* [ ] add exclude tags filtering (ie `-tags a,b,c -etags d` will list words from a,b,c tags and that dont have d tag)
* [ ] Make a proper module structure and ensure the dependencies get installed from `go.mod`
* [ ] Check if we can add loading database to memory to speed things up
* [ ] Split operations into multiple files to ease development and avoid conflicts in merge
* [ ] Add Github workflow to build binary releases