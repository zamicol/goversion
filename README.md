# goversion

Create a go module like version for your go app that includes semver with a directory hash.  

# Why?

The go modules code is internal to the go tool, meaning it cannot be imported into a go project.  

Also, the current project's go module string can only be determinded by importing it by another project.  This tool enables discovery of that string without the need to do an external import. 


# Version
Version should be in this form:

	example.com/example v0.1.0 h1:2matHcJWnMecl8bfZnhEnn/9I6POyL9mK0DJFd5WHC8=

WORKING:
- Hashing a zip file and getting base64 sum.  
- Getting the module name from go.mod text.  

TODO:
- Get tags, probably using something like this: https://github.com/golang/tools/tree/master/go/vcs
- Construct pseudo tags.  
- Zip project directories.
- Removing VCS components before zipping like go mod.  


# Other
Other files of interest:
https://github.com/golang/go/tree/master/src/cmd/go/internal/semver

For git:
https://github.com/golang/go/tree/master/src/cmd/go/internal/modfetch/codehost