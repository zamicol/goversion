# goversion


Create a go module like version for your go app that includes semver with a directory hash.  

Version should be in this form:

	gitlab.com/zamicol/jsonflag v0.1.0 h1:2matHcJWnMecl8bfZnhEnn/9I6POyL9mK0DJFd5WHC8=


TODO:
- Get tags, probably using soemthing like this: https://github.com/golang/tools/tree/master/go/vcs
- Construct pseudo tags.  
- Zip project directories.
- Removing VCS components before zipping like go mod.  


WORKING:
- Hashing a zip file and getting base64 sum.  


Other files of interest:
https://github.com/golang/go/tree/master/src/cmd/go/internal/semver

For git:
https://github.com/golang/go/tree/master/src/cmd/go/internal/modfetch/codehost