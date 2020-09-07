# 
This repository is to demonstrate a fibonacci sequence service backed by a fast "cache" represented by postgresql.

Clone the repository to your local MacOS machine

git clone 

Verify that you have the following installed:

Docker Desktop

Docker Compose

Golang 1.15

Make (gnu make)


Further Instructions are available by entering "make <enter>" from the MacOS Terminal in the cloned repo directory.

###########Developer Notes ###########################

Given the size of the project, I decided that the testing code did not require a separate package naming convention
nor split out the test code from the project code.  Larger projects I would put test code into a test subdirectory.

Having looked at dockertest, I decided that I wanted to run the demo using docker_compose and let the postresql docker container get spun up by the Makefile for unit testing purposes.  It would be good to compare notes if ory/dockertest is architected the way we would like for testing withing golang.  The os.exit() at the end of the dockertest examples is a little concerning.  I would like it to be able to have a lifetime longer than the current function scope.  Docker Desktop has an API that I skimmed, be interesting to see how this 

In the interest of time, I decided not to try to get distribution and testing onto Windows Platform.  The Makefile will likely run on Linux if you comment out the /Application/Docker.app option.  Windows may be closer to operation, but that would require some more testing.

How I formatted the Makefile was intended to make it easy and step by step, with instructions labeled from #1 to #5.

I included a couple of algorithmic versions of the fibonacci sequence generator including an aproximation one that
runs in O(1),O(1).  The memoized version is using a go map, but there are other techniques including using a straight access array along with some pre-allocation if necessary.  Also could have better memory persistence instead of building the fibonacci sequence structure with each http call.

Not anything mentioned about concurrency and thread safety in the technical evaluation requirements email.  So the only 
thing we are using is the transactional nature of Postgresql to serialize web hits and their effects.








