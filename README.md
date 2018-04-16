gltchook
========

A GitLab-to-TeamCity hook converter.

GitLab doesn't natively support TeamCity hooks in a good ways. TeamCity
requires an odd URL format that is painful to construct manually. This
project takes a default GitLab push hook and converts it into the
appropriate TeamCity API call to recheck the affected repositories.

