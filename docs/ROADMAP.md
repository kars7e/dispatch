Dispatch Project Roadmap
========================

### How should I use this document?

This document provides description of items that the project decided to prioritize. This should serve as a reference
point for Dispatch contributors to understand where the project is going, and help determine if a contribution could be
conflicting with some longer term plans.

The fact that a feature isn't listed here doesn't mean that a patch for it will automatically be refused! We are always
happy to receive patches for new cool features we haven't thought about, or didn't judge to be a priority. Please
however understand that such patches might take longer for us to review.

### How can I help?

Short term objectives are listed in [Milestones](https://github.com/vmware/dispatch/milestones) which correspond to a
montly cadence.  Generally speaking any issue which has the label
[Feature](https://github.com/vmware/dispatch/labels/feature) or
[Enhancement](https://github.com/vmware/dispatch/labels/enhancement) are roadmap items. Our goal is to split down the
workload in such way that anybody can jump in and help. Please comment on issues if you want to work on it to avoid
duplicating effort! Similarly, if a maintainer is already assigned on an issue you'd like to participate in, pinging him
on GitHub to offer your help is the best way to go.

### How can I add something to the roadmap?

The roadmap is primarily maintaned by the Dispatch maintainers. We are aiming to be as transparent as possible through
this document and labeling issues. Because roadmap items can have broad effects on the Dispatch project any items added
or changed on this document will be the result of discussions among the maintainers and the author of a proposal.

If you have a proposal which you believe belongs on the roadmap, either raise it in a Issue with the tag "proposal".
This will start a more in depth discussion.

# 1. Features and Refactoring

## 1.1 Users and Authentication

Currently Dispatch does not contain a user database.  Although there is support for a login action against an IDP (GitHub),
this effectively only ensures that the user has a GitHub account.  Authentication and authorization is a big value
proposition for Dispatch, and the first step is to maintain a database of users.  This is a precursor to full blown
roles and authorization which are also on the roadmap.

An initial implementation should simply support a database and APIs for managing users.  Then the authorization check
should simply ensure that the entity making an API request is included in that user database.  Additionally, the user
metadata should be propogated through the system (associated with the request or event) for auditability.

## 1.2 Roles and Authorization

## 1.3 Applications or Groups

## 1.4 Image Management


## 1.1 Runtime improvements

We introduced [`runC`](https://runc.io) as a standalone low-level tool for container
execution in 2015, the first stage in spinning out parts of the Engine into standalone tools.

As runC continued evolving, and the OCI specification along with it, we created
[`containerd`](https://github.com/containerd/containerd), a daemon to control and monitor `runC`.
In late 2016 this was relaunched as the `containerd` 1.0 track, aiming to provide a common runtime
for the whole spectrum of container systems, including Kubernetes, with wide community support.
This change meant that there was an increased scope for `containerd`, including image management
and storage drivers.

Moby will rely on a long-running `containerd` companion daemon for all container execution
related operations. This could open the door in the future for Engine restarts without interrupting
running containers. The switch over to containerd 1.0 is an important goal for the project, and
will result in a significant simplification of the functions implemented in this repository.

## 1.2 Internal decoupling

A lot of work has been done in trying to decouple Moby internals. This process of creating
standalone projects with a well defined function that attract a dedicated community should continue.
As well as integrating `containerd` we would like to integrate [BuildKit](https://github.com/moby/buildkit)
as the next standalone component.

We see gRPC as the natural communication layer between decoupled components.

## 1.3 Custom assembly tooling

We have been prototyping the Moby [assembly tool](https://github.com/moby/tool) which was originally
developed for LinuxKit and intend to turn it into a more generic packaging and assembly mechanism
that can build not only the default version of Moby, as distribution packages or other useful forms,
but can also build very different container systems, themselves built of cooperating daemons built in
and running in containers. We intend to merge this functionality into this repo.