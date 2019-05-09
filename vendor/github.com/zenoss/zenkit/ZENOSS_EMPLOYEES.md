zenkit for Zenoss Employees
=====

Zenoss employees, please follow these steps when setting up the repo and jenkins 
jobs.

## Build and Repo Setup

### Create GitHub Repo
The zenkit creation script initialized your microservice as a git repo and
committed the initial code already. Now you need to add it to GitHub.

To create a new repo on Github, navigate to the [new repository
page](https://github.com/organizations/zenoss/repositories/new).  Enter the
same name you used to create the repository, since it's a Go package.  Also
a helpful description would be good. Make sure the repository will be private,
disable options to automatically add files to the repo, and click the *Create
repository* button.  Follow the instructions to _Push an existing repository
from the command line_.

Add the topic `zing` to the repo's list of topics so that it can be queryable.

There are several GitHub configuration items that are required for the new repo.
Browse to the _Settings_ tab in the new repo.  In the _Collaborations & teams_ 
settings, add the following teams and permissions:

| Team | Permission |
| ---- | ---------- |
| Administrators | Admin |
| Developers | Write |
| Automation | Write |
| Employees | Read |

In the _Branches_ settings, choose the _master_ branch for protection and
ensure that the options for "Require pull request reviews before merging"
and "Require status checks to pass before merging" are enabled.

### Create Jenkins Jobs
To create the suite of Jenkins Jobs, browse to the
[Microservice Job Builder](https://jenkins-eng.zenoss.io/job/micro-services/job/job_create/build?delay=0sec),
enter the name of the service repo, and click "Build".  This will create a
new suite of jobs for your service in the
[Microservices folder](https://jenkins-eng.zenoss.io/job/micro-services/).
