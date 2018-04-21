gltchook
========

A GitLab-to-TeamCity hook converter.

GitLab doesn't natively support TeamCity hooks in a good ways. TeamCity
requires an odd URL format that is painful to construct manually. This
project takes a default GitLab push hook and converts it into the
appropriate TeamCity API call to recheck the affected repositories.

The repository URL is extracted from the hook payload. The simplest way to
set this up for all repositories on a GitLab/TeamCity instance is to install
it once on the TeamCity server and set up a system push hook in GitLab. TC
will then be notified about all pushes regardless of project and will react
accordingly.
