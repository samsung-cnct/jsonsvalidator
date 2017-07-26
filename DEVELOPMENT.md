# golang Development

###### Mikel Nelson (6/2017)
Example golang project set up

## Docs
### Rules
* Every file should have the License and Copyright at the top.
* Every pkg should have a `doc.go` file that has the `copyright - blank line - package description - then package <blah>` only. Then every other file in the package will not have the package description.
* All func, type, interfaces, etc should have a comment.
* `LICENSE` file at top level (essetially the same as the license and copyright in each go file).
* `go vet, golint, gosimple` should be run on all code and the results should be "clean".
* ***Viable, usable godoc is the goal***.

## Tests
### Rules
* Every `X.go` should have a corresponding `X_test.go` file (if applicable)
* `X_test.go` should provide full code coverage if possible.
* Do NOT create a `X_test.go` just to silent the `[no test file]` warnings.  Only create `X_test.go` with actual tests internal and a `not implemented` failure or skip.



## Dependencies
### Rules
* Dependency Management tools for golang are in flux.  Hopeful future is `go dep`.  
* For now, create and check in `./vendor` sub-ddirectory.
* Create the ./vendor dir via whatever means you like.  For exmple:
** Dependencies are managed with [Glide](https://github.com/Masterminds/glide).
** Add new golang package dependencies to `glide.yaml`.  Add package version information to `glide.lock`.
** NOTE: `glide.lock` has a timestamp that is updated whenever the dependencies are checked.  This should not cause a `git` checkin, however, editing the file should be checked in.  This is not automated.
** You should only have to rum `make dep` or `make dep-update`, and only if missing packages or they are out of date.  This is not automated.


## Development Helpers
### cobra
[cobra](https://github.com/spf13/cobra) pkg is used for the `main/cmd` packages.  New `cmd` files may be added by hand, or easier, with the `cobra add` command template add.
 
